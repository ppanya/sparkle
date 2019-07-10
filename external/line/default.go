package line

import "context"

var (
	lineGETProfileURL = "https://api.line.me/v2/profile"
)

type LineClient interface {
	GetProfile(ctx context.Context, LineAccessToken string) (p *Profile, err error)
}

type DefaultLineClient struct{}

func NewDefaultLineClient() *DefaultLineClient {
	return &DefaultLineClient{}

}
