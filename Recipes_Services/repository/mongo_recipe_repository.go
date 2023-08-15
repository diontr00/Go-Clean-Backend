package repository

import (
	"context"
	"khanhanhtr/sample/model"
	"khanhanhtr/sample/mongo"

	"github.com/hashicorp/go-multierror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mongoRecipeRepo struct {
	database   mongo.Database
	collection string
}

func (r *mongoRecipeRepo) Create(ctx context.Context, recipe *model.Recipe) (interface{}, error) {
	id, err := r.database.UseCollection(r.collection).InsertOne(ctx, recipe)
	return id, err
}

func (r *mongoRecipeRepo) Get(
	ctx context.Context,
) ([]model.Recipe, error) {
	var errors error
	var recipe []model.Recipe

	cursor, cursor_err := r.database.UseCollection(r.collection).Find(ctx, bson.M{})
	defer func() {
		close_err := cursor.Close(ctx)
		if close_err != nil {
			errors = multierror.Append(errors, close_err)
		}
	}()

	decode_err := cursor.All(ctx, &recipe)

	errors = multierror.Append(errors, cursor_err, decode_err)
	return recipe, errors.(*multierror.Error).ErrorOrNil()
}

func (r *mongoRecipeRepo) GetById(ctx context.Context, recipeId string) (model.Recipe, error) {

	var recipe model.Recipe
	id, err := primitive.ObjectIDFromHex(recipeId)
	if err != nil {
		return model.Recipe{}, err

	}
	result := r.database.UseCollection(r.collection).FindOne(ctx, bson.M{
		"_id": id,
	})
	err = result.Decode(&recipe)
	if err != nil {
		return model.Recipe{}, err
	}
	return recipe, nil
}

func (r *mongoRecipeRepo) UpdateOne(
	ctx context.Context,
	recipeId string,
	recipe *model.Recipe,
) (int64, error) {
	var errors error
	objectId, id_error := primitive.ObjectIDFromHex(recipeId)

	result, coll_err := r.database.UseCollection(r.collection).UpdateOne(
		ctx,
		bson.M{"_id": objectId},
		bson.D{
			{"$set", bson.D{
				{"name", recipe.Name},
				{"instructions", recipe.Instructions},
				{"ingredients", recipe.Ingredients},
				{"tags", recipe.Tags},
			}},
		},
	)

	errors = multierror.Append(errors, id_error, coll_err)

	return result.ModifiedCount, errors.(*multierror.Error).ErrorOrNil()
}

func (r *mongoRecipeRepo) Delete(ctx context.Context, recipeId string) (int64, error) {
	var errors error
	object_id, obj_err := primitive.ObjectIDFromHex(recipeId)
	result, result_err := r.database.UseCollection(r.collection).DeleteOne(ctx, bson.M{
		"_id": object_id,
	})

	errors = multierror.Append(errors, obj_err, result_err)

	return result, errors.(*multierror.Error).ErrorOrNil()

}

func NewMongoRecipeRepo(db mongo.Database, collection string) model.RecipeRepository {
	return &mongoRecipeRepo{database: db, collection: collection}
}
