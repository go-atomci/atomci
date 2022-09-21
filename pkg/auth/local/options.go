package local

type Option func(*Provider)

func Name(name string) Option {
	return func(p *Provider) {
		p.name = name
	}
}

func Email(email string) Option {
	return func(p *Provider) {
		p.email = email
	}
}

func User(user string) Option {
	return func(p *Provider) {
		p.user = user
	}
}

func Password(password string) Option {
	return func(p *Provider) {
		p.password = password
	}
}
