# go-ldap-client

[![codecov](https://codecov.io/gh/colynn/go-ldap-client/branch/master/graph/badge.svg)](https://codecov.io/gh/colynn/go-ldap-client)
[![GoDoc](https://godoc.org/github.com/colynn/go-ldap-client?status.svg)](https://pkg.go.dev/github.com/colynn/go-ldap-client?tab=doc) 

Simple ldap client to authenticate, retrieve basic information and groups for a user.

## Usage

The only external dependency is [gopkg.in/ldap.v3](https://pkg.go.dev/github.com/go-ldap/ldap/v3).

```
package main

import (
	"log"

	ldap "github.com/colynn/go-ldap-client/v3"
)

func main() {
	client := &ldap.Client{
		Base:         "dc=example,dc=com",
		Host:         "ldap.example.com",
		Port:         389,
		UseSSL:       false,
		BindDN:       "uid=readonlysuer,ou=People,dc=example,dc=com",
		BindPassword: "readonlypassword",
		UserFilter:   "(uid=%s)",
		GroupFilter: "(memberUid=%s)",
		Attributes:   []string{"givenName", "sn", "mail", "uid"},
	}
	// It is the responsibility of the caller to close the connection
	defer client.Close()

	ok, user, err := client.Authenticate("username", "password")
	if err != nil {
		log.Fatalf("Error authenticating user %s: %+v", "username", err)
	}
	if !ok {
		log.Fatalf("Authenticating failed for user %s", "username")
	}
	log.Printf("User: %+v", user)
	
	groups, err := client.GetGroupsOfUser("username")
	if err != nil {
		log.Fatalf("Error getting groups for user %s: %+v", "username", err)
	}
	log.Printf("Groups: %+v", groups) 
}
```

## SSL(LDAPS)
If you use SSL, you will need to pass the server name for certificate verification or skip domain name verification e.g.`client.ServerName = "ldap.example.com"`.

## Why?
Because [go-ldap-client](https://github.com/jtblin/go-ldap-client) been a long time didn't maintenance from 2017 to now.
So re-create it, make it better for everyone to use and maintain.
