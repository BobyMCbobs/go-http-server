package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

func TestHandler_serveHandlerVuejsHistoryMode(t *testing.T) {
	type fields struct {
		Error404FilePath   string
		HeaderMap          map[string][]string
		GzipEnabled        bool
		HeaderMapEnabled   bool
		TemplateMap        map[string]string
		TemplateMapEnabled bool
		VueJSHistoryMode   bool
		ServeFolder        string
	}
	tests := []struct {
		name        string
		fields      fields
		serveFolder []struct {
			name    string
			content string
		}
		tests []struct {
			path       string
			content    string
			statusCode int
			headers    map[string]string
		}
	}{
		struct {
			name        string
			fields      fields
			serveFolder []struct {
				name    string
				content string
			}
			tests []struct {
				path       string
				content    string
				statusCode int
				headers    map[string]string
			}
		}{
			name: "basic test",
			serveFolder: []struct {
				name    string
				content string
			}{
				{
					name:    "index.html",
					content: "<h1>hello world</h1>",
				},
			},
			tests: []struct {
				path       string
				content    string
				statusCode int
				headers    map[string]string
			}{
				{
					path:       "/",
					content:    "<h1>hello world</h1>",
					statusCode: http.StatusOK,
				},
				{
					path:       "/aaa",
					content:    "<h1>hello world</h1>",
					statusCode: http.StatusOK,
				},
				{
					path:       "/bbbb",
					content:    "<h1>hello world</h1>",
					statusCode: http.StatusOK,
				},
				{
					path:       "/index.html",
					content:    "<h1>hello world</h1>",
					statusCode: http.StatusOK,
				},
			},
		},
		struct {
			name        string
			fields      fields
			serveFolder []struct {
				name    string
				content string
			}
			tests []struct {
				path       string
				content    string
				statusCode int
				headers    map[string]string
			}
		}{
			name: "templating test",
			fields: fields{
				TemplateMap: map[string]string{
					"Message": "some-kinda-test",
				},
			},
			serveFolder: []struct {
				name    string
				content string
			}{
				{
					name:    "index.html",
					content: "<h1>{{.Message}}</h1>",
				},
			},
			tests: []struct {
				path       string
				content    string
				statusCode int
				headers    map[string]string
			}{
				{
					path:       "/",
					content:    "<h1>some-kinda-test</h1>",
					statusCode: http.StatusOK,
				},
			},
		},
		struct {
			name        string
			fields      fields
			serveFolder []struct {
				name    string
				content string
			}
			tests []struct {
				path       string
				content    string
				statusCode int
				headers    map[string]string
			}
		}{
			name: "templating test with bad template 1",
			fields: fields{
				TemplateMap: map[string]string{
					"Message": "some-kinda-test",
				},
			},
			serveFolder: []struct {
				name    string
				content string
			}{
				{
					name:    "index.html",
					content: "<h1>{{.Message}</h1>",
				},
			},
			tests: []struct {
				path       string
				content    string
				statusCode int
				headers    map[string]string
			}{
				{
					path: "/",
					content: `500 internal error
`,
					statusCode: http.StatusInternalServerError,
				},
			},
		},
		struct {
			name        string
			fields      fields
			serveFolder []struct {
				name    string
				content string
			}
			tests []struct {
				path       string
				content    string
				statusCode int
				headers    map[string]string
			}
		}{
			name: "templating test with bad template 2",
			fields: fields{
				TemplateMap: map[string]string{
					"Message": "some-kinda-test",
				},
			},
			serveFolder: []struct {
				name    string
				content string
			}{
				{
					name:    "index.html",
					content: "<h1>{{nil}}</h1>",
				},
			},
			tests: []struct {
				path       string
				content    string
				statusCode int
				headers    map[string]string
			}{
				{
					path: "/",
					content: `500 internal error
<h1>`,
					statusCode: http.StatusInternalServerError,
				},
			},
		},
		struct {
			name        string
			fields      fields
			serveFolder []struct {
				name    string
				content string
			}
			tests []struct {
				path       string
				content    string
				statusCode int
				headers    map[string]string
			}
		}{
			name: "header map test",
			fields: fields{
				HeaderMapEnabled: true,
				HeaderMap: map[string][]string{
					"X-Some-Header": []string{"Yes"},
				},
			},
			serveFolder: []struct {
				name    string
				content string
			}{
				{
					name:    "index.html",
					content: "<h1>hello world</h1>",
				},
			},
			tests: []struct {
				path       string
				content    string
				statusCode int
				headers    map[string]string
			}{
				{
					path:       "/",
					content:    "<h1>hello world</h1>",
					statusCode: http.StatusOK,
					headers: map[string]string{
						"X-Some-Header": "Yes",
					},
				},
			},
		},
		struct {
			name        string
			fields      fields
			serveFolder []struct {
				name    string
				content string
			}
			tests []struct {
				path       string
				content    string
				statusCode int
				headers    map[string]string
			}
		}{
			name: "basic test with disallowed path",
			serveFolder: []struct {
				name    string
				content string
			}{
				{
					name:    "index.html",
					content: "<h1>hello world</h1>",
				},
			},
			tests: []struct {
				path       string
				content    string
				statusCode int
				headers    map[string]string
			}{
				{
					path:       "/",
					content:    "<h1>hello world</h1>",
					statusCode: http.StatusOK,
				},
				{
					path:       "/.ghs.yaml",
					content:    "<h1>hello world</h1>",
					statusCode: http.StatusOK,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := &Handler{
				Error404FilePath:   tt.fields.Error404FilePath,
				HeaderMap:          tt.fields.HeaderMap,
				GzipEnabled:        tt.fields.GzipEnabled,
				HeaderMapEnabled:   tt.fields.HeaderMapEnabled,
				TemplateMap:        tt.fields.TemplateMap,
				TemplateMapEnabled: true,
				VueJSHistoryMode:   true,
				ServeFolder:        tt.fields.ServeFolder,
			}
			dir, err := os.MkdirTemp("", "serveFolder")
			if err != nil {
				t.Fatal(err)
			}
			h.ServeFolder = dir
			defer os.RemoveAll(dir)

			for _, f := range tt.serveFolder {
				file := filepath.Join(dir, f.name)
				if err := os.WriteFile(file, []byte(f.content), 0644); err != nil {
					t.Fatal(err)
				}
			}

			for _, te := range tt.tests {
				srv := httptest.NewServer(h.serveHandlerVuejsHistoryMode())
				defer srv.Close()
				tp, err := url.JoinPath(srv.URL, te.path)
				if err != nil {
					t.Fatal(err)
				}
				client := srv.Client()
				resp, err := client.Get(tp)
				if err != nil {
					t.Fatal(err)
				}

				b, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatal(err)
				}
				if got, want := resp.StatusCode, te.statusCode; got != want {
					t.Errorf("Handler.serveHandlerVuejsHistoryMode() = %v, want %v for path %v", got, want, te.path)
				}
				if got, want := string(b), te.content; got != want {
					t.Errorf("Handler.serveHandlerVuejsHistoryMode() = %v, want %v for path %v", got, want, te.path)
				}
				for hk, hv := range te.headers {
					if got, want := resp.Header.Get(hk), hv; got != want {
						t.Errorf("Handler.serveHandlerStandard() = %v, want %v for path %v", got, want, te.path)
					}
				}
			}
		})
	}
}

