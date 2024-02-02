package geoprovider

import (
	"context"
	"encoding/json"
	"errors"

	"Buran.com/4Geoservice_1/service"
	"github.com/ekomobile/dadata/v2"
	"github.com/ekomobile/dadata/v2/api/suggest"
	"github.com/ekomobile/dadata/v2/client"
)

type GeoProvider interface {
	AddressSearch(input string) (*Address, error)
	GeoCode(lat, lng string) (*Address, error)
	AddressSearchRPC(args string, reply *Address) error
	GeoCodeRPC(args []string, reply *Address) error
}

type Provider struct {
	Responder service.Responder
}

type Address struct {
	Data interface{}
}

func (p *Provider) AddressSearchRPC(args string, reply *Address) error {
	resp, err := p.AddressSearch(args)
	reply.Data = resp.Data
	if err != nil {
		return err
	}
	return nil
}

func (p *Provider) AddressSearch(input string) (*Address, error) {
	if p.Responder != nil {
		ok, resp, err := p.Responder.CheckDataExists(input, true)
		if err != nil {
			return &Address{}, err
		}
		if ok {
			return &Address{Data: resp.Data}, nil
		}
	}

	creds := client.Credentials{
		ApiKeyValue:    "2a7817aac732e00f81ebf1066b13560c4080a35d",
		SecretKeyValue: "b143f6b27785fcd56fe6f85b1961b471d40b6369",
	}
	api := dadata.NewSuggestApi(client.WithCredentialProvider(&creds))

	params := suggest.RequestParams{
		Query: input,
	}

	result, err := api.Address(context.Background(), &params)
	if err != nil {
		return &Address{}, errors.New("400: Wrong request format")
	}
	resultToAdd, err := json.Marshal(result)
	if err != nil {
		return &Address{}, errors.New("cannot marshal data to json")
	}
	
	if p.Responder != nil {
		err = p.Responder.AddData(input, "/api/address/search", resultToAdd)
	}

	if err != nil {
		return &Address{}, err
	}
	return &Address{Data: result}, nil
}

func (p *Provider) GeoCodeRPC(args []string, reply *Address) error {
	resp, err := p.GeoCode(args[0], args[1])
	reply.Data = resp.Data
	if err != nil {
		return err
	}
	return nil
}

func (p *Provider) GeoCode(lat, lng string) (*Address, error) {
	if p.Responder != nil {
		ok, resp, err := p.Responder.CheckDataExists(lat + " " + lng, true)
		if err != nil {
			return &Address{}, err
		}
		if ok {
			return &Address{Data: resp.Data}, nil
		}
	}

	creds := client.Credentials{
		ApiKeyValue:    "2a7817aac732e00f81ebf1066b13560c4080a35d",
		SecretKeyValue: "b143f6b27785fcd56fe6f85b1961b471d40b6369",
	}

	api := dadata.NewSuggestApi(client.WithCredentialProvider(&creds))

	params := suggest.RequestParams{
		Query: lat + " " + lng,
	}

	result, err := api.Address(context.Background(), &params)
	if err != nil {
		return &Address{}, errors.New("400: Wrong request format")
	}
	resultToAdd, err := json.Marshal(result)
	if err != nil {
		return &Address{}, errors.New("cannot marshal data to json")
	}

	if p.Responder != nil {
		err = p.Responder.AddData(lat + " " + lng, "/api/address/geocode", resultToAdd)
	}

	if err != nil {
		return &Address{}, err
	}
	return &Address{Data: result}, nil
}