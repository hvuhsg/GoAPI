package main

import (
	"net/http"

	"github.com/hvuhsg/goapi"
)

func main() {
	app := goapi.GoAPI("external handler", "1.0v")

	// Serve files from GoAPI app using the Include method that allow us to include external handlers in the app
	app.Include("/", http.FileServer(http.Dir(".")))

	app.Run("127.0.0.1", 8080)
}
