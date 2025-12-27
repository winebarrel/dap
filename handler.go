package dap

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	auth "github.com/abbot/go-http-auth"
	"github.com/gorilla/securecookie"
)

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK") //nolint:errcheck
}

type AuthHandler struct {
	*AuthOptions
	backend       *url.URL
	proxy         *httputil.ReverseProxy
	sc            *securecookie.SecureCookie
	cookieName    string
	authenticator *auth.DigestAuth
}

func NewAuthHandler(backend *url.URL, options *AuthOptions) http.Handler {
	proxy := httputil.NewSingleHostReverseProxy(backend)
	sc := securecookie.New([]byte(options.CookieHashKey), nil)
	cookieName := fmt.Sprintf("_dap_%s", options.Realm)

	authenticator := auth.NewDigestAuthenticator(options.Realm, func(user, realm string) string {
		if secret, ok := options.Users[user]; ok {
			return secret
		}

		return ""
	})

	return &AuthHandler{
		AuthOptions:   options,
		backend:       backend,
		proxy:         proxy,
		sc:            sc,
		cookieName:    cookieName,
		authenticator: authenticator,
	}
}

func (h *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(h.cookieName)

	if cookie != nil || err == nil {
		var authInfo auth.Info
		err = h.sc.Decode(h.cookieName, cookie.Value, &authInfo)

		if err != nil || !authInfo.Authenticated {
			if err != nil {
				log.Println(err)
			}

			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
	} else {
		ctx := h.authenticator.NewContext(context.Background(), r)
		authInfo := auth.FromContext(ctx)

		if authInfo != nil {
			authInfo.UpdateHeaders(w.Header())
		}

		if authInfo == nil || !authInfo.Authenticated {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if err := h.setCookie(w, authInfo); err != nil {
			log.Println(err)

			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	r.Host = h.backend.Host
	h.proxy.ServeHTTP(w, r)
}

func (h *AuthHandler) setCookie(w http.ResponseWriter, authInfo *auth.Info) error {
	encoded, err := h.sc.Encode(h.cookieName, *authInfo)

	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     h.cookieName,
		Value:    encoded,
		HttpOnly: true,
		Domain:   h.CookieDomain,
		Secure:   h.CookieSecure,
	}

	if h.CookieSameSite != nil {
		cookie.SameSite = h.CookieSameSite.Value()
	}

	http.SetCookie(w, cookie)

	return nil
}
