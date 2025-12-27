# dap

Digest authentication proxy that stores authentication information in a cookie.

## Usage

```
Usage: dap --backends=KEY=VALUE;... --users=KEY=VALUE;... --realm=STRING --cookie-hash-key=STRING [flags]

Flags:
  -h, --help                                 Show help.
  -b, --backends=KEY=VALUE;...               Port and backend mapping ($DAP_BACKENDS).
      --health-check="_health"               Health check path ($DAP_HEALTH_CHECK).
  -u, --users=KEY=VALUE;...                  User and secret mapping ($DAP_USERS).
  -r, --realm=STRING                         Auth realm ($DAP_REALM).
  -k, --cookie-hash-key=STRING               Hash key for cookie encryption ($DAP_COOKIE_HASH_KEY).
      --cookie-domain=STRING                 Cookie 'domain' attr ($DAP_COOKIE_DOMAIN).
      --[no-]cookie-secure                   Cookie 'secure' attr ($DAP_COOKIE_SECURE).
      --cookie-same-site=COOKIE-SAME-SITE    Cookie 'samesite' attr ($DAP_COOKIE_SAMESITE).
      --version
```

```sh
$ export DAP_COOKIE_HASH_KEY=my-secret
$ export DAP_BACKENDS='8080=https://example.com;8081=https://www.yahoo.co.jp'
$ export DAP_REALM=example.com
# echo -n 'test:example.com:hello' | sha1sum
$ export DAP_USERS='john=b98e16cbc3d01734b264adba7baa3bf9'
$ go run ./cmd/dap
```

```sh
$ curl localhost:8080
Unauthorized
$ curl localhost:8081
Unauthorized

$ curl --digest --user "john:hello" localhost:8080
...<title>Example Domain</title>...
$ curl --digest --user "john:hello" localhost:8081
...<title>Yahoo! JAPAN</title>...

$ curl -c cookie.txt --digest --user "john:hello" localhost:8080
$ curl -b cookie.txt "john:hello" localhost:8080
...<title>Example Domain</title>...
```
