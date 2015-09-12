package auth

import (
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type GithubAuth struct {
	oauthConfig *oauth2.Config
}

func NewGithubAuth(clientId, secret string) *GithubAuth {
	return &GithubAuth{
		oauthConfig: &oauth2.Config{
			ClientID:     clientId,
			ClientSecret: secret,
			Scopes:       []string{"user", "repo", "admin:repo_hook"},
			Endpoint:     github.Endpoint,
		},
	}
}

func (a *GithubAuth) AuthHandle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.Redirect(w, r, a.oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOnline), http.StatusTemporaryRedirect)
}

func (a *GithubAuth) CallbackHandle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	log.Debugf("code: %s | state: %s true: %t", code, state, code == "")

	if code == "" || state != "state" {
		log.Error("invalid code or state")
		w.Write([]byte(`invalid code or state`))
		return
	}

	token, err := a.oauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Error(err)
		w.Write([]byte(err.Error()))
		return
	}

	expire := time.Now().AddDate(0, 0, 1)
	cookie := http.Cookie{
		Name:    "github",
		Value:   token.AccessToken,
		Path:    "/",
		Expires: expire,
	}
	http.SetCookie(w, &cookie)
	log.Debugf("Received Auth Callback token: %s", token)
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
