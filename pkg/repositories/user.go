package sparklerepo

import (
	"context"
	"github.com/octofoxio/sparkle"
	commonv1 "github.com/octofoxio/sparkle/pkg/common/v1"
	entitiesv1 "github.com/octofoxio/sparkle/pkg/entities/v1"
	"github.com/rs/xid"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, input *entitiesv1.UserRecord) (ID string, err error)
	FindByID(ctx context.Context, ID string) (*entitiesv1.UserRecord, error)
	FindOne(ctx context.Context, filter interface{}) (*entitiesv1.UserRecord, error)
	UpdateByID(ctx context.Context, ID string, input *entitiesv1.UserRecord) error
}

type DefaultUserRepository struct {
	User sparkle.Collection
}

func (d *DefaultUserRepository) FindOne(ctx context.Context, filter interface{}) (*entitiesv1.UserRecord, error) {
	var result entitiesv1.UserRecord
	err := d.User.FindOne(ctx, filter, &result)
	if err == sparkle.ErrNotFound {
		return nil, err
	}
	return &result, err
}

func (d *DefaultUserRepository) UpdateByID(ctx context.Context, ID string, input *entitiesv1.UserRecord) error {
	// ensure to
	// remove ID from input
	input.ID = nil
	err := d.User.Save(ctx, ID, input)
	return err
}

func NewDefaultUserRepository(user sparkle.Collection) *DefaultUserRepository {
	return &DefaultUserRepository{User: user}
}

func (d *DefaultUserRepository) FindByID(ctx context.Context, ID string) (*entitiesv1.UserRecord, error) {
	var result entitiesv1.UserRecord
	err := d.User.FindByID(ctx, ID, &result)
	return &result, err
}

func (d *DefaultUserRepository) Create(ctx context.Context, input *entitiesv1.UserRecord) (string, error) {
	id := xid.New()
	input.CreatedAt = commonv1.NewTimestamp(time.Now())
	err := d.User.Save(ctx, id.String(), input)
	return id.String(), err
}
