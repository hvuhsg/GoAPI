package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/hvuhsg/goapi"
)

var host = flag.String("host", "127.0.0.1", "server host (127.0.0.1 / 0.0.0.0)")
var port = flag.Int("port", 8080, "server port")
var useNgrok = flag.Bool("ngrok", false, "use ngrok tunnel")
var ngrokToken = flag.String("ngroktoken", "", "ngrok token")

func main() {
	flag.Parse()

	currentDir, err := os.Getwd()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	api := goapi.GoAPI("current path serve", "1.0v")
	api.Include("/", http.FileServer(http.Dir(currentDir)))

	if *useNgrok {
		api.RunNgrok(*ngrokToken)
	} else {
		api.Run(*host, *port)
	}
}
