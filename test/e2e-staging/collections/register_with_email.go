package collections

import (
	"fmt"
	"github.com/octofoxio/sparkle/pkg/common/v1"
	"github.com/octofoxio/sparkle/pkg/svcs/v1"
	"golang.org/x/net/context"
)

func RegisterWithEmail(client svcsv1.SparkleClient) {
	output, err := client.Register(context.Background(),
		&svcsv1.RegisterInput{
			RegisterInputData: &svcsv1.RegisterInput_Email{
				Email: &svcsv1.RegisterWithEmailInput{
					Email:         commonv1.NotNullString("maixezer@gmail.com"),
					DisplayName:   commonv1.NotNullString("Mai"),
					CallbackURL:   commonv1.NotNullString("https://www.google.com"),
					PlainPassword: commonv1.NotNullString("something"),
					PhoneNumber:   commonv1.NotNullString("+66909940794"),
					FullName:      commonv1.NotNullString("Johnny Apple Seed"),
					Gender:        commonv1.Gender_Female,
				},
			},
		},
	)

	if err != nil {
		panic(err)
	}
	fmt.Println(output.Result.GetID().GetData())
}
