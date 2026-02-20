package events

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type Event interface {
	GetID() string
	GetType() string
	GetPayload() interface{}
	ToJSON() ([]byte, error)
}

type BaseEvent struct {
	ID      string      `json:"id"`
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func (e *BaseEvent) GetID() string {
	return e.ID
}

func (e *BaseEvent) GetType() string {
	return e.Type
}

func (e *BaseEvent) GetPayload() interface{} {
	return e.Payload
}

func (e *BaseEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

const (
	EventTypePolicyReload  = "policy.reload"
	EventTypePolicyCreated = "policy.created"
	EventTypePolicyUpdated = "policy.updated"
	EventTypePolicyDeleted = "policy.deleted"
	EventTypeUserLogin     = "user.login"
	EventTypeUserLogout    = "user.logout"
)

type PolicyReloadEventPayload struct {
	TriggeredBy string `json:"triggered_by"`
	Timestamp   string `json:"timestamp"`
	Reason      string `json:"reason"`
}

func NewPolicyReloadEvent(triggeredBy, reason string) Event {
	return &BaseEvent{
		ID:   uuid.New().String(),
		Type: EventTypePolicyReload,
		Payload: PolicyReloadEventPayload{
			TriggeredBy: triggeredBy,
			Timestamp:   time.Now().Format(time.RFC3339),
			Reason:      reason,
		},
	}
}

type PolicyChangedEventPayload struct {
	Action      string `json:"action"`
	Subject     string `json:"subject"`
	Object      string `json:"object"`
	ActionType  string `json:"action"`
	Owner       string `json:"owner"`
	PerformedBy string `json:"performed_by"`
	Timestamp   string `json:"timestamp"`
}

func NewPolicyCreatedEvent(subject, object, actionType, owner, performedBy string) Event {
	return &BaseEvent{
		ID:   uuid.New().String(),
		Type: EventTypePolicyCreated,
		Payload: PolicyChangedEventPayload{
			Action:      "created",
			Subject:     subject,
			Object:      object,
			ActionType:  actionType,
			Owner:       owner,
			PerformedBy: performedBy,
			Timestamp:   time.Now().Format(time.RFC3339),
		},
	}
}

func NewPolicyDeletedEvent(subject, object, actionType, owner, performedBy string) Event {
	return &BaseEvent{
		ID:   uuid.New().String(),
		Type: EventTypePolicyDeleted,
		Payload: PolicyChangedEventPayload{
			Action:      "deleted",
			Subject:     subject,
			Object:      object,
			ActionType:  actionType,
			Owner:       owner,
			PerformedBy: performedBy,
			Timestamp:   time.Now().Format(time.RFC3339),
		},
	}
}
