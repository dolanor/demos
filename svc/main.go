package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"sync"

	"dagger.io/dagger"
)

func panicerr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	ctx := context.Background()
	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	panicerr(err)

	svc := c.Container().From("golang:1.19").
		WithNewFile("./web.go", dagger.ContainerWithNewFileOpts{
			Contents: webgo,
		}).
		WithExposedPort(9999).
		WithExec([]string{"go", "run", "web.go"})

	curler := c.Container().From("alpine:3.14.6").
		WithServiceBinding("web", svc).
		WithExec([]string{"apk", "add", "curl"})

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		res, err := curler.
			WithExec([]string{"echo", "sleeping 3s"}).
			WithExec([]string{"sleep", "3s"}).
			WithExec([]string{"curl", "-s", "http://web:9999/A"}).
			Stdout(ctx)
		panicerr(err)
		fmt.Println(res)
		wg.Done()
	}()

	res, err := curler.
		WithExec([]string{"curl", "-s", "http://web:9999/C"}).
		Stdout(ctx)
	panicerr(err)
	fmt.Println(res)

	wg.Wait()

}

var webgo = `package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, r.URL.Path)
	})
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		panic(err)
	}
}`
