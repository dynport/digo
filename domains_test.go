package digo

import (
	"errors"
	"os"
	"testing"
)

func getAccount() (*Account, error) {
	apiKey := os.Getenv("DIGITAL_OCEAN_API_KEY")
	clientId := os.Getenv("DIGITAL_OCEAN_CLIENT_ID")

	if apiKey == "" {
		return nil, errors.New("Need DIGITAL_OCEAN_API_KEY")
	}

	if clientId == "" {
		return nil, errors.New("Need DIGITAL_OCEAN_CLIENT_ID")
	}

	return &Account{
		ApiKey:   apiKey,
		ClientId: clientId,
	}, nil
}

func TestDomains(t *testing.T) {
	account, e := getAccount()

	if e != nil {
		t.Fatal(e.Error())
	}

	domain, _ := account.NewDomain()
	domain.IpAddress = "1.2.3.4"
	domain.Name = "newdomain-123456.com"

	if e := domain.Save(); e != nil {
		t.Fatal("Could not create domain: " + e.Error())
	}

	domain, e = account.Domain(domain.Id)

	if e != nil {
		t.Fatal("Could not find domain: " + e.Error())
	}

	list, e := account.Domains()

	if e != nil {
		t.Fatal("Could not retrieve list of domains: " + e.Error())
	}

	found := false

	for _, d := range list {
		if d.Id == domain.Id {
			found = true
		}
	}

	if found == false {
		t.Fatal("Could not find domain")
	}

	if e := domain.Destroy(); e != nil {
		t.Fatal("Could not destroy")
	}

	list, _ = account.Domains()

	found = false

	for _, d := range list {
		if d.Id == domain.Id {
			found = true
		}
	}

	if found == true {
		t.Fatal("Expected to not find example.com!")
	}
}

func TestDomainRecords(t *testing.T) {
	account, e := getAccount()

	if e != nil {
		t.Fatal(e.Error())
	}

	domain, _ := account.NewDomain()
	domain.IpAddress = "1.2.3.4"
	domain.Name = "newdomain-123456.com"
	domain.Save()

	defaultRecords, e := domain.Records()

	if e != nil {
		t.Fatal("Could not load record list: ", e.Error())
	}

	record, _ := domain.NewRecord()
	record.Name = "sample"
	record.Data = "1.2.3.4"
	record.RecordType = "A"

	if e = record.Save(); e != nil {
		t.Fatal("Could not create DNS record: " + e.Error())
	}

	record.Name = "sample2"

	if e := record.Save(); e != nil {
		t.Fatal("Could not save update to record: " + e.Error())
	}

	if record, e := domain.Record(record.Id); e != nil {
		t.Fatal("Could not load existing record: " + e.Error())
	} else if record.Name != "sample2" {
		t.Fatal("Record name not expected!")
	}

	if records, _ := domain.Records(); records == nil || len(records) != (len(defaultRecords)+1) {
		t.Fatal("Did not have expected number of records!")
	}

	if e = record.Destroy(); e != nil {
		t.Fatal("Unable to delete DNS record: " + e.Error())
	}

	if records, _ := domain.Records(); len(records) != len(defaultRecords) {
		t.Fatal("Did not have expected number of records!")
	}

	domain.Destroy()
}
