# GoAPI

```go
package main

import goapi "github.com/hvuhsg/goapi"

app: = goapi.GoAPI("hello world api")

var ping := app.Path("/ping")
ping.Methods(goapi.GET, goapi.POST)
ping.Description("ping pong")
ping.Parameter("age", goapi.VIsInt{}, goapi.VRange{Min: 0, Max: 120})
ping.Action(
    func (request *goapi.Request) {
        age := request.GetInt("age")
        return map[string] string {
            "message": "pong", "age": strconv.Iota(timestapm)
        }
    }
)


func main() {
    app.Run("0.0.0.0", 8080)
}
```