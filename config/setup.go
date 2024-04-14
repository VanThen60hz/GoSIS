package configs

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	_ "github.com/denisenkom/go-mssqldb"
)

func ConnectMongoDB() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(EnvMongoURI()))
	if err != nil {
		log.Fatal(err)
	}

	// ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")
	return client
}

// Client instance
var MongoDB *mongo.Client = ConnectMongoDB()

// getting database collections
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("apicompany").Collection(collectionName)
	return collection
}

func ConnectSqlServerDB() *sql.DB {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := sql.Open("sqlserver", EnvSqlServerURI())
	if err != nil {
		log.Fatal(err)
	}

	// ping the database
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to SQL Server")
	return db
}

// Client instance
var SqlServerDB *sql.DB = ConnectSqlServerDB()