func TestHandler_serveHandlerStandard(t *testing.T) {
	type fields struct {
		Error404FilePath   string
		HeaderMap          map[string][]string
		GzipEnabled        bool
		HeaderMapEnabled   bool
		TemplateMap        map[string]string
		TemplateMapEnabled bool
		VueJSHistoryMode   bool
		ServeFolder        string
	}
	tests := []struct {
		name        string
		fields      fields
		serveFolder []struct {
			name    string
			content string
		}
		tests []struct {
			path       string
			content    string
			statusCode int
			headers    map[string]string
		}
	}{
		{
			name: "basic test",
			fields: fields{
				Error404FilePath: "404.html",
			},
			serveFolder: []struct {
				name    string
				content string
			}{
				{
					name:    "index.html",
					content: "<h1>hello world</h1>",
				},
				{
					name:    "404.html",
					content: "<h1>unknown page</h1>",
				},
			},
			tests: []struct {
				path       string
				content    string
				statusCode int
				headers    map[string]string
			}{
				{
					path:       "/",
					content:    "<h1>hello world</h1>",
					statusCode: http.StatusOK,
				},
				{
					path:       "/aaaaa",
					content:    `<h1>unknown page</h1>`,
					statusCode: http.StatusNotFound,
				},
				{
					path:       "/MMMMM.html",
					content:    `<h1>unknown page</h1>`,
					statusCode: http.StatusNotFound,
				},
				{
					path:       "/index.html",
					content:    "<h1>hello world</h1>",
					statusCode: http.StatusOK,
				},
			},
		},
		{
			name: "basic test missing 404 page",
			fields: fields{
				Error404FilePath: "404.html",
			},
			serveFolder: []struct {
				name    string
				content string
			}{
				{
					name:    "index.html",
					content: "<h1>hello world</h1>",
				},
			},
			tests: []struct {
				path       string
				content    string
				statusCode int
				headers    map[string]string
			}{
				{
					path:       "/",
					content:    "<h1>hello world</h1>",
					statusCode: http.StatusOK,
				},
				{
					path: "/aaaaa",
					content: `404 page not found
`,
					statusCode: http.StatusNotFound,
				},
				{
					path: "/MMMMM.html",
					content: `404 page not found
`,
					statusCode: http.StatusNotFound,
				},
				{
					path:       "/index.html",
					content:    "<h1>hello world</h1>",
					statusCode: http.StatusOK,
				},
			},
		},
		{
			name: "header map test",
			fields: fields{
				HeaderMapEnabled: true,
				HeaderMap: map[string][]string{
					"X-Some-Header": []string{"Yes"},
				},
			},
			serveFolder: []struct {
				name    string
				content string
			}{
				{
					name:    "index.html",
					content: "<h1>hello world</h1>",
				},
			},
			tests: []struct {
				path       string
				content    string
				statusCode int
				headers    map[string]string
			}{
				{
					path:       "/",
					content:    "<h1>hello world</h1>",
					statusCode: http.StatusOK,
					headers: map[string]string{
						"X-Some-Header": "Yes",
					},
				},
			},
		},
		{
			name: "basic test with disallowed paths",
			fields: fields{
				Error404FilePath: "404.html",
			},
			serveFolder: []struct {
				name    string
				content string
			}{
				{
					name:    "index.html",
					content: "<h1>hello world</h1>",
				},
				{
					name:    "404.html",
					content: "<h1>unknown page</h1>",
				},
			},
			tests: []struct {
				path       string
				content    string
				statusCode int
				headers    map[string]string
			}{
				{
					path:       "/",
					content:    "<h1>hello world</h1>",
					statusCode: http.StatusOK,
				},
				{
					path:       "/.ghs.yaml",
					content:    `<h1>unknown page</h1>`,
					statusCode: http.StatusNotFound,
				},
				{
					path:       "/MMMMM.html",
					content:    `<h1>unknown page</h1>`,
					statusCode: http.StatusNotFound,
				},
				{
					path:       "/index.html",
					content:    "<h1>hello world</h1>",
					statusCode: http.StatusOK,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := &Handler{
				Error404FilePath:   tt.fields.Error404FilePath,
				HeaderMap:          tt.fields.HeaderMap,
				GzipEnabled:        tt.fields.GzipEnabled,
				HeaderMapEnabled:   tt.fields.HeaderMapEnabled,
				TemplateMap:        tt.fields.TemplateMap,
				TemplateMapEnabled: tt.fields.TemplateMapEnabled,
				VueJSHistoryMode:   tt.fields.VueJSHistoryMode,
				ServeFolder:        tt.fields.ServeFolder,
			}
			dir, err := os.MkdirTemp("", "serveFolder")
			if err != nil {
				t.Fatal(err)
			}
			h.ServeFolder = dir
			defer os.RemoveAll(dir)

			for _, f := range tt.serveFolder {
				file := filepath.Join(dir, f.name)
				if err := os.WriteFile(file, []byte(f.content), 0644); err != nil {
					t.Fatal(err)
				}
			}

			for _, te := range tt.tests {
				srv := httptest.NewServer(h.serveHandlerStandard())
				defer srv.Close()
				tp, err := url.JoinPath(srv.URL, te.path)
				if err != nil {
					t.Fatal(err)
				}
				client := srv.Client()
				resp, err := client.Get(tp)
				if err != nil {
					t.Fatal(err)
				}
				b, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatal(err)
				}
				if got, want := resp.StatusCode, te.statusCode; got != want {
					t.Errorf("Handler.serveHandlerStandard() = %v, want %v for path %v", got, want, te.path)
				}
				if got, want := string(b), te.content; got != want {
					t.Errorf("Handler.serveHandlerStandard() = %v, want %v for path %v", got, want, te.path)
				}
				for hk, hv := range te.headers {
					if got, want := resp.Header.Get(hk), hv; got != want {
						t.Errorf("Handler.serveHandlerStandard() = %v, want %v for path %v", got, want, te.path)
					}
				}
			}
		})
	}
}

