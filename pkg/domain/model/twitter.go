package model

import (
	"context"
	"fmt"
	twitter "github.com/g8rswimmer/go-twitter/v2"
	"log"
	"net/http"
	"os"
	"strings"
)

type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

type TwitterInfo struct {
	Name     string
	Username string
}

func verifyTwitter(names []string) map[string]TwitterInfo {
	client := &twitter.Client{
		Authorizer: authorize{
			Token: os.Getenv("TwitterToken"),
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}
	opts := twitter.UserLookupOpts{}

	userResponse, err := client.UserNameLookup(context.Background(), names, opts)
	if err != nil {
		log.Panicf("user lookup error: %v", err)
	}

	dictionaries := userResponse.Raw.UserDictionaries()
	var resp = make(map[string]TwitterInfo)
	for _, d := range dictionaries {
		var newInfo TwitterInfo
		newInfo.Name = d.User.Name
		newInfo.Username = d.User.UserName
		resp[strings.ToUpper(newInfo.Username)] = newInfo
	}
	//return list of verified accounts
	return resp
}

func ValidateTwitter(data []*EventData) {
	var entries []string
	var targets []*EventData
	//load entries
	for _, d := range data {
		if d.snsTwitter.Status != NG && len(d.snsTwitter.Value) > 0 {
			fmt.Println(d.snsTwitter.Value)
			entries = append(entries, d.snsTwitter.Value)
			targets = append(targets, d)
		}
	}
	//get list of verified accounts in entries
	accounts := verifyTwitter(entries)
	var verifiedInfo []struct {
		username     string
		name         string
		eventOrgName string
	}
	for _, d := range targets {
		ac, ok := accounts[strings.ToUpper(d.snsTwitter.Value)]
		if ok {
			d.snsTwitter.Value = ac.Username
			d.snsTwitter.setVerification(Verified)
			verifiedInfo = append(verifiedInfo, struct {
				username     string
				name         string
				eventOrgName string
			}{username: ac.Username, name: ac.Name, eventOrgName: d.eventOrgName})
		} else {
			d.snsTwitter.setVerification(Error)
		}
	}
	for _, s := range verifiedInfo {
		fmt.Printf("%20s %50s %50s\n", s.username, s.name, s.eventOrgName)
	}
}
