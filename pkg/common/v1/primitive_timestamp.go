package commonv1

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"time"
)

func NewTimestampNow() *Timestamp {
	return &Timestamp{
		Seconds: int64(time.Now().Unix()),
	}
}

func NewTimestamp(t time.Time) *Timestamp {
	return &Timestamp{
		Seconds: int64(t.Unix()),
	}
}

func (m *Timestamp) GetTime() time.Time {
	return time.Unix(m.Seconds, 0)
}
func (m *Timestamp) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if m == nil {
		return bsonx.Null().MarshalBSONValue()
	}
	return bsonx.Time(time.Unix(m.Seconds, 0)).MarshalBSONValue()
}

func (m *Timestamp) UnmarshalBSONValue(b bsontype.Type, bb []byte) error {
	if b == bson.TypeDateTime {
		v := bsonx.Time(time.Unix(0, 0))
		err := v.UnmarshalBSONValue(b, bb)
		if err != nil {
			return err
		}
		m.Seconds = v.Time().Unix()
		return nil
	} else {
		return fmt.Errorf("cannot parse Timestamp receive %s but require int64", b.String())
	}
}
