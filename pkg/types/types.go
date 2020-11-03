/*
	handle all types used by API
*/

package types

import (
	"net/http"
)

// Endpoints ...
// array for endpoint fields
type Endpoints []struct {
	EndpointPath string
	HandlerFunc  http.HandlerFunc
	HTTPMethods  []string
}
