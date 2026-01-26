package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database holds the MongoDB client and database instance
type Database struct {
	Client *mongo.Client
	DB     *mongo.Database
}

// FindAll implements repositories.MongoRepository interface for MongoDB
func (d *Database) FindAll(ctx context.Context, collection string, filter interface{}, results interface{}) error {
	// Convert filter to bson.M if it's not already
	var bsonFilter bson.M
	if filter == nil {
		bsonFilter = bson.M{}
	} else if m, ok := filter.(bson.M); ok {
		bsonFilter = m
	} else {
		// Try to marshal/unmarshal if it's a different type
		bsonFilter = bson.M{}
	}

	cursor, err := d.DB.Collection(collection).Find(ctx, bsonFilter)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	return cursor.All(ctx, results)
}

// NewConnection creates a new MongoDB connection
func NewConnection() (*Database, error) {
	// Get connection string from environment or use default
	// Check DATABASE_URI first (used by config system), then MONGODB_URI for backward compatibility
	uri := os.Getenv("DATABASE_URI")
	if uri == "" {
		uri = os.Getenv("MONGODB_URI")
	}
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	// Get database name from environment or use default
	// Check DATABASE_NAME first (used by config system), then DB_NAME for backward compatibility
	dbName := os.Getenv("DATABASE_NAME")
	if dbName == "" {
		dbName = os.Getenv("DB_NAME")
	}
	if dbName == "" {
		dbName = "cronobserver"
	}

	// Create context with timeout for connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping to verify connection
	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Printf("Connected to MongoDB at %s, database: %s", uri, dbName)

	return &Database{
		Client: client,
		DB:     client.Database(dbName),
	}, nil
}

// Close gracefully closes the MongoDB connection
func (d *Database) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := d.Client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
	}

	log.Println("Disconnected from MongoDB")
	return nil
}
