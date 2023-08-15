package model

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Recipe struct {
	ID           primitive.ObjectID `json:"id"           bson:"_id"`
	Name         string             `json:"name"         bson:"name"         validate:"required"`
	Tags         []string           `json:"tags"         bson:"tags"`
	Ingredients  []string           `json:"ingredients"  bson:"ingredients"  validate:"required"`
	Instructions []string           `json:"instructions" bson:"instructions" validate:"required"`
	PublishedAt  time.Time          `json:"publishedAt"  bson:"publishedAt"`
}

type FindOptions struct {
	Page  int64
	Limit int64
}

//go:generate mockery --name RecipeUseCase
type RecipeUseCase interface {
	Create(c context.Context, receipe *Recipe) (interface{}, error)
	Get(c context.Context, options *FindOptions) ([]Recipe, error)
	GetById(c context.Context, id string) (Recipe, error)
	UpdateOne(c context.Context, id string, recipe *Recipe) (int64, error)
	Delete(c context.Context, id string) (int64, error)
}

//go:generate mockery --name RecipeRepository
type RecipeRepository interface {
	Create(c context.Context, receipe *Recipe) (interface{}, error)
	Get(c context.Context) ([]Recipe, error)
	GetById(c context.Context, id string) (Recipe, error)
	UpdateOne(c context.Context, id string, recipe *Recipe) (int64, error)
	Delete(c context.Context, id string) (int64, error)
}
