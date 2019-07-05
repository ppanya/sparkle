package commonv1

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"time"
)

func NewTime(t time.Time) *Timestamp {
	return &Timestamp{
		Seconds: int64(t.Second()),
	}
}

func (m *Timestamp) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsonx.Time(time.Unix(m.Seconds, 0)).MarshalBSONValue()
}

func (m *Timestamp) UnmarshalBSONValue(b bsontype.Type, bb []byte) error {
	if b == bson.TypeDateTime {
		v := bsonx.Time(time.Now())
		err := v.UnmarshalBSONValue(b, bb)
		if err != nil {
			return err
		}
		m.Seconds = int64(v.Time().Second())
		return nil
	} else {
		return fmt.Errorf("cannot parse Timestamp receive %s but require int64", b.String())
	}
}
