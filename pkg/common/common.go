/*
	common function calls
*/

package common

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"sigs.k8s.io/yaml"
)

// AppBuild metadata
var (
	AppName                  = "go-http-server"
	AppBuildVersion          = "0.0.0"
	AppBuildHash             = "???"
	AppBuildDate             = "???"
	AppBuildMode             = "???"
	AppServeFolderConfigName = ".ghs.yaml"
)

// GetAppEnvFile ...
// location of an env file to load
func GetAppEnvFile() (output string) {
	return GetEnvOrDefault("APP_ENV_FILE", ".env")
}

// GetAppHealthPortEnabled ...
// enable the binding of a health port
func GetAppHealthPortEnabled() (output bool) {
	return GetEnvOrDefault("APP_HEALTH_PORT_ENABLED", "true") == "true"
}

// GetAppHealthPort ...
// the port to bind the health service to
func GetAppHealthPort() (output string) {
	return GetEnvOrDefault("APP_HEALTH_PORT", ":8081")
}

// GetAppPort ...
// the port to serve web traffic on
func GetAppPort() (output string) {
	return GetEnvOrDefault("APP_PORT", ":8080")
}

// GetAppHTTPSPort ...
// The port to serve HTTPS traffic on, if enabled.
func GetAppHTTPSPort() (output string) {
	return GetEnvOrDefault("APP_HTTPS_PORT", ":8443")
}

// GetAppHTTPSCrtPath ...
// The TLS cert to use for serving HTTPS
func GetAppHTTPSCrtPath() (output string) {
	return GetEnvOrDefault("APP_HTTPS_CRT_PATH", "")
}

// GetAppHTTPSKeyPath ...
// The TLS cert to use for serving HTTPS
func GetAppHTTPSKeyPath() (output string) {
	return GetEnvOrDefault("APP_HTTPS_KEY_PATH", "")
}

// GetAppEnableHTTPS ...
// Whether to enable serving HTTPS.
func GetAppEnableHTTPS() (output bool) {
	return GetEnvOrDefault("APP_ENABLE_HTTPS", "false") == "true"
}

// GetAppMetricsPort ...
// return the port which the app should serve metrics on
func GetAppMetricsPort() (output string) {
	return GetEnvOrDefault("APP_PORT_METRICS", ":2112")
}

// GetAppMetricsEnabled ...
// serve metrics endpoint
func GetAppMetricsEnabled() (output bool) {
	return GetEnvOrDefault("APP_METRICS_ENABLED", "true") == "true"
}

// GetAppRealIPHeader ...
// the header to use instead of r.RemoteAddr
func GetAppRealIPHeader() (output string) {
	return GetEnvOrDefault("APP_HTTP_REAL_IP_HEADER", "")
}

// GetServeFolder ...
// return the path of the folder to serve
func GetServeFolder() (output string) {
	pwd, _ := os.Getwd()
	return GetEnvOrDefault("APP_SERVE_FOLDER", GetEnvOrDefault("KO_DATA_PATH", pwd))
}

// GetTemplateMapPath ...
// return the path of the template map
func GetTemplateMapPath() (output string) {
	return GetEnvOrDefault("APP_TEMPLATE_MAP_PATH", "./map.yaml")
}

// GetVuejsHistoryMode ...
// return if to use Vuejs history mode
func GetVuejsHistoryMode() (output bool) {
	return GetEnvOrDefault("APP_VUEJS_HISTORY_MODE", "false") == "true"
}

// GetEnableGZIP ...
// Return whether we should handle GZIP.
func GetEnableGZIP() (enable bool) {
	return GetEnvOrDefault("APP_HANDLE_GZIP", "false") == "true"
}

// GetHeaderSetEnable ...
// return if headers should be templated
func GetHeaderSetEnable() (output bool) {
	return GetEnvOrDefault("APP_HEADER_SET_ENABLE", "false") == "true"
}

// GetHeaderMapPath ...
// return the path of the header map
func GetHeaderMapPath() (output string) {
	return GetEnvOrDefault("APP_HEADER_MAP_PATH", "./headers.yaml")
}

// Get404PageFileName ...
// return the name of the file to serve for 404 for standard directory serving
func Get404PageFileName() (output string) {
	return GetEnvOrDefault("APP_404_PAGE_FILE_NAME", "404.html")
}

// GetRedirectRoutesEnabled ...
// return if redirecting routes should be enabled
func GetRedirectRoutesEnabled() (output bool) {
	return GetEnvOrDefault("APP_REDIRECT_ROUTES_ENABLED", "true") == "true"
}

