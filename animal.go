package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"sync"
)

var startTime time.Time
const VERSION string = "v1"
const MAXCALL int = 5
const GBIFOCCURANCESLIMIT int = 300
const GBIFSPECIESLIMIT int = 1000;

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

//gbifJSON get the spices result
type gbifJSON struct {
	Results []SpeciesJSON `json:"results"`
}

/*
Species represents extint species by country
*/
type Species struct {
	Key       int    `json:"key"`
	Kingdom   string `json:"kingdom"`
	Phylum    string `json:"phylum"`
	Order     string `json:"order"`
	Family    string `json:"family"`
	Species   string `json:"species"`
	SciName   string `json:"scientificName"`
	CanonName string `json:"canonicalName"`
	IsExtinct bool   `json:"extinct"`
}

//SpeciesJSON represents a temporay struct to get data from GBIF
type SpeciesJSON struct {
	Key       int    `json:"speciesKey"`
	Kingdom   string `json:"kingdom"`
	Phylum    string `json:"phylum"`
	Order     string `json:"order"`
	Family    string `json:"family"`
	Species   string `json:"species"`
	SciName   string `json:"scientificName"`
	CanonName string `json:"canonicalName"`
	IsExtinct bool   `json:"extinct"`
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
	Uptime        string `json:"uptime"`
}
func uptime() time.Duration {
    return time.Since(startTime)
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
func speciesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.Header.Add(w.Header(), "content-type", "application/json")
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) != 5 {
			http.Error(w, "Expecting format .../country_ios2", http.StatusBadRequest)
			return
		}
		if parts[4] == "" {
			//Get all
		} else if _, err := strconv.Atoi(parts[4]); err == nil {
			query := "http://api.gbif.org/v1/occurrence/search?speciesKey=" + parts[4]
			fmt.Println(query)
			resp, err := http.Get(query)
			if err != nil {
				//We fucked up
			}
			data, _ := ioutil.ReadAll(resp.Body)
			var speciesJSON gbifJSON
			json.Unmarshal(data, &speciesJSON)
			json.NewEncoder(w).Encode(speciesJSON)

		} else {
			fmt.Print(parts[4])
			http.Error(w, "Expecting format .../speices_key", http.StatusBadRequest)
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
		
		uptime := uptime()
		uptimeString := fmt.Sprintf("%.0f seconds",uptime.Seconds()) 
		diag := Diag{respneGbif.StatusCode, respneCon.StatusCode, VERSION, uptimeString}
		json.NewEncoder(w).Encode(diag)
		return
	}
	http.Error(w, "only get method allowed", http.StatusNotImplemented)
	return
}
func cachingSpecies(){
	//TODO CreateFile
	f, err := os.Create("Species2.txt")
	if err != nil {
		fmt.Println("We fucked up")
		//We fucked up
	}
	datachan := make(chan []byte)
	//data[MAXCALL];
	var wg sync.WaitGroup
	for i := 0; i < MAXCALL; i++ {
		wg.Add(1)
		go func() {
			limit := strconv.Itoa(GBIFSPECIESLIMIT)
			offset := strconv.Itoa(GBIFSPECIESLIMIT*i)
			query := "http://api.gbif.org/v1/species/search?limit=" + limit + "&offset=" + offset
			resp, err := http.Get(query)
			if err != nil {
				//We fucked up
			}
			data, _ := ioutil.ReadAll(resp.Body)
			datachan <- data
		}() 
		go func(){
			getdata := <- datachan
			f.Write(getdata)
			fmt.Println("We have recived data")
			wg.Done()
		}()

	} 
	wg.Wait()
	fmt.Println("We are done")
	f.Close()
}
func init(){
	startTime = time.Now()
}
func main() {
	fmt.Println("Caching API's ")
	cachingSpecies();
	fmt.Println("Starting application:")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", nilHandler)
	http.HandleFunc("/conservation/v1/country/", countryHandler)
	http.HandleFunc("/conservation/v1/species/", speciesHandler)
	http.HandleFunc("/conservation/v1/diag/", diagHandler)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
