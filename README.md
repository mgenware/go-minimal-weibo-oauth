# go-minimal-weibo-oauth
A minimal library for Weibo oauth

## Installation
```sh
go get github.com/mgenware/go-minimal-weibo-oauth
```

# Example
```go
package main

import (
	"fmt"
	"net/http"

	"github.com/mgenware/go-minimal-weibo-oauth"
)

const (
    clientID     = "{{Your Client ID}}"
    clientSecret = "{{Your Client Secret}}"
    redirectionURL  = "{{Your Redirection URL}}"
    urlState     = "{{Some random string}}"
)

var weibo *weiboOAuth.OAuth

func init() {
	var err error
	weibo, err = weiboOAuth.NewWeiboOAuth(clientID, clientSecret, redirectionURL)
	if err != nil {
		panic(err)
	}

	weiboOAuth.Logging = true
}

func oauthHandler(w http.ResponseWriter, r *http.Request) {
	urlStr, err := weibo.GetAuthorizationURL(urlState)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	http.Redirect(w, r, urlStr, http.StatusMovedPermanently)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	if code == "" {
		w.Write([]byte("Invalid code"))
		return
	}

	state := r.FormValue("state")
	if state != urlState {
		w.Write([]byte("Invalid state"))
		return
	}

	token, err := weibo.GetAccessToken(code)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	userID := token.UIDString
	profile, err := weibo.GetUserInfo(token.AccessToken, userID)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(fmt.Sprint(profile)))
}

func main() {
	http.HandleFunc("/weibo_oauth", oauthHandler)
	http.HandleFunc("/weibo_oauth_callback", callbackHandler)
	http.ListenAndServe(":3000", nil)
}

```
