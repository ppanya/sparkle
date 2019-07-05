package entitiesv1

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type UserRecord struct {
	User `bson:",omitempty,inline"`
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
