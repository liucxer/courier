package b

import (
	"context"
	"github.com/liucxer/courier/httptransport/httpx"
)

// +gengo:jsonschema=false
// +gengo:validator
type GetByID struct {
	httpx.MethodGet `path:"/xxx/:id"`
	// ID
	ID string `in:"path" name:"id"`
	// PullPolicy
	PullPolicy PullPolicy `in:"query" name:"pullPolicy,omitempty"`
	// Name
	Name string `in:"query" validate:"@string[4,]" name:"name,omitempty"`
	// Label
	Label []string `in:"query" name:"label,omitempty"`
}

func (r *GetByID) Output(context.Context) (interface{}, error) {
	return &Data{ID: r.ID}, nil
}

type Data struct {
	ID string `json:"id"`
}