// GetRedirectRoutesPath ...
// return if redirecting routes should be enabled
func GetRedirectRoutesPath() (output string) {
	return GetEnvOrDefault("APP_REDIRECT_ROUTES_PATH", "./redirects.yaml")
}

// GetEnvOrDefault ...
// given an env var return it's value, else return a default
func GetEnvOrDefault(envName string, defaultValue string) (output string) {
	output = os.Getenv(envName)
	if output == "" {
		output = defaultValue
	}
	return output
}

// EvaluateEnvFromMap ...
// evaluates environment variables from map[string]string{}
func EvaluateEnvFromMap(input map[string]string) (output map[string]string) {
	output = map[string]string{}
	for index, value := range input {
		output[index] = os.ExpandEnv(value)
	}
	return output
}

// LoadMapConfig ...
// loads map config as YAML
func LoadMapConfig(path string) (output map[string]string, err error) {
	mapBytes, err := os.ReadFile(path)
	if err != nil {
		return output, fmt.Errorf("Failed to load map file: %v", err.Error())
	}
	err = yaml.Unmarshal(mapBytes, &output)
	if err != nil {
		return map[string]string{}, err
	}
	return output, nil
}

// EvaluateEnvFromHeaderMap ...
// evaluates environment variables from map[string][]string
func EvaluateEnvFromHeaderMap(input map[string][]string) (output map[string][]string) {
	output = map[string][]string{}
	for key, value := range input {
		for _, valueSub := range value {
			output[key] = append(output[key], os.ExpandEnv(valueSub))
		}
	}
	return output
}

// LoadHeaderMapConfig ...
// loads map config as YAML
func LoadHeaderMapConfig(path string) (output map[string][]string, err error) {
	mapBytes, err := os.ReadFile(path)
	if err != nil {
		return map[string][]string{}, fmt.Errorf("Failed to load map file: %v", err.Error())
	}
	err = yaml.Unmarshal(mapBytes, &output)
	if err != nil {
		return map[string][]string{}, err
	}
	return output, nil
}

// WriteHeadersToResponse ...
// overwrites the headers in the response
func WriteHeadersToResponse(w http.ResponseWriter, headerMap map[string][]string) http.ResponseWriter {
	for key, value := range headerMap {
		for _, valueSub := range value {
			w.Header().Set(key, valueSub)
		}
	}
	return w
}

// GetRequestIP ...
// returns r.RemoteAddr unless RealIPHeader is set
func GetRequestIP(r *http.Request) (requestIP string) {
	realIPHeader := GetAppRealIPHeader()
	headerValue := r.Header.Get(realIPHeader)
	if realIPHeader == "" || headerValue == "" {
		return r.RemoteAddr
	}
	return headerValue
}

// Logging ...
// a basic middleware for logging
func Logging(next http.Handler) http.Handler {
	// log all requests
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestIP := GetRequestIP(r)
		log.Printf("%v %v %v %v %v %v %#v", r.Method, r.URL, r.Proto, r.Response, requestIP, r.RemoteAddr, r.Header)
		next.ServeHTTP(w, r)
	})
}

// DotfileConfig ...
// dotfiles found in the web root
type DotfileConfig struct {
	HeaderMap        map[string][]string `json:"headerMap"`
	HistoryMode      bool                `json:"historyMode"`
	RedirectRoutes   map[string]string   `json:"redirectRoutes"`
	TemplateMap      map[string]string   `json:"templateMap"`
	Error404FilePath string              `json:"error404FilePath"`
}

// LoadDotfileConfig ...
// loads a .ghs.yaml in the serve folder
func LoadDotfileConfig(serveFolder string) (cfg DotfileConfig, err error) {
	configPath := path.Join(serveFolder, AppServeFolderConfigName)
	file, err := os.ReadFile(configPath)
	if err != nil {
		return DotfileConfig{}, err
	}
	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return DotfileConfig{}, err
	}
	return cfg, nil
}

// LoadRedirectRoutesConfig ...
// loads map config as YAML
func LoadRedirectRoutesConfig(path string) (output map[string]string, err error) {
	mapBytes, err := os.ReadFile(path)
	if err != nil {
		return map[string]string{}, fmt.Errorf("Failed to load map file: %v", err.Error())
	}
	err = yaml.Unmarshal(mapBytes, &output)
	if err != nil {
		return map[string]string{}, err
	}
	return output, nil
}
