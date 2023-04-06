package scaleway

import (
	"context"
	"fmt"

	secret "github.com/scaleway/scaleway-sdk-go/api/secret/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func GetSecretPlaintext(ctx context.Context, config *Config, secretID, secretRevision string) (string, error) {
	scwClient, err := scw.NewClient(
		scw.WithProfile(config),
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

type Config = scw.Profile

func LoadConfig() *Config {
	return scw.LoadEnvProfile()
}
