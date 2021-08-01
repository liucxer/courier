package routes

import (
	"github.com/liucxer/courier/courier"
	"github.com/liucxer/courier/httptransport"
	"github.com/liucxer/courier/httptransport/openapi"
)

var RootRouter = courier.NewRouter(httptransport.BasePath("/demo"))

func init() {
	RootRouter.Register(openapi.OpenAPIRouter)
}
