package profile

import "time"

// AuditEventType describes the kind of change made to a profile.
type AuditEventType string

const (
	AuditEventCreated AuditEventType = "created"
	AuditEventUpdated AuditEventType = "updated"
	AuditEventDeleted AuditEventType = "deleted"
	AuditEventCopied  AuditEventType = "copied"
	AuditEventRenamed AuditEventType = "renamed"
)

// AuditEvent records a single change to a profile in the store.
type AuditEvent struct {
	Timestamp   time.Time      `json:"timestamp"`
	EventType   AuditEventType `json:"event_type"`
	ProfileName string         `json:"profile_name"`
	Actor       string         `json:"actor,omitempty"`
	Detail      string         `json:"detail,omitempty"`
}

// AuditLog holds an ordered list of audit events.
type AuditLog struct {
	Events []AuditEvent `json:"events"`
}

// Append adds a new event to the log.
func (a *AuditLog) Append(event AuditEvent) {
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}
	a.Events = append(a.Events, event)
}

// FilterByProfile returns all events for a specific profile name.
func (a *AuditLog) FilterByProfile(name string) []AuditEvent {
	var out []AuditEvent
	for _, e := range a.Events {
		if e.ProfileName == name {
			out = append(out, e)
		}
	}
	return out
}

// FilterByType returns all events of a given type.
func (a *AuditLog) FilterByType(t AuditEventType) []AuditEvent {
	var out []AuditEvent
	for _, e := range a.Events {
		if e.EventType == t {
			out = append(out, e)
		}
	}
	return out
}

// Last returns the most recent event, or nil if the log is empty.
func (a *AuditLog) Last() *AuditEvent {
	if len(a.Events) == 0 {
		return nil
	}
	e := a.Events[len(a.Events)-1]
	return &e
}
