package sparkleuc

import sparklecrypto "github.com/octofoxio/sparkle/pkg/crypto"

type SessionPayload struct {
	UserID string
}

func NewSession(signer sparklecrypto.TokenSigner, UserID string) (string, error) {
	return signer.Sign(SessionPayload{UserID: UserID})
}
