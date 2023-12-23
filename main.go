package main

import (
	"log"

	"github.com/joho/godotenv"

	common "gitlab.com/BobyMCbobs/go-http-server/pkg/common"
	ghs "gitlab.com/BobyMCbobs/go-http-server/pkg/httpserver"
)

func main() {
	log.Printf("%v (%v, %v, %v, %v)\n", common.AppName, common.AppBuildVersion, common.AppBuildHash, common.AppBuildMode, common.AppBuildDate)
	if common.AppBuildMode == "development" {
		_ = godotenv.Load(".env")
	}
	ws := ghs.NewWebServer()
	log.Printf("Configuration: %#v\n", ws)
	ws.Listen()
}
