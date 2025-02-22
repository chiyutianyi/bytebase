package api

import (
	"context"
	"encoding/json"
)

// Activity type
type ActivityType string

const (
	// Issue related
	ActivityIssueCreate              ActivityType = "bb.issue.create"
	ActivityIssueCommentCreate       ActivityType = "bb.issue.comment.create"
	ActivityIssueFieldUpdate         ActivityType = "bb.issue.field.update"
	ActivityIssueStatusUpdate        ActivityType = "bb.issue.status.update"
	ActivityPipelineTaskStatusUpdate ActivityType = "bb.pipeline.task.status.update"

	// Member related
	ActivityMemberCreate     ActivityType = "bb.member.create"
	ActivityMemberRoleUpdate ActivityType = "bb.member.role.update"
	ActivityMemberActivate   ActivityType = "bb.member.activate"
	ActivityMemberDeactivate ActivityType = "bb.member.deactivate"
)

func (e ActivityType) String() string {
	switch e {
	case ActivityIssueCreate:
		return "bb.issue.create"
	case ActivityIssueCommentCreate:
		return "bb.issue.comment.create"
	case ActivityIssueFieldUpdate:
		return "bb.issue.field.update"
	case ActivityIssueStatusUpdate:
		return "bb.issue.status.update"
	case ActivityPipelineTaskStatusUpdate:
		return "bb.pipeline.task.status.update"
	case ActivityMemberCreate:
		return "bb.member.create"
	case ActivityMemberRoleUpdate:
		return "bb.member.role.update"
	case ActivityMemberActivate:
		return "bb.member.activate"
	case ActivityMemberDeactivate:
		return "bb.member.deactivate"
	}
	return "bb.activity.unknown"
}

type ActivityLevel string

const (
	ACTIVITY_INFO    ActivityLevel = "INFO"
	ACTIVITY_WARNING ActivityLevel = "WARNING"
	ACTIVITY_ERROR   ActivityLevel = "ERROR"
)

func (e ActivityLevel) String() string {
	switch e {
	case ACTIVITY_INFO:
		return "INFO"
	case ACTIVITY_WARNING:
		return "WARNING"
	case ACTIVITY_ERROR:
		return "ERROR"
	}
	return "UNKNOWN"
}

// These payload types are only used when marshalling to the json format for saving into the database.
// So we annotate with json tag using camelCase naming which is consistent with normal
// json naming convention. More importantly, frontend code can simply use JSON.parse to
// convert to the expected struct there.
type ActivityIssueCreatePayload struct {
	// Used by inbox to display info without paying the join cost
	IssueName string `json:"issueName"`
}

type ActivityIssueCommentCreatePayload struct {
	// Used by inbox to display info without paying the join cost
	IssueName string `json:"issueName"`
}

type ActivityIssueFieldUpdatePayload struct {
	FieldId  IssueFieldId `json:"fieldId"`
	OldValue string       `json:"oldValue,omitempty"`
	NewValue string       `json:"newValue,omitempty"`
	// Used by inbox to display info without paying the join cost
	IssueName string `json:"issueName"`
}

type ActivityIssueStatusUpdatePayload struct {
	OldStatus IssueStatus `json:"oldStatus,omitempty"`
	NewStatus IssueStatus `json:"newStatus,omitempty"`
	// Used by inbox to display info without paying the join cost
	IssueName string `json:"issueName"`
}

type ActivityPipelineTaskStatusUpdatePayload struct {
	TaskId    int        `json:"taskId"`
	OldStatus TaskStatus `json:"oldStatus,omitempty"`
	NewStatus TaskStatus `json:"newStatus,omitempty"`
	// Used by inbox to display info without paying the join cost
	IssueName string `json:"issueName"`
	TaskName  string `json:"taskName"`
}

type ActivityMemberCreatePayload struct {
	PrincipalId    int          `json:"principalId"`
	PrincipalName  string       `json:"principalName"`
	PrincipalEmail string       `json:"principalEmail"`
	MemberStatus   MemberStatus `json:"memberStatus"`
	Role           Role         `json:"role"`
}

type ActivityMemberRoleUpdatePayload struct {
	PrincipalId    int    `json:"principalId"`
	PrincipalName  string `json:"principalName"`
	PrincipalEmail string `json:"principalEmail"`
	OldRole        Role   `json:"oldRole"`
	NewRole        Role   `json:"newRole"`
}

type ActivityMemberActivateDeactivatePayload struct {
	PrincipalId    int    `json:"principalId"`
	PrincipalName  string `json:"principalName"`
	PrincipalEmail string `json:"principalEmail"`
	Role           Role   `json:"role"`
}

type Activity struct {
	ID int `jsonapi:"primary,activity"`

	// Standard fields
	CreatorId int
	Creator   *Principal `jsonapi:"attr,creator"`
	CreatedTs int64      `jsonapi:"attr,createdTs"`
	UpdaterId int
	Updater   *Principal `jsonapi:"attr,updater"`
	UpdatedTs int64      `jsonapi:"attr,updatedTs"`
	// Related fields
	// The object where this activity belongs
	// e.g if Type is "bb.issue.xxx", then this field refers to the corresponding issue's id.
	ContainerId int `jsonapi:"attr,containerId"`

	// Domain specific fields
	Type    ActivityType  `jsonapi:"attr,actionType"`
	Level   ActivityLevel `jsonapi:"attr,level"`
	Comment string        `jsonapi:"attr,comment"`
	Payload string        `jsonapi:"attr,payload"`
}

type ActivityCreate struct {
	// Standard fields
	// Value is assigned from the jwt subject field passed by the client.
	CreatorId int

	// Domain specific fields
	ContainerId int          `jsonapi:"attr,containerId"`
	Type        ActivityType `jsonapi:"attr,actionType"`
	Level       ActivityLevel
	Comment     string `jsonapi:"attr,comment"`
	Payload     string `jsonapi:"attr,payload"`
}

type ActivityFind struct {
	ID *int

	// Domain specific fields
	ContainerId *int
}

func (find *ActivityFind) String() string {
	str, err := json.Marshal(*find)
	if err != nil {
		return err.Error()
	}
	return string(str)
}

type ActivityPatch struct {
	ID int `jsonapi:"primary,activityPatch"`

	// Standard fields
	// Value is assigned from the jwt subject field passed by the client.
	UpdaterId int

	// Domain specific fields
	Comment *string `jsonapi:"attr,comment"`
}

type ActivityDelete struct {
	ID int

	// Standard fields
	// Value is assigned from the jwt subject field passed by the client.
	DeleterId int
}

type ActivityService interface {
	CreateActivity(ctx context.Context, create *ActivityCreate) (*Activity, error)
	FindActivityList(ctx context.Context, find *ActivityFind) ([]*Activity, error)
	FindActivity(ctx context.Context, find *ActivityFind) (*Activity, error)
	PatchActivity(ctx context.Context, patch *ActivityPatch) (*Activity, error)
	DeleteActivity(ctx context.Context, delete *ActivityDelete) error
}