func TestHandler_ServeHandler(t *testing.T) {
	type fields struct {
		Error404FilePath   string
		HeaderMap          map[string][]string
		GzipEnabled        bool
		HeaderMapEnabled   bool
		TemplateMap        map[string]string
		TemplateMapEnabled bool
		VueJSHistoryMode   bool
		ServeFolder        string
	}
	tests := []struct {
		name        string
		fields      fields
		wantHandler http.Handler
	}{
		{
			name: "handler default",
		},
		{
			name: "handler vuejs history mode",
			fields: fields{
				VueJSHistoryMode: true,
			},
		},
		{
			name: "handler default with gzip",
			fields: fields{
				GzipEnabled: true,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := &Handler{
				Error404FilePath:   tt.fields.Error404FilePath,
				HeaderMap:          tt.fields.HeaderMap,
				GzipEnabled:        tt.fields.GzipEnabled,
				HeaderMapEnabled:   tt.fields.HeaderMapEnabled,
				TemplateMap:        tt.fields.TemplateMap,
				TemplateMapEnabled: tt.fields.TemplateMapEnabled,
				VueJSHistoryMode:   tt.fields.VueJSHistoryMode,
				ServeFolder:        tt.fields.ServeFolder,
			}
			// TODO test this better
			if gotHandler := h.ServeHandler(); gotHandler == nil {
				t.Errorf("Handler.ServeHandler() = %v, want not nil", gotHandler)
			}
		})
	}
}

