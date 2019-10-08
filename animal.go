package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var startTime time.Time

//VERSION The current api version
const VERSION string = "v1"

//MAXCALL How may calls we are going to make to the API
const MAXCALL int = 5

//GBIFOCCURANCESLIMIT What the limit for Gbif occurances is
const GBIFOCCURANCESLIMIT int = 300

//GBIFSPECIESLIMIT What the limit for Gbif Speces is
const GBIFSPECIESLIMIT int = 1000

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

//GbifJSON get the spices result
type GbifJSON struct {
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
	Key       int    `json:"key"`
	Kingdom   string `json:"kingdom"`
	Phylum    string `json:"phylum"`
	Order     string `json:"order"`
	Family    string `json:"family"`
	Species   string `json:"species"`
	SciName   string `json:"scientificName"`
	CanonName string `json:"canonicalName"`
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
	Gbif          int    `json:"gbif"`
	Restcountries int    `json:"restcountries"`
	Version       string `json:"version"`
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
			http.Error(w, "Expecting format .../specieskey", http.StatusBadRequest)
			return
		}
		if parts[4] == "" {
			//Get all
		} else if _, err := strconv.Atoi(parts[4]); err == nil {
			file, err := os.Open("Species.txt")
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			data, err := ioutil.ReadAll(file)
			i := 0
			var species []SpeciesJSON
			for {
				var doc SpeciesJSON
				err := json.Unmarshal(data, &doc)
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Fatal(err)
				}
				species[i] = doc
				i++
			}
			json.NewEncoder(w).Encode(species)

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
		uptimeString := fmt.Sprintf("%.0f seconds", uptime.Seconds())
		diag := Diag{respneGbif.StatusCode, respneCon.StatusCode, VERSION, uptimeString}
		json.NewEncoder(w).Encode(diag)
		return
	}
	http.Error(w, "only get method allowed", http.StatusNotImplemented)
	return
}
func cachingSpecies() {
	//TODO CreateFile

	file, err := os.Open("Species.txt")
	if err != nil {
		//Create the OS
		file, err = os.Create("Species.txt")

	} else {
		//ut exsist and we need to see how old it is
		info, err := os.Stat("Species.txt")
		if err != nil {
			//Send help
		}
		mtime := info.ModTime()
		fmt.Println("Last changed:", mtime)
		timenow := time.Now()
		//Does not care for timezones btw
		if timenow.Hour() > mtime.Hour()+24 {
			fmt.Println("Cache is old, Run Update")

		} else {
			return
		}
	}

	if err != nil {
		fmt.Println("We fucked up")
		//We fucked up
	}
	datachan := make(chan []byte)
	//data[MAXCALL];
	var wg sync.WaitGroup
	for i := 0; i < MAXCALL; i++ {
		wg.Add(1)
		go func(i int) {
			limit := strconv.Itoa(GBIFSPECIESLIMIT)
			offset := strconv.Itoa(GBIFSPECIESLIMIT * i)
			query := "http://api.gbif.org/v1/species/search?limit=" + limit + "&offset=" + offset
			resp, err := http.Get(query)
			if err != nil {
				//We fucked up
			}
			data, _ := ioutil.ReadAll(resp.Body)
			var results GbifJSON
			//Get it in correct format
			json.Unmarshal(data, &results)
			data, _ = json.Marshal(results.Results)
			datachan <- data
		}(i)
		go func() {
			getdata := <-datachan
			file.Write(getdata)
			fmt.Println("We have recived data")
			wg.Done()
		}()

	}
	wg.Wait()
	fmt.Println("We are done")
	file.Close()
}
func init() {
	startTime = time.Now()
}
func main() {
	fmt.Println("Caching API's ")
	cachingSpecies()
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
