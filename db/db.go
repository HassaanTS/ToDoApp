package db

import (
	"ToDoApp/todos"
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func BuildURI() string {
	dbHost := os.Getenv("MONGODB_HOSTNAME")
	dbPort := os.Getenv("MONGODB_PORT")
	dbName := os.Getenv("MONGODB_DBNAME")
	uri := fmt.Sprintf("mongodb://%s:%s/%s", dbHost, dbPort, dbName)
	return uri
}

func ConnectDB(uri string) (*mongo.Client, context.Context, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(),
		30*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return client, ctx, cancel, fmt.Errorf("error while connecting to mongo...%w", err)
	}
	return client, ctx, cancel, nil
}

func DisconnectDB(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {
	defer cancel()
	defer func() error {
		err := client.Disconnect(ctx)
		if err != nil {
			return fmt.Errorf("error while disconnecting from mongo...%w", err)
		}
		return nil
	}()
}

// INSERT
func InsertRecord(client *mongo.Client, ctx context.Context, dataBase, col string, doc interface{}) (*mongo.InsertOneResult, error) {
	collection := client.Database(dataBase).Collection(col)
	result, err := collection.InsertOne(ctx, doc)
	return result, err
}

// GET
func GetRecords(client *mongo.Client, ctx context.Context, dataBase, col string) ([]todos.ToDo, error) {
	var records []todos.ToDo
	var err error
	query, field := bson.D{}, bson.D{}
	collection := client.Database(dataBase).Collection(col)
	result, err := collection.Find(ctx, query, options.Find().SetProjection(field))
	if err != nil {
		return records, fmt.Errorf("error while querying for records...%w", err)
	}

	// get all records from cursor
	if err = result.All(ctx, &records); err != nil {
		return records, fmt.Errorf("error while extracting records from cursor...%w", err)
	}

	return records, nil
}

// UPDATE
func UpdateRecord(client *mongo.Client, ctx context.Context, dataBase, col string, filter, update interface{}) (result *mongo.UpdateResult, err error) {
	collection := client.Database(dataBase).Collection(col)
	result, err = collection.UpdateOne(ctx, filter, update)
	return
}

// DELETE
func DeleteRecord(client *mongo.Client, ctx context.Context, dataBase, col string, query interface{}) (result *mongo.DeleteResult, err error) {
	collection := client.Database(dataBase).Collection(col)
	result, err = collection.DeleteOne(ctx, query)
	return
}
