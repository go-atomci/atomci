package ldap

type Option func(*Provider)

func BaseDN(baseDN string) Option {
	return func(p *Provider) {
		p.baseDN = baseDN
	}
}

func Host(host string) Option {
	return func(p *Provider) {
		p.host = host
	}
}

func Port(port int) Option {
	return func(p *Provider) {
		p.port = port
	}
}

func BindDN(bindDN string) Option {
	return func(p *Provider) {
		p.bindDN = bindDN
	}
}

func BindPassword(bindPassword string) Option {
	return func(p *Provider) {
		p.bindPassword = bindPassword
	}
}

func UserFilter(userFilter string) Option {
	return func(p *Provider) {
		p.userFilter = userFilter
	}
}
