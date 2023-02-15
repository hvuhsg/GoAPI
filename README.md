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
    func (request *goapi.Request) {
        timestamp := request.GetInt("timestamp")
        return map[string] string {
            "message": "pong", "timestamp": strconv.Iota(timestapm)
        }
    }
)


func main() {
    app.Run("0.0.0.0", 8080)
}
```