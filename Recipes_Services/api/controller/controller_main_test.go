package controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"khanhanhtr/sample/api/controller"
	"khanhanhtr/sample/config"
	"khanhanhtr/sample/model"
	"khanhanhtr/sample/model/mocks"
	"khanhanhtr/sample/usecase"
	"time"

	"github.com/go-playground/validator"
	"github.com/stretchr/testify/require"

	redismock "khanhanhtr/sample/redis/mocks"

	logStack "github.com/diontr00/logstack"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type global struct {
	recipecontroller   *controller.RecipeController
	redis_mock         *redismock.Client
	repo_mock          *mocks.RecipeRepository
	sample_recipe      model.Recipe
	sample_user_signin model.UserSigninRequest
	useruc_mock        *mocks.UserUseCase
	userController     *controller.UserController
}

var g *global

func TestMain(m *testing.M) {

	g = &global{}
	logger := logStack.DefaultLogger()
	validator := validator.New()
	redis := redismock.NewClient(&testing.T{})
	trans := config.NewTrans()
	reciperepo := mocks.NewRecipeRepository(&testing.T{})

	recipeusecase := usecase.NewRecipeUsecase(reciperepo, logger, time.Second*2, redis)
	userusecase := mocks.NewUserUseCase(&testing.T{})

	sample_recipe := model.Recipe{
		ID:           primitive.NewObjectID(),
		Name:         "Com tam",
		Tags:         []string{"vietnamese food"},
		Ingredients:  []string{"something1", "something2"},
		Instructions: []string{"instruction1", "instruction2"},
		PublishedAt:  time.Now(),
	}

	signin_request := model.UserSigninRequest{
		Username: "user1",
		Password: "123",
	}

	recipecontroller := &controller.RecipeController{
		Logger:        logger,
		Validator:     validator,
		Translator:    trans,
		RecipeUseCase: recipeusecase,
	}

	userControler := &controller.UserController{
		Logger:      logger,
		Validator:   validator,
		Translator:  trans,
		UserUseCase: userusecase,
	}

	g.recipecontroller = recipecontroller
	g.sample_recipe = sample_recipe
	g.sample_user_signin = signin_request
	g.redis_mock = redis
	g.repo_mock = reciperepo
	g.useruc_mock = userusecase
	g.userController = userControler

	m.Run()
}

func newHttpRequest(data interface{}, method, path string) *http.Request {
	t := &testing.T{}
	b, err := json.Marshal(data)
	require.NoError(t, err)

	reader := bytes.NewReader(b)

	req := httptest.NewRequest(method, path, reader)
	req.Header.Add("Content-Type", "application/json")
	return req
}
