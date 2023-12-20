package httpserver

import (
	"crypto/tls"
	_ "embed"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"

	"gitlab.com/BobyMCbobs/go-http-server/pkg/handlers"
	"gitlab.com/BobyMCbobs/go-http-server/pkg/metrics"
)

var (
	//go:embed testdata/server-cert.pem
	tlsPublic string
	//go:embed testdata/server-key.pem
	tlsPrivate string
)

func TestNewWebServer(t *testing.T) {
	tests := []struct {
		name                 string
		env                  map[string]string
		dotfileContent       string
		setServeFolderToTemp bool
		findValue            func(*WebServer) any
		want                 any
	}{
		{
			name: "set app port",
			env: map[string]string{
				"APP_PORT": ":8123",
			},
			findValue: func(ws *WebServer) any {
				return ws.AppPort
			},
			want: ":8123",
		},
		{
			name: "set app port",
			env: map[string]string{
				"APP_PORT": ":8123",
			},
			findValue: func(ws *WebServer) any {
				return ws.AppPort
			},
			want: ":8123",
		},
		{
			name: "use settings from dotfile",
			dotfileContent: `---
headerMap:
  X-Abc:
    - Hello
templateMap:
  AAAA: BBBB
redirectRoutes:
  /abc: /cba
  /example: http://example.com
`,
			setServeFolderToTemp: true,
			findValue: func(ws *WebServer) any {
				return []any{
					ws.HeaderMap, ws.TemplateMap, ws.RedirectRoutes,
				}
			},
			want: []any{
				map[string][]string{"X-Abc": {"Hello"}},
				map[string]string{"AAAA": "BBBB"},
				map[string]string{"/abc": "/cba", "/example": "http://example.com"},
			},
		},
		{
			name:                 "dotfile overrides env",
			setServeFolderToTemp: true,
			env: map[string]string{
				"APP_404_PAGE_FILE_NAME": "404-from-env.html",
			},
			dotfileContent: `---
error404FilePath: 404-from-dotfile.html
`,
			findValue: func(ws *WebServer) any {
				return ws.Error404FilePath
			},
			want: "404-from-dotfile.html",
		},
		{
			name:                 "dotfile template and header env aren't evaluated",
			setServeFolderToTemp: true,
			dotfileContent: `---
headerMap:
  X-Abc:
    - ${HOST}
templateMap:
  AAAA: ${HOST}
`,
			findValue: func(ws *WebServer) any {
				return []any{ws.HeaderMap, ws.TemplateMap}
			},
			want: []any{map[string][]string{"X-Abc": {"${HOST}"}}, map[string]string{"AAAA": "${HOST}"}},
		},
		{
			name: "use error 404 page from dotfile",
			dotfileContent: `---
error404FilePath: 404-lol.html
`,
			setServeFolderToTemp: true,
			findValue: func(ws *WebServer) any {
				return ws.Error404FilePath
			},
			want: "404-lol.html",
		},
		{
			name:                 "can't read dotfile",
			dotfileContent:       `%(*$#**)`,
			setServeFolderToTemp: true,
			findValue: func(ws *WebServer) any {
				return nil
			},
			want: nil,
		},
		{
			// TODO needs some thought
			name: "set middleware",
			findValue: func(*WebServer) any {
				return true
			},
			want: true,
		},
		// TODO redirect routes + failure
		// {
		// 	name: "redirect routes",
		// },
		// TODO extrahandlers + / root fail
		// TODO tls
		{
			name: "tls enabled",
			env: map[string]string{
				"APP_ENABLE_HTTPS": "true",
			},
			findValue: func(ws *WebServer) any {
				return ws.serverTLS != nil
			},
			want: true,
		},
		{
			name: "bad http origins",
			env: map[string]string{
				"APP_HTTP_ALLOWED_ORIGINS": "%&*#exam???ple.com",
			},
			findValue: func(ws *WebServer) any {
				return len(ws.HTTPAllowedOrigins)
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		tt := tt
		// NOTE sets env and cannot be parallelised
		t.Run(tt.name, func(t *testing.T) {
			prevEnv := map[string]string{}
			dir, err := os.MkdirTemp("", "load-newwebserver")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(dir)
			if tt.setServeFolderToTemp {
				if tt.env == nil {
					tt.env = map[string]string{}
				}
				tt.env["APP_SERVE_FOLDER"] = dir
			}
			for k, v := range tt.env {
				prevEnv[k] = os.Getenv(k)
				os.Setenv(k, v)
			}
			defer func() {
				for k, v := range prevEnv {
					os.Setenv(k, v)
				}
			}()
			if tt.dotfileContent != "" {
				if err := os.WriteFile(path.Join(dir, ".ghs.yaml"), []byte(tt.dotfileContent), 0644); err != nil {
					t.Fatalf("error: failed to write .ghs.yaml: %v", err)
				}
			}
			if got := tt.findValue(NewWebServer()); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWebServer() = %+v (%v), want %+v (%v)", got, reflect.TypeOf(got), tt.want, reflect.TypeOf(tt.want))
			}
		})
	}
}

func TestWebServer_SetServeFolder(t *testing.T) {
	type fields struct {
		AppPort               string
		HTTPAllowedOrigins    []string
		EnvFile               string
		Error404FilePath      string
		ExtraHandlers         []*ExtraHandler
		ExtraMiddleware       []func(http.Handler) http.Handler
		GzipEnabled           bool
		HTTPPort              string
		HTTPSPort             string
		HTTPSPortEnabled      bool
		HeaderMap             map[string][]string
		HeaderMapEnabled      bool
		HeaderMapPath         string
		HealthPort            string
		HealthPortEnabled     bool
		MetricsPort           string
		MetricsPortEnabled    bool
		RealIPHeader          string
		RedirectRoutes        map[string]string
		RedirectRoutesEnabled bool
		RedirectRoutesPath    string
		ServeFolder           string
		TLSCertPath           string
		TLSConfig             *tls.Config
		TLSKeyPath            string
		TemplateMap           map[string]string
		TemplateMapEnabled    bool
		TemplateMapPath       string
		VueJSHistoryMode      bool
		handler               *handlers.Handler
		server                *http.Server
		serverTLS             *http.Server
	}
	type args struct {
		path string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *WebServer
	}{
		{
			name: "basic set serve folder",
			args: args{
				path: "./public",
			},
			want: &WebServer{
				ServeFolder: "./public",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			w := &WebServer{
				AppPort:               tt.fields.AppPort,
				HTTPAllowedOrigins:    tt.fields.HTTPAllowedOrigins,
				EnvFile:               tt.fields.EnvFile,
				Error404FilePath:      tt.fields.Error404FilePath,
				ExtraHandlers:         tt.fields.ExtraHandlers,
				ExtraMiddleware:       tt.fields.ExtraMiddleware,
				GzipEnabled:           tt.fields.GzipEnabled,
				HTTPPort:              tt.fields.HTTPPort,
				HTTPSPort:             tt.fields.HTTPSPort,
				HTTPSPortEnabled:      tt.fields.HTTPSPortEnabled,
				HeaderMap:             tt.fields.HeaderMap,
				HeaderMapEnabled:      tt.fields.HeaderMapEnabled,
				HeaderMapPath:         tt.fields.HeaderMapPath,
				HealthPort:            tt.fields.HealthPort,
				HealthPortEnabled:     tt.fields.HealthPortEnabled,
				MetricsPort:           tt.fields.MetricsPort,
				MetricsPortEnabled:    tt.fields.MetricsPortEnabled,
				RealIPHeader:          tt.fields.RealIPHeader,
				RedirectRoutes:        tt.fields.RedirectRoutes,
				RedirectRoutesEnabled: tt.fields.RedirectRoutesEnabled,
				RedirectRoutesPath:    tt.fields.RedirectRoutesPath,
				ServeFolder:           tt.fields.ServeFolder,
				TLSCertPath:           tt.fields.TLSCertPath,
				TLSConfig:             tt.fields.TLSConfig,
				TLSKeyPath:            tt.fields.TLSKeyPath,
				TemplateMap:           tt.fields.TemplateMap,
				TemplateMapEnabled:    tt.fields.TemplateMapEnabled,
				TemplateMapPath:       tt.fields.TemplateMapPath,
				VueJSHistoryMode:      tt.fields.VueJSHistoryMode,
				handler:               tt.fields.handler,
				server:                tt.fields.server,
				serverTLS:             tt.fields.serverTLS,
			}
			if got := w.SetServeFolder(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WebServer.SetServeFolder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebServer_SetExtraHandlers(t *testing.T) {
	type fields struct {
		AppPort               string
		HTTPAllowedOrigins    []string
		EnvFile               string
		Error404FilePath      string
		ExtraHandlers         []*ExtraHandler
		ExtraMiddleware       []func(http.Handler) http.Handler
		GzipEnabled           bool
		HTTPPort              string
		HTTPSPort             string
		HTTPSPortEnabled      bool
		HeaderMap             map[string][]string
		HeaderMapEnabled      bool
		HeaderMapPath         string
		HealthPort            string
		HealthPortEnabled     bool
		MetricsPort           string
		MetricsPortEnabled    bool
		RealIPHeader          string
		RedirectRoutes        map[string]string
		RedirectRoutesEnabled bool
		RedirectRoutesPath    string
		ServeFolder           string
		TLSCertPath           string
		TLSConfig             *tls.Config
		TLSKeyPath            string
		TemplateMap           map[string]string
		TemplateMapEnabled    bool
		TemplateMapPath       string
		VueJSHistoryMode      bool
		handler               *handlers.Handler
		server                *http.Server
		serverTLS             *http.Server
	}
	type args struct {
		hs []*ExtraHandler
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *WebServer
	}{
		{
			name: "basic set extra handler",
			args: args{
				hs: []*ExtraHandler{
					{
						Path: "/aaa",
						HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
							t.Log("request hereeeee")
						},
						HTTPMethods: []string{http.MethodGet},
					},
				},
			},
			want: &WebServer{
				ExtraHandlers: []*ExtraHandler{
					{
						Path: "/aaa",
						HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
							t.Log("request hereeeee")
						},
						HTTPMethods: []string{http.MethodGet},
					},
				},
			},
		},
		{
			name: "no handlers",
			want: &WebServer{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			w := &WebServer{
				AppPort:               tt.fields.AppPort,
				HTTPAllowedOrigins:    tt.fields.HTTPAllowedOrigins,
				EnvFile:               tt.fields.EnvFile,
				Error404FilePath:      tt.fields.Error404FilePath,
				ExtraHandlers:         tt.fields.ExtraHandlers,
				ExtraMiddleware:       tt.fields.ExtraMiddleware,
				GzipEnabled:           tt.fields.GzipEnabled,
				HTTPPort:              tt.fields.HTTPPort,
				HTTPSPort:             tt.fields.HTTPSPort,
				HTTPSPortEnabled:      tt.fields.HTTPSPortEnabled,
				HeaderMap:             tt.fields.HeaderMap,
				HeaderMapEnabled:      tt.fields.HeaderMapEnabled,
				HeaderMapPath:         tt.fields.HeaderMapPath,
				HealthPort:            tt.fields.HealthPort,
				HealthPortEnabled:     tt.fields.HealthPortEnabled,
				MetricsPort:           tt.fields.MetricsPort,
				MetricsPortEnabled:    tt.fields.MetricsPortEnabled,
				RealIPHeader:          tt.fields.RealIPHeader,
				RedirectRoutes:        tt.fields.RedirectRoutes,
				RedirectRoutesEnabled: tt.fields.RedirectRoutesEnabled,
				RedirectRoutesPath:    tt.fields.RedirectRoutesPath,
				ServeFolder:           tt.fields.ServeFolder,
				TLSCertPath:           tt.fields.TLSCertPath,
				TLSConfig:             tt.fields.TLSConfig,
				TLSKeyPath:            tt.fields.TLSKeyPath,
				TemplateMap:           tt.fields.TemplateMap,
				TemplateMapEnabled:    tt.fields.TemplateMapEnabled,
				TemplateMapPath:       tt.fields.TemplateMapPath,
				VueJSHistoryMode:      tt.fields.VueJSHistoryMode,
				handler:               tt.fields.handler,
				server:                tt.fields.server,
				serverTLS:             tt.fields.serverTLS,
			}
			type checked struct {
				Path        string
				HTTPMethods []string
			}
			var expected []checked
			for _, r := range tt.want.ExtraHandlers {
				expected = append(expected, checked{
					Path:        r.Path,
					HTTPMethods: r.HTTPMethods,
				})
			}
			var result []checked
			got := w.SetExtraHandlers(tt.args.hs...)
			for _, r := range got.ExtraHandlers {
				result = append(result, checked{
					Path:        r.Path,
					HTTPMethods: r.HTTPMethods,
				})
			}
			if !reflect.DeepEqual(result, expected) {
				t.Errorf("WebServer.SetExtraHandlers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebServer_SetExtraMiddleware(t *testing.T) {
	type fields struct {
		AppPort               string
		HTTPAllowedOrigins    []string
		EnvFile               string
		Error404FilePath      string
		ExtraHandlers         []*ExtraHandler
		ExtraMiddleware       []func(http.Handler) http.Handler
		GzipEnabled           bool
		HTTPPort              string
		HTTPSPort             string
		HTTPSPortEnabled      bool
		HeaderMap             map[string][]string
		HeaderMapEnabled      bool
		HeaderMapPath         string
		HealthPort            string
		HealthPortEnabled     bool
		MetricsPort           string
		MetricsPortEnabled    bool
		RealIPHeader          string
		RedirectRoutes        map[string]string
		RedirectRoutesEnabled bool
		RedirectRoutesPath    string
		ServeFolder           string
		TLSCertPath           string
		TLSConfig             *tls.Config
		TLSKeyPath            string
		TemplateMap           map[string]string
		TemplateMapEnabled    bool
		TemplateMapPath       string
		VueJSHistoryMode      bool
		handler               *handlers.Handler
		server                *http.Server
		serverTLS             *http.Server
	}
	type args struct {
		m []func(http.Handler) http.Handler
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *WebServer
	}{
		{
			name: "basic set middleware",
			args: args{
				m: []func(h http.Handler) http.Handler{},
			},
			want: &WebServer{
				ExtraMiddleware: []func(http.Handler) http.Handler{},
			},
		},
		{
			name: "no middleware",
			args: args{},
			want: &WebServer{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			w := &WebServer{
				AppPort:               tt.fields.AppPort,
				HTTPAllowedOrigins:    tt.fields.HTTPAllowedOrigins,
				EnvFile:               tt.fields.EnvFile,
				Error404FilePath:      tt.fields.Error404FilePath,
				ExtraHandlers:         tt.fields.ExtraHandlers,
				ExtraMiddleware:       tt.fields.ExtraMiddleware,
				GzipEnabled:           tt.fields.GzipEnabled,
				HTTPPort:              tt.fields.HTTPPort,
				HTTPSPort:             tt.fields.HTTPSPort,
				HTTPSPortEnabled:      tt.fields.HTTPSPortEnabled,
				HeaderMap:             tt.fields.HeaderMap,
				HeaderMapEnabled:      tt.fields.HeaderMapEnabled,
				HeaderMapPath:         tt.fields.HeaderMapPath,
				HealthPort:            tt.fields.HealthPort,
				HealthPortEnabled:     tt.fields.HealthPortEnabled,
				MetricsPort:           tt.fields.MetricsPort,
				MetricsPortEnabled:    tt.fields.MetricsPortEnabled,
				RealIPHeader:          tt.fields.RealIPHeader,
				RedirectRoutes:        tt.fields.RedirectRoutes,
				RedirectRoutesEnabled: tt.fields.RedirectRoutesEnabled,
				RedirectRoutesPath:    tt.fields.RedirectRoutesPath,
				ServeFolder:           tt.fields.ServeFolder,
				TLSCertPath:           tt.fields.TLSCertPath,
				TLSConfig:             tt.fields.TLSConfig,
				TLSKeyPath:            tt.fields.TLSKeyPath,
				TemplateMap:           tt.fields.TemplateMap,
				TemplateMapEnabled:    tt.fields.TemplateMapEnabled,
				TemplateMapPath:       tt.fields.TemplateMapPath,
				VueJSHistoryMode:      tt.fields.VueJSHistoryMode,
				handler:               tt.fields.handler,
				server:                tt.fields.server,
				serverTLS:             tt.fields.serverTLS,
			}
			if got := w.SetExtraMiddleware(tt.args.m...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WebServer.SetExtraMiddleware() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebServer_LoadTLS(t *testing.T) {
	type fields struct {
		AppPort               string
		HTTPAllowedOrigins    []string
		EnvFile               string
		Error404FilePath      string
		ExtraHandlers         []*ExtraHandler
		ExtraMiddleware       []func(http.Handler) http.Handler
		GzipEnabled           bool
		HTTPPort              string
		HTTPSPort             string
		HTTPSPortEnabled      bool
		HeaderMap             map[string][]string
		HeaderMapEnabled      bool
		HeaderMapPath         string
		HealthPort            string
		HealthPortEnabled     bool
		MetricsPort           string
		MetricsPortEnabled    bool
		RealIPHeader          string
		RedirectRoutes        map[string]string
		RedirectRoutesEnabled bool
		RedirectRoutesPath    string
		ServeFolder           string
		TLSCertPath           string
		TLSConfig             *tls.Config
		TLSKeyPath            string
		TemplateMap           map[string]string
		TemplateMapEnabled    bool
		TemplateMapPath       string
		VueJSHistoryMode      bool
		handler               *handlers.Handler
		server                *http.Server
		serverTLS             *http.Server
	}
	tests := []struct {
		name         string
		fields       fields
		publicKey    string
		privateKey   string
		errorMessage string
		want         int
	}{
		// TODO test this better
		{
			name:       "basic",
			publicKey:  tlsPublic,
			privateKey: tlsPrivate,
			want:       1,
		},
		{
			name:         "no keys",
			errorMessage: "no such file or directory",
			want:         1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WebServer{
				AppPort:               tt.fields.AppPort,
				HTTPAllowedOrigins:    tt.fields.HTTPAllowedOrigins,
				EnvFile:               tt.fields.EnvFile,
				Error404FilePath:      tt.fields.Error404FilePath,
				ExtraHandlers:         tt.fields.ExtraHandlers,
				ExtraMiddleware:       tt.fields.ExtraMiddleware,
				GzipEnabled:           tt.fields.GzipEnabled,
				HTTPPort:              tt.fields.HTTPPort,
				HTTPSPort:             tt.fields.HTTPSPort,
				HTTPSPortEnabled:      tt.fields.HTTPSPortEnabled,
				HeaderMap:             tt.fields.HeaderMap,
				HeaderMapEnabled:      tt.fields.HeaderMapEnabled,
				HeaderMapPath:         tt.fields.HeaderMapPath,
				HealthPort:            tt.fields.HealthPort,
				HealthPortEnabled:     tt.fields.HealthPortEnabled,
				MetricsPort:           tt.fields.MetricsPort,
				MetricsPortEnabled:    tt.fields.MetricsPortEnabled,
				RealIPHeader:          tt.fields.RealIPHeader,
				RedirectRoutes:        tt.fields.RedirectRoutes,
				RedirectRoutesEnabled: tt.fields.RedirectRoutesEnabled,
				RedirectRoutesPath:    tt.fields.RedirectRoutesPath,
				ServeFolder:           tt.fields.ServeFolder,
				TLSCertPath:           tt.fields.TLSCertPath,
				TLSConfig:             tt.fields.TLSConfig,
				TLSKeyPath:            tt.fields.TLSKeyPath,
				TemplateMap:           tt.fields.TemplateMap,
				TemplateMapEnabled:    tt.fields.TemplateMapEnabled,
				TemplateMapPath:       tt.fields.TemplateMapPath,
				VueJSHistoryMode:      tt.fields.VueJSHistoryMode,
				handler:               tt.fields.handler,
				server:                tt.fields.server,
				serverTLS:             tt.fields.serverTLS,
			}
			dir, err := os.MkdirTemp("", "load-tls")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(dir)
			publicKeyPath := path.Join(dir, "tls.cert")
			privateKeyPath := path.Join(dir, "tls.key")
			if tt.publicKey != "" {
				if err := os.WriteFile(publicKeyPath, []byte(tt.publicKey), 0644); err != nil {
					t.Fatal(err)
				}
				w.TLSCertPath = publicKeyPath
			}
			if tt.privateKey != "" {
				if err := os.WriteFile(privateKeyPath, []byte(tt.privateKey), 0644); err != nil {
					t.Fatal(err)
				}
				w.TLSKeyPath = privateKeyPath
			}
			got, err := w.LoadTLS()
			if err != nil && tt.errorMessage == "" {
				t.Fatalf("error loading tls cert: %v", err)
			} else if err != nil && !strings.Contains(err.Error(), tt.errorMessage) {
				t.Fatalf("unexpected error loading tls cert: %v", err)
			}
			if !reflect.DeepEqual(len(got.TLSConfig.Certificates), tt.want) {
				t.Errorf("WebServer.LoadTLS() = %v, want %v", len(got.TLSConfig.Certificates), tt.want)
			}
		})
	}
}

func TestWebServer_LoadTemplateMap(t *testing.T) {
	type fields struct {
		AppPort               string
		HTTPAllowedOrigins    []string
		EnvFile               string
		Error404FilePath      string
		ExtraHandlers         []*ExtraHandler
		ExtraMiddleware       []func(http.Handler) http.Handler
		GzipEnabled           bool
		HTTPPort              string
		HTTPSPort             string
		HTTPSPortEnabled      bool
		HeaderMap             map[string][]string
		HeaderMapEnabled      bool
		HeaderMapPath         string
		HealthPort            string
		HealthPortEnabled     bool
		MetricsPort           string
		MetricsPortEnabled    bool
		RealIPHeader          string
		RedirectRoutes        map[string]string
		RedirectRoutesEnabled bool
		RedirectRoutesPath    string
		ServeFolder           string
		TLSCertPath           string
		TLSConfig             *tls.Config
		TLSKeyPath            string
		TemplateMap           map[string]string
		TemplateMapEnabled    bool
		TemplateMapPath       string
		VueJSHistoryMode      bool
		handler               *handlers.Handler
		server                *http.Server
		serverTLS             *http.Server
		dotfileLoaded         bool
	}
	tests := []struct {
		name                string
		fields              fields
		templateFileContent string
		env                 map[string]string
		errorMessage        string
		want                map[string]string
	}{
		{
			name: "basic with env",
			fields: fields{
				dotfileLoaded: false,
			},
			env: map[string]string{
				"THINGY": "THIS",
			},
			templateFileContent: `---
X-Abc: "123"
Something: "${THINGY}"
`,
			want: map[string]string{
				"X-Abc":     "123",
				"Something": "THIS",
			},
		},
		{
			name: "using dotfiles",
			fields: fields{
				dotfileLoaded: true,
				TemplateMap: map[string]string{
					"X-123": "Abc",
				},
			},
			want: map[string]string{
				"X-123": "Abc",
			},
		},
		{
			name: "file doesn't exist",
			fields: fields{
				dotfileLoaded: false,
			},
			errorMessage: "error: template map file not found",
			want:         nil,
		},
		{
			name: "fails to load templatemap config",
			fields: fields{
				dotfileLoaded: false,
			},
			templateFileContent: `%#%@?#:`,
			errorMessage:        "error converting YAML to JSON",
			want:                nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			w := &WebServer{
				AppPort:               tt.fields.AppPort,
				HTTPAllowedOrigins:    tt.fields.HTTPAllowedOrigins,
				EnvFile:               tt.fields.EnvFile,
				Error404FilePath:      tt.fields.Error404FilePath,
				ExtraHandlers:         tt.fields.ExtraHandlers,
				ExtraMiddleware:       tt.fields.ExtraMiddleware,
				GzipEnabled:           tt.fields.GzipEnabled,
				HTTPPort:              tt.fields.HTTPPort,
				HTTPSPort:             tt.fields.HTTPSPort,
				HTTPSPortEnabled:      tt.fields.HTTPSPortEnabled,
				HeaderMap:             tt.fields.HeaderMap,
				HeaderMapEnabled:      tt.fields.HeaderMapEnabled,
				HeaderMapPath:         tt.fields.HeaderMapPath,
				HealthPort:            tt.fields.HealthPort,
				HealthPortEnabled:     tt.fields.HealthPortEnabled,
				MetricsPort:           tt.fields.MetricsPort,
				MetricsPortEnabled:    tt.fields.MetricsPortEnabled,
				RealIPHeader:          tt.fields.RealIPHeader,
				RedirectRoutes:        tt.fields.RedirectRoutes,
				RedirectRoutesEnabled: tt.fields.RedirectRoutesEnabled,
				RedirectRoutesPath:    tt.fields.RedirectRoutesPath,
				ServeFolder:           tt.fields.ServeFolder,
				TLSCertPath:           tt.fields.TLSCertPath,
				TLSConfig:             tt.fields.TLSConfig,
				TLSKeyPath:            tt.fields.TLSKeyPath,
				TemplateMap:           tt.fields.TemplateMap,
				TemplateMapEnabled:    tt.fields.TemplateMapEnabled,
				TemplateMapPath:       tt.fields.TemplateMapPath,
				VueJSHistoryMode:      tt.fields.VueJSHistoryMode,
				handler:               tt.fields.handler,
				server:                tt.fields.server,
				serverTLS:             tt.fields.serverTLS,
				dotfileLoaded:         tt.fields.dotfileLoaded,
			}
			prevEnv := map[string]string{}
			for k, v := range tt.env {
				prevEnv[k] = os.Getenv(k)
				os.Setenv(k, v)
			}
			defer func() {
				for k, v := range prevEnv {
					os.Setenv(k, v)
				}
			}()
			dir, err := os.MkdirTemp("", "load-template-map")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(dir)
			w.ServeFolder = dir
			if tt.templateFileContent != "" {
				templatesFilePath := path.Join(dir, "templates.yaml")
				if err := os.WriteFile(templatesFilePath, []byte(tt.templateFileContent), 0644); err != nil {
					t.Fatalf("failed to write headers file: %v", err)
				}
				w.TemplateMapPath = templatesFilePath
			}
			w.handler = w.newHandlerForWebServer()
			w, err = w.LoadTemplateMap()
			if err != nil && tt.errorMessage == "" {
				t.Fatalf("error loading templatemap: %v", err)
			} else if err != nil && !strings.Contains(err.Error(), tt.errorMessage) {
				t.Fatalf("unexpected error loading templatemap: %v", err)
			}
			if got := w.handler.TemplateMap; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WebServer.LoadTemplateMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebServer_SetTemplateMap(t *testing.T) {
	type fields struct {
		AppPort               string
		HTTPAllowedOrigins    []string
		EnvFile               string
		Error404FilePath      string
		ExtraHandlers         []*ExtraHandler
		ExtraMiddleware       []func(http.Handler) http.Handler
		GzipEnabled           bool
		HTTPPort              string
		HTTPSPort             string
		HTTPSPortEnabled      bool
		HeaderMap             map[string][]string
		HeaderMapEnabled      bool
		HeaderMapPath         string
		HealthPort            string
		HealthPortEnabled     bool
		MetricsPort           string
		MetricsPortEnabled    bool
		RealIPHeader          string
		RedirectRoutes        map[string]string
		RedirectRoutesEnabled bool
		RedirectRoutesPath    string
		ServeFolder           string
		TLSCertPath           string
		TLSConfig             *tls.Config
		TLSKeyPath            string
		TemplateMap           map[string]string
		TemplateMapEnabled    bool
		TemplateMapPath       string
		VueJSHistoryMode      bool
		handler               *handlers.Handler
		server                *http.Server
		serverTLS             *http.Server
		dotfileLoaded         bool
	}
	type args struct {
		input map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		env    map[string]string
		want   map[string]string
	}{
		{
			name: "basic",
			args: args{
				input: map[string]string{
					"X-Something-Here": "Value",
				},
			},
			want: map[string]string{
				"X-Something-Here": "Value",
			},
		},
		{
			name: "eval env",
			env: map[string]string{
				"COOL_123": "It is",
			},
			args: args{
				input: map[string]string{
					"X-Something-Here": "${COOL_123}",
				},
			},
			want: map[string]string{
				"X-Something-Here": "It is",
			},
		},
		{
			name: "dotfile no eval env",
			fields: fields{
				dotfileLoaded: true,
			},
			env: map[string]string{
				"COOL_123": "It is",
			},
			args: args{
				input: map[string]string{
					"X-Something-Here": "${COOL_123}",
				},
			},
			want: map[string]string{
				"X-Something-Here": "${COOL_123}",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		// NOTE cannot be parallel since using env
		t.Run(tt.name, func(t *testing.T) {
			w := &WebServer{
				AppPort:               tt.fields.AppPort,
				HTTPAllowedOrigins:    tt.fields.HTTPAllowedOrigins,
				EnvFile:               tt.fields.EnvFile,
				Error404FilePath:      tt.fields.Error404FilePath,
				ExtraHandlers:         tt.fields.ExtraHandlers,
				ExtraMiddleware:       tt.fields.ExtraMiddleware,
				GzipEnabled:           tt.fields.GzipEnabled,
				HTTPPort:              tt.fields.HTTPPort,
				HTTPSPort:             tt.fields.HTTPSPort,
				HTTPSPortEnabled:      tt.fields.HTTPSPortEnabled,
				HeaderMap:             tt.fields.HeaderMap,
				HeaderMapEnabled:      tt.fields.HeaderMapEnabled,
				HeaderMapPath:         tt.fields.HeaderMapPath,
				HealthPort:            tt.fields.HealthPort,
				HealthPortEnabled:     tt.fields.HealthPortEnabled,
				MetricsPort:           tt.fields.MetricsPort,
				MetricsPortEnabled:    tt.fields.MetricsPortEnabled,
				RealIPHeader:          tt.fields.RealIPHeader,
				RedirectRoutes:        tt.fields.RedirectRoutes,
				RedirectRoutesEnabled: tt.fields.RedirectRoutesEnabled,
				RedirectRoutesPath:    tt.fields.RedirectRoutesPath,
				ServeFolder:           tt.fields.ServeFolder,
				TLSCertPath:           tt.fields.TLSCertPath,
				TLSConfig:             tt.fields.TLSConfig,
				TLSKeyPath:            tt.fields.TLSKeyPath,
				TemplateMap:           tt.fields.TemplateMap,
				TemplateMapEnabled:    tt.fields.TemplateMapEnabled,
				TemplateMapPath:       tt.fields.TemplateMapPath,
				VueJSHistoryMode:      tt.fields.VueJSHistoryMode,
				handler:               tt.fields.handler,
				server:                tt.fields.server,
				serverTLS:             tt.fields.serverTLS,
				dotfileLoaded:         tt.fields.dotfileLoaded,
			}
			prevEnv := map[string]string{}
			for k, v := range tt.env {
				prevEnv[k] = os.Getenv(k)
				os.Setenv(k, v)
			}
			defer func() {
				for k, v := range prevEnv {
					os.Setenv(k, v)
				}
			}()
			w.handler = w.newHandlerForWebServer()
			if got := w.SetTemplateMap(tt.args.input).handler.TemplateMap; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WebServer.SetTemplateMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebServer_LoadHeaderMap(t *testing.T) {
	type fields struct {
		AppPort               string
		HTTPAllowedOrigins    []string
		EnvFile               string
		Error404FilePath      string
		ExtraHandlers         []*ExtraHandler
		ExtraMiddleware       []func(http.Handler) http.Handler
		GzipEnabled           bool
		HTTPPort              string
		HTTPSPort             string
		HTTPSPortEnabled      bool
		HeaderMap             map[string][]string
		HeaderMapEnabled      bool
		HeaderMapPath         string
		HealthPort            string
		HealthPortEnabled     bool
		MetricsPort           string
		MetricsPortEnabled    bool
		RealIPHeader          string
		RedirectRoutes        map[string]string
		RedirectRoutesEnabled bool
		RedirectRoutesPath    string
		ServeFolder           string
		TLSCertPath           string
		TLSConfig             *tls.Config
		TLSKeyPath            string
		TemplateMap           map[string]string
		TemplateMapEnabled    bool
		TemplateMapPath       string
		VueJSHistoryMode      bool
		handler               *handlers.Handler
		server                *http.Server
		serverTLS             *http.Server
		dotfileLoaded         bool
	}
	tests := []struct {
		name              string
		fields            fields
		headerFileContent string
		env               map[string]string
		errorMessage      string
		want              map[string][]string
	}{
		{
			name: "basic with env",
			fields: fields{
				dotfileLoaded: false,
			},
			env: map[string]string{
				"THINGY": "THIS",
			},
			headerFileContent: `---
X-Abc:
  - "123"
Something:
  - "${THINGY}"
`,
			want: map[string][]string{
				"X-Abc":     {"123"},
				"Something": {"THIS"},
			},
		},
		{
			name: "using dotfiles",
			fields: fields{
				dotfileLoaded: true,
				HeaderMap: map[string][]string{
					"X-123": {"Abc"},
				},
			},
			want: map[string][]string{
				"X-123": {"Abc"},
			},
		},
		{
			name: "file doesn't exist",
			fields: fields{
				dotfileLoaded: false,
			},
			errorMessage: "error: header template not found",
			want:         nil,
		},
		{
			name: "fails to load headermap config",
			fields: fields{
				dotfileLoaded: false,
			},
			headerFileContent: `%#%@?#:`,
			errorMessage:      "error converting YAML to JSON",
			want:              nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		// NOTE cannot be parallel because using env
		t.Run(tt.name, func(t *testing.T) {
			w := &WebServer{
				AppPort:               tt.fields.AppPort,
				HTTPAllowedOrigins:    tt.fields.HTTPAllowedOrigins,
				EnvFile:               tt.fields.EnvFile,
				Error404FilePath:      tt.fields.Error404FilePath,
				ExtraHandlers:         tt.fields.ExtraHandlers,
				ExtraMiddleware:       tt.fields.ExtraMiddleware,
				GzipEnabled:           tt.fields.GzipEnabled,
				HTTPPort:              tt.fields.HTTPPort,
				HTTPSPort:             tt.fields.HTTPSPort,
				HTTPSPortEnabled:      tt.fields.HTTPSPortEnabled,
				HeaderMap:             tt.fields.HeaderMap,
				HeaderMapEnabled:      tt.fields.HeaderMapEnabled,
				HeaderMapPath:         tt.fields.HeaderMapPath,
				HealthPort:            tt.fields.HealthPort,
				HealthPortEnabled:     tt.fields.HealthPortEnabled,
				MetricsPort:           tt.fields.MetricsPort,
				MetricsPortEnabled:    tt.fields.MetricsPortEnabled,
				RealIPHeader:          tt.fields.RealIPHeader,
				RedirectRoutes:        tt.fields.RedirectRoutes,
				RedirectRoutesEnabled: tt.fields.RedirectRoutesEnabled,
				RedirectRoutesPath:    tt.fields.RedirectRoutesPath,
				ServeFolder:           tt.fields.ServeFolder,
				TLSCertPath:           tt.fields.TLSCertPath,
				TLSConfig:             tt.fields.TLSConfig,
				TLSKeyPath:            tt.fields.TLSKeyPath,
				TemplateMap:           tt.fields.TemplateMap,
				TemplateMapEnabled:    tt.fields.TemplateMapEnabled,
				TemplateMapPath:       tt.fields.TemplateMapPath,
				VueJSHistoryMode:      tt.fields.VueJSHistoryMode,
				handler:               tt.fields.handler,
				server:                tt.fields.server,
				serverTLS:             tt.fields.serverTLS,
				dotfileLoaded:         tt.fields.dotfileLoaded,
			}
			prevEnv := map[string]string{}
			for k, v := range tt.env {
				prevEnv[k] = os.Getenv(k)
				os.Setenv(k, v)
			}
			defer func() {
				for k, v := range prevEnv {
					os.Setenv(k, v)
				}
			}()
			dir, err := os.MkdirTemp("", "load-header-map")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(dir)
			w.ServeFolder = dir
			if tt.headerFileContent != "" {
				headersFilePath := path.Join(dir, "headers.yaml")
				if err := os.WriteFile(headersFilePath, []byte(tt.headerFileContent), 0644); err != nil {
					t.Fatalf("failed to write headers file: %v", err)
				}
				w.HeaderMapPath = headersFilePath
			}
			w.handler = w.newHandlerForWebServer()
			w, err = w.LoadHeaderMap()
			if err != nil && tt.errorMessage == "" {
				t.Fatalf("error loading headermap: %v", err)
			} else if err != nil && !strings.Contains(err.Error(), tt.errorMessage) {
				t.Fatalf("unexpected error loading headermap: %v", err)
			}
			if got := w.handler.HeaderMap; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WebServer.LoadHeaderMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebServer_SetHeaderMap(t *testing.T) {
	type fields struct {
		AppPort               string
		HTTPAllowedOrigins    []string
		EnvFile               string
		Error404FilePath      string
		ExtraHandlers         []*ExtraHandler
		ExtraMiddleware       []func(http.Handler) http.Handler
		GzipEnabled           bool
		HTTPPort              string
		HTTPSPort             string
		HTTPSPortEnabled      bool
		HeaderMap             map[string][]string
		HeaderMapEnabled      bool
		HeaderMapPath         string
		HealthPort            string
		HealthPortEnabled     bool
		MetricsPort           string
		MetricsPortEnabled    bool
		RealIPHeader          string
		RedirectRoutes        map[string]string
		RedirectRoutesEnabled bool
		RedirectRoutesPath    string
		ServeFolder           string
		TLSCertPath           string
		TLSConfig             *tls.Config
		TLSKeyPath            string
		TemplateMap           map[string]string
		TemplateMapEnabled    bool
		TemplateMapPath       string
		VueJSHistoryMode      bool
		handler               *handlers.Handler
		server                *http.Server
		serverTLS             *http.Server
		dotfileLoaded         bool
	}
	type args struct {
		input map[string][]string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		env    map[string]string
		want   map[string][]string
	}{
		{
			name: "basic",
			args: args{
				input: map[string][]string{
					"X-Something-Here": {"Value"},
				},
			},
			want: map[string][]string{
				"X-Something-Here": {"Value"},
			},
		},
		{
			name: "eval env",
			env: map[string]string{
				"COOL_123": "It is",
			},
			args: args{
				input: map[string][]string{
					"X-Something-Here": {"${COOL_123}"},
				},
			},
			want: map[string][]string{
				"X-Something-Here": {"It is"},
			},
		},
		{
			name: "dotfile no eval env",
			fields: fields{
				dotfileLoaded: true,
			},
			env: map[string]string{
				"COOL_123": "It is",
			},
			args: args{
				input: map[string][]string{
					"X-Something-Here": {"${COOL_123}"},
				},
			},
			want: map[string][]string{
				"X-Something-Here": {"${COOL_123}"},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			w := &WebServer{
				AppPort:               tt.fields.AppPort,
				HTTPAllowedOrigins:    tt.fields.HTTPAllowedOrigins,
				EnvFile:               tt.fields.EnvFile,
				Error404FilePath:      tt.fields.Error404FilePath,
				ExtraHandlers:         tt.fields.ExtraHandlers,
				ExtraMiddleware:       tt.fields.ExtraMiddleware,
				GzipEnabled:           tt.fields.GzipEnabled,
				HTTPPort:              tt.fields.HTTPPort,
				HTTPSPort:             tt.fields.HTTPSPort,
				HTTPSPortEnabled:      tt.fields.HTTPSPortEnabled,
				HeaderMap:             tt.fields.HeaderMap,
				HeaderMapEnabled:      tt.fields.HeaderMapEnabled,
				HeaderMapPath:         tt.fields.HeaderMapPath,
				HealthPort:            tt.fields.HealthPort,
				HealthPortEnabled:     tt.fields.HealthPortEnabled,
				MetricsPort:           tt.fields.MetricsPort,
				MetricsPortEnabled:    tt.fields.MetricsPortEnabled,
				RealIPHeader:          tt.fields.RealIPHeader,
				RedirectRoutes:        tt.fields.RedirectRoutes,
				RedirectRoutesEnabled: tt.fields.RedirectRoutesEnabled,
				RedirectRoutesPath:    tt.fields.RedirectRoutesPath,
				ServeFolder:           tt.fields.ServeFolder,
				TLSCertPath:           tt.fields.TLSCertPath,
				TLSConfig:             tt.fields.TLSConfig,
				TLSKeyPath:            tt.fields.TLSKeyPath,
				TemplateMap:           tt.fields.TemplateMap,
				TemplateMapEnabled:    tt.fields.TemplateMapEnabled,
				TemplateMapPath:       tt.fields.TemplateMapPath,
				VueJSHistoryMode:      tt.fields.VueJSHistoryMode,
				handler:               tt.fields.handler,
				server:                tt.fields.server,
				serverTLS:             tt.fields.serverTLS,
				dotfileLoaded:         tt.fields.dotfileLoaded,
			}
			prevEnv := map[string]string{}
			for k, v := range tt.env {
				prevEnv[k] = os.Getenv(k)
				os.Setenv(k, v)
			}
			defer func() {
				for k, v := range prevEnv {
					os.Setenv(k, v)
				}
			}()
			w.handler = w.newHandlerForWebServer()
			if got := w.SetHeaderMap(tt.args.input).handler.HeaderMap; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WebServer.SetHeaderMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebServer_newHandlerForWebServer(t *testing.T) {
	type fields struct {
		AppPort               string
		HTTPAllowedOrigins    []string
		EnvFile               string
		Error404FilePath      string
		ExtraHandlers         []*ExtraHandler
		ExtraMiddleware       []func(http.Handler) http.Handler
		GzipEnabled           bool
		HTTPPort              string
		HTTPSPort             string
		HTTPSPortEnabled      bool
		HeaderMap             map[string][]string
		HeaderMapEnabled      bool
		HeaderMapPath         string
		HealthPort            string
		HealthPortEnabled     bool
		MetricsPort           string
		MetricsPortEnabled    bool
		RealIPHeader          string
		RedirectRoutes        map[string]string
		RedirectRoutesEnabled bool
		RedirectRoutesPath    string
		ServeFolder           string
		TLSCertPath           string
		TLSConfig             *tls.Config
		TLSKeyPath            string
		TemplateMap           map[string]string
		TemplateMapEnabled    bool
		TemplateMapPath       string
		VueJSHistoryMode      bool
		handler               *handlers.Handler
		server                *http.Server
		serverTLS             *http.Server
	}
	tests := []struct {
		name   string
		fields fields
		want   *handlers.Handler
	}{
		{
			name: "basic",
			fields: fields{
				ServeFolder:        "./folder",
				VueJSHistoryMode:   false,
				TemplateMapEnabled: true,
				HeaderMapEnabled:   true,
				Error404FilePath:   "./404.html",
				GzipEnabled:        true,
				HeaderMap: map[string][]string{
					"A": {"B"},
				},
				TemplateMap: map[string]string{"A": "B"},
			},
			want: &handlers.Handler{
				ServeFolder:        "./folder",
				VueJSHistoryMode:   false,
				TemplateMapEnabled: true,
				HeaderMapEnabled:   true,
				Error404FilePath:   "./404.html",
				GzipEnabled:        true,
				HeaderMap: map[string][]string{
					"A": {"B"},
				},
				TemplateMap: map[string]string{"A": "B"},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			w := &WebServer{
				AppPort:               tt.fields.AppPort,
				HTTPAllowedOrigins:    tt.fields.HTTPAllowedOrigins,
				EnvFile:               tt.fields.EnvFile,
				Error404FilePath:      tt.fields.Error404FilePath,
				ExtraHandlers:         tt.fields.ExtraHandlers,
				ExtraMiddleware:       tt.fields.ExtraMiddleware,
				GzipEnabled:           tt.fields.GzipEnabled,
				HTTPPort:              tt.fields.HTTPPort,
				HTTPSPort:             tt.fields.HTTPSPort,
				HTTPSPortEnabled:      tt.fields.HTTPSPortEnabled,
				HeaderMap:             tt.fields.HeaderMap,
				HeaderMapEnabled:      tt.fields.HeaderMapEnabled,
				HeaderMapPath:         tt.fields.HeaderMapPath,
				HealthPort:            tt.fields.HealthPort,
				HealthPortEnabled:     tt.fields.HealthPortEnabled,
				MetricsPort:           tt.fields.MetricsPort,
				MetricsPortEnabled:    tt.fields.MetricsPortEnabled,
				RealIPHeader:          tt.fields.RealIPHeader,
				RedirectRoutes:        tt.fields.RedirectRoutes,
				RedirectRoutesEnabled: tt.fields.RedirectRoutesEnabled,
				RedirectRoutesPath:    tt.fields.RedirectRoutesPath,
				ServeFolder:           tt.fields.ServeFolder,
				TLSCertPath:           tt.fields.TLSCertPath,
				TLSConfig:             tt.fields.TLSConfig,
				TLSKeyPath:            tt.fields.TLSKeyPath,
				TemplateMap:           tt.fields.TemplateMap,
				TemplateMapEnabled:    tt.fields.TemplateMapEnabled,
				TemplateMapPath:       tt.fields.TemplateMapPath,
				VueJSHistoryMode:      tt.fields.VueJSHistoryMode,
				handler:               tt.fields.handler,
				server:                tt.fields.server,
				serverTLS:             tt.fields.serverTLS,
			}
			if got := w.newHandlerForWebServer(); !reflect.DeepEqual(*got, *tt.want) {
				t.Errorf("WebServer.newHandlerForWebServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebServer_SetHandler(t *testing.T) {
	type fields struct {
		AppPort               string
		HTTPAllowedOrigins    []string
		EnvFile               string
		Error404FilePath      string
		ExtraHandlers         []*ExtraHandler
		ExtraMiddleware       []func(http.Handler) http.Handler
		GzipEnabled           bool
		HTTPPort              string
		HTTPSPort             string
		HTTPSPortEnabled      bool
		HeaderMap             map[string][]string
		HeaderMapEnabled      bool
		HeaderMapPath         string
		HealthPort            string
		HealthPortEnabled     bool
		MetricsPort           string
		MetricsPortEnabled    bool
		RealIPHeader          string
		RedirectRoutes        map[string]string
		RedirectRoutesEnabled bool
		RedirectRoutesPath    string
		ServeFolder           string
		TLSCertPath           string
		TLSConfig             *tls.Config
		TLSKeyPath            string
		TemplateMap           map[string]string
		TemplateMapEnabled    bool
		TemplateMapPath       string
		VueJSHistoryMode      bool
		handler               *handlers.Handler
		server                *http.Server
		serverTLS             *http.Server
	}
	type args struct {
		input *handlers.Handler
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *handlers.Handler
	}{
		{
			name: "basic",
			args: args{
				input: &handlers.Handler{
					VueJSHistoryMode: true,
				},
			},
			want: &handlers.Handler{
				VueJSHistoryMode: true,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			w := &WebServer{
				AppPort:               tt.fields.AppPort,
				HTTPAllowedOrigins:    tt.fields.HTTPAllowedOrigins,
				EnvFile:               tt.fields.EnvFile,
				Error404FilePath:      tt.fields.Error404FilePath,
				ExtraHandlers:         tt.fields.ExtraHandlers,
				ExtraMiddleware:       tt.fields.ExtraMiddleware,
				GzipEnabled:           tt.fields.GzipEnabled,
				HTTPPort:              tt.fields.HTTPPort,
				HTTPSPort:             tt.fields.HTTPSPort,
				HTTPSPortEnabled:      tt.fields.HTTPSPortEnabled,
				HeaderMap:             tt.fields.HeaderMap,
				HeaderMapEnabled:      tt.fields.HeaderMapEnabled,
				HeaderMapPath:         tt.fields.HeaderMapPath,
				HealthPort:            tt.fields.HealthPort,
				HealthPortEnabled:     tt.fields.HealthPortEnabled,
				MetricsPort:           tt.fields.MetricsPort,
				MetricsPortEnabled:    tt.fields.MetricsPortEnabled,
				RealIPHeader:          tt.fields.RealIPHeader,
				RedirectRoutes:        tt.fields.RedirectRoutes,
				RedirectRoutesEnabled: tt.fields.RedirectRoutesEnabled,
				RedirectRoutesPath:    tt.fields.RedirectRoutesPath,
				ServeFolder:           tt.fields.ServeFolder,
				TLSCertPath:           tt.fields.TLSCertPath,
				TLSConfig:             tt.fields.TLSConfig,
				TLSKeyPath:            tt.fields.TLSKeyPath,
				TemplateMap:           tt.fields.TemplateMap,
				TemplateMapEnabled:    tt.fields.TemplateMapEnabled,
				TemplateMapPath:       tt.fields.TemplateMapPath,
				VueJSHistoryMode:      tt.fields.VueJSHistoryMode,
				handler:               tt.fields.handler,
				server:                tt.fields.server,
				serverTLS:             tt.fields.serverTLS,
			}
			if got := w.SetHandler(tt.args.input).handler; !reflect.DeepEqual(*got, *tt.want) {
				t.Errorf("WebServer.SetHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebServer_NewMetricsFromWebServer(t *testing.T) {
	type fields struct {
		AppPort               string
		HTTPAllowedOrigins    []string
		EnvFile               string
		Error404FilePath      string
		ExtraHandlers         []*ExtraHandler
		ExtraMiddleware       []func(http.Handler) http.Handler
		GzipEnabled           bool
		HTTPPort              string
		HTTPSPort             string
		HTTPSPortEnabled      bool
		HeaderMap             map[string][]string
		HeaderMapEnabled      bool
		HeaderMapPath         string
		HealthPort            string
		HealthPortEnabled     bool
		MetricsPort           string
		MetricsPortEnabled    bool
		RealIPHeader          string
		RedirectRoutes        map[string]string
		RedirectRoutesEnabled bool
		RedirectRoutesPath    string
		ServeFolder           string
		TLSCertPath           string
		TLSConfig             *tls.Config
		TLSKeyPath            string
		TemplateMap           map[string]string
		TemplateMapEnabled    bool
		TemplateMapPath       string
		VueJSHistoryMode      bool
		handler               *handlers.Handler
		server                *http.Server
		serverTLS             *http.Server
	}
	tests := []struct {
		name   string
		fields fields
		want   *metrics.Metrics
	}{
		{
			name: "settings set",
			fields: fields{
				MetricsPortEnabled: true,
				MetricsPort:        ":91234",
			},
			want: &metrics.Metrics{
				Enabled: true,
				Port:    ":91234",
			},
		},
		{
			name: "disabled",
			want: &metrics.Metrics{
				Enabled: false,
			},
		},
		{
			name: "not set",
			want: &metrics.Metrics{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			w := &WebServer{
				AppPort:               tt.fields.AppPort,
				HTTPAllowedOrigins:    tt.fields.HTTPAllowedOrigins,
				EnvFile:               tt.fields.EnvFile,
				Error404FilePath:      tt.fields.Error404FilePath,
				ExtraHandlers:         tt.fields.ExtraHandlers,
				ExtraMiddleware:       tt.fields.ExtraMiddleware,
				GzipEnabled:           tt.fields.GzipEnabled,
				HTTPPort:              tt.fields.HTTPPort,
				HTTPSPort:             tt.fields.HTTPSPort,
				HTTPSPortEnabled:      tt.fields.HTTPSPortEnabled,
				HeaderMap:             tt.fields.HeaderMap,
				HeaderMapEnabled:      tt.fields.HeaderMapEnabled,
				HeaderMapPath:         tt.fields.HeaderMapPath,
				HealthPort:            tt.fields.HealthPort,
				HealthPortEnabled:     tt.fields.HealthPortEnabled,
				MetricsPort:           tt.fields.MetricsPort,
				MetricsPortEnabled:    tt.fields.MetricsPortEnabled,
				RealIPHeader:          tt.fields.RealIPHeader,
				RedirectRoutes:        tt.fields.RedirectRoutes,
				RedirectRoutesEnabled: tt.fields.RedirectRoutesEnabled,
				RedirectRoutesPath:    tt.fields.RedirectRoutesPath,
				ServeFolder:           tt.fields.ServeFolder,
				TLSCertPath:           tt.fields.TLSCertPath,
				TLSConfig:             tt.fields.TLSConfig,
				TLSKeyPath:            tt.fields.TLSKeyPath,
				TemplateMap:           tt.fields.TemplateMap,
				TemplateMapEnabled:    tt.fields.TemplateMapEnabled,
				TemplateMapPath:       tt.fields.TemplateMapPath,
				VueJSHistoryMode:      tt.fields.VueJSHistoryMode,
				handler:               tt.fields.handler,
				server:                tt.fields.server,
				serverTLS:             tt.fields.serverTLS,
			}
			if got := w.NewMetricsFromWebServer(); !reflect.DeepEqual(*got, *tt.want) {
				t.Errorf("WebServer.NewMetricsFromWebServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebServer_Listen(t *testing.T) {
	type fields struct {
		AppPort               string
		HTTPAllowedOrigins    []string
		EnvFile               string
		Error404FilePath      string
		ExtraHandlers         []*ExtraHandler
		ExtraMiddleware       []func(http.Handler) http.Handler
		GzipEnabled           bool
		HTTPPort              string
		HTTPSPort             string
		HTTPSPortEnabled      bool
		HeaderMap             map[string][]string
		HeaderMapEnabled      bool
		HeaderMapPath         string
		HealthPort            string
		HealthPortEnabled     bool
		MetricsPort           string
		MetricsPortEnabled    bool
		RealIPHeader          string
		RedirectRoutes        map[string]string
		RedirectRoutesEnabled bool
		RedirectRoutesPath    string
		ServeFolder           string
		TLSCertPath           string
		TLSConfig             *tls.Config
		TLSKeyPath            string
		TemplateMap           map[string]string
		TemplateMapEnabled    bool
		TemplateMapPath       string
		VueJSHistoryMode      bool
		handler               *handlers.Handler
		server                *http.Server
		serverTLS             *http.Server
	}
	tests := []struct {
		name          string
		noQuitChannel bool
		fields        fields
	}{
		{
			name: "basic",
			fields: fields{
				AppPort:     fmt.Sprintf(":%v", rand.Intn(65000-50000)+50000),
				MetricsPort: fmt.Sprintf(":%v", rand.Intn(65000-50000)+50000),
			},
		},
		{
			name: "tls",
			fields: fields{
				AppPort:          fmt.Sprintf(":%v", rand.Intn(65000-50000)+50000),
				HTTPSPort:        fmt.Sprintf(":%v", rand.Intn(65000-50000)+50000),
				HTTPSPortEnabled: true,
				MetricsPort:      fmt.Sprintf(":%v", rand.Intn(65000-50000)+50000),
			},
		},
		// {
		// 	name: "no quit channel",
		// 	fields: fields{
		// 		AppPort:     fmt.Sprintf(":%v", rand.Intn(65000-50000)+50000),
		// 		MetricsPort: fmt.Sprintf(":%v", rand.Intn(65000-50000)+50000),
		// 	},
		// 	noQuitChannel: true,
		// },
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			w := &WebServer{
				AppPort:               tt.fields.AppPort,
				HTTPAllowedOrigins:    tt.fields.HTTPAllowedOrigins,
				EnvFile:               tt.fields.EnvFile,
				Error404FilePath:      tt.fields.Error404FilePath,
				ExtraHandlers:         tt.fields.ExtraHandlers,
				ExtraMiddleware:       tt.fields.ExtraMiddleware,
				GzipEnabled:           tt.fields.GzipEnabled,
				HTTPPort:              tt.fields.HTTPPort,
				HTTPSPort:             tt.fields.HTTPSPort,
				HTTPSPortEnabled:      tt.fields.HTTPSPortEnabled,
				HeaderMap:             tt.fields.HeaderMap,
				HeaderMapEnabled:      tt.fields.HeaderMapEnabled,
				HeaderMapPath:         tt.fields.HeaderMapPath,
				HealthPort:            tt.fields.HealthPort,
				HealthPortEnabled:     tt.fields.HealthPortEnabled,
				MetricsPort:           tt.fields.MetricsPort,
				MetricsPortEnabled:    tt.fields.MetricsPortEnabled,
				RealIPHeader:          tt.fields.RealIPHeader,
				RedirectRoutes:        tt.fields.RedirectRoutes,
				RedirectRoutesEnabled: tt.fields.RedirectRoutesEnabled,
				RedirectRoutesPath:    tt.fields.RedirectRoutesPath,
				ServeFolder:           tt.fields.ServeFolder,
				TLSCertPath:           tt.fields.TLSCertPath,
				TLSConfig:             tt.fields.TLSConfig,
				TLSKeyPath:            tt.fields.TLSKeyPath,
				TemplateMap:           tt.fields.TemplateMap,
				TemplateMapEnabled:    tt.fields.TemplateMapEnabled,
				TemplateMapPath:       tt.fields.TemplateMapPath,
				VueJSHistoryMode:      tt.fields.VueJSHistoryMode,
				handler:               tt.fields.handler,
				server:                tt.fields.server,
				serverTLS:             tt.fields.serverTLS,
			}
			dir, err := os.MkdirTemp("", "load-listen")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(dir)
			w.server = &http.Server{
				Addr: w.AppPort,
			}
			if w.HTTPSPortEnabled {
				publicKeyPath := path.Join(dir, "tls.cert")
				privateKeyPath := path.Join(dir, "tls.key")
				if err := os.WriteFile(publicKeyPath, []byte(tlsPublic), 0644); err != nil {
					t.Fatal(err)
				}
				w.TLSCertPath = publicKeyPath
				if err := os.WriteFile(privateKeyPath, []byte(tlsPrivate), 0644); err != nil {
					t.Fatal(err)
				}
				w.TLSKeyPath = privateKeyPath
				if _, err := w.LoadTLS(); err != nil {
					t.Fatalf("error: failed to load TLS: %v\n", err)
				}
				w.serverTLS = &http.Server{
					Addr:      w.HTTPSPort,
					TLSConfig: w.TLSConfig,
				}
			}
			w.newHandlerForWebServer()
			quitChan := make(chan bool, 1)
			if tt.noQuitChannel {
				quitChan = nil
			}
			go w.Listen(quitChan)
			if !tt.noQuitChannel {
				defer func() {
					quitChan <- true
				}()
			}
		})
	}
}
