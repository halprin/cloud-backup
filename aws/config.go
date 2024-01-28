package aws

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"log"
	"path/filepath"
)

const credentialFile = "credentials"
const configFile = "config"

func GetConfig(awsCredentialConfigPath string, awsProfile string) (aws.Config, error) {

	cfgOptions := []func(*config.LoadOptions) error{
		config.WithSharedConfigProfile(awsProfile),
	}

	if awsCredentialConfigPath != "" {
		cfgOptions = append(cfgOptions, config.WithSharedCredentialsFiles([]string{filepath.Join(awsCredentialConfigPath, credentialFile)}))
		cfgOptions = append(cfgOptions, config.WithSharedConfigFiles([]string{filepath.Join(awsCredentialConfigPath, configFile)}))
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), cfgOptions...)

	if err != nil {
		log.Println("Initial AWS config failed")
		return aws.Config{}, err
	}

	return cfg, nil
}
