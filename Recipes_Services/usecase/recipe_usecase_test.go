package usecase_test

import (
	"context"
	"errors"
	"khanhanhtr/sample/model"
	"khanhanhtr/sample/model/mocks"
	"khanhanhtr/sample/usecase"
	"testing"
	"time"

	redismocks "khanhanhtr/sample/redis/mocks"

	logStack "github.com/diontr00/logstack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	sample_recipe *model.Recipe
)

func init() {
	sample_recipe = &model.Recipe{
		ID:           primitive.NewObjectID(),
		Name:         "Vietnamese Pho",
		Tags:         []string{"Vietname", "Asian Food"},
		Instructions: []string{"Cook Beef for 2 hours", "Make noodle"},
		Ingredients:  []string{"Beef", "Chilli"},
		PublishedAt:  time.Now(),
	}

}

func TestRecipeCreate_Unit(t *testing.T) {
	recipeRepo := mocks.NewRecipeRepository(t)
	redis := redismocks.NewClient(t)
	timeout := time.Second

	recipeUsecase := usecase.NewRecipeUsecase(recipeRepo, logStack.DefaultLogger(), timeout, redis)

	t.Run("Success", func(t *testing.T) {
		recipeRepo.On(
			"Create",
			mock.AnythingOfType("*context.timerCtx"),
			mock.AnythingOfType("*model.Recipe"),
		).Return(sample_recipe.ID, nil).Once()

		redis.On("Rpush", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("string"), mock.Anything).
			Return(nil).
			Once()

		id, err := recipeUsecase.Create(context.Background(), sample_recipe)
		assert.NoError(t, err)
		assert.NotEmpty(t, id)
		assert.Equal(t, id, sample_recipe.ID)
		recipeRepo.AssertExpectations(t)
	})

	t.Run("Failure", func(t *testing.T) {

		recipeRepo.On(
			"Create",
			mock.AnythingOfType("*context.timerCtx"),
			mock.AnythingOfType("*model.Recipe"),
		).Return("", errors.New("")).Once()

		id, err := recipeUsecase.Create(context.Background(), sample_recipe)
		assert.Error(t, err)
		assert.Empty(t, id)
		recipeRepo.AssertExpectations(t)
	})
}

func SkipTestRecipeGetAll_Unit(t *testing.T) {
	recipeRepo := mocks.NewRecipeRepository(t)
	redismocks := redismocks.NewClient(t)

	timeout := time.Second

	recipeUsecase := usecase.NewRecipeUsecase(recipeRepo, nil, timeout, redismocks)
	t.Run("Success", func(t *testing.T) {
		recipeRepo.On(
			"Get",
			mock.AnythingOfType("*context.timerCtx"),
		).Return([]model.Recipe{*sample_recipe}, nil).Once()

		recipe, err := recipeUsecase.Get(context.Background(), nil)
		assert.NoError(t, err)
		assert.NotEmpty(t, recipe)
		recipeRepo.AssertExpectations(t)
	})

	t.Run("Failure", func(t *testing.T) {
		recipeRepo.On(
			"Get",
			mock.AnythingOfType("*context.timerCtx"),
		).Return([]model.Recipe{}, errors.New("")).Once()

		recipe, err := recipeUsecase.Get(context.Background(), nil)
		assert.Error(t, err)
		assert.Empty(t, recipe)
		recipeRepo.AssertExpectations(t)
	})
}

func TestGetById_Unit(t *testing.T) {
	recipeRepo := mocks.NewRecipeRepository(t)

	timeout := time.Second

	recipeUsecase := usecase.NewRecipeUsecase(recipeRepo, nil, timeout, nil)
	t.Run("Success", func(t *testing.T) {
		recipeRepo.On(
			"GetById",
			mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("string"),
		).Return(*sample_recipe, nil).Once()

		recipe, err := recipeUsecase.GetById(context.Background(), sample_recipe.ID.String())
		assert.NoError(t, err)
		assert.NotEmpty(t, recipe)
		recipeRepo.AssertExpectations(t)
	})

	t.Run("Failure", func(t *testing.T) {
		recipeRepo.On(
			"GetById",
			mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("string"),
		).Return(*sample_recipe, errors.New("")).Once()

		recipe, err := recipeUsecase.GetById(context.Background(), sample_recipe.ID.String())
		assert.Error(t, err)
		assert.Empty(t, recipe)
		recipeRepo.AssertExpectations(t)
	})

}
