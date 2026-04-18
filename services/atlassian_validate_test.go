package services

import (
	"reflect"
	"testing"
)

func TestValidateAtlassianEnv(t *testing.T) {
	tests := []struct {
		name      string
		host      string
		mail      string
		token     string
		pat       string
		wantMiss  []string
	}{
		{
			name:     "all empty — host plus auth missing",
			wantMiss: []string{"ATLASSIAN_HOST", "ATLASSIAN_PAT or (ATLASSIAN_EMAIL + ATLASSIAN_TOKEN)"},
		},
		{
			name:     "host only — auth missing",
			host:     "https://x.atlassian.net",
			wantMiss: []string{"ATLASSIAN_PAT or (ATLASSIAN_EMAIL + ATLASSIAN_TOKEN)"},
		},
		{
			name:     "PAT only — valid for Server/DC",
			host:     "https://jira.company.com",
			pat:      "pat-xyz",
			wantMiss: nil,
		},
		{
			name:     "email + token — valid for Cloud",
			host:     "https://x.atlassian.net",
			mail:     "a@b.c",
			token:    "tkn",
			wantMiss: nil,
		},
		{
			name:     "email without token — auth missing",
			host:     "https://x.atlassian.net",
			mail:     "a@b.c",
			wantMiss: []string{"ATLASSIAN_PAT or (ATLASSIAN_EMAIL + ATLASSIAN_TOKEN)"},
		},
		{
			name:     "token without email — auth missing",
			host:     "https://x.atlassian.net",
			token:    "tkn",
			wantMiss: []string{"ATLASSIAN_PAT or (ATLASSIAN_EMAIL + ATLASSIAN_TOKEN)"},
		},
		{
			name:     "PAT wins even when email/token also present",
			host:     "https://jira.company.com",
			mail:     "a@b.c",
			token:    "tkn",
			pat:      "pat-xyz",
			wantMiss: nil,
		},
		{
			name:     "host missing but auth present — host still reported",
			pat:      "pat-xyz",
			wantMiss: []string{"ATLASSIAN_HOST"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateAtlassianEnv(tt.host, tt.mail, tt.token, tt.pat)
			if !reflect.DeepEqual(got, tt.wantMiss) {
				t.Errorf("ValidateAtlassianEnv = %v, want %v", got, tt.wantMiss)
			}
		})
	}
}
