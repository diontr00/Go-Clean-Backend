package config

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	vault "github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type (
	Env struct {
		App   AppEnv
		Vault VaultEnv
		Mongo MongoEnv
		Mq    MqEnv
	}

	MongoEnv struct {
		DBHost        string `env:"DB_HOST"`
		DBName        string `env:"DB_NAME"`
		MongoPort     string `env:"MONGO_PORT"`
		DBUri         string `env:"DB_URI"`
		mongoRootName string
		mongoRootPass string
		MongoPath     string
		ServiceName   string `env:"MONGO_SERVICE_NAME,default=recipemongo"`
	}

	VaultEnv struct {
		Vault_addr string `env:"VAULT_ADDRESS"`
		Vault_path string `env:"VAULT_PATH"`
	}

	AppEnv struct {
		RSSCert        []byte
		RSSKey         []byte
		Env            string        `env:"APP_ENV,default=dev"`
		ServerPort     string        `env:"SERVER_PORT,default=:8001"`
		ContextTimeout time.Duration `env:"CONTEXT_TIMEOUT,default=2s"`
		LogLocation    string        `env:"LOG_LOCATION"`
	}
	MqEnv struct {
		MqHost         string `env:"MQ_HOST,default=localhost:5672"`
		MqVhost        string `env:"MQ_VHOST"`
		MqUri          string
		NoRetries      int           `env:"MQ_RETRY_NO,default=5"`
		ConfirmMode    bool          `env:"MQ_CONFIRM_MODE,default=true"`
		ReInitDelay    time.Duration `env:"MQ_REINIT_DELAY,default=2s"`
		ReConnectDelay time.Duration `env:"MQ_RECONNECT_DELAY,default=2s"`
	}
)

//go:embed .env.sample.dev
var env_byte []byte

func NewEnv() *Env {
	env_file, err := envFile()
	if err != nil {
		log.Fatal(err.Error())

	}
	err = godotenv.Load(env_file)
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	defer os.Remove(env_file)

	env := &Env{}
	err = envconfig.Process(context.Background(), env)
	if err != nil {
		log.Fatal(err.Error())
	}

	if env.App.Env == "dev" {
		fmt.Println("Running App in development env ")
	}

	secret_engine, err := env.setupSecret()
	if err != nil {
		log.Fatal(err.Error())
	}

	var err_ch = make(chan error, 3)
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()

		err = env.loadMongo(ctx, secret_engine)
		if err != nil {
			err_ch <- err
		}
	}()

	go func() {
		defer wg.Done()
		err = env.loadApp(ctx, secret_engine)
		if err != nil {
			err_ch <- err
		}
	}()

	go func() {
		defer wg.Done()
		err = env.loadRabbit(ctx, secret_engine)
		if err != nil {
			err_ch <- err
		}
	}()

	wg.Wait()
	close(err_ch)

	for err := range err_ch {
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	return env
}

func (e *Env) loadMongo(ctx context.Context, secret_engine *vault.Client) error {
	mongo_root_name, err := getSecret(
		ctx,
		secret_engine,
		"MONGO_INITDB_ROOT_USERNAME",
		"web-server",
	)
	if err != nil {
		return err
	}
	e.Mongo.mongoRootName = mongo_root_name

	mongo_root_pass, err := getSecret(
		ctx,
		secret_engine,
		"MONGO_INITDB_ROOT_PASSWORD",
		"web-server",
	)
	if err != nil {
		return err
	}
	e.Mongo.mongoRootPass = mongo_root_pass

	mongo_path := fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s?authSource=admin",
		e.Mongo.ServiceName,
		e.Mongo.mongoRootName,
		e.Mongo.mongoRootPass,
		e.Mongo.MongoPort,
		e.Mongo.DBName,
	)

	e.Mongo.MongoPath = mongo_path

	return nil
}

func (e *Env) loadApp(ctx context.Context, secret_engine *vault.Client) error {

	rss_cert, err := getSecret(ctx, secret_engine, "CERT", "web-server/cert")
	if err != nil {
		return err
	}

	rss_key, err := getSecret(ctx, secret_engine, "KEY", "web-server/cert")
	if err != nil {
		return err
	}

	e.App.RSSCert = []byte(rss_cert)
	e.App.RSSKey = []byte(rss_key)
	return nil
}

func (e *Env) loadRabbit(ctx context.Context, secret_engine *vault.Client) error {
	username, err := getSecret(ctx, secret_engine, "RABBIT_MQ_USERNAME", "web-server")
	if err != nil {
		return err
	}
	password, err := getSecret(ctx, secret_engine, "RABBIT_MQ_PASSWORD", "web-server")
	if err != nil {
		return err
	}

	e.Mq.MqUri = fmt.Sprintf("amqp://%s:%s@%s/%s", username, password, e.Mq.MqHost, e.Mq.MqVhost)
	return nil

}

// Create temp file for godotenv to load embedded .env
func envFile() (string, error) {
	var err error
	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, ".env.tmp")
	err = os.WriteFile(tempFile, env_byte, 0600)

	if err != nil {
		return "", err
	}

	return tempFile, err
}

func (e *Env) setupSecret() (*vault.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	client, err := vault.New(
		vault.WithAddress(e.Vault.Vault_addr),
		vault.WithRequestTimeout(10*time.Second),
	)
	if err != nil {
		return nil, err
	}

	resp, err := client.Auth.UserpassLogin(
		ctx,
		os.Getenv("VAULT_USERNAME"),
		schema.UserpassLoginRequest{Password: os.Getenv("VAULT_PASSWORD")},
	)

	if err != nil {
		return nil, err
	}

	if err = client.SetToken(resp.Auth.ClientToken); err != nil {
		return nil, err

	}

	return client, nil
}

func getSecret(
	ctx context.Context,
	client *vault.Client,
	secret_key string,
	secret_path string,
) (string, error) {
	secret, err := client.Secrets.KvV2Read(ctx, secret_path, vault.WithMountPath("secret"))
	if err != nil {
		return "", err
	}

	data := secret.Data.Data
	value, ok := data[secret_key]
	if !ok {
		return "", fmt.Errorf("Cannot get %s from secret", secret_key)
	}

	strValue, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("Invalid format for %s: %v", secret_key, value)
	}

	return strValue, nil
}
