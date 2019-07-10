package sparkleuc

import (
	"github.com/octofoxio/foundation"
	"github.com/octofoxio/sparkle"
	"github.com/octofoxio/sparkle/external/line"
	sparklecrypto "github.com/octofoxio/sparkle/pkg/crypto"
	sparklerepo "github.com/octofoxio/sparkle/pkg/repositories"
)

type RegisterUseCase struct {
	identity sparklerepo.IdentityRepository
	user     sparklerepo.UserRepository
	session  sparklerepo.SessionRepository

	signer     sparklecrypto.TokenSigner
	lineClient line.LineClient

	sparkle.EmailSender
	fs foundation.FileSystem
}

func NewRegisterUseCase(signer sparklecrypto.TokenSigner, session sparklerepo.SessionRepository, identity sparklerepo.IdentityRepository, user sparklerepo.UserRepository, emailSender sparkle.EmailSender, fs foundation.FileSystem, client line.LineClient) *RegisterUseCase {
	return &RegisterUseCase{identity: identity, user: user, EmailSender: emailSender, fs: fs, signer: signer, session: session, lineClient: client}
}
