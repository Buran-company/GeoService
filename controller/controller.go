package controller

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	geoprovider "Buran.com/4Geoservice_1/provider"
	service "Buran.com/4Geoservice_1/service"
	"github.com/golang-jwt/jwt"
)

type Controllerer interface {
	VerifyJWT(endpointHandler func(writer http.ResponseWriter, request *http.Request)) http.HandlerFunc
	Register(writer http.ResponseWriter, r *http.Request)
	LogIn(w http.ResponseWriter, r *http.Request)
	AddressSearch(w http.ResponseWriter, r *http.Request)
	AddressGeocode(w http.ResponseWriter, r *http.Request)
}

type UserController struct {
	Responder    service.Responder
	Prov geoprovider.Provider
	AddressQuery string `json:"queryaddr"`
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

var (
	requestsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "myapp_requests_total",
		Help: "Total number of requests",
	})

	requestDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "myapp_request_duration_seconds",
		Help:    "Request duration in seconds",
		Buckets: prometheus.LinearBuckets(0.1, 0.1, 10),
	})
)

func NewUserController(responder service.Responder, queryaddr, lat, lng string) *UserController {
	prometheus.MustRegister(requestsTotal)
	prometheus.MustRegister(requestDuration)
	return &UserController{
		Responder:    responder,
		Prov: geoprovider.Provider{Responder: responder},
		AddressQuery: queryaddr,
		Lat: lat,
		Lng: lng,
	}
}

func (u *UserController) VerifyJWT(endpointHandler func(writer http.ResponseWriter, request *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		startTime := time.Now()
		requestsTotal.Inc()
		token, _ := logging()
		request.Header.Set("Authorization", token)
		var keyfunc jwt.Keyfunc = func(token *jwt.Token) (interface{}, error) {
			return []byte("SecretYouShouldHide"), nil
		}

		parsed, err := jwt.Parse(request.Header.Get("Authorization"), keyfunc)
		if err != nil {
			u.Responder.ErrorUnauthorized(writer, err)
			duration := time.Since(startTime).Seconds()
			requestDuration.Observe(duration)
			return
		}

		if !parsed.Valid {
			u.Responder.ErrorUnauthorized(writer, errors.New("token is not valid"))
			duration := time.Since(startTime).Seconds()
			requestDuration.Observe(duration)
			return
		}
		u.Responder.OutputJSON(writer, service.Response{
			Success: true,
			Message: "Token is valid.",
			Data:    nil,
		})
		endpointHandler(writer, request)
		duration := time.Since(startTime).Seconds()
		requestDuration.Observe(duration)
	})
}

func generateJWT() (string, error) {
	var sampleSecretKey = []byte("SecretYouShouldHide")

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(24 * time.Hour).Unix()
	claims["authorized"] = true
	claims["user"] = "username"
	tokenString, err := token.SignedString(sampleSecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Register godoc
// @Summary Registering new User!
// @Router /api/register [get]
func (u *UserController) Register(writer http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	requestsTotal.Inc()
	token, err := generateJWT()
	if err != nil {
		u.Responder.ErrorInternal(writer, err)
		duration := time.Since(startTime).Seconds()
		requestDuration.Observe(duration)
	} else {
		u.Responder.OutputJSON(writer, service.Response{
			Success: true,
			Message: "200 OK",
			Data:    nil,
		})
		os.Truncate("/home/hexedchild1/Kata/Repository/go-kata/course4Geoservice_1/keys.txt", 0)
		os.WriteFile("/home/hexedchild1/Kata/Repository/go-kata/course4Geoservice_1/keys.txt", []byte(token), 0644)
		duration := time.Since(startTime).Seconds()
		requestDuration.Observe(duration)
	}
}

func logging() (string, error) {
	data, err := os.ReadFile("/home/hexedchild1/Kata/Repository/go-kata/course4Geoservice_1/keys.txt")
	if err != nil {
		log.Fatal(err.Error())
	}
	return string(data), nil
}

// LogIn godoc
// @Summary Just Logging In
// @Router /api/login [get]
func (u *UserController) LogIn(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	requestsTotal.Inc()
	token, _ := logging()
	r.Header.Set("Authorization", token)
	u.Responder.OutputJSON(w, service.Response{
		Success: true,
		Message: "You've logged in successfully!",
		Data:    nil,
	})
	duration := time.Since(startTime).Seconds()
	requestDuration.Observe(duration)
}

// AddressSearch godoc
// @Summary Retrieves all possible Info based on given Adress
// @Produce json
// @Router /api/address/search [get]
func (u *UserController) AddressSearch(w http.ResponseWriter, r *http.Request) {
	result, err := u.Prov.AddressSearch(u.AddressQuery)
	if err != nil {
		u.Responder.ErrorInternal(w, err)
		return
	}
	u.Responder.OutputJSON(w, service.Response{
		Success: true,
		Message: "",
		Data:    result.Data,
	})
}

// AddressGeocode godoc
// @Summary Retrieves all possible Info based on given IP
// @Produce json
// @Router /api/address/geocode [get]
func (u *UserController) AddressGeocode(w http.ResponseWriter, r *http.Request) {
	result, err := u.Prov.GeoCode(u.Lat, u.Lng)
	if err != nil {
		u.Responder.ErrorInternal(w, err)
		return
	}
	u.Responder.OutputJSON(w, service.Response{
		Success: true,
		Message: "",
		Data:    result.Data,
	})
}
