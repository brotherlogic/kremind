package main

import (
	ghbclient "github.com/brotherlogic/githubridge/client"
	rstore_client "github.com/brotherlogic/rstore/client"
)

type Server struct {
	rclient  *rstore_client.RStoreClient
	ghclient *ghbclient.GithubridgeClient
}
