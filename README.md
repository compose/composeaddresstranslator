ComposeAddressTranslator for gocql
=

This package implements the [AddressTranslator](https://github.com/gocql/gocql/blob/1f5574155a67802fea8f94d486fe95ea76178242/address_translators.go#L8) interface for the [gocql/gocql](https://github.com/gocql/gocql) package.

It is intended to be used with [Compose Hosted Scylla](https://www.compose.com/scylladb) to properly map internal ip addresses to public-facing portals.

```go
package main

import (
	"log"

	"github.com/gocql/gocql"
	"github.com/compose/composeaddresstranslator"
)

func main() {
    // Create address translator with info from deployment dashboard
	cat := composeaddresstranslator.New(map[string]string{
		"10.49.168.5:9042": "some.portal.host:15992",
		"10.49.168.6:9042": "another.portal.host:15993",
		"10.49.168.7:9042": "more.portal.hosts:15994",
	})
    // cat := composeaddresstranslator.NewFromJSONString(`
    // "10.49.168.5:9042": "some.portal.host:15992",
    // "10.49.168.6:9042": "another.portal.host:15993",
    // "10.49.168.7:9042": "more.portal.hosts:15994"
    // }`)

    // Initialize cluster with the contact points from the address translator.
	cluster := gocql.NewCluster(cat.ContactPoints()...)  
    // Set your username and password
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: "{REDACTED}",
		Password: "{REDACTED}",
	}
    // Set the address translator for this cluster to our ComposeAddressTranslator
	cluster.AddressTranslator = cat
	// This seems to be necessary to avoid gocql attempting to fetch info from cluster about 
	// the seed node's external IP which fails with a warning.
	cluster.IgnorePeerAddr = true

    cluster.Keyspace = "test"
	session, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	// Insert a tweet
	if err := session.Query(`INSERT INTO tweet (timeline, id, text) VALUES (?, ?, ?)`,
		"me", gocql.TimeUUID(), "hello hello hello").Exec(); err != nil {
		log.Fatal(err)
	}
}
```
