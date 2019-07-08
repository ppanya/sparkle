package entitiesv1

type IdentityRecord struct {
	Identity `bson:"inline"`
	UserID   string `bson:",omitempty"`
	SiteName string `bson:",omitempty"`
}
