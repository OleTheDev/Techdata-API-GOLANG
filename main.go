package main

import (
	"fmt"
	"log"
	"net/http"

	"TTApi/Config"
	"TTApi/Models/Techdata"
	"TTApi/Models/Techdata/Handlers"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

/*
	Setup the server
*/
func main() {
	//Define config
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error with config file: %s", err)
	}

	if Config.LiveMode {
		fmt.Print("Server Launched in Production Environment!")
	} else {
		fmt.Print("Server Launched in Development/Quality Environment!")
	}

	Config.ApplyConfig()

	//Create routes
	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/techdata/auth", Techdata.AuthLogin).Methods(http.MethodGet)
	api.HandleFunc("/techdata/products", getProductsTechData).Methods(http.MethodGet)
	api.HandleFunc("/techdata/orders", getOrdersTechData).Methods(http.MethodGet)
	api.HandleFunc("/techdata/customers", getCustomersTechData).Methods(http.MethodGet)

	//Launch on port in config.yml, default 8080
	log.Fatal(http.ListenAndServe(viper.GetString("server.port"), r))
}

/*
	Get products from techdata
*/
func getProductsTechData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	Config.ApplyConfig()

	results := Techdata.ApiRequest(Config.Token, Config.SOIN, Config.Signature, Config.Timestamp, "http://eu-uat-papi.tdmarketplace.net/catalog/products/1", "GET")

	w.Write([]byte(results))
}

/*
	Get order history
*/
func getOrdersTechData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	Config.ApplyConfig()

	results := Techdata.ApiRequest(Config.Token, Config.SOIN, Config.Signature, Config.Timestamp, "https://partnerapi.tdstreamone.eu/order/subscriptions/1", "GET")

	w.Write([]byte(results))
}

func getCustomersTechData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	Config.ApplyConfig()

	results := Techdata.ApiRequest(Config.Token, Config.SOIN, Config.Signature, Config.Timestamp, "https://partnerapi.tdstreamone.eu/endCustomer/details/1", "GET")

	TechCustomers.CustomerList([]byte(results))

	w.Write([]byte(results))
}
