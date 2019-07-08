package sparklerepo

import (
	"context"
	"github.com/octofoxio/sparkle"
	entitiesv1 "github.com/octofoxio/sparkle/pkg/entities/v1"
	"github.com/rs/xid"
)

type SessionRepository interface {
	Create(ctx context.Context, input *entitiesv1.SessionRecord) (id string, err error)
	FindOne(ctx context.Context, filter interface{}) (*entitiesv1.SessionRecord, error)
	Update()
	Delete()
}

type DefaultSessionRepository struct {
	Session sparkle.Collection
}

func (d *DefaultSessionRepository) Create(ctx context.Context, input *entitiesv1.SessionRecord) (ID string, err error) {
	ID = xid.New().String()
	err = d.Session.Save(ctx, ID, input)
	return ID, err
}

func (d *DefaultSessionRepository) FindOne(ctx context.Context, filter interface{}) (*entitiesv1.SessionRecord, error) {
	var result entitiesv1.SessionRecord
	err := d.Session.FindOne(ctx, filter, &result)
	return &result, err
}

func (d *DefaultSessionRepository) Update() {
	panic("implement me")
}

func (d *DefaultSessionRepository) Delete() {
	panic("implement me")
}

func NewDefaultSessionRepository(db sparkle.Collection) *DefaultSessionRepository {
	return &DefaultSessionRepository{
		Session: db,
	}
}
