package utiles

import (
	"encoding/json"
	"net/http"
)

type WBody struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
}

func (p *WBody) ResBodyStatusOK(data interface{}, message string, success bool, w http.ResponseWriter) {

	p.Data = data
	p.Message = message
	p.Success = success
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(*p)
}
func (p *WBody) ResBodyInternalServerError(data interface{}, message string, success bool, w http.ResponseWriter) {

	p.Data = data
	p.Message = message
	p.Success = success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(*p)

}

func (p *WBody) ResBodyBadRequest(data interface{}, message string, success bool, w http.ResponseWriter) {

	p.Data = data
	p.Message = message
	p.Success = success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(*p)
}
