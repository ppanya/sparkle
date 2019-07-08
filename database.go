package sparkle

import "context"

type TransactionalProvider interface {
	Begin(ctx context.Context, fn func(context TransactionalContext) error) error
}
type TransactionalContext interface {
	context.Context
	Rollback(ctx context.Context) error
	Commit(ctx context.Context) error
}

type Collection interface {
	FindByID(ctx context.Context, ID string, value interface{}) error
	FindOne(ctx context.Context, filter, value interface{}) error
	Save(ctx context.Context, ID string, entity interface{}) error
	DeleteByID(ctx context.Context, ID string) error
}

type Database interface {
	FindByID(ctx context.Context, Collection, ID string, value interface{}) error
	Save(ctx context.Context, Collection, ID string, entity interface{}) error
	DeleteByID(ctx context.Context, Collection, ID string) error
	Collection(name string) Collection
}
