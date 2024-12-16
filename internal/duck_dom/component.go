package duckdom

type Component struct {
	Position
	Styles
	Content string
	Buffer  string
	// NOTE: we should really think about it
	// maybe it would be better if we just made a bunch of functions
	// that take *Component as an input and does some actions with it
	Action func()
}

func (self *Component) ExecuteAction() {
	self.Action()
}
