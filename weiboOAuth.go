package weiboOAuth

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	AuthURL        = "https://api.weibo.com/oauth2/authorize"
	AccessTokenURL = "https://api.weibo.com/oauth2/access_token"
	UserInfoURL    = "https://api.weibo.com/2/users/show.json"
)

type OAuth struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type OAuthToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	RemindIn    string `json:"remind_in"`
	UIDString   string `json:"uid"`

	Error        string `json:"error"`
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_description"`
}

type UserInfo struct {
	UID         int64  `json:"id"`
	Name        string `json:"name"`
	Location    string `json:"location"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

func NewWeiboOAuth(clientID, clientSecret, redirectURL string) (*OAuth, error) {
	if clientID == "" {
		return nil, errors.New("clientID cannot be empty")
	}
	if clientSecret == "" {
		return nil, errors.New("clientSecret cannot be empty")
	}
	if redirectURL == "" {
		return nil, errors.New("redirectURL cannot be empty")
	}

	oauth := &OAuth{}
	oauth.ClientID = clientID
	oauth.ClientSecret = clientSecret
	oauth.RedirectURL = redirectURL
	return oauth, nil
}

func (oauth *OAuth) GetAuthorizeURL() string {
	qs := url.Values{"client_id": {oauth.ClientID},
		"redirect_uri": {oauth.RedirectURL}}
	urlStr := AuthURL + "?" + qs.Encode()
	return urlStr
}

func (oauth *OAuth) GetAccessToken(code string) (*OAuthToken, error) {
	if code == "" {
		return nil, errors.New("code cannot be empty")
	}
	resp, err := http.PostForm(AccessTokenURL,
		url.Values{"client_id": {oauth.ClientID},
			"client_secret": {oauth.ClientSecret},
			"grant_type":    {"authorization_code"},
			"code":          {code},
			"redirect_uri":  {oauth.RedirectURL}})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	token := &OAuthToken{}
	err = json.Unmarshal(body, token)
	if err != nil {
		return nil, err
	}
	if token.ErrorCode != 0 {
		return nil, errors.New(token.ErrorMessage)
	}
	return token, err
}

func (oauth *OAuth) GetUserInfo(token *OAuthToken, uid string) (*UserInfo, error) {
	if token == nil {
		return nil, errors.New("token cannot be nil")
	}
	qs := url.Values{"access_token": {token.AccessToken},
		"uid": {uid}}
	urlStr := UserInfoURL + "?" + qs.Encode()

	resp, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ret := &UserInfo{}
	err = json.Unmarshal(body, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
