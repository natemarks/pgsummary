package aws

import (
	"github.com/natemarks/easyaws/secrets"
	"github.com/natemarks/pgsummary/util"
	"github.com/rs/zerolog"
)

type GetCredentialInput struct {
	SecretId, UsernameKey, PasswordKey string
}

type GetCredentialOutput struct {
	Username, Password string
}

func GetCredentialsFromAWSSM(input GetCredentialInput, log *zerolog.Logger) GetCredentialOutput {

	docInput := secrets.GetSecretJSONInput{AWSSMSecretID: input.SecretId}
	// Get the AWS Secrets JSON Doc that contains the database credentials
	jsonDoc, err := secrets.GetSecretJSON(docInput, log)
	util.CheckError(err, log)

	// Get the username from the JSON Doc
	username, err := secrets.LookupJSONKey(input.UsernameKey, jsonDoc)
	util.CheckError(err, log)

	// Get the password from the JSON Doc
	password, err := secrets.LookupJSONKey(input.PasswordKey, jsonDoc)
	util.CheckError(err, log)
	res := GetCredentialOutput{
		Username: username,
		Password: password,
	}
	return res
}
