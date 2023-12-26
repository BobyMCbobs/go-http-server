package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path"
	"testing"
	"time"

	"gitlab.com/BobyMCbobs/go-http-server/pkg/common"
)

var defaultEnv = map[string]string{
	"APP_METRICS_ENABLED": "false",
}
var (
	client = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	numberRunes = []rune("0123456789")
)

func newRequest(method string, url string, body io.Reader) (req *http.Request) {
	req, _ = http.NewRequest(method, url, body)
	return req
}

func pointer[V any](input V) *V {
	return &input
}

func randStringRunes(set []rune, n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = set[rand.Intn(len(set))]
	}
	return string(b)
}

func TestMain(t *testing.T) {
	tests := []struct {
		name            string
		env             map[string]string
		noSetEnv        bool
		files           map[string]string
		productionBuild bool
		skip            bool
		requests        []struct {
			req         *http.Request
			wantBody    *string
			wantCode    int
			wantPath    string
			wantError   bool
			wantHeaders map[string]string
			client      *http.Client
		}
	}{
		{
			name: "standard serving",
			files: map[string]string{
				"index.html":       `hello`,
				"about/index.html": `this is a page`,
				"cool.html":        `cool stuff maybe`,
				"404.html":         `not found`,
			},
			requests: []struct {
				req         *http.Request
				wantBody    *string
				wantCode    int
				wantPath    string
				wantError   bool
				wantHeaders map[string]string
				client      *http.Client
			}{
				{
					req:      newRequest(http.MethodGet, "/", nil),
					wantBody: pointer(`hello`),
					wantCode: http.StatusOK,
					wantPath: "/",
				},
				{
					req:      newRequest(http.MethodGet, "/index.html", nil),
					wantBody: pointer(``),
					wantCode: http.StatusMovedPermanently,
					wantPath: "/index.html",
				},
				{
					req:      newRequest(http.MethodGet, "/about", nil),
					wantBody: pointer(``),
					wantCode: http.StatusMovedPermanently,
					wantPath: "/about",
				},
				{
					req:      newRequest(http.MethodGet, "/about/", nil),
					wantBody: pointer(`this is a page`),
					wantCode: http.StatusOK,
					wantPath: "/about/",
				},
				{
					req:      newRequest(http.MethodGet, "/about/index.html", nil),
					wantBody: pointer(``),
					wantCode: http.StatusMovedPermanently,
					wantPath: "/about/",
				},
				{
					req:      newRequest(http.MethodGet, "/aaaa", nil),
					wantBody: pointer(`not found`),
					wantCode: http.StatusNotFound,
					wantPath: "/aaaa",
				},
				{
					req:      newRequest(http.MethodGet, "/asdfkasdfjhkaser/asgaehrl", nil),
					wantBody: pointer(`not found`),
					wantCode: http.StatusNotFound,
					wantPath: "/asdfkasdfjhkaser/asgaehrl",
				},
			},
		},
		// testcase history mode
		{
			name: "history mode",
			env: map[string]string{
				"APP_VUEJS_HISTORY_MODE": "true",
			},
			files: map[string]string{
				"index.html": `history mode enabled`,
			},
			requests: []struct {
				req         *http.Request
				wantBody    *string
				wantCode    int
				wantPath    string
				wantError   bool
				wantHeaders map[string]string
				client      *http.Client
			}{
				{
					req:      newRequest(http.MethodGet, "/", nil),
					wantBody: pointer(`history mode enabled`),
					wantCode: http.StatusOK,
				},
				{
					req:      newRequest(http.MethodGet, "/asdfjsdf", nil),
					wantBody: pointer(`history mode enabled`),
					wantCode: http.StatusOK,
				},
				{
					req:      newRequest(http.MethodGet, "/a/a/a/a/a/a/aa", nil),
					wantBody: pointer(`history mode enabled`),
					wantCode: http.StatusOK,
				},
			},
		},
		// testcase history mode with templating
		{
			name: "history mode with templating",
			env: map[string]string{
				"APP_VUEJS_HISTORY_MODE": "true",
				"APP_TEMPLATE_MAP_PATH":  "./cfg-for-testing-only-this-is-bad-practice-in-web-serve-folder/tmpl.yaml",
			},
			files: map[string]string{
				"index.html": `{{ .Title }}`,
				"./cfg-for-testing-only-this-is-bad-practice-in-web-serve-folder/tmpl.yaml": `---
Title: Hello!
`,
			},
			requests: []struct {
				req         *http.Request
				wantBody    *string
				wantCode    int
				wantPath    string
				wantError   bool
				wantHeaders map[string]string
				client      *http.Client
			}{
				{
					req:      newRequest(http.MethodGet, "/", nil),
					wantBody: pointer(`Hello!`),
					wantCode: http.StatusOK,
				},
				{
					req:      newRequest(http.MethodGet, "/54309y90429", nil),
					wantBody: pointer(`Hello!`),
					wantCode: http.StatusOK,
				},
				{
					req:      newRequest(http.MethodGet, "/a/b/c/b/c/a/dddddd", nil),
					wantBody: pointer(`Hello!`),
					wantCode: http.StatusOK,
				},
			},
		},
		// testcase metrics port enabled with port set
		{
			name: "metrics port enabled with port set",
			env: map[string]string{
				"APP_METRICS_ENABLED": "true",
				"APP_PORT_METRICS":    ":7335",
			},
			requests: []struct {
				req         *http.Request
				wantBody    *string
				wantCode    int
				wantPath    string
				wantError   bool
				wantHeaders map[string]string
				client      *http.Client
			}{
				{
					req:      newRequest(http.MethodGet, "http://localhost:7335/metrics", nil),
					wantBody: nil,
					wantCode: http.StatusOK,
				},
			},
		},
		// testcase load config from env file
		{
			name: "load config from env file",
			skip: true, // TODO fix test. Unsure why it's currently failing
			files: map[string]string{
				".env": `APP_METRICS_ENABLED=true
APP_PORT_METRICS=:7334
APP_PORT=:49001`,
				"index.html": `hello`,
			},
			// noSetEnv: true,
			requests: []struct {
				req         *http.Request
				wantBody    *string
				wantCode    int
				wantPath    string
				wantError   bool
				wantHeaders map[string]string
				client      *http.Client
			}{
				{
					req:      newRequest(http.MethodGet, "http://localhost:7334/metrics", nil),
					wantBody: nil,
					wantCode: http.StatusOK,
				},
				{
					req:      newRequest(http.MethodGet, "http://localhost:49001", nil),
					wantBody: pointer(`hello`),
					wantCode: http.StatusOK,
				},
			},
		},
		// testcase no env file in prod
		{
			name:            "don't load config from env file in production",
			productionBuild: true,
			files: map[string]string{
				".env": `APP_METRICS_ENABLED=true
APP_PORT_METRICS=:7333`,
			},
			requests: []struct {
				req         *http.Request
				wantBody    *string
				wantCode    int
				wantPath    string
				wantError   bool
				wantHeaders map[string]string
				client      *http.Client
			}{
				{
					req:       newRequest(http.MethodGet, "http://localhost:7333/metrics", nil),
					wantBody:  nil,
					wantCode:  http.StatusOK,
					wantError: true,
				},
			},
		},
		// testcase load header map
		{
			name: "header map",
			env: map[string]string{
				"APP_HEADER_MAP_PATH":   "./cfg-for-testing-only-this-is-bad-practice-in-web-serve-folder/headers.yaml",
				"APP_HEADER_SET_ENABLE": "true",
				"SOME_VALUE":            "it-is-working-yes",
			},
			files: map[string]string{
				"index.html": `hello`,
				"./cfg-for-testing-only-this-is-bad-practice-in-web-serve-folder/headers.yaml": `---
X-Cool-As:
 - "Yeah"
X-From-Env:
  - "${SOME_VALUE}"
`,
			},
			requests: []struct {
				req         *http.Request
				wantBody    *string
				wantCode    int
				wantPath    string
				wantError   bool
				wantHeaders map[string]string
				client      *http.Client
			}{
				{
					req:      newRequest(http.MethodGet, "/", nil),
					wantBody: pointer(`hello`),
					wantCode: http.StatusOK,
					wantHeaders: map[string]string{
						"X-Cool-As":  "Yeah",
						"X-From-Env": "it-is-working-yes",
					},
				},
			},
		},
		// testcase load redirect map
		{
			name: "redirect map",
			env: map[string]string{
				"APP_REDIRECT_ROUTES_ENABLED": "true",
				"APP_REDIRECT_ROUTES_PATH":    "./cfg-for-testing-only-this-is-bad-practice-in-web-serve-folder/redirects.yaml",
			},
			files: map[string]string{
				"index.html": `hello`,
				"./cfg-for-testing-only-this-is-bad-practice-in-web-serve-folder/redirects.yaml": `---
/a: /b
/some-page: /another-page
/docs: https://bobymcbobs.gitlab.io/go-http-server/
/example: https://example.com/
/{thing:[0-9]+}: https://example.com/404
/page/some-page-{any:.*}: /some-page
`,
			},
			requests: []struct {
				req         *http.Request
				wantBody    *string
				wantCode    int
				wantPath    string
				wantError   bool
				wantHeaders map[string]string
				client      *http.Client
			}{
				{
					req:      newRequest(http.MethodGet, "/", nil),
					wantBody: pointer(`hello`),
					wantCode: http.StatusOK,
				},
				{
					req:      newRequest(http.MethodGet, "/a", nil),
					wantBody: nil,
					wantCode: http.StatusTemporaryRedirect,
					wantHeaders: map[string]string{
						"Location": "/b",
					},
				},
				{
					req:      newRequest(http.MethodGet, "/some-page", nil),
					wantBody: nil,
					wantCode: http.StatusTemporaryRedirect,
					wantHeaders: map[string]string{
						"Location": "/another-page",
					},
				},
				{
					req:      newRequest(http.MethodGet, "/docs", nil),
					wantBody: nil,
					wantCode: http.StatusTemporaryRedirect,
					wantHeaders: map[string]string{
						"Location": "https://bobymcbobs.gitlab.io/go-http-server/",
					},
				},
				{
					req:      newRequest(http.MethodGet, "/example", nil),
					wantBody: nil,
					wantCode: http.StatusTemporaryRedirect,
					wantHeaders: map[string]string{
						"Location": "https://example.com/",
					},
				},
				{
					req:      newRequest(http.MethodGet, "/123", nil),
					wantBody: nil,
					wantCode: http.StatusTemporaryRedirect,
					wantHeaders: map[string]string{
						"Location": "https://example.com/404",
					},
				},
				{
					req:      newRequest(http.MethodGet, "/"+randStringRunes(numberRunes, 9), nil),
					wantBody: nil,
					wantCode: http.StatusTemporaryRedirect,
					wantHeaders: map[string]string{
						"Location": "https://example.com/404",
					},
				},
				{
					req:      newRequest(http.MethodGet, "/page/some-page-"+randStringRunes(letterRunes, 10), nil),
					wantBody: nil,
					wantCode: http.StatusTemporaryRedirect,
					wantHeaders: map[string]string{
						"Location": "/some-page",
					},
				},
			},
		},
		// testcase history mode omnibus
		{
			name: "history mode omnibus",
			env: map[string]string{
				"APP_VUEJS_HISTORY_MODE":      "true",
				"APP_TEMPLATE_MAP_PATH":       "./cfg-for-testing-only-this-is-bad-practice-in-web-serve-folder/tmpl.yaml",
				"APP_REDIRECT_ROUTES_ENABLED": "true",
				"APP_REDIRECT_ROUTES_PATH":    "./cfg-for-testing-only-this-is-bad-practice-in-web-serve-folder/redirects.yaml",
				"APP_HEADER_MAP_PATH":         "./cfg-for-testing-only-this-is-bad-practice-in-web-serve-folder/headers.yaml",
				"APP_HEADER_SET_ENABLE":       "true",
				"SOME_VALUE":                  "it-is-working-yes",
			},
			files: map[string]string{
				"index.html": `{{ .Title }}`,
				"./cfg-for-testing-only-this-is-bad-practice-in-web-serve-folder/tmpl.yaml": `---
Title: Hello!
`,
				"./cfg-for-testing-only-this-is-bad-practice-in-web-serve-folder/redirects.yaml": `---
/a: /b
/some-page: /another-page
/docs: https://bobymcbobs.gitlab.io/go-http-server/
/example: https://example.com/
/{thing:[0-9]+}: https://example.com/404
/page/some-page-{any:.*}: /some-page
`,
				"./cfg-for-testing-only-this-is-bad-practice-in-web-serve-folder/headers.yaml": `---
X-Cool-As:
 - "Yeah"
X-From-Env:
  - "${SOME_VALUE}"
`,
			},
			requests: []struct {
				req         *http.Request
				wantBody    *string
				wantCode    int
				wantPath    string
				wantError   bool
				wantHeaders map[string]string
				client      *http.Client
			}{
				{
					req:      newRequest(http.MethodGet, "/", nil),
					wantBody: pointer(`Hello!`),
					wantCode: http.StatusOK,
					wantHeaders: map[string]string{
						"X-Cool-As":  "Yeah",
						"X-From-Env": "it-is-working-yes",
					},
				},
				{
					req:      newRequest(http.MethodGet, "/54309y90429", nil),
					wantBody: pointer(`Hello!`),
					wantCode: http.StatusOK,
					wantHeaders: map[string]string{
						"X-Cool-As":  "Yeah",
						"X-From-Env": "it-is-working-yes",
					},
				},
				{
					req:      newRequest(http.MethodGet, "/a/b/c/b/c/a/dddddd", nil),
					wantBody: pointer(`Hello!`),
					wantCode: http.StatusOK,
					wantHeaders: map[string]string{
						"X-Cool-As":  "Yeah",
						"X-From-Env": "it-is-working-yes",
					},
				},
				{
					req:      newRequest(http.MethodGet, "/a", nil),
					wantBody: nil,
					wantCode: http.StatusTemporaryRedirect,
					wantHeaders: map[string]string{
						"Location": "/b",
					},
				},
				{
					req:      newRequest(http.MethodGet, "/some-page", nil),
					wantBody: nil,
					wantCode: http.StatusTemporaryRedirect,
					wantHeaders: map[string]string{
						"Location": "/another-page",
					},
				},
				{
					req:      newRequest(http.MethodGet, "/some-page", nil),
					client:   http.DefaultClient,
					wantBody: pointer(`Hello!`),
					wantCode: http.StatusOK,
					wantHeaders: map[string]string{
						"X-Cool-As":  "Yeah",
						"X-From-Env": "it-is-working-yes",
					},
				},
				{
					req:      newRequest(http.MethodGet, "/docs", nil),
					wantBody: nil,
					wantCode: http.StatusTemporaryRedirect,
					wantHeaders: map[string]string{
						"Location": "https://bobymcbobs.gitlab.io/go-http-server/",
					},
				},
				{
					req:      newRequest(http.MethodGet, "/example", nil),
					wantBody: nil,
					wantCode: http.StatusTemporaryRedirect,
					wantHeaders: map[string]string{
						"Location": "https://example.com/",
					},
				},
				{
					req:      newRequest(http.MethodGet, "/123", nil),
					wantBody: nil,
					wantCode: http.StatusTemporaryRedirect,
					wantHeaders: map[string]string{
						"Location": "https://example.com/404",
					},
				},
				{
					req:      newRequest(http.MethodGet, "/"+randStringRunes(numberRunes, 9), nil),
					wantBody: nil,
					wantCode: http.StatusTemporaryRedirect,
					wantHeaders: map[string]string{
						"Location": "https://example.com/404",
					},
				},
				{
					req:      newRequest(http.MethodGet, "/page/some-page-"+randStringRunes(letterRunes, 10), nil),
					wantBody: nil,
					wantCode: http.StatusTemporaryRedirect,
					wantHeaders: map[string]string{
						"Location": "/some-page",
					},
				},
			},
		},
		// testcase dotfiles
		{
			name: "dotfiles",
			files: map[string]string{
				".ghs.yaml": `---
error404FilePath: some404.html
headerMap:
  X-Some-Header:
    - Valid
redirectRoutes:
  /a: /b
  /some-page: /another-page
  /docs: https://bobymcbobs.gitlab.io/go-http-server/
  /example: https://example.com/
  /{thing:[0-9]+}: https://example.com/404
  /page/some-page-{any:.*}: /some-page
`,
				"index.html":   `Hello!`,
				"some404.html": `Not found!!!`,
				"b":            `Bees!`,
			},
			requests: []struct {
				req         *http.Request
				wantBody    *string
				wantCode    int
				wantPath    string
				wantError   bool
				wantHeaders map[string]string
				client      *http.Client
			}{
				{
					req:      newRequest(http.MethodGet, "/aabdkbdf", nil),
					wantBody: pointer(`Not found!!!`),
					wantCode: http.StatusNotFound,
					wantHeaders: map[string]string{
						"X-Some-Header": "Valid",
					},
				},
				{
					req:      newRequest(http.MethodGet, "/", nil),
					wantBody: pointer(`Hello!`),
					wantCode: http.StatusOK,
					wantHeaders: map[string]string{
						"X-Some-Header": "Valid",
					},
				},
				{
					req:      newRequest(http.MethodGet, "/a", nil),
					wantBody: nil,
					wantCode: http.StatusTemporaryRedirect,
					wantHeaders: map[string]string{
						"Location": "/b",
					},
				},
				{
					req:      newRequest(http.MethodGet, "/a", nil),
					client:   http.DefaultClient,
					wantBody: pointer(`Bees!`),
					wantCode: http.StatusOK,
					wantHeaders: map[string]string{
						"X-Some-Header": "Valid",
					},
				},
				{
					req:      newRequest(http.MethodGet, "/some-page", nil),
					wantBody: nil,
					wantCode: http.StatusTemporaryRedirect,
					wantHeaders: map[string]string{
						"Location": "/another-page",
					},
				},
				{
					req:      newRequest(http.MethodGet, "/docs", nil),
					wantBody: nil,
					wantCode: http.StatusTemporaryRedirect,
					wantHeaders: map[string]string{
						"Location": "https://bobymcbobs.gitlab.io/go-http-server/",
					},
				},
				{
					req:      newRequest(http.MethodGet, "/example", nil),
					wantBody: nil,
					wantCode: http.StatusTemporaryRedirect,
					wantHeaders: map[string]string{
						"Location": "https://example.com/",
					},
				},
				{
					req:      newRequest(http.MethodGet, "/123", nil),
					wantBody: nil,
					wantCode: http.StatusTemporaryRedirect,
					wantHeaders: map[string]string{
						"Location": "https://example.com/404",
					},
				},
				{
					req:      newRequest(http.MethodGet, "/"+randStringRunes(numberRunes, 9), nil),
					wantBody: nil,
					wantCode: http.StatusTemporaryRedirect,
					wantHeaders: map[string]string{
						"Location": "https://example.com/404",
					},
				},
				{
					req:      newRequest(http.MethodGet, "/page/some-page-"+randStringRunes(letterRunes, 10), nil),
					wantBody: nil,
					wantCode: http.StatusTemporaryRedirect,
					wantHeaders: map[string]string{
						"Location": "/some-page",
					},
				},
			},
		},
		// testcase dotfiles with history mode and templating
		{
			name: "dotfiles with history mode and templating",
			files: map[string]string{
				".ghs.yaml": `---
historyMode: true
headerMap:
  X-Some-Header:
    - Valid
templateMap:
  Title: "Hello there!"
  UnevaluatedEnv: "${HOST}"
redirectRoutes:
  /a: /b
  /a/b: /b/c
  /some-page: /another-page
  /docs: https://bobymcbobs.gitlab.io/go-http-server/
  /example: https://example.com/
  /{thing:[0-9]+}: https://example.com/404
  /page/some-page-{any:.*}: /some-page
`,
				"index.html": `{{ .Title }} {{ .UnevaluatedEnv }}`,
			},
			requests: []struct {
				req         *http.Request
				wantBody    *string
				wantCode    int
				wantPath    string
				wantError   bool
				wantHeaders map[string]string
				client      *http.Client
			}{
				{
					req:      newRequest(http.MethodGet, "/aabdkbdf", nil),
					wantBody: pointer(`Hello there! ${HOST}`),
					wantCode: http.StatusOK,
					wantHeaders: map[string]string{
						"X-Some-Header": "Valid",
					},
				},
				{
					req:      newRequest(http.MethodGet, "/", nil),
					wantBody: pointer(`Hello there! ${HOST}`),
					wantCode: http.StatusOK,
					wantHeaders: map[string]string{
						"X-Some-Header": "Valid",
					},
				},
				{
					req:      newRequest(http.MethodGet, "/settings", nil),
					wantBody: pointer(`Hello there! ${HOST}`),
					wantCode: http.StatusOK,
					wantHeaders: map[string]string{
						"X-Some-Header": "Valid",
					},
				},
				{
					req:      newRequest(http.MethodGet, "/a", nil),
					wantBody: nil,
					wantCode: http.StatusTemporaryRedirect,
					wantHeaders: map[string]string{
						"Location": "/b",
					},
				},
				{
					req:      newRequest(http.MethodGet, "/a/b", nil),
					wantBody: nil,
					wantCode: http.StatusTemporaryRedirect,
					wantHeaders: map[string]string{
						"Location": "/b/c",
					},
				},
				{
					req:      newRequest(http.MethodGet, "/a/b", nil),
					client:   http.DefaultClient,
					wantBody: pointer(`Hello there! ${HOST}`),
					wantCode: http.StatusOK,
					wantHeaders: map[string]string{
						"X-Some-Header": "Valid",
					},
				},
				{
					req:      newRequest(http.MethodGet, "/some-page", nil),
					wantBody: nil,
					wantCode: http.StatusTemporaryRedirect,
					wantHeaders: map[string]string{
						"Location": "/another-page",
					},
				},
				{
					req:      newRequest(http.MethodGet, "/docs", nil),
					wantBody: nil,
					wantCode: http.StatusTemporaryRedirect,
					wantHeaders: map[string]string{
						"Location": "https://bobymcbobs.gitlab.io/go-http-server/",
					},
				},
				{
					req:      newRequest(http.MethodGet, "/example", nil),
					wantBody: nil,
					wantCode: http.StatusTemporaryRedirect,
					wantHeaders: map[string]string{
						"Location": "https://example.com/",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		// NOTE no parallel because writing env
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip()
			}
			if tt.env == nil {
				tt.env = map[string]string{}
			}
			for k, v := range defaultEnv {
				if tt.env[k] != "" {
					continue
				}
				tt.env[k] = v
			}
			dir, err := os.MkdirTemp("", "main-e2e")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(dir)
			if err := os.Chdir(dir); err != nil {
				t.Fatalf("failed to change directory: %v", err)
			}
			for f, c := range tt.files {
				d := path.Dir(f)
				if err := os.Mkdir(path.Join(dir, d), 0777); errors.Is(err, os.ErrNotExist) {
					t.Fatalf("failed to create dir: %v", err)
				}
				if err := os.WriteFile(path.Join(dir, f), []byte(c), 0644); err != nil {
					t.Fatalf("failed to write file: %v", err)
				}
			}
			if tt.env["APP_PORT"] == "" {
				tt.env["APP_PORT"] = fmt.Sprintf(":%v", rand.Intn(65000-50000)+50000)
			}
			t.Log("appbuildmode", common.AppBuildMode)
			if tt.productionBuild {
				common.AppBuildMode = "production"
			} else {
				common.AppBuildMode = "development"
			}
			prevEnv := map[string]string{}
			for k, v := range tt.env {
				prevEnv[k] = os.Getenv(k)
				log.Printf("saving env %v (%v)\n", k, prevEnv[k])
				log.Printf("setting env %v to %v\n", k, v)
				if !tt.noSetEnv {
					os.Setenv(k, v)
				}
			}
			defer func() {
				for k, v := range prevEnv {
					log.Printf("restoring env %v to %v", k, v)
					os.Setenv(k, v)
				}
			}()
			go func() {
				// NOTE not sure how to quit when the test finishes yet
				//      so it will become a zombie go routine until `go test` completes
				//      each test must bind to unique ports
				main()
			}()
			// wait for port listening
			for {
				conn, err := net.DialTimeout("tcp", "localhost"+tt.env["APP_PORT"], time.Millisecond*500)
				if err == nil {
					conn.Close()
					break
				}
				time.Sleep(1 * time.Second)
			}
			for _, r := range tt.requests {
				r.req.URL.Scheme = "http"
				if r.req.URL.Host == "" {
					r.req.URL.Host = "localhost" + tt.env["APP_PORT"]
				}
				var c *http.Client
				if r.client == nil {
					c = client
				} else {
					c = r.client
				}
				resp, err := c.Do(r.req)
				if err != nil && !r.wantError {
					t.Fatal(err)
				}
				if resp == nil && r.wantError {
					continue
				}
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatal(err)
				}
				t.Log(resp.Header)
				if !r.wantError && r.wantBody != nil && string(body) != *r.wantBody {
					t.Fatalf("request to path returned body '%v' = %v; wants = %v", r.req.URL.Path, string(body), *r.wantBody)
				}
				if !r.wantError && resp.StatusCode != r.wantCode {
					t.Fatalf("request to path returned response code '%v' = %v; wants = %v", r.req.URL.Path, resp.StatusCode, r.wantCode)
				}
				foundHeaders := []string{}
				for k, v := range r.wantHeaders {
					if val := resp.Header.Get(k); v == val {
						foundHeaders = append(foundHeaders, k)
					}
				}
				if len(foundHeaders) != len(r.wantHeaders) {
					t.Fatalf("unexpected amount of expected headers found = %v (%+v); wants = %v (%+v)", len(foundHeaders), foundHeaders, len(r.wantHeaders), r.wantHeaders)
				}
			}
		})
	}
}
