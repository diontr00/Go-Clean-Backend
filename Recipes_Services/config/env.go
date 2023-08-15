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
		Redis RedisEnv
		Auth0 Auth0Env
	}

	Auth0Env struct {
		Domain        string `env:"AUTH_DOMAIN"`
		Client_ID     string
		Client_Secret string
		AUTH_DB_NAME  string
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

	RedisEnv struct {
		RedisPath     string `env:"REDIS_PATH"`
		RedisPassword string `env:"REDIS_PASS"`
		RedisDB       int    `env:"REDIS_DB"`
	}

	VaultEnv struct {
		Vault_addr string `env:"VAULT_ADDRESS"`
		Vault_path string `env:"VAULT_PATH"`
	}

	AppEnv struct {
		RecipeCert       []byte
		RecipeKey        []byte
		EncryptKey       string
		CookieEncryptKey string
		Env              string        `env:"APP_ENV,default=dev"`
		ServerPort       string        `env:"SERVER_PORT,default=:8000"`
		ContextTimeout   time.Duration `env:"CONTEXT_TIMEOUT,default=2s"`
		LogLocation      string        `env:"LOG_LOCATION"`
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
		err = env.loadAuth0(ctx, secret_engine)
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

	encrypt_key, err := getSecret(ctx, secret_engine, "JWT_ENCRYPTION_KEY", "web-server")
	if err != nil {
		return err
	}

	e.App.EncryptKey = encrypt_key
	cookie_encrypt_key, err := getSecret(ctx, secret_engine, "COOKIE_ENCRYPTION_KEY", "web-server")
	if err != nil {
		return err
	}
	e.App.CookieEncryptKey = cookie_encrypt_key

	recipe_cert, err := getSecret(ctx, secret_engine, "CERT", "web-server/cert")
	if err != nil {
		return err
	}

	recipe_key, err := getSecret(ctx, secret_engine, "KEY", "web-server/cert")
	if err != nil {
		return err
	}

	e.App.RecipeCert = []byte(recipe_cert)
	e.App.RecipeKey = []byte(recipe_key)
	return nil
}

func (e *Env) loadAuth0(ctx context.Context, secret_engine *vault.Client) error {

	client_id, err := getSecret(ctx, secret_engine, "CLIENTID", "web-server/auth")
	if err != nil {
		return err
	}
	client_secret, err := getSecret(
		ctx,
		secret_engine,
		"CLIENTSECRET",
		"web-server/auth",
	)
	if err != nil {
		return err
	}

	auth_db_name, err := getSecret(ctx, secret_engine, "AUTH_DBNAME", "web-server/auth")
	if err != nil {
		return err
	}

	e.Auth0.Client_ID = client_id
	e.Auth0.Client_Secret = client_secret
	e.Auth0.AUTH_DB_NAME = auth_db_name
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
