package mongo

import (
	"context"
	"errors"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//go:generate mockery --name Client
type Client interface {
	UseDatabase(string) Database
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	StartSession() (mongo.Session, error)
	UseSession(ctx context.Context, fn func(mongo.SessionContext) error) error
	Ping(ctx context.Context) error
}

//go:generate mockery --name Database
type Database interface {
	UseCollection(collName string) Collection
	Client() Client
}

//go:generate mockery --name Collection
type Collection interface {
	FindOne(ctx context.Context, filter interface{}) SingleResult
	InsertOne(ctx context.Context, document interface{}) (insertID interface{}, err error)
	InsertMany(ctx context.Context, documents []interface{}) (insertIDs []interface{}, err error)
	DeleteOne(ctx context.Context, filter interface{}) (deleteCount int64, err error)
	DeleteMany(ctx context.Context, filter interface{}) (deleteCount int64, err error)
	Find(ctx context.Context, filter interface{}, options ...*options.FindOptions) (Cursor, error)
	CountDocuments(
		ctx context.Context,
		filter interface{},
		options ...*options.CountOptions,
	) (int64, error)
	Aggregate(ctx context.Context, pipeline interface{}) (Cursor, error)
	UpdateOne(
		ctx context.Context,
		filter interface{},
		update interface{},
		options ...*options.UpdateOptions,
	) (*mongo.UpdateResult, error)
	UpdateMany(
		ctx context.Context, filter interface{},
		update interface{},
		options ...*options.UpdateOptions,
	) (*mongo.UpdateResult, error)
}

// go:generate mockery --name SingleResult
type SingleResult interface {
	Decode(val interface{}) error
	Err() error
}

// go:generate mockery --name Cursor
type Cursor interface {
	Close(ctx context.Context) error
	Next(ctx context.Context) bool
	Decode(val interface{}) error
	// Iterate the cursor and decode each document into the result
	All(ctx context.Context, result interface{}) error
}

// MongoDB Wrapper Concrete
type mongoClient struct {
	client *mongo.Client
}

type mongoDatabase struct {
	database *mongo.Database
}
type mongoCollection struct {
	collection *mongo.Collection
}

type mongoSingleResult struct {
	result *mongo.SingleResult
}

type mongoCursor struct {
	cursor *mongo.Cursor
}

type mongoSession struct {
	mongo.Session
}

// Handle dynamic values  and provide generic approach for setting the zero value of Go type
type nullawareDecoder struct {
	// hold reference to the default decoder provided by BSON codec lib
	defDecoder bsoncodec.ValueDecoder
	// hold the vlaue of specific go type when not explicitly set
	zeroValue reflect.Value
}

func (d *nullawareDecoder) DecodeValue(
	dctx bsoncodec.DecodeContext,
	vr bsonrw.ValueReader,
	val reflect.Value,
) error {
	if vr.Type() != bson.TypeNull {
		// If value is not null delegate to the default BSON codec
		return d.defDecoder.DecodeValue(dctx, vr, val)

	}
	if !val.CanSet() {
		return errors.New("Val not setable")
	}

	val.Set(d.zeroValue)

	return nil
}

func NewBsonNullableDecoder(def bsoncodec.ValueDecoder, zeroValue reflect.Value) *nullawareDecoder {
	return &nullawareDecoder{defDecoder: def, zeroValue: zeroValue}
}

func CheckMongoError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case mongo.IsDuplicateKeyError(err):
		return ErrKeyExist{err: err}
	case mongo.IsTimeout(err):
		return ErrTimeout{err: err}
	case mongo.IsDuplicateKeyError(err):
		return ErrKeyExist{err: err}
	case errors.Is(err, mongo.ErrEmptySlice):
		return ErrEmptyParams{err: err}
	default:
		return ErrInternal{err: err}
	}
}

// Client implementation
// ---------------------------------------------------------------------------------------------------------------------------------------

// Create a client and intialize it with connect
func NewClient(ctx context.Context, connection string) (Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connection))
	return &mongoClient{client: client}, err
}

func (mc *mongoClient) Ping(ctx context.Context) error {
	return mc.client.Ping(ctx, readpref.Primary())
}

func (mc *mongoClient) UseDatabase(dbName string) Database {
	db := mc.client.Database(dbName)
	return &mongoDatabase{database: db}
}

