package templates

// MongoDB Database Template
const MongoTemplate = `package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connection wraps MongoDB client and database
type Connection struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// NewConnection creates a new MongoDB database connection
func NewConnection(mongoURL string) (*Connection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	// Extract database name from URL or use default
	dbName := "{{.Database.Name}}"
	if dbName == "" {
		dbName = "app"
	}

	return &Connection{
		Client:   client,
		Database: client.Database(dbName),
	}, nil
}

// Disconnect closes the MongoDB connection
func (c *Connection) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return c.Client.Disconnect(ctx)
}

// Collection returns a collection from the database
func (c *Connection) Collection(name string) *mongo.Collection {
	return c.Database.Collection(name)
}
`
