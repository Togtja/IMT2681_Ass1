package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

/*
Country represents extint species by country
*/
type Country struct {
	Code        string   `json:"code"`
	CountryName string   `json:"countryname"`
	CountryFlag string   `json:"countryflag"`
	Species     []string `json:"species"`
	SpeciesKey  []int    `json:"speciesKey"`
}

//CountryJSON represents a temporay struct to get data from restcountry
type CountryJSON struct {
	Code        string `json:"alpha2Code"`
	CountryName string `json:"name"`
	CountryFlag string `json:"flag"`
}

/*
Diag represents some cool comments
	{
		"gbif": <value>, e.g "200"
		"restcountries": <value>, e.g "200"
		"version": <value>, e.g "v1"
		"uptime": <value>, e.g "50.5"
	}
*/
type Diag struct {
	Gbif          int     `json:"gbif"`
	Restcountries int     `json:"restcountries"`
	Version       string  `json:"version"`
	Uptime        float64 `json:"uptime"`
}

func nilHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Default Handler: Invalid request received.")
	http.Error(w, "Invalid request", http.StatusBadRequest)
}

func countryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.Header.Add(w.Header(), "content-type", "application/json")
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) != 5 {
			http.Error(w, "Expecting format .../country_ios2", http.StatusBadRequest)
			return
		}
		if parts[4] == "" {
			//Get all
		} else if len(parts[4]) == 2 {
			query := "https://restcountries.eu/rest/v2/alpha/" + parts[4] + "?fields=name;flag;alpha2Code"
			fmt.Println(query)
			resp, err := http.Get(query)
			if err != nil {
				//We fucked up
			}
			data, _ := ioutil.ReadAll(resp.Body)
			var countryJSON CountryJSON
			json.Unmarshal(data, &countryJSON)
			json.NewEncoder(w).Encode(countryJSON)

		} else {
			fmt.Print(parts[4])
			http.Error(w, "Expecting country ios2 code", http.StatusBadRequest)
			return
		}
		return
	}
	http.Error(w, "only get method allowed", http.StatusNotImplemented)
	return
}
func diagHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.Header.Add(w.Header(), "content-type", "application/json")
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) != 5 {
			fmt.Println("Length is: ", len(parts))
			http.Error(w, "Expecting format .../", http.StatusBadRequest)
			return
		}
		respneGbif, err := http.Get("http://api.gbif.org/v1/")
		if err != nil {
			http.Error(w, "The HTTP request failed with error", http.StatusInternalServerError)
			fmt.Printf("The HTTP request failed with error %s\n", err)
			return
		}

		respneCon, err := http.Get("https://restcountries.eu/rest/v2/")
		if err != nil {
			http.Error(w, "The HTTP request failed with error", http.StatusInternalServerError)
			fmt.Printf("The HTTP request failed with error %s\n", err)
			return
		}
		//TODO get correct uptime
		diag := Diag{respneGbif.StatusCode, respneCon.StatusCode, "v1", 1.5}
		json.NewEncoder(w).Encode(diag)
		return
	}
	http.Error(w, "only get method allowed", http.StatusNotImplemented)
	return
}
func main() {
	fmt.Println("Starting application:")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", nilHandler)
	http.HandleFunc("/conservation/v1/country/", countryHandler)
	//http.HandleFunc("/conservation/v1/species/", helloHandler)
	http.HandleFunc("/conservation/v1/diag/", diagHandler)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
