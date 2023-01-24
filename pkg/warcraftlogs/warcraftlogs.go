package warcraftlogs

import (
	"context"

	"github.com/farwydi/divot/pkg/database"
	"github.com/hasura/go-graphql-client"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type WarcraftLogs struct {
	client *graphql.Client
	db     database.Database

	saveAnonymous bool
}

func NewWarcraftLogs(
	clientID, clientSecret string,
	db database.Database,
) *WarcraftLogs {
	cfg := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     "https://www.warcraftlogs.com/oauth/token",
	}
	httpClient := oauth2.NewClient(context.Background(), cfg.TokenSource(context.Background()))

	client := graphql.NewClient("https://www.warcraftlogs.com/api/v2/client", httpClient)

	return &WarcraftLogs{
		client: client,
		db:     db,
	}
}
