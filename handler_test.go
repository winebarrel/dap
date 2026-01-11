package dap_test

import (
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"testing"

	dac "github.com/Snawoot/go-http-digest-auth-client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/dap"
)

func TestHandler_OK(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	var server *httptest.Server

	{
		u, _ := url.Parse("http://example.com")
		h := dap.NewAuthHandler(u, &dap.AuthOptions{
			Users:         map[string]string{"john": "b98e16cbc3d01734b264adba7baa3bf9"},
			Realm:         "example.com",
			CookieHashKey: "my-secret",
		})

		server = httptest.NewServer(h)
		t.Cleanup(server.Close)
	}

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Transport: dac.NewDigestTransport("john", "hello", http.DefaultTransport),
		Jar:       jar,
	}

	// Digest auth
	resp, err := client.Get(server.URL)
	require.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	// Cookie auth
	client.Transport = http.DefaultTransport
	resp, err = client.Get(server.URL)
	require.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	// No cookie
	client.Jar = nil
	resp, err = client.Get(server.URL)
	require.NoError(err)
	assert.Equal(http.StatusUnauthorized, resp.StatusCode)
}

func TestHandler_UnauthorizedUser(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	var server *httptest.Server

	{
		u, _ := url.Parse("http://example.com")
		h := dap.NewAuthHandler(u, &dap.AuthOptions{
			Users:         map[string]string{"john": "b98e16cbc3d01734b264adba7baa3bf9"},
			Realm:         "example.com",
			CookieHashKey: "my-secret",
		})

		server = httptest.NewServer(h)
		t.Cleanup(server.Close)
	}

	client := &http.Client{
		Transport: dac.NewDigestTransport("johnx", "hello", http.DefaultTransport),
	}
	resp, err := client.Get(server.URL)
	require.NoError(err)

	assert.Equal(http.StatusUnauthorized, resp.StatusCode)
}

func TestHandler_UnauthorizedPassword(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	var server *httptest.Server

	{
		u, _ := url.Parse("http://example.com")
		h := dap.NewAuthHandler(u, &dap.AuthOptions{
			Users:         map[string]string{"john": "b98e16cbc3d01734b264adba7baa3bf9"},
			Realm:         "example.com",
			CookieHashKey: "my-secret",
		})

		server = httptest.NewServer(h)
		t.Cleanup(server.Close)
	}

	client := &http.Client{
		Transport: dac.NewDigestTransport("john", "wrongpassword", http.DefaultTransport),
	}
	resp, err := client.Get(server.URL)
	require.NoError(err)

	assert.Equal(http.StatusUnauthorized, resp.StatusCode)
}

func TestHandler_UnauthorizedRealm(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	var server *httptest.Server

	{
		u, _ := url.Parse("http://example.com")
		h := dap.NewAuthHandler(u, &dap.AuthOptions{
			Users:         map[string]string{"john": "b98e16cbc3d01734b264adba7baa3bf9"},
			Realm:         "x.example.com",
			CookieHashKey: "my-secret",
		})

		server = httptest.NewServer(h)
		t.Cleanup(server.Close)
	}

	client := &http.Client{
		Transport: dac.NewDigestTransport("john", "hello", http.DefaultTransport),
	}
	resp, err := client.Get(server.URL)
	require.NoError(err)

	assert.Equal(http.StatusUnauthorized, resp.StatusCode)
}
