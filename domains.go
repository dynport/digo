package digo

import (
	"errors"
	"fmt"
	"net/url"
)

type DomainResponse struct {
	Status string  `json:"status"`
	Domain *Domain `json:"domain"`
}

type DomainsResponse struct {
	Status  string    `json:"status"`
	Domains []*Domain `json:"domains"`
}

type DomainRecordResponse struct {
	Status       string        `json:"status"`
	DomainRecord *DomainRecord `json:"record"`
}

type DomainRecordsResponse struct {
	Status        string          `json:"status"`
	DomainRecords []*DomainRecord `json:"records"`
}

type Domain struct {
	Id                int    `json:"id"`
	Name              string `json:"name"`
	TTL               int    `json:"ttl"`
	LiveZoneFile      string `json:"live_zone_file"`
	Error             string `json:"error"`
	ZoneFileWithError string `json:"zone_file_with_error"`
	IpAddress         string
	account           *Account
}

type DomainRecord struct {
	Id         int    `json:"id"`
	DomainId   int    `json:"domain_id"`
	RecordType string `json:"record_type"`
	Name       string `json:"name"`
	Data       string `json:"data"`
	Priority   int    `json:"priority,omitempty"`
	Port       int    `json:"port,omitempty"`
	Weight     int    `json:"weight,omitempty"`
	domain     *Domain
}

func (account *Account) Domains() ([]*Domain, error) {
	r := &DomainsResponse{}

	if e := account.loadResource("/domains", r); e != nil {
		return r.Domains, e
	}

	for _, d := range r.Domains {
		d.account = account
	}

	return r.Domains, nil
}

func (account *Account) Domain(id int) (*Domain, error) {
	r := &DomainResponse{}

	if e := account.loadResource(fmt.Sprintf("/domains/%d", id), r); e != nil {
		return nil, e
	}

	r.Domain.account = account

	return r.Domain, nil
}

func (account *Account) NewDomain() (*Domain, error) {
	return &Domain{
		account: account,
	}, nil
}

func (domain *Domain) Save() error {
	qs := &url.Values{}

	if domain.IpAddress == "" {
		return errors.New("IP Address must not be blank")
	}

	if domain.Name == "" {
		return errors.New("Name must not be blank")
	}

	qs.Set("name", domain.Name)
	qs.Set("ip_address", domain.IpAddress)

	dr := &DomainResponse{}

	e := domain.account.loadResource(fmt.Sprintf("/domains/new?%s", qs.Encode()), dr)

	if e == nil {
		domain.Id = dr.Domain.Id
	}

	return e
}

func (domain *Domain) Destroy() error {
	return domain.account.loadResource(fmt.Sprintf("/domains/%d/destroy", domain.Id), &EventResponse{})
}

func (domain *Domain) Records() ([]*DomainRecord, error) {
	r := &DomainRecordsResponse{}

	if e := domain.account.loadResource(fmt.Sprintf("/domains/%d/records", domain.Id), r); e != nil {
		return r.DomainRecords, e
	}

	for _, dr := range r.DomainRecords {
		dr.domain = domain
	}

	return r.DomainRecords, nil
}

func (domain *Domain) Record(id int) (*DomainRecord, error) {
	r := &DomainRecordResponse{}

	if e := domain.account.loadResource(fmt.Sprintf("/domains/%d/records/%d", domain.Id, id), r); e != nil {
		return nil, e
	}

	r.DomainRecord.domain = domain

	return r.DomainRecord, nil
}

func (domain *Domain) NewRecord() (*DomainRecord, error) {
	return &DomainRecord{
		domain: domain,
	}, nil
}

func (record *DomainRecord) Save() error {
	qs := &url.Values{}

	qs.Set("domain_id", fmt.Sprintf("%d", record.DomainId))
	qs.Set("record_type", record.RecordType)
	qs.Set("name", record.Name)
	qs.Set("data", record.Data)

	if record.Priority > 0 {
		qs.Set("priority", fmt.Sprintf("%d", record.Priority))
	}

	if record.Port > 0 {
		qs.Set("port", fmt.Sprintf("%d", record.Port))
	}

	if record.Weight > 0 {
		qs.Set("weight", fmt.Sprintf("%d", record.Weight))
	}

	drr := &DomainRecordResponse{}

	var e error

	if record.Id > 0 {
		e = record.domain.account.loadResource(fmt.Sprintf("/domains/%d/records/%d/edit?%s", record.domain.Id, record.Id, qs.Encode()), drr)

	} else {
		e = record.domain.account.loadResource(fmt.Sprintf("/domains/%d/records/new?%s", record.domain.Id, qs.Encode()), drr)

	}

	if e != nil {
		return e
	}

	record.Id = drr.DomainRecord.Id

	return nil
}

func (record *DomainRecord) Destroy() error {
	return record.domain.account.loadResource(fmt.Sprintf("/domains/%d/records/%d/destroy", record.domain.Id, record.Id), &EventResponse{})
}