func (mc *mongoClient) StartSession() (mongo.Session, error) {
	session, err := mc.client.StartSession()
	return &mongoSession{session}, err
}

// Depreciate Newclient already initialize the connection
func (mc *mongoClient) Connect(ctx context.Context) error {
	return nil
}

func (mc *mongoClient) Disconnect(ctx context.Context) error {
	return mc.client.Disconnect(ctx)
}

func (mc *mongoClient) UseSession(ctx context.Context, fn func(mongo.SessionContext) error) error {
	return mc.client.UseSession(ctx, fn)
}

// Database implementation
// -------------------------------------------------------------------------------------------------------------------------------------

func (md *mongoDatabase) UseCollection(colName string) Collection {
	collection := md.database.Collection(colName)
	return &mongoCollection{collection: collection}
}

func (md *mongoDatabase) Client() Client {
	client := md.database.Client()
	return &mongoClient{client: client}
}

// Collection Implementation
// --------------------------------------------------------------------------------------------------------------------------------------

func (mc *mongoCollection) FindOne(ctx context.Context, filter interface{}) SingleResult {
	singleResult := mc.collection.FindOne(ctx, filter)
	return &mongoSingleResult{result: singleResult}
}

func (mc *mongoCollection) UpdateOne(
	ctx context.Context,
	filter interface{},
	update interface{},
	opts ...*options.UpdateOptions,
) (*mongo.UpdateResult, error) {
	return mc.collection.UpdateOne(ctx, filter, update, opts[:]...)
}

func (mc *mongoCollection) InsertOne(
	ctx context.Context,
	document interface{},
) (interface{}, error) {
	id, err := mc.collection.InsertOne(ctx, document)
	err = CheckMongoError(err)
	return id.InsertedID, err
}

func (mc *mongoCollection) InsertMany(
	ctx context.Context,
	document []interface{},
) ([]interface{}, error) {
	res, err := mc.collection.InsertMany(ctx, document)
	err = CheckMongoError(err)
	return res.InsertedIDs, err
}

func (mc *mongoCollection) DeleteOne(ctx context.Context, filter interface{}) (int64, error) {
	count, err := mc.collection.DeleteOne(ctx, filter)
	err = CheckMongoError(err)
	return count.DeletedCount, err
}

func (mc *mongoCollection) DeleteMany(ctx context.Context, filter interface{}) (int64, error) {
	count, err := mc.collection.DeleteMany(ctx, filter)
	err = CheckMongoError(err)
	return count.DeletedCount, err
}

func (mc *mongoCollection) Find(
	ctx context.Context,
	filter interface{},
	opts ...*options.FindOptions,
) (Cursor, error) {
	findResult, err := mc.collection.Find(ctx, filter, opts[:]...)
	err = CheckMongoError(err)
	return &mongoCursor{cursor: findResult}, err
}

func (mc *mongoCollection) Aggregate(ctx context.Context, pipeline interface{}) (Cursor, error) {
	aggregateResult, err := mc.collection.Aggregate(ctx, pipeline)

	return &mongoCursor{cursor: aggregateResult}, err
}

func (mc *mongoCollection) UpdateMany(
	ctx context.Context,
	filter interface{},
	update interface{},
	opts ...*options.UpdateOptions,
) (*mongo.UpdateResult, error) {

	return mc.collection.UpdateMany(ctx, filter, update, opts[:]...)
}

func (mc *mongoCollection) CountDocuments(
	ctx context.Context,
	filter interface{},
	opts ...*options.CountOptions,
) (int64, error) {
	return mc.collection.CountDocuments(ctx, filter, opts...)
}

// SingleResult Implementation
// ----------------------------------------------------------------------------------------------------------------------------------

func (sr *mongoSingleResult) Decode(v interface{}) error {
	return sr.result.Decode(v)
}
func (sr *mongoSingleResult) Err() error {
	return sr.result.Err()
}

// Cursor Implementation
// -----------------------------------------------------------------------------------------------------------------------------------

func (mc *mongoCursor) Close(ctx context.Context) error {
	return mc.cursor.Close(ctx)
}

func (mc *mongoCursor) Next(ctx context.Context) bool {
	return mc.cursor.Next(ctx)
}

func (mc *mongoCursor) Decode(v interface{}) error {
	return mc.cursor.Decode(v)
}

func (mc *mongoCursor) All(ctx context.Context, result interface{}) error {
	return mc.cursor.All(ctx, result)
}
