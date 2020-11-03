/*
	initialise the API
*/

package main

import (
	"log"

	gohttpserver "gitlab.com/safesurfer/go-http-server/cmd/gohttpserver"
	common "gitlab.com/safesurfer/go-http-server/pkg/common"
)

func main() {
	// initialise the app
	log.Printf("%v (%v, %v, %v, %v)\n", common.AppName, common.AppBuildVersion, common.AppBuildHash, common.AppBuildMode, common.AppBuildDate)
	gohttpserver.HandleWebserver()
}
