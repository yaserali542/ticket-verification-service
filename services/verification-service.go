package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"github.com/yaserali542/ticket-verification-service/encryption"
	"github.com/yaserali542/ticket-verification-service/models"
	"github.com/yaserali542/ticket-verification-service/repository"
)

type VerficationService struct {
	Repository *repository.Repository
}

func (*VerficationService) GetUserDetails(username, jwtToken string) (*models.BasicFields, error) {
	url := fmt.Sprintf("%v/basic-details", viper.GetViper().GetString("account_service_url"))

	requestBytes, _ := json.Marshal(models.BasicFieldsRequest{UserName: username})
	reqBody := bytes.NewReader(requestBytes)

	signature := encryption.CalculateHmac(requestBytes)

	req, _ := http.NewRequest(
		"POST",
		url,
		reqBody,
	)

	// add a request header
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("Signature", signature)
	req.Header.Add("token", jwtToken)

	// send an HTTP using `req` object
	res, err := http.DefaultClient.Do(req)
	// check for response error
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid return status from :%v", url)
	}
	var basicAccountDetails models.BasicFields

	// read response data
	if err := json.NewDecoder(res.Body).Decode(&basicAccountDetails); err != nil {
		return nil, err
	}

	// close response body
	res.Body.Close()

	return &basicAccountDetails, nil
}

func (s *VerficationService) VerifyQRCode(jwtToken string, request *models.VerifyQRCodeRequest) (*models.QRCodeDetails, bool, error) {

	url := fmt.Sprintf("%v/qrscan-verify", viper.GetViper().GetString("service_management_url"))

	_, notFound, err := s.Repository.GetVerificationQRData(request.DigitalSignature)

	if err != nil {
		return nil, true, err
	}
	if !notFound {
		return nil, false, nil
	}

	requestBytes, _ := json.Marshal(request)
	reqBody := bytes.NewReader(requestBytes)

	signature := encryption.CalculateHmac(requestBytes)

	req, _ := http.NewRequest(
		"POST",
		url,
		reqBody,
	)
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	req.Header.Add("Signature", signature)
	req.Header.Add("token", jwtToken)
	res, err := http.DefaultClient.Do(req)
	// check for response error
	if err != nil {
		return nil, true, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, true, fmt.Errorf("invalid return status from :%v", url)
	}
	var qrCodeDetails models.QRCodeDetails

	// read response data
	if err := json.NewDecoder(res.Body).Decode(&qrCodeDetails); err != nil {
		return nil, true, err
	}

	claims := &models.Claims{}
	jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		jwtKey := viper.GetViper().GetString("jwt.key")
		return []byte(jwtKey), nil
	})

	verifiedQRCode := &models.VerifiedQRCodes{
		BookingId:        qrCodeDetails.BookingDetails.ID,
		DigitalSignature: request.DigitalSignature,
		UserName:         qrCodeDetails.UserDetails.UserName,
		VerifierUserName: claims.Username,
	}

	if err := s.Repository.CreateVerificationQRData(verifiedQRCode); err != nil {
		return nil, true, err
	}
	return &qrCodeDetails, true, nil
}

func (s *VerficationService) ConfirmVerfication(jwtToken string, request *models.VerifyQRCodeRequest) (*models.VerifiedQRCodes, bool, error) {

	qrDetails, notFound, err := s.VerifyQRCode(jwtToken, request)
	if !notFound || err != nil {
		return nil, notFound, err
	}

	claims := &models.Claims{}
	jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		jwtKey := viper.GetViper().GetString("jwt.key")
		return []byte(jwtKey), nil
	})

	verifiedQRCode := &models.VerifiedQRCodes{
		BookingId:        qrDetails.BookingDetails.ID,
		DigitalSignature: request.DigitalSignature,
		UserName:         qrDetails.UserDetails.UserName,
		VerifierUserName: claims.Username,
	}

	if err := s.Repository.CreateVerificationQRData(verifiedQRCode); err != nil {
		return nil, true, err
	}
	return verifiedQRCode, true, nil

}
