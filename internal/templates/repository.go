package templates

// RepositoryTemplate remains unchanged
const RepositoryTemplate = `package repository

import (
	"context"

	"{{.Config.Project.Module}}/internal/domain/entity"
)

// {{.Repository.Name}} defines the interface for {{.Repository.Entity}} {{.Repository.Name}}
type {{.Repository.Name}} interface {
	Create(ctx context.Context, {{.Repository.Entity | lower}} *entity.{{.Repository.Entity}}) error
	GetByID(ctx context.Context, id string) (*entity.{{.Repository.Entity}}, error)
	Update(ctx context.Context, {{.Repository.Entity | lower}} *entity.{{.Repository.Entity}}) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*entity.{{.Repository.Entity}}, error)
}
`

// PostgreSQL repository implementation template
const PostgresRepositoryTemplate = `package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
	
	"{{.Config.Project.Module}}/internal/domain/entity"
	"{{.Config.Project.Module}}/internal/domain/repository"
)

type {{.Repository.Name | lower}} struct {
	pg *Postgres
}

// New{{.Repository.Name}} creates a new PostgreSQL {{.Repository.Name}}
func New{{.Repository.Name}}(pg *Postgres) repository.{{.Repository.Name}} {
	return &{{.Repository.Name | lower}}{pg: pg}
}

func (r *{{.Repository.Name | lower}}) Create(ctx context.Context, {{.Repository.Entity | lower}} *entity.{{.Repository.Entity}}) error {
	err := r.pg.WithSchema(ctx, "public", func(conn *pgx.Conn) error {
		return nil
	})
	return err
}

func (r *{{.Repository.Name | lower}}) GetByID(ctx context.Context, id string) (*entity.{{.Repository.Entity}}, error) {
	{{.Repository.Entity | lower}} := &entity.{{.Repository.Entity}}{}
	err := r.pg.WithSchema(ctx, "public", func(conn *pgx.Conn) error {
		return nil
	})
	return {{.Repository.Entity | lower}}, err
}

func (r *{{.Repository.Name | lower}}) Update(ctx context.Context, {{.Repository.Entity | lower}} *entity.{{.Repository.Entity}}) error {
	err := r.pg.WithSchema(ctx, "public", func(conn *pgx.Conn) error {
		return nil
	})
	return err
}

func (r *{{.Repository.Name | lower}}) Delete(ctx context.Context, id string) error {
	err := r.pg.WithSchema(ctx, "public", func(conn *pgx.Conn) error {
		return nil
	})
	return err
}

func (r *{{.Repository.Name | lower}}) List(ctx context.Context, limit, offset int) ([]*entity.{{.Repository.Entity}}, error) {
	err := r.pg.WithSchema(ctx, "public", func(conn *pgx.Conn) error {
		return nil
	})
	return nil, err
}
`

// Updated MongoDB repository implementation template
const MongoRepositoryTemplate = `package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"{{.Config.Project.Module}}/internal/domain/entity"
	"{{.Config.Project.Module}}/internal/domain/repository"
)

type {{.Repository.Name | lower}} struct {
	collection *mongo.Collection
}

// New{{.Repository.Name}} creates a new MongoDB {{.Repository.Name}}
func New{{.Repository.Name}}(db *mongo.Database) repository.{{.Repository.Name}} {
	return &{{.Repository.Name | lower}}{
		collection: db.Collection("{{.Repository.Entity | snakeCase}}s"),
	}
}

func (r *{{.Repository.Name | lower}}) Create(ctx context.Context, {{.Repository.Entity | lower}} *entity.{{.Repository.Entity}}) error {
	{{.Repository.Entity | lower}}.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(ctx, {{.Repository.Entity | lower}})
	return err
}

func (r *{{.Repository.Name | lower}}) GetByID(ctx context.Context, id string) (*entity.{{.Repository.Entity}}, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	
	{{.Repository.Entity | lower}} := &entity.{{.Repository.Entity}}{}
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode({{.Repository.Entity | lower}})
	if err != nil {
		return nil, err
	}
	return {{.Repository.Entity | lower}}, nil
}

func (r *{{.Repository.Name | lower}}) Update(ctx context.Context, {{.Repository.Entity | lower}} *entity.{{.Repository.Entity}}) error {
	filter := bson.M{"_id": {{.Repository.Entity | lower}}.ID}
	update := bson.M{"$set": {{.Repository.Entity | lower}}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *{{.Repository.Name | lower}}) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (r *{{.Repository.Name | lower}}) List(ctx context.Context, limit, offset int) ([]*entity.{{.Repository.Entity}}, error) {
	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset))
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var {{.Repository.Entity | lower}}s []*entity.{{.Repository.Entity}}
	for cursor.Next(ctx) {
		{{.Repository.Entity | lower}} := &entity.{{.Repository.Entity}}{}
		if err := cursor.Decode({{.Repository.Entity | lower}}); err != nil {
			return nil, err
		}
		{{.Repository.Entity | lower}}s = append({{.Repository.Entity | lower}}s, {{.Repository.Entity | lower}})
	}
	return {{.Repository.Entity | lower}}s, nil
}
`
