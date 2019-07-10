package sparkleuc

import (
	"context"
	"errors"
	"github.com/octofoxio/sparkle"
	commonv1 "github.com/octofoxio/sparkle/pkg/common/v1"
	entitiesv1 "github.com/octofoxio/sparkle/pkg/entities/v1"
	sparklerepo "github.com/octofoxio/sparkle/pkg/repositories"
)

type ProfileUseCase struct {
	identity sparklerepo.IdentityRepository
	user     sparklerepo.UserRepository
}

func NewProfileUseCase(identity sparklerepo.IdentityRepository, user sparklerepo.UserRepository) *ProfileUseCase {
	return &ProfileUseCase{identity: identity, user: user}
}

func (p *ProfileUseCase) MustGetIdentity(ctx context.Context, UserID, SiteName string) (identity *entitiesv1.IdentityRecord) {
	identity, err := p.identity.FindOne(ctx, &entitiesv1.IdentityRecord{
		SiteName: SiteName,
		UserID:   UserID,
	})
	if err != nil {
		panic(err)
	}
	return identity
}
func (p *ProfileUseCase) MustGetDefaultIdentity(ctx context.Context, UserID string) (identity *entitiesv1.IdentityRecord) {
	identity, err := p.identity.FindOne(ctx, &entitiesv1.IdentityRecord{
		SiteName: "default",
		UserID:   UserID,
	})
	if err != nil {
		if err == sparkle.ErrNotFound {
			panic(commonv1.T_DefaultIdentityNotFound)
		}
		panic(err)
	}
	return identity
}

func (p *ProfileUseCase) PutIdentity(ctx context.Context, UserID, SiteName string, input *entitiesv1.Identity) (identity *entitiesv1.IdentityRecord, err error) {

	if SiteName == "default" {
		return nil, errors.New("default identity is read only")
	}

	ID := UserID + "::" + SiteName
	err = p.identity.UpdateByID(ctx, ID, &entitiesv1.IdentityRecord{
		Identity: *input,
		SiteName: SiteName,
		UserID:   UserID,
	})
	if err != nil {
		return nil, err
	}

	identity = p.MustGetIdentity(ctx, UserID, SiteName)

	return identity, err

}

func (p *ProfileUseCase) GetIdentity(ctx context.Context, UserID, SiteName string) (identity *entitiesv1.IdentityRecord, err error) {

	identity, err = p.identity.FindOne(ctx, &entitiesv1.IdentityRecord{
		SiteName: SiteName,
		UserID:   UserID,
	})

	if err == sparkle.ErrNotFound {
		return p.MustGetDefaultIdentity(ctx, UserID), nil
	}

	if err != nil {
		return nil, err
	}

	return identity, nil

}
