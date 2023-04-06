package main

import (
	"context"
	"log"
	"os"

	"dagger.io/dagger"

	"github.com/dolanor/demos/secret/gcp"
	"github.com/dolanor/demos/secret/scaleway"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	secretID := "my-personal-secret-id"

	secretPlaintext, err := scwSecret(ctx)
	// secretPlaintext, err := gcpSecret(ctx)
	if err != nil {
		panic(err)
	}
	log.Println("plaintext:", secretPlaintext)

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		panic(err)
	}
	defer c.Close()

	secret := c.SetSecret(secretID, string(secretPlaintext))

	cont := c.Container().From("alpine:3.17.2")

	echo := cont.
		WithSecretVariable("SECRET", secret).
		// try to echo the secret
		WithExec([]string{"sh", "-c", "echo secret: $SECRET"})

	out, err := echo.Stdout(ctx)
	if err != nil {
		panic(err)
	}
	log.Println("echo:", out)

	cat := cont.
		WithMountedSecret("/my/secret/file", secret).
		WithExec([]string{"sh", "-c", "cat /my/secret/file > /tmp/cleartext-file"}).
		// try to cat the secret file
		WithExec([]string{"sh", "-c", "cat /my/secret/file"})
	out, err = cat.Stdout(ctx)
	if err != nil {
		panic(err)
	}
	log.Println("cat:", out)

	ok, err := cat.File("/tmp/cleartext-file").
		Export(ctx, "cleartext.file")

	if err != nil || !ok {
		panic(err)
	}
}

func scwSecret(ctx context.Context) (string, error) {
	scwSecretID := "8cfafcab-7309-4625-a9c9-175644c008ca"
	secretRevision := "2"

	scwConfig := scaleway.LoadConfig()

	// get the secret from Scaleway Secret Manager
	return scaleway.GetSecretPlaintext(ctx, scwConfig, scwSecretID, secretRevision)
}

func gcpSecret(ctx context.Context) (string, error) {
	// You need to have your GOOGLE_APPLICATION_CREDENTIALS env var set to point to your json credentials
	projectID := "test-dagger-io"
	gcpSecretID := "my-personal-secret-id"
	secretRevision := "2"

	// get the secret from Google Cloud Secret Manager
	return gcp.GetSecretPlaintext(ctx, projectID, gcpSecretID, secretRevision)
}
