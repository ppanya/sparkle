package sparklerepo

import (
	"context"
	"github.com/octofoxio/sparkle"
	entitiesv1 "github.com/octofoxio/sparkle/pkg/entities/v1"
	"github.com/rs/xid"
)

const (
	UserCollectionName = "user"
)

type UserRepository interface {
	Create(ctx context.Context, input *entitiesv1.UserRecord) (ID string, err error)
	FindByID(ctx context.Context, ID string) (*entitiesv1.UserRecord, error)
	UpdateByID(ctx context.Context, ID string, input *entitiesv1.UserRecord) error
}

type DefaultUserRepository struct {
	User sparkle.Collection
}

func (d *DefaultUserRepository) UpdateByID(ctx context.Context, ID string, input *entitiesv1.UserRecord) error {
	// ensure to
	// remove ID from input
	input.ID = nil
	err := d.User.Save(ctx, ID, input)
	return err
}

func NewDefaultUserRepository(db sparkle.Database) *DefaultUserRepository {
	User := db.Collection(UserCollectionName)
	return &DefaultUserRepository{User: User}
}

func (d *DefaultUserRepository) FindByID(ctx context.Context, ID string) (*entitiesv1.UserRecord, error) {
	var result entitiesv1.UserRecord
	err := d.User.FindByID(ctx, ID, &result)
	return &result, err
}

func (d *DefaultUserRepository) Create(ctx context.Context, input *entitiesv1.UserRecord) (string, error) {
	id := xid.New()
	err := d.User.Save(ctx, id.String(), input)
	return id.String(), err
}
