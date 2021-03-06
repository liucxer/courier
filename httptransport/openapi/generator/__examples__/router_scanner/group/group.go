package group

import (
	"context"

	"github.com/liucxer/courier/courier"
	"github.com/liucxer/courier/httptransport/httpx"

	"github.com/liucxer/courier/httptransport"
)

var GroupRouter = courier.NewRouter(httptransport.Group("/group"))
var HeathRouter = courier.NewRouter(&Health{})

func init() {
	GroupRouter.Register(HeathRouter)
}

type Health struct {
	httpx.MethodHead
}

func (Health) MiddleOperators() courier.MiddleOperators {
	return courier.MiddleOperators{
		httptransport.Group("/health"),
	}
}

func (*Health) Output(ctx context.Context) (result interface{}, err error) {
	return
}
