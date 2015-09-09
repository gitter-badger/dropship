package github

import (
	"encoding/json"
	"net/http"

	"github.com/ChrisMcKenzie/dropship/logging"
	"github.com/google/go-github/github"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"
)

var log = logging.GetLogger()

func AddHook(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	c, _ := r.Cookie("github")
	log.Debug(c.Value)
	if c.Value == "" {
		return
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.Value},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	name := "web"

	hook, _, err := client.Repositories.CreateHook(
		p.ByName("repo_owner"),
		p.ByName("repo_name"),
		&github.Hook{
			Name:   &name,
			Events: []string{"deployment"},
			Config: map[string]interface{}{
				"url": "http://joog.chrismckenzie.io/deploy/github.com/" + p.ByName("repo_owner") + "/" + p.ByName("repo_name"),
			},
		},
	)

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	log.Debugf("Hook Created: %v", hook)
	json, _ := json.Marshal(hook)
	w.Write(json)
}

func GetRepos(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	c, _ := r.Cookie("github")
	log.Debug(c.Value)
	if c.Value == "" {
		return
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.Value},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	repos, _, err := client.Repositories.List("", &github.RepositoryListOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	})
	if err != nil {
		return
	}

	re, err := json.Marshal(repos)
	if err != nil {
		return
	}

	w.Write(re)
}