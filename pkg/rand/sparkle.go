package rand

import (
	"github.com/brianvoe/gofakeit"
	commonv1 "github.com/octofoxio/sparkle/pkg/common/v1"
	svcsv1 "github.com/octofoxio/sparkle/pkg/svcs/v1"
)

func RegisterWithEmailInput() *svcsv1.RegisterWithEmailInput {

	return &svcsv1.RegisterWithEmailInput{
		Email:         commonv1.NotNullString(gofakeit.Email()),
		FullName:      commonv1.NotNullString(gofakeit.FirstName()),
		DisplayName:   commonv1.NotNullString(gofakeit.Username()),
		CallbackURL:   commonv1.NotNullString(gofakeit.URL()),
		Gender:        commonv1.Gender_Female,
		PhoneNumber:   commonv1.NotNullString(gofakeit.Phone()),
		PlainPassword: commonv1.NotNullString(gofakeit.Password(true, true, true, true, false, 10)),
	}

}
