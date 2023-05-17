package main

import (
	"fmt"
	"log"

	"github.com/cli/go-gh"
)

// get the registration token from gh api
func GetToken(repo string, org string, ent string, remove bool) (token string) {
	var tokenType string
	if remove {
		tokenType = "remove-token"
	} else {
		tokenType = "registration-token"
	}

	location := fmt.Sprintf("repos/%s/", repo)
	if org != "" {
		location = fmt.Sprintf("orgs/%s/", org)
	}
	if ent != "" {
		location = fmt.Sprintf("enterprises/%s/", ent)
	}

	tokenGenCall := fmt.Sprintf("%s/actions/runners/%s", location, tokenType)

	value, stdErr, err := gh.Exec("api", tokenGenCall, "--jq", ".token")
	if err != nil {
		log.Fatal(stdErr.String()+"/n", err)
	}
	return value.String()
}
