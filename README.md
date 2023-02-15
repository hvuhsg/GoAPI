# GoAPI

```go
package main

import goapi "github.com/hvuhsg/goapi"

app: = goapi.GoApp("hello world api")

var ping := app.Path("/ping")
ping.Methods(goapi.GET, goapi.POST)
ping.Description("ping pong")
ping.Parameter("timestamp")
ping.Action(
    func (request) {
        timestamp := request.get("timestamp")
        return map[string] string {
            "message": "pong"
        }
    }
)


func main() {
    app.Run("0.0.0.0", 8080)
}
```