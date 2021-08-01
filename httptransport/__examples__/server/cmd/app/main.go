package main

import (
	"github.com/liucxer/courier/courier"
	"github.com/liucxer/courier/httptransport/__examples__/server/cmd/app/routes"

	"github.com/liucxer/courier/httptransport"
)

func main() {
	ht := &httptransport.HttpTransport{
		Port: 8080,
	}
	ht.SetDefaults()

	courier.Run(routes.RootRouter, ht)
}
