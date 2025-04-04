package services

import (
	"log"
	"sync"

	jira "github.com/ctreminiom/go-atlassian/jira/v2"
	"github.com/pkg/errors"
)

var JiraClient = sync.OnceValue[*jira.Client](func() *jira.Client {
	host, mail, token := loadAtlassianCredentials()

	if host == "" || mail == "" || token == "" {
		log.Fatal("ATLASSIAN_HOST, ATLASSIAN_EMAIL, ATLASSIAN_TOKEN are required")
	}

	instance, err := jira.New(nil, host)
	if err != nil {
		log.Fatal(errors.WithMessage(err, "failed to create jira client"))
	}

	instance.Auth.SetBasicAuth(mail, token)

	return instance
})