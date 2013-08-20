package main

import (
	"fmt"
	"github.com/dynport/gologger"
	"github.com/dynport/gocli"
	"os"
	"strconv"
	"strings"
)

var (
	cli    = &gocli.Router{}
	logger = gologger.New()
	account *Account
)

func CurrentAccount() *Account {
	if account == nil {
		var e error
		account, e = AccountFromEnv()
		if e != nil {
			ExitWith("unable to load account from env: " + e.Error())
		}
		logger.Debugf("using account %+v", account)
	}
	return account
}

func ExitWith(err interface{}) {
	logger.Error(err)
	os.Exit(1)
}

func init() {
	logger.Start()
	if os.Getenv("DEBUG") == "true" {
		logger.LogLevel = gologger.DEBUG
	}
}

const (
	DIGITAL_OCEAN_CLIENT_ID         = "DIGITAL_OCEAN_CLIENT_ID"
	DIGITAL_OCEAN_API_KEY           = "DIGITAL_OCEAN_API_KEY"
	DIGITAL_OCEAN_DEFAULT_REGION_ID = "DIGITAL_OCEAN_DEFAULT_REGION_ID"
	DIGITAL_OCEAN_DEFAULT_SIZE_ID   = "DIGITAL_OCEAN_DEFAULT_SIZE_ID"
	DIGITAL_OCEAN_DEFAULT_IMAGE_ID  = "DIGITAL_OCEAN_DEFAULT_IMAGE_ID"
	DIGITAL_OCEAN_DEFULT_SSH_KEY    = "DIGITAL_OCEAN_DEFULT_SSH_KEY"
)

func AccountFromEnv() (*Account, error) {
	account := &Account{}
	account.ClientId = os.Getenv(DIGITAL_OCEAN_CLIENT_ID)
	account.ApiKey = os.Getenv(DIGITAL_OCEAN_API_KEY)
	account.RegionId, _ = strconv.Atoi(os.Getenv(DIGITAL_OCEAN_DEFAULT_REGION_ID))
	account.SizeId, _ = strconv.Atoi(os.Getenv(DIGITAL_OCEAN_DEFAULT_SIZE_ID))
	account.ImageId, _ = strconv.Atoi(os.Getenv(DIGITAL_OCEAN_DEFAULT_IMAGE_ID))
	account.SshKey, _ = strconv.Atoi(os.Getenv(DIGITAL_OCEAN_DEFULT_SSH_KEY))

	allErrors := []string{}

	if account.ClientId == "" {
		allErrors = append(allErrors, fmt.Sprintf("%s must be set in env", DIGITAL_OCEAN_CLIENT_ID))
	}
	if account.ApiKey == "" {
		allErrors = append(allErrors, fmt.Sprintf("%s must be set in env", DIGITAL_OCEAN_API_KEY))
	}
	if len(allErrors) > 0 {
		return nil, fmt.Errorf(strings.Join(allErrors, "\n"))
	}
	return account, nil
}

func init() {
	if os.Getenv("DEBUG") == "true" {
		logger.LogLevel = gologger.DEBUG
	}
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			ExitWith(r)
		}
	}()
	if e := cli.Handle(os.Args); e != nil {
		ExitWith(e.Error())
	}
}
