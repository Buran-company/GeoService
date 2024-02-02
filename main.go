package main

import (
	"encoding/gob"
	"log"
	"net"
	"net/http"
	"net/rpc"

	"Buran.com/4Geoservice_1/controller"
	geoprovider "Buran.com/4Geoservice_1/provider"
	"Buran.com/4Geoservice_1/service"
	"github.com/ekomobile/dadata/v2/api/suggest"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ptflp/godecoder"
	swaggerFiles "github.com/swaggo/files"
	"go.uber.org/zap"
)

// @title My First API
// @version 0.0.0.0.0.0.0.1
// @description 'Cause today... is a day... of forgiveness! And to you... We can only say... A simple answer... Oh, hey, hey, hey!
func main() {
	gob.Register([]*suggest.AddressSuggestion{})
	gob.Register(geoprovider.Address{})
	gob.Register(new(interface{}))
	uController := controller.NewUserController(service.NewResponder(godecoder.NewDecoder(), &zap.Logger{}), "odintsovo makovskogo 2", "55.878", "37.653")

	http.HandleFunc("/api/register", uController.Register)
	http.HandleFunc("/api/login", uController.LogIn)
	http.HandleFunc("/api/address/geocode", uController.VerifyJWT(uController.AddressGeocode))
	http.HandleFunc("/api/address/search", uController.VerifyJWT(uController.AddressSearch))

	http.Handle("/swagger/*any", swaggerFiles.Handler)
	http.Handle("/metrics", promhttp.Handler())

	prov := new(geoprovider.Provider)
	rpc.Register(prov)

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("Ошибка при запуске сервера:", err)
	}

	log.Println("Сервер запущен на порту 8080")

	rpc.Accept(l)
}
