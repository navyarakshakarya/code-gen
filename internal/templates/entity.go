package templates

const PostgresEntityTemplate = `package entity

import (
	"time"
)

// {{.Entity}} represents the {{.Entity}} entity for PostgreSQL
type {{.Entity}} struct {
	ID        int64     ` + "`json:\"id\" db:\"id\"`" + `
	Name      string    ` + "`json:\"name\" db:\"name\"`" + `
	CreatedAt time.Time ` + "`json:\"created_at\" db:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`json:\"updated_at\" db:\"updated_at\"`" + `
}

// Validate validates the {{.Entity}} entity
func (e *{{.Entity}}) Validate() error {
	// TODO: Implement validation logic
	return nil
}
`

const MongoEntityTemplate = `package entity

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// {{.Entity}} represents the {{.Entity}} entity for MongoDB
type {{.Entity}} struct {
	ID        primitive.ObjectID ` + "`json:\"id\" bson:\"_id,omitempty\"`" + `
	Name      string    ` + "`json:\"name\" bson:\"name\"`" + `
	CreatedAt time.Time          ` + "`json:\"created_at\" bson:\"created_at\"`" + `
	UpdatedAt time.Time          ` + "`json:\"updated_at\" bson:\"updated_at\"`" + `
}

// Validate validates the {{.Entity}} entity
func (e *{{.Entity}}) Validate() error {
	// TODO: Implement validation logic
	return nil
}
`
