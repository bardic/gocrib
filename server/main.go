package main

import "github.com/bardic/gocrib/server/router"

// Cribbage API
//
//	@title						cribbage server
//	@version					0.0.4
//	@description				cribbage rest server
//
//	@host						localhost
//	@securityDefinitions.basic	BasicAuth
//	@host						localhost:1323
//	@BasePath					/v1
func main() {

	r := &router.Router{}
	r.New()
}
