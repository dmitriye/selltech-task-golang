package app

const (
	S_EMPTY    = 0
	S_UPDATING = 1
	S_OK       = 2
)

type AppState struct {
	state int
}

var names = map[int]string{
	S_EMPTY:    "empty",
	S_UPDATING: "updating",
	S_OK:       "ok",
}

func (o *AppState) SetState(s int) {
	if _, ok := names[s]; !ok {
		panic("can not set unknown state")
	}
	o.state = s
}

func (o *AppState) GetState() int {
	return o.state
}

func (o *AppState) GetName() string {
	return names[o.state]
}

func (o *AppState) IsEmpty() bool {
	return o.state == S_EMPTY
}

func (o *AppState) IsUpdating() bool {
	return o.state == S_UPDATING
}

func (o *AppState) IsOk() bool {
	return o.state == S_OK
}

func NewAppState() *AppState {
	return &AppState{S_EMPTY}
}
