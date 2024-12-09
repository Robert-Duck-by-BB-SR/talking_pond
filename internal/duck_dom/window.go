package duckdom

type Window struct {
	ActiveChildId int
	Children      []Renderable
	Pos           Position
}

func (self *Window) Render() string {
	return ""
}
func (self *Window) SetStyle(string) {}

func (self *Window) Active() Renderable { return self.Children[self.ActiveChildId] }
func (self *Window) SetActive(id int)   { self.ActiveChildId = id }
func (self *Window) ActiveIndex() int   { return self.ActiveChildId }

func (self *Window) GetPos() Position { return self.Pos }
