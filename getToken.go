package main

import (
	"fmt"
	"log"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/repository"
)

// get the registration token from gh api
func GetToken(repo repository.Repository, org string, ent string, remove bool) (token string) {
	var tokenType string
	if remove {
		tokenType = "remove-token"
	} else {
		tokenType = "registration-token"
	}

	location := fmt.Sprintf("repos/%s/%s", repo.Owner, repo.Name)
	if org != "" {
		location = fmt.Sprintf("orgs/%s", org)
	}
	if ent != "" {
		location = fmt.Sprintf("enterprises/%s", ent)
	}

	tokenGenCall := fmt.Sprintf("%s/actions/runners/%s", location, tokenType)
	ghRest, err := api.DefaultRESTClient()
	if err != nil {
		log.Fatal(err)
	}
	tokencontainer := struct {
		Token string
	}{}
	//post with an empty body
	err = ghRest.Post(tokenGenCall, nil, &tokencontainer)
	if err != nil {
		log.Fatal(err)
	}
	if tokencontainer.Token == "" {
		log.Fatal("Unable to get a valid token")
	}
	return tokencontainer.Token
}
