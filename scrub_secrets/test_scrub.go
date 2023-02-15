package main

import (
	"context"
	"log"
	"os"
	"time"

	"dagger.io/dagger"
)

func main() {
	err := os.Setenv("MY_SECRET_ID", "my secret value")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))

	secret := c.Host().EnvVariable("MY_SECRET_ID").Secret()

	stdout, err := c.Container().
		From("alpine:latest").
		WithEnvVariable("CACHEBUSTER", time.Now().String()).
		WithSecretVariable("MY_SECRET_ID", secret).
		WithExec([]string{"sh", "-c", "echo super secret: $MY_SECRET_ID; date;"}).
		Stdout(ctx)

	log.Println("LOGSTDOUT:", stdout)
}
