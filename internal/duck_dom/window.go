package duckdom

type Window struct {
	ActiveChildId int
	Children      []Renderable
	Pos           Position
}

func (self *Window) Render() string {
	return ""
}

func (self *Window) SetPos(p Position) {
	self.Pos = p
}

func (self *Window) GetPos() Position {
	return self.Pos
}
func (self *Window) SetStyle(string) {}
