# GoAPI

```go
package main

import goapi "github.com/hvuhsg/goapi"

app: = goapi.GoApp(title = "hello world api")

app.Path("/ping").Methods([]int{GET}).Description("ping pong").Parameter("timestamp").Action(
        func ping(request) {
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