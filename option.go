package servant

type option struct {
	name string
}

type Option func(o *option)

func Name(name string) Option {
	return func(o *option) {
		o.name = name
	}
}
