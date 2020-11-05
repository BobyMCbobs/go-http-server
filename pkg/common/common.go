/*
	common function calls
*/

package common

import (
	"fmt"
	"net/http"
	"os"
	"log"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

// AppBuild metadata
const (
	AppName         = "go-http-server"
	AppBuildVersion = "0.0.0"
	AppBuildHash    = "???"
	AppBuildDate    = "???"
	AppBuildMode    = "???"
)

// GetAppEnvFile ...
// location of an env file to load
func GetAppEnvFile() (output string) {
	return GetEnvOrDefault("APP_ENV_FILE", ".env")
}

// GetAppHealthPortEnabled ...
// enable the binding of a health port
func GetAppHealthPortEnabled() (output string) {
	return GetEnvOrDefault("APP_HEALTH_PORT_ENABLED", "true")
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

// GetAppMetricsPort ...
// return the port which the app should serve metrics on
func GetAppMetricsPort() (output string) {
	return GetEnvOrDefault("APP_PORT_METRICS", ":2112")
}

// GetAppMetricsEnabled ...
// serve metrics endpoint
func GetAppMetricsEnabled() (output string) {
	return GetEnvOrDefault("APP_METRICS_ENABLED", "true")
}

// GetAppRealIPHeader ...
// the header to use instead of r.RemoteAddr
func GetAppRealIPHeader() (output string) {
	return GetEnvOrDefault("APP_HTTP_REAL_IP_HEADER", "")
}

// GetServeFolder ...
// return the path of the folder to serve
func GetServeFolder() (output string) {
	return GetEnvOrDefault("APP_SERVE_FOLDER", "./dist")
}

// GetTemplateMapPath ...
// return the path of the template map
func GetTemplateMapPath() (output string) {
	return GetEnvOrDefault("APP_TEMPLATE_MAP_PATH", "./map.yaml")
}

// GetVuejsHistoryMode ...
// return if to use Vuejs history mode
func GetVuejsHistoryMode() (output string) {
	return GetEnvOrDefault("APP_VUEJS_HISTORY_MODE", "false")
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
	mapBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return output, fmt.Errorf("Failed to load map file: %v", err.Error())
	}
	err = yaml.Unmarshal(mapBytes, &output)
	return output, err
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