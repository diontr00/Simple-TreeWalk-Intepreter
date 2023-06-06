package object

// Keep track of variable
type Tracker struct {
	store map[string]Object
	outer *Tracker
}

func (t *Tracker) Get(name string) (Object, bool) {
	obj, ok := t.store[name]
	// Take in account clousure effect , where the enclosed function can access the environment of the outer function (scopes)
	if !ok && t.outer != nil {
		obj, ok = t.outer.Get(name)
	}
	return obj, ok
}
func (t *Tracker) Set(name string, val Object) Object {
	t.store[name] = val
	return val
}

func NewTracker() *Tracker {
	s := make(map[string]Object)
	return &Tracker{store: s, outer: nil}
}

// keeping track of enclosing environment
func NewEnclosedTracker(outer *Tracker) *Tracker {
	tracker := NewTracker()
	tracker.outer = outer
	return tracker
}
