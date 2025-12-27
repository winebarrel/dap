package dap

import (
	"net/http"
	"net/url"
	"strings"
)

type Options struct {
	Backends map[uint]*url.URL `short:"b" required:"" env:"DAP_BACKENDS" help:"Port and backend mapping."`
	AuthOptions
}

type CookieSameSite string

func (v *CookieSameSite) Value() http.SameSite {
	switch strings.ToLower(string(*v)) {
	case "lax":
		return http.SameSiteLaxMode
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteDefaultMode
	}
}

type AuthOptions struct {
	Users          map[string]string `short:"u" required:"" env:"DAP_USERS" help:"User and secret mapping."`
	Realm          string            `short:"r" required:"" env:"DAP_REALM" help:"Auth realm."`
	CookieHashKey  string            `short:"k" required:"" env:"DAP_COOKIE_HASH_KEY" help:"Hash key for cookie encryption."`
	CookieDomain   string            `env:"DAP_COOKIE_DOMAIN" help:"Cookie 'domain' attr."`
	CookieSecure   bool              `negatable:"" env:"DAP_COOKIE_SECURE" help:"Cookie 'secure' attr."`
	CookieSameSite *CookieSameSite   `enum:"lax,strict,none" env:"DAP_COOKIE_SAMESITE" help:"Cookie 'samesite' attr."`
}
