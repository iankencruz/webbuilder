package oidc

import (
	"context"
	"fmt"

	gooidc "github.com/coreos/go-oidc/v3/oidc"
	"github.com/iankencruz/webbuilder/internal/config"
	"golang.org/x/oauth2"
)

// var standardOIDC = map[string]bool{
// 	"rauthy":  true,
// 	"google":  true,
// 	"zitadel": true,
// }

func NewRegistry(ctx context.Context, cfg *config.Config) (*Registry, error) {
	r := &Registry{providers: make(map[string]*Provider)}
	for name, pc := range cfg.OIDCProvider {
		p, err := buildProvider(ctx, pc)
		if err != nil {
			return nil, fmt.Errorf("oidc provider %s: %w", name, err)
		}
		r.Register(name, p)
	}
	return r, nil
}

func buildProvider(ctx context.Context, pc config.OIDCProvider) (*Provider, error) {
	p, err := gooidc.NewProvider(ctx, pc.IssuerURL)
	if err != nil {
		return nil, err
	}
	return &Provider{
		Config: &oauth2.Config{
			ClientID:     pc.ClientID,
			ClientSecret: pc.ClientSecret,
			RedirectURL:  pc.RedirectURL,
			Endpoint:     p.Endpoint(),
			Scopes:       pc.Scopes,
		},
		Verifier:     p.Verifier(&gooidc.Config{ClientID: pc.ClientID}),
		PostLoginURL: pc.PostLoginURL,
	}, nil
}
