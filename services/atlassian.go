package services

import (
	"log"
	"os"
	"sync"

	"github.com/ctreminiom/go-atlassian/jira/agile"
	"github.com/pkg/errors"
)

func loadAtlassianCredentials() (host, mail, token string) {
	host = os.Getenv("ATLASSIAN_HOST")
	mail = os.Getenv("ATLASSIAN_EMAIL")
	token = os.Getenv("ATLASSIAN_TOKEN")

	if host == "" || mail == "" || token == "" {
		log.Fatal("ATLASSIAN_HOST, ATLASSIAN_EMAIL, ATLASSIAN_TOKEN are required, please set it in MCP Config")
	}

	return host, mail, token
}

var AgileClient = sync.OnceValue[*agile.Client](func() *agile.Client {
	host, mail, token := loadAtlassianCredentials()

	instance, err := agile.New(nil, host)
	if err != nil {
		log.Fatal(errors.WithMessage(err, "failed to create agile client"))
	}

	instance.Auth.SetBasicAuth(mail, token)

	return instance
})
