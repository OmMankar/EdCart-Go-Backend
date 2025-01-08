package database

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connection(uri string) (Client *mongo.Client, err error) {

	if uri == "" {
		return nil, errors.New("set your 'MONGODB_URI' environment variable")
	}

	Client, err = mongo.Connect(context.TODO(), options.Client().
		ApplyURI(uri))
	if err != nil {
		fmt.Println("Error in connecting with Db")

		return nil, err
	}
	fmt.Println("Database connection successful")

	return Client, nil
}
