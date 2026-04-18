package services

import (
	"log"
	"sync"

	jira "github.com/ctreminiom/go-atlassian/jira/v3"
	"github.com/pkg/errors"
)

var JiraClient = sync.OnceValue[*jira.Client](func() *jira.Client {
	host, mail, token, pat := loadAtlassianCredentials()

	instance, err := jira.New(nil, host)
	if err != nil {
		log.Fatal(errors.WithMessage(err, "failed to create jira client"))
	}

	if pat != "" {
		instance.Auth.SetBearerToken(pat)
	} else {
		instance.Auth.SetBasicAuth(mail, token)
	}

	return instance
})