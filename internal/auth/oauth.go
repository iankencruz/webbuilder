package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

type OAuth2Provider struct {
	oauth2Config *oauth2.Config
	userInfoURL  string
	postLoginURL string
}

func NewOAuth2Provider(name, clientID, clientSecret, redirectURI, issuerURI, postLoginURI string, scopes []string) *OAuth2Provider {
	// default to OIDC discovery-style endpoints
	authURL := issuerURI + "/oauth/v2/authorize"
	tokenURL := issuerURI + "/oauth/v2/token"
	userInfoURL := issuerURI + "/oidc/v1/userinfo"

	// provider-specific endpoint overrides
	switch name {
	case "google":
		authURL = "https://accounts.google.com/o/oauth2/v2/auth"
		tokenURL = "https://oauth2.googleapis.com/token"
		userInfoURL = "https://www.googleapis.com/oauth2/v3/userinfo"
	case "github":
		authURL = "https://github.com/login/oauth/authorize"
		tokenURL = "https://github.com/login/oauth/access_token"
		userInfoURL = "https://api.github.com/user"
	case "discord":
		authURL = "https://discord.com/api/oauth2/authorize"
		tokenURL = "https://discord.com/api/oauth2/token"
		userInfoURL = "https://discord.com/api/users/@me"
	}

	cfg := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURI,
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:   authURL,
			TokenURL:  tokenURL,
			AuthStyle: oauth2.AuthStyleInHeader,
		},
	}

	return &OAuth2Provider{
		oauth2Config: cfg,
		userInfoURL:  userInfoURL,
		postLoginURL: postLoginURI,
	}
}

func (o *OAuth2Provider) AuthenticationURL(state string) string {
	return o.oauth2Config.AuthCodeURL(state)
}

func (o *OAuth2Provider) RedirectURI() string {
	return o.oauth2Config.RedirectURL
}

func (o *OAuth2Provider) PostLoginURL() string {
	return o.postLoginURL
}

func (o *OAuth2Provider) ExchangeCode(ctx context.Context, code string) (*UserProfile, error) {
	token, err := o.oauth2Config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("oauth2 exchange failed: %w", err)
	}

	client := o.oauth2Config.Client(ctx, token)
	resp, err := client.Get(o.userInfoURL)
	if err != nil {
		return nil, fmt.Errorf("failed fetching user info from provider: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user info endpoint returned bad HTTP status: %d", resp.StatusCode)
	}

	var claims struct {
		Subject    string `json:"sub"`
		Email      string `json:"email"`
		GivenName  string `json:"given_name"`
		FamilyName string `json:"family_name"`
		Avatar     string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&claims); err != nil {
		return nil, fmt.Errorf("failed decoding user info payload: %w", err)
	}

	return &UserProfile{
		Subject:    claims.Subject,
		Email:      claims.Email,
		GivenName:  claims.GivenName,
		FamilyName: claims.FamilyName,
		Avatar:     claims.Avatar,
	}, nil
}
