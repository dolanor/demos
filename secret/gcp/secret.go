package gcp

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

func GetSecretPlaintext(ctx context.Context, projectID, secretID, secretRevision string) (string, error) {
	secretVersion := fmt.Sprintf("projects/%s/secrets/%s/versions/%s", projectID, secretID, secretRevision)

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
