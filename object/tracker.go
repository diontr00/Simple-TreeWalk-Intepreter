package object

// Keep track of variable
type Tracker struct {
	store map[string]Object
}

func (t *Tracker) Get(name string) (Object, bool) {
	obj, ok := t.store[name]
	return obj, ok
}
func (t *Tracker) Set(name string, val Object) Object {
	t.store[name] = val
	return val
}

func NewTracker() *Tracker {
	s := make(map[string]Object)
	return &Tracker{store: s}
}
