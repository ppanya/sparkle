package entitiesv1

import commonv1 "github.com/octofoxio/sparkle/pkg/common/v1"

type IdentityRecord struct {
	ID       *commonv1.String `bson:"_id,omitempty"`
	Identity `bson:"inline"`
	UserID   string `bson:",omitempty"`
	SiteName string `bson:",omitempty"`
}
