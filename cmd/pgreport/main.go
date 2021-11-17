package main

import (
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/natemarks/pgsummary/aws"
	"github.com/natemarks/pgsummary/pg"
	"github.com/natemarks/pgsummary/util"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"time"
)

func outputInstanceReport(hostname string, json []byte) {
	JSONString := string(json)
	fmt.Print(JSONString)
	t := time.Now()
	filename := fmt.Sprintf("%s-%s.json", hostname, t.Format("20060102150405"))
	_ = ioutil.WriteFile(filename, json, 0644)
}

func main() {

	logger := log.With().Str("test_key", "test_value").Logger()

	hostPtr := flag.String("host", "localhost", "Postgres instance FQDN")
	portPtr := flag.Int("port", 5432, "Postgres instance port")
	dbnamePtr := flag.String("dbname", "some_dbname", "Postgres database name")
	secretIdPtr := flag.String("secretId", "myenv/mysecret", "AWS Secret ID")
	secretUsernameKeyPtr := flag.String("secretUsernameKey", "some_username_key", "AWS Secret JSON doc username key")
	secretPasswordKeyPtr := flag.String("secretPasswordKey", "some_password_key", "AWS Secret JSON doc password key")

	flag.Parse()

	gci := aws.GetCredentialInput{
		SecretId:    *secretIdPtr,
		UsernameKey: *secretUsernameKeyPtr,
		PasswordKey: *secretPasswordKeyPtr,
	}
	gco := aws.GetCredentialsFromAWSSM(gci, &logger)

	cData := pg.ConnData{
		Host:     *hostPtr,
		Port:     *portPtr,
		Username: gco.Username,
		Password: gco.Password,
		DbName:   *dbnamePtr,
	}

	pg.ValidateCredentials(cData, &logger)

	ir, err := pg.GetInstanceReport(cData, &logger)
	util.CheckError(err, &logger)
	// convert object to JSON []bytes
	jb, err := json.MarshalIndent(ir, "", "  ")
	util.CheckError(err, &logger)
	outputInstanceReport(*hostPtr, jb)
}
