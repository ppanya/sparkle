package entitiesv1

type IdentityRecord struct {
	Identity
	UserID   string `bson:",omitempty"`
	SiteName string `bson:",omitempty"`
}
