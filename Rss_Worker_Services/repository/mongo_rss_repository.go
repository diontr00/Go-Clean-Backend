package repository

import (
	"context"
	"khanhanhtr/sample2/model"
	"khanhanhtr/sample2/mongo"
)

type mongoRssRepo struct {
	database   mongo.Database
	collection string
}

func (r *mongoRssRepo) InsertEntries(ctx context.Context, entries []*model.RSSEntry) (int, error) {

	collection := r.database.UseCollection(r.collection)
	var interfaceSlice []interface{}
	for _, entry := range entries {
		interfaceSlice = append(interfaceSlice, entry)
	}
	count, err := collection.InsertMany(ctx, interfaceSlice)
	if err != nil {
		return 0, err
	}
	return len(count), err
}

func NewMongoRssRepo(db mongo.Database, collection string) model.RssRepository {
	return &mongoRssRepo{
		database:   db,
		collection: collection,
	}
}
