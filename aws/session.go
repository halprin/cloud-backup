package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
	"path/filepath"
)

const credentialFile = "credentials"
const configFile = "config"

func GetSession(awsCredentialConfigPath string, awsProfile string) (*session.Session, error) {
	options := session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           awsProfile,
	}

	if awsCredentialConfigPath != "" {
		options.Config = aws.Config{
			Credentials: credentials.NewSharedCredentials(filepath.Join(awsCredentialConfigPath, credentialFile), awsProfile),
		}
		options.SharedConfigFiles = []string{filepath.Join(awsCredentialConfigPath, configFile)}
	}

	awsSession, err := session.NewSessionWithOptions(options)

	if err != nil {
		log.Println("Initial AWS session failed")
		return nil, err
	}

	return awsSession, nil
}
