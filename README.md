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
		ctx.JSON(rest.R().OK())
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
		ctx.JSON(rest.R().OK().Data(
			"id", id,
			"path", path,
		))
	})

	shack.Run(":8080", r)
}
```

### Querystring parameters
```go
func main() {
	r := shack.NewRouter()
	r.GET("/example", func(ctx *shack.Context) {
		foo := ctx.Query("foo").Int()
		bar := ctx.Query("bar", "defaultBar").Value()
		ctx.String(strconv.Itoa(foo), bar)
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
		data := ctx.Form("data").Value()
		
		forms := &forms{}
		ctx.Forms().Bind(forms, "json")
		
		ctx.JSON(rest.R().Data(
			"data", data,
			"forms", forms,
		))
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
		r.GET("/login", loginHandler)
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
		Output("/logs").
		Enable()
	
	shack.Log.Error("some error",
		"timestamp", time.Now().Unix(),
	)
}
```