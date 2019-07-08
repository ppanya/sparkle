package endpoints

import (
	"context"
	"errors"
	"github.com/octofoxio/sparkle"
	commonv1 "github.com/octofoxio/sparkle/pkg/common/v1"
	"github.com/octofoxio/sparkle/pkg/svcs/v1"
	sparkleuc "github.com/octofoxio/sparkle/pkg/usecases"
)

func NewSparkleEndpoints(register *sparkleuc.RegisterUseCase, login *sparkleuc.LoginUseCase) *SparkleEndpoints {
	return &SparkleEndpoints{
		register: register,
		login:    login,
	}
}

type SparkleEndpoints struct {
	register *sparkleuc.RegisterUseCase
	login    *sparkleuc.LoginUseCase
}

func (i *SparkleEndpoints) ValidateAccessToken(c context.Context, in *svcsv1.ValidateAccessTokenInput) (*svcsv1.ValidateAccessTokenOutput, error) {

	_, err := i.login.ValidateSession(c, in.AccessToken.GetData())
	if err != nil {
		return &svcsv1.ValidateAccessTokenOutput{
			Result: &svcsv1.ValidateAccessTokenOutput_SessionStatus{
				IsValid: false,
				Message: commonv1.NotNullString(err.Error()),
			},
		}, nil
	}

	return &svcsv1.ValidateAccessTokenOutput{
		Result: &svcsv1.ValidateAccessTokenOutput_SessionStatus{
			IsValid: true,
		},
	}, nil

}

func (i *SparkleEndpoints) GetMyProfile(c context.Context, in *svcsv1.GetMyProfileInput) (*svcsv1.GetMyProfileOutput, error) {

	user, ok := sparkle.GetUserProfileFromContext(c)
	if !ok {
		panic(commonv1.T_MissingCredentials)
	}

	return &svcsv1.GetMyProfileOutput{
		Result: &user.User,
	}, nil

}

func (i *SparkleEndpoints) Login(c context.Context, in *svcsv1.LoginInput) (o *svcsv1.LoginOutput, err error) {

	defer func() {
		if err != nil {
			o = &svcsv1.LoginOutput{}
		}
	}()

	switch v := in.GetLoginInputData().(type) {
	case *svcsv1.LoginInput_Email:
		output, err := i.login.LoginWithEmail(c, v.Email)
		if err != nil {
			return nil, err
		}
		return &svcsv1.LoginOutput{
			Result: &output.Session,
		}, err
	case *svcsv1.LoginInput_Facebook:
		return nil, errors.New("implement me")
	}
	return nil, errors.New("unknown error")
}

func (d *SparkleEndpoints) RegisterWithEmail(c context.Context, in *svcsv1.RegisterWithEmailInput) (*svcsv1.RegisterWithEmailOutput, error) {
	o, e := d.register.RegisterWithEmail(c, in)
	return o, e
}

func (s *SparkleEndpoints) ConfirmEmailByAccessTokenHandler(c context.Context, accessToken, callbackURL string) error {
	if accessToken == "" {
		return errors.New("invalid token")
	}
	if callbackURL == "" {
		return errors.New("invalid callback URL")
	}
	return s.register.ConfirmEmailHandler(c, accessToken)
}
