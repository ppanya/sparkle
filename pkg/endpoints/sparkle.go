package endpoints

import (
	"context"
	"errors"
	"github.com/octofoxio/sparkle"
	commonv1 "github.com/octofoxio/sparkle/pkg/common/v1"
	entitiesv1 "github.com/octofoxio/sparkle/pkg/entities/v1"
	"github.com/octofoxio/sparkle/pkg/svcs/v1"
	sparkleuc "github.com/octofoxio/sparkle/pkg/usecases"
)

func NewSparkleEndpoints(register *sparkleuc.RegisterUseCase, login *sparkleuc.LoginUseCase, profile *sparkleuc.ProfileUseCase) *SparkleEndpoints {
	return &SparkleEndpoints{
		register: register,
		login:    login,
		profile:  profile,
	}
}

type SparkleEndpoints struct {
	register *sparkleuc.RegisterUseCase
	login    *sparkleuc.LoginUseCase
	profile  *sparkleuc.ProfileUseCase
}

func (d *SparkleEndpoints) Register(c context.Context, in *svcsv1.RegisterInput) (*svcsv1.RegisterOutput, error) {

	switch data := in.GetRegisterInputData().(type) {
	case *svcsv1.RegisterInput_Email:
		o, e := d.register.RegisterWithEmail(c, data.Email)
		return &svcsv1.RegisterOutput{Result: o}, e
	case *svcsv1.RegisterInput_Line:
		o, e := d.register.RegisterWithLine(c, data.Line)
		return &svcsv1.RegisterOutput{Result: o}, e
	}
	panic("invalid register provider")
}

func (i *SparkleEndpoints) PutIdentity(c context.Context, in *svcsv1.PutIdentityInput) (*svcsv1.PutIdentityOutput, error) {
	user, ok := sparkle.GetUserProfileFromContext(c)
	if !ok {
		panic(commonv1.T_MissingCredentials)
	}
	output, err := i.profile.PutIdentity(c, user.GetID().GetData(), in.GetSiteName().GetData(), in.Data)
	if err != nil {
		return &svcsv1.PutIdentityOutput{}, nil
	}
	return &svcsv1.PutIdentityOutput{
		Result: &output.Identity,
	}, err

}

func (i *SparkleEndpoints) GetIdentity(c context.Context, in *svcsv1.GetIdentityInput) (*svcsv1.GetIdentityOutput, error) {

	user, ok := sparkle.GetUserProfileFromContext(c)
	if !ok {
		panic(commonv1.T_MissingCredentials)
	}

	out, err := i.profile.GetIdentity(c, user.ID.GetData(), in.GetSiteName().GetData())

	if err != nil {
		panic(err)
	}

	return &svcsv1.GetIdentityOutput{
		Result: &out.Identity,
	}, nil

}

func (i *SparkleEndpoints) ValidateAccessToken(c context.Context, in *svcsv1.ValidateAccessTokenInput) (*svcsv1.ValidateAccessTokenOutput, error) {

	s, _, err := i.login.ValidateSession(c, in.AccessToken.GetData())
	if err != nil {
		return &svcsv1.ValidateAccessTokenOutput{
			Result: &svcsv1.ValidateAccessTokenOutput_SessionStatus{
				IsValid: false,
				Message: commonv1.NotNullString(err.Error()),
				Session: s,
			},
		}, nil
	}

	return &svcsv1.ValidateAccessTokenOutput{
		Result: &svcsv1.ValidateAccessTokenOutput_SessionStatus{
			IsValid: true,
			Session: s,
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

	var session *entitiesv1.SessionRecord
	switch v := in.GetLoginInputData().(type) {
	case *svcsv1.LoginInput_Email:
		session, err = i.login.LoginWithEmail(c, v.Email)
		if err != nil {
			return nil, err
		}
	case *svcsv1.LoginInput_Line:
		session, err = i.login.LoginWithLine(c, v.Line)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unknown error")
	}

	return &svcsv1.LoginOutput{
		Result: &session.Session,
	}, nil

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