func TestHandler_ServeStandardRedirect(t *testing.T) {
	type fields struct {
		Error404FilePath   string
		HeaderMap          map[string][]string
		GzipEnabled        bool
		HeaderMapEnabled   bool
		TemplateMap        map[string]string
		TemplateMapEnabled bool
		VueJSHistoryMode   bool
		ServeFolder        string
	}
	type args struct {
		from string
		to   string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   struct {
			location   string
			statusCode int
			errorBody  string
		}
	}{
		struct {
			name   string
			fields fields
			args   args
			want   struct {
				location   string
				statusCode int
				errorBody  string
			}
		}{
			name: "basic path redirect",
			args: args{
				from: "/a",
				to:   "/b",
			},
			want: struct {
				location   string
				statusCode int
				errorBody  string
			}{
				location: "/b",
			},
		},
		struct {
			name   string
			fields fields
			args   args
			want   struct {
				location   string
				statusCode int
				errorBody  string
			}
		}{
			name: "basic path redirect to url",
			args: args{
				from: "/a",
				to:   "http://example.com/",
			},
			want: struct {
				location   string
				statusCode int
				errorBody  string
			}{
				location: "http://example.com/",
			},
		},
		struct {
			name   string
			fields fields
			args   args
			want   struct {
				location   string
				statusCode int
				errorBody  string
			}
		}{
			name: "bad url",
			args: args{
				from: "/a",
				to:   ":90%?88**",
			},
			want: struct {
				location   string
				statusCode int
				errorBody  string
			}{
				errorBody:  "fatal: unable to redirect to destination URL",
				statusCode: http.StatusInternalServerError,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		if tt.want.statusCode == 0 {
			tt.want.statusCode = http.StatusTemporaryRedirect
		}
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := &Handler{
				Error404FilePath:   tt.fields.Error404FilePath,
				HeaderMap:          tt.fields.HeaderMap,
				GzipEnabled:        tt.fields.GzipEnabled,
				HeaderMapEnabled:   tt.fields.HeaderMapEnabled,
				TemplateMap:        tt.fields.TemplateMap,
				TemplateMapEnabled: tt.fields.TemplateMapEnabled,
				VueJSHistoryMode:   tt.fields.VueJSHistoryMode,
				ServeFolder:        tt.fields.ServeFolder,
			}
			req := httptest.NewRequest("GET", "http://example.com", nil)
			w := httptest.NewRecorder()
			h.ServeStandardRedirect(tt.args.from, tt.args.to)(w, req)

			// if tt.error != ""
			if code := w.Result().StatusCode; code != tt.want.statusCode {
				t.Errorf("Handler.ServeStandardRedirect() = %v, want %v", code, http.StatusTemporaryRedirect)
			}
			if got := w.Result().Header.Get("Location"); got != tt.want.location {
				t.Errorf("Handler.ServeStandardRedirect() = %v, want %v", got, tt.want)
			}
		})
	}
}
