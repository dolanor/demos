package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"dagger.io/dagger"
	"github.com/dolanor/demos/secret/scaleway"
	secret "github.com/scaleway/scaleway-sdk-go/api/secret/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func main() {
	ctx := context.Background()

	secretID := "8cfafcab-7309-4625-a9c9-175644c008ca"
	secretRevision := "1"

	scwProfile := scaleway.LoadConfig()

	// get the secret from Scaleway Secret Manager
	secretPlaintext, err := scaleway.GetSecretPlaintext(ctx, scwProfile, secretID, secretRevision)
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

func scwGetSecretPlaintext(ctx context.Context, scwProfile *scw.Profile, secretID, secretRevision string) (string, error) {
	scwClient, err := scw.NewClient(
		scw.WithProfile(scwProfile),
	)
	if err != nil {
		return "", fmt.Errorf("scw.NewClient: %w", err)
	}
	secReq := &secret.AccessSecretVersionRequest{
		Region:   scw.RegionFrPar,
		SecretID: secretID,
		Revision: "1",
	}

	secAPI := secret.NewAPI(scwClient)
	if err != nil {
		return "", fmt.Errorf("secret.NewAPI: %w", err)
	}

	res, err := secAPI.AccessSecretVersion(secReq, scw.WithContext(ctx))
	if err != nil {
		return "", fmt.Errorf("AccessSecretVersion: %w", err)
	}

	secretPlaintext := res.Data

	return string(secretPlaintext), nil
}
