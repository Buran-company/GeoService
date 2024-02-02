package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"Buran.com/4Geoservice_1/repository"
	"github.com/ekomobile/dadata/v2/api/model"
	"github.com/ekomobile/dadata/v2/api/suggest"
	"github.com/ptflp/godecoder"

	"go.uber.org/zap"
)

//go:generate easytags $GOFILE
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type Responder interface {
	OutputJSON(w http.ResponseWriter, responseData interface{})

	ErrorUnauthorized(w http.ResponseWriter, err error)
	ErrorBadRequest(w http.ResponseWriter, err error)
	ErrorForbidden(w http.ResponseWriter, err error)
	ErrorInternal(w http.ResponseWriter, err error)

	CheckDataExists(address string, isAddress bool) (bool, Response, error)
	AddData(query string, request string, response []byte) error
}

type Respond struct {
	log *zap.Logger
	godecoder.Decoder
	DataBaseHandler repository.DataBaseHandler
}

func NewResponder(decoder godecoder.Decoder, logger *zap.Logger) Responder {
	r := &Respond{log: logger, Decoder: decoder, DataBaseHandler: repository.DataBaseHandler{}}
	err := r.DataBaseHandler.ConnectToDB()
	if err != nil {
		return nil
	}
	return r
}

func (r *Respond) AddData(query string, request string, response []byte) error {
	err := r.DataBaseHandler.Create(repository.Repo{
		Query:    query,
		Request:  request,
		Response: response,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *Respond) CheckDataExists(address string, isAddress bool) (bool, Response, error) {
	history, err := r.DataBaseHandler.List(-1, -1)
	if err != nil {
		return false, Response{}, err
	}
	for _, search := range history {
		if search.Query == address {
			var resp1 []*model.Address
			var resp2 *suggest.GeoIPResponse
			var err error 
			if isAddress {
				err = json.Unmarshal(search.Response, &resp1)
			} else {
				err = json.Unmarshal(search.Response, &resp2)
			}
			if err != nil {
				return false, Response{}, err
			}
			if isAddress {
				return true, Response{
					Success: true,
					Message: "",
					Data:    resp1,
				}, nil
			} else {
				return true, Response{
					Success: true,
					Message: "",
					Data:    resp2,
				}, nil
			}
		}
	}
	return false, Response{}, nil
}

func (r *Respond) OutputJSON(w http.ResponseWriter, responseData interface{}) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	if err := r.Encode(w, responseData); err != nil {
		r.log.Error("responder json encode error", zap.Error(err))
	}
}

func (r *Respond) ErrorBadRequest(w http.ResponseWriter, err error) {
	r.log.Info("http response bad request status code", zap.Error(err))
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	if err := r.Encode(w, Response{
		Success: false,
		Message: err.Error(),
		Data:    nil,
	}); err != nil {
		r.log.Info("response writer error on write", zap.Error(err))
	}
}

func (r *Respond) ErrorForbidden(w http.ResponseWriter, err error) {
	r.log.Warn("http resposne forbidden", zap.Error(err))
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusForbidden)
	if err := r.Encode(w, Response{
		Success: false,
		Message: err.Error(),
		Data:    nil,
	}); err != nil {
		r.log.Error("response writer error on write", zap.Error(err))
	}
}

func (r *Respond) ErrorUnauthorized(w http.ResponseWriter, err error) {
	r.log.Warn("http resposne Unauthorized", zap.Error(err))
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	if err := r.Encode(w, Response{
		Success: false,
		Message: err.Error(),
		Data:    nil,
	}); err != nil {
		r.log.Error("response writer error on write", zap.Error(err))
	}
}

func (r *Respond) ErrorInternal(w http.ResponseWriter, err error) {
	if errors.Is(err, context.Canceled) {
		return
	}
	r.log.Error("http response internal error", zap.Error(err))
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	if err := r.Encode(w, Response{
		Success: false,
		Message: err.Error(),
		Data:    nil,
	}); err != nil {
		r.log.Error("response writer error on write", zap.Error(err))
	}
}
