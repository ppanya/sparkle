package sparklerepo

import (
	"context"
	"github.com/octofoxio/sparkle"
	entitiesv1 "github.com/octofoxio/sparkle/pkg/entities/v1"
	"github.com/rs/xid"
)

type IdentityRepository interface {
	Create(ctx context.Context, input *entitiesv1.IdentityRecord) (ID string, err error)
	FindByID(ctx context.Context, ID string) (*entitiesv1.IdentityRecord, error)
	UpdateByID(ctx context.Context, ID string, input *entitiesv1.IdentityRecord) error
}

type DefaultIdentityRepository struct {
	Identity sparkle.Collection
}

func (d *DefaultIdentityRepository) Create(ctx context.Context, input *entitiesv1.IdentityRecord) (ID string, err error) {
	ID = xid.New().String()
	err = d.Identity.Save(ctx, ID, input)
	return ID, err
}

func (d *DefaultIdentityRepository) FindByID(ctx context.Context, ID string) (*entitiesv1.IdentityRecord, error) {
	var result entitiesv1.IdentityRecord
	err := d.Identity.FindByID(ctx, ID, &result)
	return &result, err
}

func (d *DefaultIdentityRepository) UpdateByID(ctx context.Context, ID string, input *entitiesv1.IdentityRecord) error {
	// ensure to
	// remove ID from input
	input.ID = nil
	err := d.Identity.Save(ctx, ID, input)
	return err
}
