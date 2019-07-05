package sparkleuc

import (
	"github.com/octofoxio/sparkle"
	sparklerepo "github.com/octofoxio/sparkle/pkg/repositories"
)

type SparkleUseCase struct {
	identity sparklerepo.IdentityRepository
	user     sparklerepo.UserRepository
	sparkle.EmailSender
}
