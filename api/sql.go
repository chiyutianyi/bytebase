package api

import (
	"context"

	"github.com/bytebase/bytebase/plugin/db"
)

type ConnectionInfo struct {
	DBType     db.Type `jsonapi:"attr,dbType"`
	Host       string  `jsonapi:"attr,host"`
	Port       string  `jsonapi:"attr,port"`
	Username   string  `jsonapi:"attr,username"`
	Password   string  `jsonapi:"attr,password"`
	InstanceId *int    `jsonapi:"attr,instanceId"`
}

type SqlSyncSchema struct {
	InstanceId int `jsonapi:"attr,instanceId"`
}

type SqlResultSet struct {
	// SQL operation may fail for connection issue and there is no proper http status code for it, so we return error in the response body.
	Error string `jsonapi:"attr,error"`
}

type SqlService interface {
	Ping(ctx context.Context, config *ConnectionInfo) (*SqlResultSet, error)
}
