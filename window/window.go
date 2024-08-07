package window

type Window struct {
	EventList chan *UserEvent
}

type UserEvent struct {
	UserID string
	Event  string `json:event`
	Key    string `json:key`
}

func NewWindow() *Window {
	return &Window{
		EventList: make(chan *UserEvent),
	}
}
