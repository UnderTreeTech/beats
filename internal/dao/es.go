package dao

import (
	"context"
)

var mapping = `
	{
		"settings":{
		"number_of_shards":1,
		"number_of_replicas":0
		},
		"mappings":{
			"properties":{
					"trace_id":{
						"type":"keyword"
					},
					"level":{
						"type":"keyword"
					},
					"ts":{
						"type":"keyword"
					},
					"app":{
						"type":"keyword"
					},
					"method":{
						"type":"keyword"
					},
					"req":{
						"type":"keyword"
					},
					"log":{
						"type":"text",
						"index": false
					}
			}
		}
	}
	`

// CreateIndex create an index
func (d *dao) CreateIndex(ctx context.Context, index string) (err error) {
	err = d.es.CreateIndex(ctx, index, mapping)
	return
}

// CreateIndex create an index
func (d *dao) CreateDocs(ctx context.Context, index string, docs []interface{}) (err error) {
	_, err = d.es.CreateDocs(ctx, index, docs)
	return
}

// ExistIndex check index exists
func (d *dao) ExistIndex(ctx context.Context, index string) (exist bool, err error) {
	return d.es.ExistIndex(ctx, index)
}
