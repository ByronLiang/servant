package servant

type Servant struct {
	opt option
}

func NewServant()  {
	o := option{name: "s"}
	s := Servant{opt:o}
	s.Name()
}

func (s *Servant) Name() string {
	return s.opt.name
}
