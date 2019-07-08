package sparkleuc

import (
	"bytes"
	"context"
	"github.com/octofoxio/sparkle/pkg/common"
	commonv1 "github.com/octofoxio/sparkle/pkg/common/v1"
	entitiesv1 "github.com/octofoxio/sparkle/pkg/entities/v1"
	svcsv1 "github.com/octofoxio/sparkle/pkg/svcs/v1"
	"html/template"
	"time"
)

func (s *RegisterUseCase) RegisterWithEmail(c context.Context, in *svcsv1.RegisterWithEmailInput) (out *svcsv1.RegisterWithEmailOutput, err error) {
	defer func() {

		if err != nil && out == nil {
			out = &svcsv1.RegisterWithEmailOutput{}
		}

	}()
	user := entitiesv1.User{
		Status:      entitiesv1.UserStatus_WaitingForEmailVerification,
		Email:       in.Email,
		FullName:    in.FullName,
		Gender:      in.Gender,
		PhoneNumber: in.PhoneNumber,
		CreatedAt:   commonv1.NewTimestamp(time.Now()),
	}

	record := &entitiesv1.UserRecord{
		User: user,
	}
	err = record.SetPassword(in.PlainPassword.GetData())
	if err != nil {
		return nil, err
	}
	ID, err := s.user.Create(c, record)
	user.ID = commonv1.NotNullString(ID)

	if err != nil {
		return nil, err
	}

	// สร้าง default identity
	// เพื่อใช้เป็น identity
	// แรกสำหรับทุกๆ platform
	_, err = s.identity.Create(c, &entitiesv1.IdentityRecord{
		UserID:   ID,
		SiteName: "default",
		Identity: entitiesv1.Identity{
			DisplayName: in.DisplayName,
		},
	})

	if err != nil {
		return nil, err
	}

	accessToken, err := NewSession(s.signer, ID)
	session := &entitiesv1.SessionRecord{
		Session: entitiesv1.Session{
			UserID:          commonv1.NotNullString(ID),
			AccessToken:     commonv1.NotNullString(accessToken),
			CreatedAt:       commonv1.NewTimestamp(time.Now()),
			LatestVisitedAt: commonv1.NewTimestamp(time.Now()),
		},
	}
	_, err = s.session.Create(c, session)

	if err != nil {
		return nil, err
	}

	config := common.GetConfigFromContext(c)
	confirmationURL := config.GetHost()
	confirmationURL.Path = "/c"

	confirmationURLQuery := confirmationURL.Query()
	confirmationURLQuery.Set("token", session.AccessToken.GetData())
	confirmationURLQuery.Set("callbackURL", in.CallbackURL.GetData())
	confirmationURL.RawQuery = confirmationURLQuery.Encode()

	tmpl, err := template.New("emailTemplate").Parse(config.DefaultEmailConfirmationTemplate)
	if err != nil {
		return nil, err
	}
	b := bytes.NewBuffer(nil)
	var curl = confirmationURL.String()

	err = tmpl.Execute(b, map[string]interface{}{
		"ConfirmUrl": template.HTML(curl),
	})
	if err != nil {
		return nil, err
	}

	err = s.EmailSender.Send(in.GetEmail().GetData(), config.DefaultEmailSenderAddress, b.String())
	if err != nil {
		return nil, err
	}

	return &svcsv1.RegisterWithEmailOutput{
		Result: &user,
	}, nil

}
