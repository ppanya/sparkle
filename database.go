package sparkle

import "context"

type Database interface {
	FindByID(ctx context.Context, Collection, ID string, value interface{}) error
	Save(ctx context.Context, Collection, ID string, entity interface{}) error
	DeleteByID(ctx context.Context, Collection, ID string) error
}
