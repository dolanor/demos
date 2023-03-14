package main

import (
	"context"
	"fmt"
	"io"
	"log"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"dagger.io/dagger"
)

func main() {
	ctx := context.Background()

	// You need to have your GOOGLE_APPLICATION_CREDENTIALS env var set to point to your json credentials
	projectID := "test-dagger-io"
	secretID := "my-personal-secret-id"

	// get the secret from Google Cloud Secret Manager
	secretPlaintext, err := gcpGetSecretPlaintext(ctx, projectID, secretID)
	if err != nil {
		panic(err)
	}
	log.Println("plaintext:", secretPlaintext)

	c, err := dagger.Connect(ctx, dagger.WithLogOutput(io.Discard))
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

func gcpGetSecretPlaintext(ctx context.Context, projectID, secretID string) (string, error) {
	version := 1
	secretVersion := fmt.Sprintf("projects/%s/secrets/%s/versions/%d", projectID, secretID, version)

	gcpClient, err := secretmanager.NewClient(ctx)
	if err != nil {
		panic(err)
	}
	defer gcpClient.Close()

	secReq := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretVersion,
	}

	res, err := gcpClient.AccessSecretVersion(ctx, secReq)
	if err != nil {
		panic(err)
	}

	secretPlaintext := res.Payload.Data

	return string(secretPlaintext), nil
}
