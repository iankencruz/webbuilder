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
	cfg := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURI,
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  issuerURI + "/protocol/openid-connect/auth",
			TokenURL: issuerURI + "/protocol/openid-connect/token",
		},
	}

	userInfo := issuerURI + "/protocol/openid-connect/userinfo"

	// Fallback override for providers that don't follow the OIDC spec
	switch name {
	case "google":
		userInfo = "https://www.googleapis.com/oauth2/v3/userinfo"
	case "github":
		userInfo = "https://api.github.com/user"
	case "discord":
		userInfo = "https://discord.com/api/users/@me"
	}

	return &OAuth2Provider{
		oauth2Config: cfg,
		userInfoURL:  userInfo,
		postLoginURL: postLoginURI,
	}
}

func (o *OAuth2Provider) AuthenticationURL(state string) string {
	return o.oauth2Config.AuthCodeURL(state)
}

func (o *OAuth2Provider) RedirectURI() string {
	return o.oauth2Config.RedirectURL
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
