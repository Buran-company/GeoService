package main

import (
	"encoding/gob"
	"encoding/json"
	"log"
	"net/rpc"

	geoprovider "Buran.com/4Geoservice_1/provider"
	"github.com/ekomobile/dadata/v2/api/suggest"
)

func main() {
	gob.Register([]*suggest.AddressSuggestion{})
	gob.Register(geoprovider.Address{})
	gob.Register(new(interface{}))
	client, err := rpc.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Ошибка при подключении к серверу:", err)
	}
	
	address := "odintsovo marshala biryuzova 14"
	var reply geoprovider.Address
	err = client.Call("Provider.AddressSearchRPC", address, &reply)
	if err != nil {
		log.Fatal("Ошибка при вызове удаленного метода:", err)
	}
	result, err := json.Marshal(reply.Data)
	if err != nil {
		log.Println(err)
	}
	log.Println("Результат:", string(result))

	args := []string{"55.803", "37.409"}
	err = client.Call("Provider.GeoCodeRPC", args, &reply)
	if err != nil {
		log.Fatal("Ошибка при вызове удаленного метода:", err)
	}
	result, err = json.Marshal(reply.Data)
	if err != nil {
		log.Println(err)
	}

	log.Println("Результат:", string(result))
}