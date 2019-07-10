package line

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Profile struct {
	DisplayName   string `json:"displayName"`
	UserID        string `json:"userId"`
	PictureURL    string `json:"pictureUrl"`
	StatusMessage string `json:"statusMessage"`
	Email         string `json:"email"`
}

func (l *DefaultLineClient) GetProfile(ctx context.Context, LineAccessToken string) (p *Profile, err error) {

	req, _ := http.NewRequest("GET", lineGETProfileURL, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", LineAccessToken))
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = res.Body.Close()
	}()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var pp Profile
	err = json.Unmarshal(body, &pp)
	return &pp, err

}
