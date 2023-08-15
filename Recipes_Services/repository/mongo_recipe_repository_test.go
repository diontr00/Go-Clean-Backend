package repository_test

import (
	"context"
	"errors"
	"khanhanhtr/sample/model"
	"khanhanhtr/sample/mongo"
	"khanhanhtr/sample/mongo/mocks"
	"khanhanhtr/sample/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var database_mocks *mocks.Database
var collection_mocks *mocks.Collection

type RecipeRepoTestSuite struct {
	suite.Suite
	database      mongo.Database
	collection    mongo.Collection
	sample_recipe *model.Recipe
}

func (r *RecipeRepoTestSuite) SetupSuite() {

	client := mocks.NewClient(r.T())
	client.EXPECT().UseDatabase("test_database").Return(mocks.NewDatabase(r.T())).Once()

	r.database = client.UseDatabase("test_database")

	_, ok := r.database.(*mocks.Database)
	require.True(r.Suite.T(), ok)
	database_mocks = r.database.(*mocks.Database)

	r.sample_recipe = &model.Recipe{
		ID:           primitive.NewObjectID(),
		Name:         "Com Tam Suon Bi Cha Trung",
		Tags:         []string{"Vietname", "Asian Food"},
		Instructions: []string{"Cook Beef for 2 hours", "Make noodle"},
		Ingredients:  []string{"Beef", "Chilli"},
		PublishedAt:  time.Now(),
	}

	require.NotEmpty(r.T(), r.database)
	client.AssertExpectations(r.T())
}

func (r *RecipeRepoTestSuite) TestCreate_Unit() {
	collection_name := "test_collection"

	database_mocks.EXPECT().
		UseCollection(collection_name).
		Return(mocks.NewCollection(r.T())).
		Times(3)

	r.collection = r.database.UseCollection(collection_name)

	_, ok := r.collection.(*mocks.Collection)
	require.True(r.T(), ok)
	collection_mocks = r.collection.(*mocks.Collection)

	recipe_repo := repository.NewMongoRecipeRepo(r.database, collection_name)

	r.T().Run("Success", func(t *testing.T) {

		collection_mocks.EXPECT().
			InsertOne(context.Background(), r.sample_recipe).
			Return(r.sample_recipe.ID, nil).Once()

		id, err := recipe_repo.Create(context.Background(), r.sample_recipe)
		assert.NoError(t, err)
		assert.NotEmpty(t, id)
	})

	r.T().Run("Failure", func(t *testing.T) {

		collection_mocks.EXPECT().
			InsertOne(context.Background(), r.sample_recipe).
			Return(nil, errors.New("")).Once()

		id, err := recipe_repo.Create(context.Background(), r.sample_recipe)
		assert.Error(t, err)
		assert.Empty(t, id)

	})
}

func TestRepository_Unit(t *testing.T) {
	suite.Run(t, new(RecipeRepoTestSuite))
}
