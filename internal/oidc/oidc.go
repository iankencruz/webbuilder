package oidc

import (
	gooidc "github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type Provider struct {
	Config       *oauth2.Config
	Verifier     *gooidc.IDTokenVerifier
	PostLoginURL string
}

type Registry struct {
	providers map[string]*Provider
}

func (r *Registry) Get(name string) (*Provider, bool) {
	p, ok := r.providers[name]
	return p, ok
}

func (r *Registry) Register(name string, p *Provider) {
	r.providers[name] = p
}
