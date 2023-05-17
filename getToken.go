package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/repository"
)

// get the registration token from gh api
func GetToken(repo repository.Repository, org string, ent string, remove bool) (token string) {
	var tokenType string
	if remove {
		tokenType = "remove-token"
	} else {
		tokenType = "registration-token"
	}

	location := fmt.Sprintf("repos/%s/%s", repo.Owner(), repo.Name())
	if org != "" {
		location = fmt.Sprintf("orgs/%s/", org)
	}
	if ent != "" {
		location = fmt.Sprintf("enterprises/%s/", ent)
	}

	tokenGenCall := fmt.Sprintf("%s/actions/runners/%s", location, tokenType)

	value, stdErr, err := gh.Exec("api", "-X", "POST", tokenGenCall, "--jq", ".token")
	if err != nil {
		log.Fatal(stdErr.String()+"/n", err)
	}
	return strings.Split(value.String(), "\n")[0]
}
