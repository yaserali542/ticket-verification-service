package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"github.com/yaserali542/ticket-verification-service/controllers"
	"github.com/yaserali542/ticket-verification-service/repository"
	"github.com/yaserali542/ticket-verification-service/services"
)

func main() {

	var db *gorm.DB
	var err error
	v := initializeViperConfig()
	if db, err = repository.InitializeDatabase(v); err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	if err = repository.MigrateVerifiedQRCodesTable(db); err != nil {
		log.Fatal(err)
		return
	}

	rep := repository.Repository{Db: db}
	service := services.VerficationService{Repository: &rep}
	controller := controllers.Controllers{Services: service}
	middleware := controllers.Middleware{Services: service}

	r := mux.NewRouter()
	r.Use(middleware.ValidateRequest)
	r.Use(middleware.ValidateJwtToken)
	r.HandleFunc("/verify-qr-code", controller.VerifyQRCode).Methods("POST")
	//r.HandleFunc("/events", controller.GetEvents).Methods("GET")
	//r.HandleFunc("/qrscan-verify", controller.GetQRCodeDetails).Methods("POST")
	//r.HandleFunc("/booking/{id}", controller.GetBookingFromId).Methods("GET")
	//r.HandleFunc("/booking", controller.GetUserBookings).Methods("GET")
	log.Fatal(http.ListenAndServe(":9000", r))

}
func initializeViperConfig() *viper.Viper {
	viper.SetConfigType("json")
	viper.SetConfigFile("./config/config.json")
	fmt.Printf("Using config: %s\n", viper.ConfigFileUsed())
	viper.ReadInConfig()
	return viper.GetViper()
}
