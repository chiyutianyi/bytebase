package api

import (
	"context"
	"encoding/json"

	"github.com/bytebase/bytebase/plugin/db"
)

type Instance struct {
	ID int `jsonapi:"primary,instance"`

	// Standard fields
	RowStatus RowStatus `jsonapi:"attr,rowStatus"`
	CreatorId int
	Creator   *Principal `jsonapi:"attr,creator"`
	CreatedTs int64      `jsonapi:"attr,createdTs"`
	UpdaterId int
	Updater   *Principal `jsonapi:"attr,updater"`
	UpdatedTs int64      `jsonapi:"attr,updatedTs"`

	// Related fields
	EnvironmentId int
	Environment   *Environment `jsonapi:"relation,environment"`

	// Domain specific fields
	Name         string  `jsonapi:"attr,name"`
	Engine       db.Type `jsonapi:"attr,engine"`
	ExternalLink string  `jsonapi:"attr,externalLink"`
	Host         string  `jsonapi:"attr,host"`
	Port         string  `jsonapi:"attr,port"`
	Username     string  `jsonapi:"attr,username"`
	// Password is not returned to the client
	Password string
}

type InstanceCreate struct {
	// Standard fields
	// Value is assigned from the jwt subject field passed by the client.
	CreatorId int

	// Related fields
	EnvironmentId int `jsonapi:"attr,environmentId"`

	// Domain specific fields
	Name         string  `jsonapi:"attr,name"`
	Engine       db.Type `jsonapi:"attr,engine"`
	ExternalLink string  `jsonapi:"attr,externalLink"`
	Host         string  `jsonapi:"attr,host"`
	Port         string  `jsonapi:"attr,port"`
	Username     string  `jsonapi:"attr,username"`
	Password     string  `jsonapi:"attr,password"`
}

type InstanceFind struct {
	ID *int

	// Standard fields
	RowStatus *RowStatus
}

func (find *InstanceFind) String() string {
	str, err := json.Marshal(*find)
	if err != nil {
		return err.Error()
	}
	return string(str)
}

type InstancePatch struct {
	ID int `jsonapi:"primary,instancePatch"`

	// Standard fields
	RowStatus *string `jsonapi:"attr,rowStatus"`
	// Value is assigned from the jwt subject field passed by the client.
	UpdaterId int

	// Domain specific fields
	Name         *string `jsonapi:"attr,name"`
	ExternalLink *string `jsonapi:"attr,externalLink"`
	Host         *string `jsonapi:"attr,host"`
	Port         *string `jsonapi:"attr,port"`
	Username     *string `jsonapi:"attr,username"`
	Password     *string `jsonapi:"attr,password"`
}

// Instance migration schema status
type InstanceMigrationSchemaStatus string

const (
	InstanceMigrationSchemaUnknown  InstanceMigrationSchemaStatus = "UNKNOWN"
	InstanceMigrationSchemaOK       InstanceMigrationSchemaStatus = "OK"
	InstanceMigrationSchemaNotExist InstanceMigrationSchemaStatus = "NOT_EXIST"
)

func (e InstanceMigrationSchemaStatus) String() string {
	switch e {
	case InstanceMigrationSchemaUnknown:
		return "UNKNOWN"
	case InstanceMigrationSchemaOK:
		return "OK"
	case InstanceMigrationSchemaNotExist:
		return "NOT_EXIST"
	}
	return "UNKNOWN"
}

type InstanceMigration struct {
	Status InstanceMigrationSchemaStatus `jsonapi:"attr,status"`
	Error  string                        `jsonapi:"attr,error"`
}

// MigrationHistory is stored in the instance instead of our own data file, so the field
// format is a bit different from the standard format
type MigrationHistory struct {
	ID int `jsonapi:"primary,migrationHistory"`

	// Standard fields
	Creator   string `jsonapi:"attr,creator"`
	CreatedTs int64  `jsonapi:"attr,createdTs"`
	Updater   string `jsonapi:"attr,updater"`
	UpdatedTs int64  `jsonapi:"attr,updatedTs"`

	// Domain specific fields
	Database          string             `jsonapi:"attr,database"`
	Engine            db.MigrationEngine `jsonapi:"attr,engine"`
	Type              db.MigrationType   `jsonapi:"attr,type"`
	Version           string             `jsonapi:"attr,version"`
	Description       string             `jsonapi:"attr,description"`
	Statement         string             `jsonapi:"attr,statement"`
	ExecutionDuration int                `jsonapi:"attr,executionDuration"`
	// This is a string instead of int as the issue id may come from other issue tracking system in the future
	IssueId string `jsonapi:"attr,issueId"`
	Payload string `jsonapi:"attr,payload"`
}

type InstanceService interface {
	// CreateInstance should also create the * database and the admin data source.
	CreateInstance(ctx context.Context, create *InstanceCreate) (*Instance, error)
	FindInstanceList(ctx context.Context, find *InstanceFind) ([]*Instance, error)
	FindInstance(ctx context.Context, find *InstanceFind) (*Instance, error)
	PatchInstance(ctx context.Context, patch *InstancePatch) (*Instance, error)
}
