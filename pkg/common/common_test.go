/*
	common function calls
*/

package common

import (
	"net/http"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"
)

type responseWriter struct {
	Headers    http.Header
	Body       []byte
	StatusCode int
}

func newResponseWriter() responseWriter {
	return responseWriter{
		Headers:    http.Header{},
		Body:       nil,
		StatusCode: 0,
	}
}

func (w responseWriter) Header() http.Header {
	return w.Headers
}

func (w responseWriter) Write(i []byte) (int, error) {
	copy(w.Body, i)
	return 0, nil
}

func (w responseWriter) WriteHeader(statusCode int) {
	//nolint:staticcheck
	w.StatusCode = statusCode
}

func TestGetAppHealthPortEnabled(t *testing.T) {
	tests := []struct {
		name       string
		env        map[string]string
		wantOutput bool
	}{
		{
			name:       "basic",
			env:        nil,
			wantOutput: true,
		},
		{
			name:       "set env",
			env:        map[string]string{"APP_HEALTH_PORT_ENABLED": "false"},
			wantOutput: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
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

			if gotOutput := GetAppHealthPortEnabled(); gotOutput != tt.wantOutput {
				t.Errorf("GetAppHealthPortEnabled() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestGetAppHealthPort(t *testing.T) {
	tests := []struct {
		name       string
		env        map[string]string
		wantOutput string
	}{
		{
			name:       "basic",
			env:        nil,
			wantOutput: ":8081",
		},
		{
			name:       "set env",
			env:        map[string]string{"APP_HEALTH_PORT": ":2024"},
			wantOutput: ":2024",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
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

			if gotOutput := GetAppHealthPort(); gotOutput != tt.wantOutput {
				t.Errorf("GetAppHealthPort() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestGetAppPort(t *testing.T) {
	tests := []struct {
		name       string
		env        map[string]string
		wantOutput string
	}{
		{
			name:       "basic",
			env:        nil,
			wantOutput: ":8080",
		},
		{
			name:       "set env",
			env:        map[string]string{"APP_PORT": ":7334"},
			wantOutput: ":7334",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
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

			if gotOutput := GetAppPort(); gotOutput != tt.wantOutput {
				t.Errorf("GetAppPort() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestGetAppHTTPSPort(t *testing.T) {
	tests := []struct {
		name       string
		env        map[string]string
		wantOutput string
	}{
		{
			name:       "basic",
			env:        nil,
			wantOutput: ":8443",
		},
		{
			name:       "set env",
			env:        map[string]string{"APP_HTTPS_PORT": ":10000"},
			wantOutput: ":10000",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
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

			if gotOutput := GetAppHTTPSPort(); gotOutput != tt.wantOutput {
				t.Errorf("GetAppHTTPSPort() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestGetAppHTTPSCrtPath(t *testing.T) {
	tests := []struct {
		name       string
		env        map[string]string
		wantOutput string
	}{
		{
			name:       "basic",
			env:        nil,
			wantOutput: "",
		},
		{
			name:       "set env",
			env:        map[string]string{"APP_HTTPS_CRT_PATH": "/tmp/a.crt"},
			wantOutput: "/tmp/a.crt",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
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

			if gotOutput := GetAppHTTPSCrtPath(); gotOutput != tt.wantOutput {
				t.Errorf("GetAppHTTPSCrtPath() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestGetAppHTTPSKeyPath(t *testing.T) {
	tests := []struct {
		name       string
		env        map[string]string
		wantOutput string
	}{
		{
			name:       "basic",
			env:        nil,
			wantOutput: "",
		},
		{
			name:       "set env",
			env:        map[string]string{"APP_HTTPS_KEY_PATH": "/tmp/a.key"},
			wantOutput: "/tmp/a.key",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
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

			if gotOutput := GetAppHTTPSKeyPath(); gotOutput != tt.wantOutput {
				t.Errorf("GetAppHTTPSKeyPath() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestGetAppEnableHTTPS(t *testing.T) {
	tests := []struct {
		name       string
		env        map[string]string
		wantOutput bool
	}{
		{
			name:       "basic",
			env:        nil,
			wantOutput: false,
		},
		{
			name:       "set env",
			env:        map[string]string{"APP_ENABLE_HTTPS": "true"},
			wantOutput: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
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

			if gotOutput := GetAppEnableHTTPS(); gotOutput != tt.wantOutput {
				t.Errorf("GetAppEnableHTTPS() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestGetAppMetricsPort(t *testing.T) {
	tests := []struct {
		name       string
		env        map[string]string
		wantOutput string
	}{
		{
			name:       "basic",
			env:        nil,
			wantOutput: ":2112",
		},
		{
			name:       "set env",
			env:        map[string]string{"APP_PORT_METRICS": ":7373"},
			wantOutput: ":7373",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
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

			if gotOutput := GetAppMetricsPort(); gotOutput != tt.wantOutput {
				t.Errorf("GetAppMetricsPort() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestGetAppMetricsEnabled(t *testing.T) {
	tests := []struct {
		name       string
		env        map[string]string
		wantOutput bool
	}{
		{
			name:       "basic",
			env:        nil,
			wantOutput: true,
		},
		{
			name:       "set env",
			env:        map[string]string{"APP_METRICS_ENABLED": "false"},
			wantOutput: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
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

			if gotOutput := GetAppMetricsEnabled(); gotOutput != tt.wantOutput {
				t.Errorf("GetAppMetricsEnabled() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestGetAppRealIPHeader(t *testing.T) {
	tests := []struct {
		name       string
		env        map[string]string
		wantOutput string
	}{
		{
			name:       "basic",
			env:        nil,
			wantOutput: "",
		},
		{
			name:       "set env",
			env:        map[string]string{"APP_HTTP_REAL_IP_HEADER": "X-Real-Ip"},
			wantOutput: "X-Real-Ip",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
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

			if gotOutput := GetAppRealIPHeader(); gotOutput != tt.wantOutput {
				t.Errorf("GetAppRealIPHeader() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestGetServeFolder(t *testing.T) {
	tests := []struct {
		name       string
		env        map[string]string
		wantOutput string
	}{
		{
			name: "basic",
			env:  nil,
			wantOutput: func() string {
				pwd, _ := os.Getwd()
				return pwd
			}(),
		},
		{
			name:       "set env",
			env:        map[string]string{"APP_SERVE_FOLDER": "/var/run/cool"},
			wantOutput: "/var/run/cool",
		},
		{
			name:       "set env with KO_DATA_PATH",
			env:        map[string]string{"KO_DATA_PATH": "/var/run/ko"},
			wantOutput: "/var/run/ko",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
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

			if gotOutput := GetServeFolder(); gotOutput != tt.wantOutput {
				t.Errorf("GetServeFolder() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestGetTemplateMapPath(t *testing.T) {
	tests := []struct {
		name       string
		env        map[string]string
		wantOutput string
	}{
		{
			name:       "basic",
			env:        nil,
			wantOutput: "./template-map.yaml",
		},
		{
			name:       "set env",
			env:        map[string]string{"APP_TEMPLATE_MAP_PATH": "/app/inputs.yaml"},
			wantOutput: "/app/inputs.yaml",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
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

			if gotOutput := GetTemplateMapPath(); gotOutput != tt.wantOutput {
				t.Errorf("GetTemplateMapPath() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestGetVuejsHistoryMode(t *testing.T) {
	tests := []struct {
		name       string
		env        map[string]string
		wantOutput bool
	}{
		{
			name:       "basic",
			env:        nil,
			wantOutput: false,
		},
		{
			name:       "set env",
			env:        map[string]string{"APP_VUEJS_HISTORY_MODE": "true"},
			wantOutput: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
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

			if gotOutput := GetVuejsHistoryMode(); gotOutput != tt.wantOutput {
				t.Errorf("GetVuejsHistoryMode() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestGetEnableGZIP(t *testing.T) {
	tests := []struct {
		name       string
		env        map[string]string
		wantOutput bool
	}{
		{
			name:       "basic",
			env:        nil,
			wantOutput: true,
		},
		{
			name:       "set env",
			env:        map[string]string{"APP_HANDLE_GZIP": "false"},
			wantOutput: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
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

			if gotEnable := GetEnableGZIP(); gotEnable != tt.wantOutput {
				t.Errorf("GetEnableGZIP() = %v, want %v", gotEnable, tt.wantOutput)
			}
		})
	}
}

func TestGetHeaderSetEnable(t *testing.T) {
	tests := []struct {
		name       string
		env        map[string]string
		wantOutput bool
	}{
		{
			name:       "basic",
			env:        nil,
			wantOutput: false,
		},
		{
			name:       "set env",
			env:        map[string]string{"APP_HEADER_SET_ENABLE": "true"},
			wantOutput: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
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

			if gotOutput := GetHeaderSetEnable(); gotOutput != tt.wantOutput {
				t.Errorf("GetHeaderSetEnable() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestGetHeaderMapPath(t *testing.T) {
	tests := []struct {
		name       string
		env        map[string]string
		wantOutput string
	}{
		{
			name:       "basic",
			env:        nil,
			wantOutput: "./headers.yaml",
		},
		{
			name:       "set env",
			env:        map[string]string{"APP_HEADER_MAP_PATH": "/app/headers.yaml"},
			wantOutput: "/app/headers.yaml",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
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

			if gotOutput := GetHeaderMapPath(); gotOutput != tt.wantOutput {
				t.Errorf("GetHeaderMapPath() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestGet404PageFileName(t *testing.T) {
	tests := []struct {
		name       string
		env        map[string]string
		wantOutput string
	}{
		{
			name:       "basic",
			env:        nil,
			wantOutput: "404.html",
		},
		{
			name:       "set env",
			env:        map[string]string{"APP_404_PAGE_FILE_NAME": "404-lollll.html"},
			wantOutput: "404-lollll.html",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
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

			if gotOutput := Get404PageFileName(); gotOutput != tt.wantOutput {
				t.Errorf("Get404PageFileName() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestGetRedirectRoutesEnabled(t *testing.T) {
	tests := []struct {
		name       string
		env        map[string]string
		wantOutput bool
	}{
		{
			name:       "basic",
			env:        nil,
			wantOutput: true,
		},
		{
			name:       "set env",
			env:        map[string]string{"APP_REDIRECT_ROUTES_ENABLED": "false"},
			wantOutput: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
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

			if gotOutput := GetRedirectRoutesEnabled(); gotOutput != tt.wantOutput {
				t.Errorf("GetRedirectRoutesEnabled() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestGetRedirectRoutesPath(t *testing.T) {
	tests := []struct {
		name       string
		env        map[string]string
		wantOutput string
	}{
		{
			name:       "basic",
			env:        nil,
			wantOutput: "./redirects.yaml",
		},
		{
			name:       "set env",
			env:        map[string]string{"APP_REDIRECT_ROUTES_PATH": "/app/redirects.yaml"},
			wantOutput: "/app/redirects.yaml",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
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

			if gotOutput := GetRedirectRoutesPath(); gotOutput != tt.wantOutput {
				t.Errorf("GetRedirectRoutesPath() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}

}

func TestGetHTTPAllowedOrigins(t *testing.T) {
	tests := []struct {
		name         string
		env          map[string]string
		errorMessage string
		wantOrigins  []string
	}{
		{
			name:        "basic",
			env:         map[string]string{"APP_HTTP_ALLOWED_ORIGINS": "example.com"},
			wantOrigins: []string{"example.com"},
		},
		{
			name:        "multiple",
			env:         map[string]string{"APP_HTTP_ALLOWED_ORIGINS": "example.com,islive.xyz"},
			wantOrigins: []string{"example.com", "islive.xyz"},
		},
		{
			name:        "empty",
			wantOrigins: []string{"*"},
		},
		{
			name:        "traling slash",
			env:         map[string]string{"APP_HTTP_ALLOWED_ORIGINS": "example.com,islive.xyz,"},
			wantOrigins: []string{"example.com", "islive.xyz"},
		},
		{
			name:         "bad url",
			env:          map[string]string{"APP_HTTP_ALLOWED_ORIGINS": "%&*#exam???ple.com"},
			wantOrigins:  nil,
			errorMessage: "invalid URL escape",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			gotOrigins, err := GetHTTPAllowedOrigins()
			if err != nil && tt.errorMessage == "" {
				t.Fatalf("error: %v", err)
			} else if err != nil && !strings.Contains(err.Error(), tt.errorMessage) {
				t.Fatalf("unexpected err: %v", err)
			}
			if !reflect.DeepEqual(gotOrigins, tt.wantOrigins) {
				t.Errorf("GetHTTPAllowedOrigins() = %v, want %v", gotOrigins, tt.wantOrigins)
			}
		})
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	type args struct {
		envName      string
		defaultValue string
	}
	tests := []struct {
		name       string
		env        map[string]string
		args       args
		wantOutput string
	}{
		{
			name: "basic",
			args: args{
				envName:      "AAAAAAAAAAA",
				defaultValue: "a",
			},
			wantOutput: "a",
		},
		{
			name: "basic",
			env: map[string]string{
				"AAAAAAAAAAA": "b",
			},
			args: args{
				envName:      "AAAAAAAAAAA",
				defaultValue: "a",
			},
			wantOutput: "b",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prevEnv := map[string]string{}
			for k, v := range tt.env {
				prevEnv[k] = os.Getenv(k)
				os.Setenv(k, v)
			}
			prevEnv[tt.args.envName] = os.Getenv(tt.args.envName)
			defer func() {
				for k, v := range prevEnv {
					os.Setenv(k, v)
				}
			}()
			if gotOutput := GetEnvOrDefault(tt.args.envName, tt.args.defaultValue); gotOutput != tt.wantOutput {
				t.Errorf("GetEnvOrDefault(%v, %v) = %v, want %v", tt.args.envName, tt.args.defaultValue, gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestEvaluateEnvFromMap(t *testing.T) {
	type args struct {
		input   map[string]string
		fromEnv bool
	}
	tests := []struct {
		name       string
		args       args
		env        map[string]string
		wantOutput map[string]string
	}{
		{
			name: "basic",
			env: map[string]string{
				"BBBBBBBB": "bbbb",
			},
			args: args{
				input: map[string]string{
					"AAA": "aaa",
					"BBB": "${BBBBBBBB}",
				},
				fromEnv: true,
			},
			wantOutput: map[string]string{
				"AAA": "aaa",
				"BBB": "bbbb",
			},
		},
		{
			name: "no env",
			env: map[string]string{
				"BBBBBBBB": "bbbb",
			},
			args: args{
				input: map[string]string{
					"AAA": "aaa",
					"BBB": "${BBBBBBBB}",
				},
			},
			wantOutput: map[string]string{
				"AAA": "aaa",
				"BBB": "${BBBBBBBB}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			if gotOutput := EvaluateEnvFromMap(tt.args.input, tt.args.fromEnv); !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("EvaluateEnvFromMap(%v, %v) = %v, want %v", tt.args.input, tt.args.fromEnv, gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestLoadTemplateMapConfig(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name       string
		args       args
		files      map[string]string
		wantOutput map[string]string
		wantErr    bool
	}{
		{
			name: "basic",
			args: args{
				path: "template.yaml",
			},
			files: map[string]string{
				"template.yaml": `---
A: B
B: "123"
SomeKey: Hello!
`,
			},
			wantOutput: map[string]string{
				"A":       "B",
				"B":       "123",
				"SomeKey": "Hello!",
			},
		},
		{
			name:    "no path given",
			wantErr: true,
		},
		{
			name: "bad config",
			args: args{
				path: "template.yaml",
			},
			files: map[string]string{
				"template.yaml": `$#*$*&<<>><!#>`,
			},
			wantErr:    true,
			wantOutput: map[string]string{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dir, err := os.MkdirTemp("", "load-template-map-config")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(dir)
			for f, c := range tt.files {
				if err := os.WriteFile(path.Join(dir, f), []byte(c), 0644); err != nil {
					t.Fatalf("failed to write file: %v", err)
				}
			}
			gotOutput, err := LoadTemplateMapConfig(path.Join(dir, tt.args.path))
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadTemplateMapConfig(%v) error = %v, wantErr %v", tt.args.path, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("LoadTemplateMapConfig(%v) = %v, want %v", tt.args.path, gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestEvaluateEnvFromHeaderMap(t *testing.T) {
	type args struct {
		input   map[string][]string
		fromEnv bool
	}
	tests := []struct {
		name       string
		args       args
		env        map[string]string
		wantOutput map[string][]string
	}{
		{
			name: "basic",
			env:  map[string]string{"BBB": "123"},
			args: args{
				input: map[string][]string{
					"X-Thingy": {"AAA", "${BBB}"},
				},
				fromEnv: true,
			},
			wantOutput: map[string][]string{
				"X-Thingy": {"AAA", "123"},
			},
		},
		{
			name: "no env",
			env:  map[string]string{"BBB": "123"},
			args: args{
				input: map[string][]string{
					"X-Thingy": {"AAA", "${BBB}"},
				},
				fromEnv: false,
			},
			wantOutput: map[string][]string{
				"X-Thingy": {"AAA", "${BBB}"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			if gotOutput := EvaluateEnvFromHeaderMap(tt.args.input, tt.args.fromEnv); !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("EvaluateEnvFromHeaderMap(%v, %v) = %v, want %v", tt.args.input, tt.args.fromEnv, gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestLoadHeaderMapConfig(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name       string
		args       args
		files      map[string]string
		wantOutput map[string][]string
		wantErr    bool
	}{
		{
			name: "basic",
			args: args{
				path: "headers.yaml",
			},
			files: map[string]string{
				"headers.yaml": `---
A:
  - B
B:
  - "123"
SomeKey:
  - Hello!
`,
			},
			wantOutput: map[string][]string{
				"A":       {"B"},
				"B":       {"123"},
				"SomeKey": {"Hello!"},
			},
		},
		{
			name:       "no path given",
			wantErr:    true,
			wantOutput: map[string][]string{},
		},
		{
			name: "no config",
			args: args{
				path: "headers.yaml",
			},
			wantOutput: nil,
		},
		{
			name: "bad config",
			args: args{
				path: "headers.yaml",
			},
			files: map[string]string{
				"headers.yaml": `$#*$*&<<>><!#>`,
			},
			wantErr:    true,
			wantOutput: map[string][]string{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dir, err := os.MkdirTemp("", "load-template-map-config")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(dir)
			for f, c := range tt.files {
				if err := os.WriteFile(path.Join(dir, f), []byte(c), 0644); err != nil {
					t.Fatalf("failed to write file: %v", err)
				}
			}
			gotOutput, err := LoadHeaderMapConfig(path.Join(dir, tt.args.path))
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadHeaderMapConfig(%v) error = %v, wantErr %v", tt.args.path, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("LoadHeaderMapConfig(%v) = %v, want %v", tt.args.path, gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestWriteHeadersToResponse(t *testing.T) {
	type args struct {
		w         http.ResponseWriter
		headerMap map[string][]string
	}
	tests := []struct {
		name string
		args args
		want map[string][]string
	}{
		{
			name: "basic",
			args: args{
				headerMap: map[string][]string{
					"X-Very-Cool": {"Yes indeed"},
				},
				w: newResponseWriter(),
			},
			want: map[string][]string{
				"X-Very-Cool": {"Yes indeed"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteHeadersToResponse(tt.args.w, tt.args.headerMap)
			foundHeaders := 0
			for k, v := range tt.want {
				if val := tt.args.w.Header().Get(k); val == strings.Join(v, " ") {
					foundHeaders++
				}
			}
			if foundHeaders != len(tt.want) {
				t.Fatalf("unable to find all headers")
			}
		})
	}
}

func TestGetRequestIP(t *testing.T) {
	RealIPHeader := "X-Some-Real-Ip-Header"
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name          string
		args          args
		wantRequestIP string
	}{
		{
			name: "basic",
			args: args{
				r: func() *http.Request {
					req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
					// TODO use from httpserver.WebServer
					os.Setenv("APP_HTTP_REAL_IP_HEADER", RealIPHeader)
					req.Header.Set(RealIPHeader, "123.456.789.12")
					return req
				}(),
			},
			wantRequestIP: "123.456.789.12",
		},
	}
	defer func() {
		os.Unsetenv(RealIPHeader)
	}()
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if gotRequestIP := GetRequestIP(tt.args.r); gotRequestIP != tt.wantRequestIP {
				t.Errorf("GetRequestIP(%v) = %v, want %v", tt.args.r, gotRequestIP, tt.wantRequestIP)
			}
		})
	}
}

func Test_statusRecorder_WriteHeader(t *testing.T) {
	type fields struct {
		ResponseWriter http.ResponseWriter
		Status         int
	}
	type args struct {
		status int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "basic",
			fields: fields{
				ResponseWriter: responseWriter{},
				Status:         http.StatusOK,
			},
			args: args{
				status: http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &statusRecorder{
				ResponseWriter: tt.fields.ResponseWriter,
				Status:         tt.fields.Status,
			}
			r.WriteHeader(tt.args.status)
		})
	}
}

func TestLogging(t *testing.T) {
	type args struct {
		next http.Handler
	}
	tests := []struct {
		name string
		args args
		want http.Handler
	}{
		{
			name: "basic",
			args: args{
				next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
			if err != nil {
				t.Fatal(err)
			}
			w := responseWriter{}
			Logging(tt.args.next).ServeHTTP(w, req)
			// TODO add test logic
			// if got := Logging(tt.args.next); !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("Logging(%v) = %v, want %v", tt.args.next, got, tt.want)
			// }
		})
	}
}

func TestLoadDotfileConfig(t *testing.T) {
	type args struct {
		serveFolder string
	}
	tests := []struct {
		name    string
		files   map[string]string
		args    args
		wantCfg *DotfileConfig
		wantErr bool
	}{
		{
			name: "basic",
			files: map[string]string{
				".ghs.yaml": `---
error404FilePath: 404-lol.html
headerMap:
  X-Cool:
    - "Yes"
historyMode: true
redirectRoutes:
  /aaa: http://example.com
templateMap:
  AAA: BBB
`,
			},
			wantCfg: &DotfileConfig{
				Error404FilePath: "404-lol.html",
				HeaderMap: map[string][]string{
					"X-Cool": {"Yes"},
				},
				HistoryMode: true,
				RedirectRoutes: map[string]string{
					"/aaa": "http://example.com",
				},
				TemplateMap: map[string]string{
					"AAA": "BBB",
				},
			},
		},
		{
			name:    "no dotfile",
			wantErr: false,
		},
		{
			name:    "fail loading dotfile",
			wantErr: true,
			files: map[string]string{
				".ghs.yaml": `%&*#exam???ple.com`,
			},
			wantCfg: &DotfileConfig{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dir, err := os.MkdirTemp("", "load-dotfile")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(dir)
			for f, c := range tt.files {
				if err := os.WriteFile(path.Join(dir, f), []byte(c), 0644); err != nil {
					t.Fatalf("failed to write file: %v", err)
				}
			}
			gotCfg, err := LoadDotfileConfig(dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadDotfileConfig(%v) error = %v, wantErr %v", tt.args.serveFolder, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCfg, tt.wantCfg) {
				t.Errorf("LoadDotfileConfig(%v) = %+v, want %+v", tt.args.serveFolder, gotCfg, tt.wantCfg)
			}
		})
	}
}

func TestLoadRedirectRoutesConfig(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name       string
		args       args
		files      map[string]string
		wantOutput map[string]string
		wantErr    bool
	}{
		{
			name: "basic",
			args: args{
				path: "redirects.yaml",
			},
			files: map[string]string{
				"redirects.yaml": `---
/aaa: /bbb`,
			},
			wantOutput: map[string]string{
				"/aaa": "/bbb",
			},
		},
		{
			name: "bad config",
			args: args{
				path: "redirects.yaml",
			},
			files: map[string]string{
				"redirects.yaml": `@%&*40<<<>>>3`,
			},
			wantOutput: map[string]string{},
			wantErr:    true,
		},
		{
			name: "no config",
			args: args{
				path: "redirects.yaml",
			},
			wantOutput: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dir, err := os.MkdirTemp("", "load-redirect-routes-config")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(dir)
			for f, c := range tt.files {
				if err := os.WriteFile(path.Join(dir, f), []byte(c), 0644); err != nil {
					t.Fatalf("failed to write file: %v", err)
				}
			}
			gotOutput, err := LoadRedirectRoutesConfig(path.Join(dir, tt.args.path))
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadRedirectRoutesConfig(%v) error = %v, wantErr %v", tt.args.path, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("LoadRedirectRoutesConfig(%v) = %v, want %v", tt.args.path, gotOutput, tt.wantOutput)
			}
		})
	}
}
