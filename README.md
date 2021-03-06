# Shack
Shack is a simple web framework written in Go.

## Installation
```bash
go get -u github.com/ichxxx/shack
```

## Examples

### Quick start
```go
import (
	"github.com/ichxxx/shack"
	"github.com/ichxxx/shack/middleware"
	"github.com/ichxxx/shack/rest"
)


func main() {
    r := shack.NewRouter()
    r.GET("/example", func(ctx *shack.Context) {
        rest.Resp(ctx).OK()
    }).With(middleware.Recovery())

    shack.Run(":8080", r)
    // or
    // http.ListenAndServe(":8080", r)
}
```

### Parameters in path
```go
func main() {
    r := shack.NewRouter()
    r.GET("/example/:id/*path", func(ctx *shack.Context) {
        id := ctx.Param("id")
        path := ctx.Param("path")
        ctx.JSON(shack.M{"id": id, "path": path})
        // or
        // rest.Resp(ctx).Data("id", id, "path", path).OK()
    })

    shack.Run(":8080", r)
}
```

### Querystring parameters
```go
type query struct {
    Foo int    `json:"foo"`
    Bar string `json:"bar"`
}

func main() {
    r := shack.NewRouter()
    r.GET("/example", func(ctx *shack.Context) {
        foo := ctx.QueryFlow("foo").Int()
        bar := ctx.Query("bar", "defaultBar")
        
        query := &query{}
        ctx.RawQueryFlow().Bind(query, "json")
        
        rest.Resp(ctx).Data(
            "foo", foo,
            "bar", bar,
            "query", query,
        )).OK()
    })

    shack.Run(":8080", r)
}
```

### Json Body
```go
type query struct {
    Foo int    `json:"foo"`
    Bar string `json:"bar"`
}

func main() {
    r := shack.NewRouter()
    r.GET("/example", func(ctx *shack.Context) {		
        query := &query{}
        ctx.BodyFlow().BindJson(query)
        
        rest.Resp(ctx).Data(
            "query", query,
        ).OK()
    })

    shack.Run(":8080", r)
}
```

### Multipart/Urlencoded Form
```go
type forms struct {
    Code int               `json:"code"`
    Msg  string            `json:"msg"`
    Data map[string]string `json:"data"`
}

func main() {
    r := shack.NewRouter()
    r.POST("/example", func(ctx *shack.Context) {
        data := ctx.Form("data")
        
        forms := &forms{}
        ctx.FormsFlow().Bind(forms, "json")
        
        rest.Resp(ctx).Data(
            "data", data,
            "forms", forms,
        ).OK()
    })

    shack.Run(":8080", r)
}
```

### Router group and middleware
```go
func main() {
    r := shack.NewRouter()
    r.Use(forAll)
    r.GET("/example", exampleHandler).With(middleware.AccessLog())
    r.Group("/api", func(r *shack.Router) {
        r.Use(onlyForApi)
        r.Handle("/article", articleHandler)
        
        r.Add(func(r *shack.Router) {
            r.Handle("/user", userHandler, http.MethodGet, http.MethodPost)
        })
    })
    
    shack.Run(":8080", r)
}
```

### Mount router
```go
func main() {
    r := shack.NewRouter()
    r.Mount("/api", apiRouter())
    
    shack.Run(":8080", r)
}

func apiRouter() *shack.Router {
    return shack.NewRouter()
}
```


### Logger
```go
func main() {
    shack.Logger("example").
        Level(shack.ErrorLevel).
        Encoding("json").
        Output("./logs").
        Enable()
    
    shack.Log.Error("some error",
        "timestamp", time.Now().Unix(),
    )
}
```

### Config
Shack has a simple toml config parser built in.

You can use it as follows:

```toml
# test_config.toml
[test]
name = "shack"
port = 8080
foo_bar = ["1", "2"]
```

```go
var TestConfig = &struct{
    // To use automatic parsing,
    // you have to combine shack.BaseConfig
    // in a struct.
    shack.BaseConfig
    Name    string
    Port    string
    FooBar  []int  `config:"foo_bar"`

}

func init() {
    // The second args `test` is the section's name in the toml file.
    shack.Config.Add(TestConfig, "test")
}

func main() {
    shack.Config.File("test_config.toml").Load()

    // After that, shack will parse the config automatically.
    // You can just use it.
    fmt.Println(TestConfig.Name)
    fmt.Println(TestConfig.Port)
    fmt.Println(TestConfig.FooBar)
}
```