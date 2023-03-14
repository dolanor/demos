package main

import (
	"context"
	"os"

	"dagger.io/dagger"
)

func main() {
	ctx := context.Background()
	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		panic(err)
	}
	defer c.Close()

	src := c.Host().Directory(".")

	goBuilder := c.Container().
		From("golang:1.20.2-alpine3.17").
		WithDirectory("/app", src).
		WithWorkdir("/app").
		WithExec([]string{"go", "build", "."})

	imageName := "dolanor/go-snyk:0.1.0"
	_, err = goBuilder.Publish(ctx, imageName)
	if err != nil {
		panic(err)
	}

	snykToken := c.SetSecret("snyk-token", os.Getenv("SNYK_TOKEN"))

	log, err := c.Container().
		From("snyk/snyk:docker").
		WithSecretVariable("SNYK_TOKEN", snykToken).
		WithExec([]string{"snyk", "container", "test", "--docker", imageName}).
		Stderr(ctx)

	if err != nil {
		panic(err)
	}

	println(log)
}
