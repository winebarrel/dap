package dap

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	auth "github.com/abbot/go-http-auth"
	"github.com/gorilla/securecookie"
)

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

type AuthHandler struct {
	*AuthOptions
	Proxy func(http.ResponseWriter, *http.Request)
}

func NewAuthHandler(ctx context.Context, backend *url.URL, options *AuthOptions) http.Handler {
	proxy := httputil.NewSingleHostReverseProxy(backend)
	sc := securecookie.New([]byte(options.CookieHashKey), nil)
	cookieName := fmt.Sprintf("_dap_%s", options.Realm)

	authenticator := auth.NewDigestAuthenticator(options.Realm, func(user, realm string) string {
		if secret, ok := options.Users[user]; ok {
			return secret
		}

		return ""
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(cookieName)

		if cookie != nil || err == nil {
			var authInfo auth.Info
			err = sc.Decode(cookieName, cookie.Value, &authInfo)

			if err != nil || !authInfo.Authenticated {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
		} else {
			ctx := authenticator.NewContext(ctx, r)
			authInfo := auth.FromContext(ctx)

			if authInfo != nil {
				authInfo.UpdateHeaders(w.Header())
			}

			if authInfo == nil || !authInfo.Authenticated {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			encoded, err := sc.Encode(cookieName, *authInfo)

			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			cookie := &http.Cookie{
				Name:     cookieName,
				Value:    encoded,
				HttpOnly: true,
				Domain:   options.CookieDomain,
				Secure:   options.CookieSecure,
			}

			if options.CookieSameSite != nil {
				cookie.SameSite = options.CookieSameSite.Value()
			}

			http.SetCookie(w, cookie)
		}

		r.Host = backend.Host
		proxy.ServeHTTP(w, r)
	})
}
