package tcell

import "time"

const (
	PasteBegin = "\x1b[200~"
	PasteEnd   = "\x1b[201~"
)

type EventKeyPasteBegin struct {
	t time.Time
}

var _ Event = (*EventKeyPasteBegin)(nil)

// When returns the time when this EventKeyPasteBegin was created.
func (ev *EventKeyPasteBegin) When() time.Time {
	return ev.t
}

func NewEventKeyPasteBegin() *EventKeyPasteBegin {
	return &EventKeyPasteBegin{t: time.Now()}
}

type EventKeyPasteEnd struct {
	t time.Time
}

var _ Event = (*EventKeyPasteEnd)(nil)

// When returns the time when this EventKeyPasteEnd was created.
func (ev *EventKeyPasteEnd) When() time.Time {
	return ev.t
}

func NewEventKeyPasteEnd() *EventKeyPasteEnd {
	return &EventKeyPasteEnd{t: time.Now()}
}
