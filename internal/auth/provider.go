package auth

import "context"

// create a user profile structure that holds  user information obtained from
// the OIDC providers: zitadel(priority), google, github, discord
type UserProfile struct {
	ID         string `json:"id"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Email      string `json:"email"`
	Avatar     string `json:"avatar"`
	Subject    string `json:"subject"`
}

type AuthProvider interface {
	AuthenticationURL(state string) string
	ExchangeCode(ctx context.Context, code string) (*UserProfile, error)
	RedirectURI() string
}

type Registry struct {
	providers map[string]AuthProvider
}

// NewRegistry initializes an empty, flexible registration map.
func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[string]AuthProvider),
	}
}

func (r *Registry) Get(name string) (AuthProvider, bool) {
	p, ok := r.providers[name]
	return p, ok
}

func (r *Registry) Register(name string, p AuthProvider) {
	r.providers[name] = p
}
