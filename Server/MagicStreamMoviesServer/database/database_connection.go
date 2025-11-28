package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Connect() *mongo.Client {
	//Read environment variable file
	err := godotenv.Load(".env")

	//Check if we can read/load env file
	if err != nil {
		log.Println("WARNING: Unable to find .env file")
	}

	//Read URI (Connection String)
	MongoDB := os.Getenv("MONGODB_URI")

	if MongoDB == "" {
		log.Fatal("MONGODB_URI not set.")
	}

	fmt.Println("MongoDB URI: ", MongoDB)

	//Set Client
	clientOptions := options.Client().ApplyURI(MongoDB)

	//Connect to MONGODB
	client, err := mongo.Connect(clientOptions)

	if err != nil {
		return nil
	}

	return client
}

func OpenCollection(collectionName string, client *mongo.Client) *mongo.Collection {

	//read .env database name variable
	databaseName := os.Getenv("DATABASE_NAME")

	fmt.Println("DATABASE_NAME: ", databaseName)

	collection := client.Database(databaseName).Collection(collectionName)

	if collection == nil {
		return nil
	}

	return collection
}
