package main

import (
	"context"
	"log"
	"os"
	"time"

	"dagger.io/dagger"
)

func main() {
	const hostFilePath = "supersecretfile"
	cleanup := setDemoHostEnv(hostFilePath)
	defer cleanup()

	// Dagger logic
	ctx := context.Background()
	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		panic(err)
	}

	secretEnv := c.Host().EnvVariable("MY_SECRET_ID").Secret()
	secretFile := c.Host().Directory(".").File(hostFilePath).Secret()

	stdout, err := c.Container().
		From("alpine:latest").
		WithEnvVariable("CACHEBUSTER", time.Now().String()).
		WithSecretVariable("MY_SECRET_ID", secretEnv).
		WithMountedSecret("/mysecret", secretFile).
		WithExec([]string{"sh", "-c", `echo -e "super secret: $MY_SECRET_ID || super secret file: "; cat /mysecret; echo " || "; date;`}).
		Stdout(ctx)
	if err != nil {
		panic(err)
	}

	log.Println("Stdout Log:", stdout)
}

func setDemoHostEnv(hostFilePath string) (cleanup func() error) {
	// Setting a test host environment variable
	const envName = "MY_SECRET_ID"
	err := os.Setenv(envName, "my secret value")
	if err != nil {
		panic(err)
	}

	// Setting a test host file
	const hostFileContent = "super secret file content"
	err = os.WriteFile(hostFilePath, []byte(hostFileContent), 0o644)
	if err != nil {
		panic(err)
	}

	return func() error {
		return os.Remove(hostFilePath)
	}
}
