package entitiesv1

import (
	"fmt"
	commonv1 "github.com/octofoxio/sparkle/pkg/common/v1"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"golang.org/x/crypto/bcrypt"
)

type UserRecord struct {
	User `bson:",omitempty,inline"`

	RegisterProvider commonv1.RegisterProvider `bson:",omitempty"`

	FacebookID *commonv1.String `bson:",omitempty"`
	LineID     *commonv1.String `bson:",omitempty"`

	// use a simple Bcrypt
	// will upgrade to another
	// complexity method later
	EncryptedPassword string `bson:",omitempty"`
}

func (x *UserRecord) SetPassword(plainTextPassword string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), bcrypt.MinCost)
	if err != nil {
		return err
	}
	x.EncryptedPassword = string(hashed)
	return nil
}

func (x UserRecord) ValidatePassword(value string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(x.EncryptedPassword), []byte(value)); err != nil {
		return false
	}
	return true
}

func (x UserStatus) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsonx.String(x.String()).MarshalBSONValue()
}

func (x *UserStatus) UnmarshalBSONValue(t bsontype.Type, b []byte) error {
	if t == bsontype.String {
		v := bsonx.String("")
		err := v.UnmarshalBSONValue(t, b)
		if err != nil {
			return err
		}
		d := UserStatus(UserStatus_value[v.String()])
		*x = d
		return nil
	} else {
		return fmt.Errorf("cannot marshal UserStatus, required string but receive %d", t.String())
	}
}
