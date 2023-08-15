package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"khanhanhtr/sample/model"
	"khanhanhtr/sample/mongo"
	"khanhanhtr/sample/redis"
	"time"

	logStack "github.com/diontr00/logstack"
)

var ErrInternal error = errors.New("Internal Error")
var ErrInvalidPageParm error = errors.New("Page and limit have to be greater than 0")
var ErrorParsingRecipes error = errors.New("Recipes Broken please try again")

type recipeUsecase struct {
	logger         *logStack.Logger
	repo           model.RecipeRepository
	contextTimeout time.Duration
	cache          redis.Client
}

func NewRecipeUsecase(
	recipeRepo model.RecipeRepository,
	logger *logStack.Logger,
	timeout time.Duration,
	cache redis.Client,
) model.RecipeUseCase {
	if logger == nil {
		logger = logStack.DefaultLogger()
	}

	return &recipeUsecase{
		repo:           recipeRepo,
		logger:         logger,
		contextTimeout: timeout,
		cache:          cache,
	}
}

func (ru *recipeUsecase) Get(
	c context.Context,
	opts *model.FindOptions,
) ([]model.Recipe, error) {
	ctx, cancel := context.WithTimeout(c, ru.contextTimeout)
	defer cancel()

	var start int64
	var stop int64
	if opts != nil {
		if opts.Limit <= 0 {
			stop = -1

		}
		if opts.Page <= 0 {
			start = 0
		}

		if opts.Page > 0 && opts.Limit > 0 {
			start = opts.Limit * (opts.Page - 1)
			stop = opts.Limit
		}

	} else {
		opts.Page = 0
		opts.Limit = -1
	}

	cache_key := "recipes"

	val, err := ru.cache.LRange(c, cache_key, start, stop).Result()

	if errors.Is(err, redis.NilErr) || len(val) == 0 {
		fmt.Println("Dev : Cache Miss")

		recipes, err := ru.repo.Get(ctx)

		if err != nil {
			err = ru.logError(err)
			return nil, err
		}

		for i := range recipes {

			recipeJSON, err := json.Marshal(recipes[i])
			if err != nil {
				ru.logger.Error(err.Error())
				continue

			}

			ru.cache.Rpush(c, cache_key, recipeJSON)
		}

		if stop == -1 {
			stop = int64(len(recipes))
		}
		return recipes[start:stop], nil
	}

	var recipes []model.Recipe
	for i := range val {
		var recipe model.Recipe
		data := []byte(val[i])

		if err := json.Unmarshal(data, &recipe); err != nil {
			ru.logger.Error(err.Error())
			continue
		}
		recipes = append(recipes, recipe)
	}

	ru.cache.Expire(c, cache_key, time.Minute*5)
	return recipes, nil
}

func (ru *recipeUsecase) Create(c context.Context, recipe *model.Recipe) (interface{}, error) {
	ctx, cancel := context.WithTimeout(c, ru.contextTimeout)
	defer cancel()

	id, err := ru.repo.Create(ctx, recipe)

	if err != nil {
		return "", ru.logError(err)
	}

	cache_key := "recipes"
	recipes, err := json.Marshal(recipe)
	if err != nil {
		ru.logger.Error(err.Error())
		return id, nil
	}

	ru.cache.Rpush(ctx, cache_key, recipes)

	return id, nil
}

func (ru *recipeUsecase) GetById(c context.Context, id string) (model.Recipe, error) {
	ctx, cancel := context.WithTimeout(c, ru.contextTimeout)
	defer cancel()
	recipe, err := ru.repo.GetById(ctx, id)
	if err != nil {
		err = ru.logError(err)
		return model.Recipe{}, err
	}
	return recipe, nil
}

func (ru *recipeUsecase) UpdateOne(
	c context.Context,
	id string,
	recipe *model.Recipe,
) (int64, error) {
	ctx, cancel := context.WithTimeout(c, ru.contextTimeout)
	defer cancel()
	count, err := ru.repo.UpdateOne(ctx, id, recipe)
	if err != nil {

		err = ru.logError(err)
		return 0, err
	}
	cache_key := "recipes"

	ru.cache.Expire(c, cache_key, time.Second*0)
	return count, nil

}

func (ru *recipeUsecase) Delete(c context.Context, id string) (int64, error) {
	ctx, cancel := context.WithTimeout(c, ru.contextTimeout)
	defer cancel()

	count, err := ru.repo.Delete(ctx, id)

	if err != nil {

		err = ru.logError(err)
		return 0, err
	}
	cache_key := "recipes"

	ru.cache.Expire(c, cache_key, time.Second*0)
	return count, nil

}

func (ru *recipeUsecase) logError(err error) error {
	var ErrKeyExist mongo.ErrKeyExist
	var ErrEmptyParam mongo.ErrEmptyParams
	var ErrTimeout mongo.ErrTimeout
	var ErrNilDocument mongo.ErrNoDocuments
	switch {
	case errors.As(err, &ErrKeyExist):
		return err
	case errors.As(err, &ErrEmptyParam):
		return err
	case errors.As(err, &ErrTimeout):
		return err
	case errors.As(err, &ErrNilDocument):
		return err

	default:
		ru.logger.Error(err.Error())
		return ErrInternal
	}
}
