package services

import (
	"log"
	"os"
	"sync"

	"github.com/ctreminiom/go-atlassian/jira/agile"
	"github.com/pkg/errors"
)

func loadAtlassianCredentials() (host, mail, token, pat string) {
	host = os.Getenv("ATLASSIAN_HOST")
	mail = os.Getenv("ATLASSIAN_EMAIL")
	token = os.Getenv("ATLASSIAN_TOKEN")
	pat = os.Getenv("ATLASSIAN_PAT")

	if host == "" {
		log.Fatal("ATLASSIAN_HOST is required, please set it in MCP Config")
	}

	if pat == "" && (mail == "" || token == "") {
		log.Fatal("Authentication required: set ATLASSIAN_PAT (for Jira Server/Data Center) or both ATLASSIAN_EMAIL and ATLASSIAN_TOKEN (for Jira Cloud)")
	}

	return host, mail, token, pat
}

var AgileClient = sync.OnceValue[*agile.Client](func() *agile.Client {
	host, mail, token, pat := loadAtlassianCredentials()

	instance, err := agile.New(nil, host)
	if err != nil {
		log.Fatal(errors.WithMessage(err, "failed to create agile client"))
	}

	if pat != "" {
		instance.Auth.SetBearerToken(pat)
	} else {
		instance.Auth.SetBasicAuth(mail, token)
	}

	return instance
})
