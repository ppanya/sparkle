package testutils

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/url"
)

type LineSession struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	IDToken     string `json:"id_token"`
}

func LineLogin(channelID, channelSecret string, result chan LineSession) {
	r := &mux.Router{}
	r.Path("/c").Methods("GET").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html")
		writer.WriteHeader(http.StatusOK)

		code := request.FormValue("code")
		fmt.Printf("authorization code: %s\n", code)
		lineAuthURL, _ := url.Parse("https://api.line.me/oauth2/v2.1/token")
		resp, err := http.PostForm(lineAuthURL.String(), map[string][]string{
			"grant_type":    {"authorization_code"},
			"code":          {code},
			"redirect_uri":  {"http://localhost:3099/c"},
			"client_id":     {channelID},
			"client_secret": {channelSecret},
		})
		if err != nil {
			panic(err)
		}
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		var s LineSession
		err = json.Unmarshal(b, &s)
		if err != nil {
			panic(err)
		}

		result <- s
		_, _ = writer.Write([]byte("<html><body>Ok, you can close this window</body></html>"))
	})

	lineLoginURL, _ := url.Parse("https://access.line.me/oauth2/v2.1/authorize")
	q := lineLoginURL.Query()
	q.Set("response_type", "code")
	q.Set("client_id", channelID)
	q.Set("redirect_uri", "http://localhost:3099/c")
	q.Set("state", "91821")
	q.Set("scope", "profile openid")
	lineLoginURL.RawQuery = q.Encode()

	fmt.Println("====== Click on link to login with Line for testing =======")
	fmt.Println(lineLoginURL.String())
	err := http.ListenAndServe("localhost:3099", r)
	if err != nil {
		panic(err)
	}
}
