package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/yaserali542/ticket-verification-service/models"
	"github.com/yaserali542/ticket-verification-service/services"
)

type Controllers struct {
	Services services.VerficationService
}

func (c *Controllers) VerifyQRCode(w http.ResponseWriter, r *http.Request) {

	var request models.VerifyQRCodeRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "incorrect data", http.StatusBadRequest)
		return
	}
	jwtToken := r.Header.Get("token")

	response, isNotFound, err := c.Services.VerifyQRCode(jwtToken, &request)

	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isNotFound {
		http.Error(w, "already verified/redeemed", http.StatusAlreadyReported)
	}

	w.Header().Add("Content-type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)

}
