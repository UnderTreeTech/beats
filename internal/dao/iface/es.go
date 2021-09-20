package iface

import (
	"context"
)

// interface Es
type Es interface {
	CreateIndex(ctx context.Context, index string) error
	CreateDocs(ctx context.Context, index string, docs []interface{}) (err error)
	ExistIndex(ctx context.Context, index string) (exist bool, err error)
}
