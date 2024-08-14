package main

// @title cribbage server
// @version 0.0.2
// @description cribbage rest server

// @host localhost
// @securityDefinitions.basic BasicAuth
// @BasePath /v1
func main() {
	r := &Router{}
	r.New()
}
