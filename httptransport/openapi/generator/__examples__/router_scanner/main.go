package main

import (
	"context"

	"github.com/liucxer/courier/courier"
	"github.com/liucxer/courier/httptransport/httpx"
	"github.com/liucxer/courier/httptransport/openapi/generator/__examples__/router_scanner/auth"
	"github.com/liucxer/courier/httptransport/openapi/generator/__examples__/router_scanner/group"

	"github.com/liucxer/courier/httptransport"
)

type Get struct {
	httpx.MethodGet `path:"/:id"`

	ID string `name:"id" in:"path"`
}

func (get Get) Output(ctx context.Context) (result interface{}, err error) {
	return
}

var Router = courier.NewRouter(httptransport.Group("/root"))

func main() {
	Router.Register(group.GroupRouter)
	Router.Register(courier.NewRouter(auth.Auth{}, Get{}))

	ht := &httptransport.HttpTransport{
		Port: 8080,
	}
	ht.SetDefaults()

	courier.Run(Router, ht)
}
