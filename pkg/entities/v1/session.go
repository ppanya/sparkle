package entitiesv1

import (
	"errors"
	"fmt"
	"github.com/octofoxio/sparkle/pkg/crypto"
	"time"
)

type SessionRecord struct {
	Session `bson:",inline"`
}

func (s *SessionRecord) IsValid() error {
	if s.GetLatestVisitedAt() == nil {
		return errors.New("session is expired")
	}
	var (
		latestVisited          = s.GetLatestVisitedAt().GetTime()
		timeSinceLatestVisited = time.Now().Sub(latestVisited)
	)

	fmt.Println(time.Now().Format(time.RFC1123Z))
	fmt.Println(latestVisited.Format(time.RFC1123Z))

	if timeSinceLatestVisited > time.Hour*24 {
		return errors.New("session is expired (24hrs)")
	}
	return nil
}

type SessionPayload struct {
	UserID string
}

func NewSession(signer sparklecrypto.TokenSigner, UserID string) (string, error) {
	return signer.Sign(SessionPayload{UserID: UserID})
}
