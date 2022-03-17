package main

import (
	"log"

	common "gitlab.com/safesurfer/go-http-server/pkg/common"
	ghs "gitlab.com/safesurfer/go-http-server/pkg/httpserver"
)

func main() {
	// initialise the app
	log.Printf("%v (%v, %v, %v, %v)\n", common.AppName, common.AppBuildVersion, common.AppBuildHash, common.AppBuildMode, common.AppBuildDate)
	ws := ghs.NewWebServer()
	log.Printf("Configuration: %#v\n", ws)
	ws.Listen()
}
