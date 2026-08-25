package audit

import (
	"time"

	"github.com/google/uuid"
)

const (
	KindStart       = "start"
	KindStop        = "stop"
	KindInterlock   = "interlock"
	KindMaintenance = "maintenance"
	KindQuota       = "quota"
)

type Event struct {
	ID     string
	At     time.Time
	Kind   string
	BeltID string
	Detail string
}

func NewEvent(kind, beltID, detail string) Event {
	return Event{
		ID:     uuid.NewString(),
		At:     time.Now().UTC(),
		Kind:   kind,
		BeltID: beltID,
		Detail: detail,
	}
}
