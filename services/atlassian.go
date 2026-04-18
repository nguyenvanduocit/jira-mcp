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

	// Defense-in-depth: main.go already validates these with the same rule,
	// but clients used by tests (or future callers that bypass main) still
	// benefit from a crash-early guard. Share the rule via ValidateAtlassianEnv
	// so both layers stay in lock-step.
	if missing := ValidateAtlassianEnv(host, mail, token, pat); len(missing) > 0 {
		log.Fatalf("Atlassian credentials incomplete — missing: %v", missing)
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
