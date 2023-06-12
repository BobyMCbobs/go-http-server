package main

import (
	"log"

	"github.com/joho/godotenv"

	common "gitlab.com/BobyMCbobs/go-http-server/pkg/common"
	ghs "gitlab.com/BobyMCbobs/go-http-server/pkg/httpserver"
)

func main() {
	// initialise the app
	log.Printf("%v (%v, %v, %v, %v)\n", common.AppName, common.AppBuildVersion, common.AppBuildHash, common.AppBuildMode, common.AppBuildDate)
	ws := ghs.NewWebServer()
	_ = godotenv.Load(ws.EnvFile)
	log.Printf("Configuration: %#v\n", ws)
	ws.Listen()
}
