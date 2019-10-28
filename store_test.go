package fstorage

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"github.com/stretchr/testify/assert"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/models"
	"log"
	"os"
	"testing"
)

var c *firestore.Client

func TestMain(m *testing.M) {
	project, ok := os.LookupEnv("PROJECT_ID")
	if !ok {
		log.Println("skipping tests as PROJECT_ID env variable is missing")
		return
	}
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: project}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalln(err)
	}
	c, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	os.Exit(func() int {
		defer c.Close()
		return m.Run()
	}())
}

func TestStoreClient(t *testing.T) {
	client := New(c, "tests")
	type holder struct {
		key string
		get func(string) (oauth2.TokenInfo, error)
		del func(string) error
	}
	tokens := map[*models.Token]holder{
		{Access: "access"}:   {key: "access", get: client.GetByAccess, del: client.RemoveByAccess},
		{Code: "code"}:       {key: "code", get: client.GetByCode, del: client.RemoveByCode},
		{Refresh: "refresh"}: {key: "refresh", get: client.GetByRefresh, del: client.RemoveByRefresh},
	}
	for i, h := range tokens {
		err := client.Create(i)
		assert.Nil(t, err)

		tok, err := h.get(h.key)
		assert.Nil(t, err)
		assert.Equal(t, i, tok)

		err = h.del(h.key)
		assert.Nil(t, err)

		_, err = h.get(h.key)
		assert.NotNil(t, err)

		err = h.del(h.key)
		assert.NotNil(t, err)
	}
}
