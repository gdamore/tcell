package tcell

import "time"

type EventRaw struct {
	t    time.Time
	code string // The escape code
}

// When returns the time when this EventMouse was created.
func (ev *EventRaw) When() time.Time {
	return ev.t
}

func (ev *EventRaw) EscapeCode() string {
	return ev.code
}

func NewEventRaw(code string) *EventRaw {
	return &EventRaw{
		t:    time.Now(),
		code: code,
	}
}
