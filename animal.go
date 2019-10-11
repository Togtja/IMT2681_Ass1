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
	"sync"
	"sort"
	"time"
)

//VERSION The current api version
const VERSION string = "v1"

//MAXCALL How may calls we are going to make to the API
const MAXCALL int = 2

//GBIFOCCURANCESLIMIT What the limit for Gbif occurances is
const GBIFOCCURANCESLIMIT int = 300

//GBIFSPECIESLIMIT What the limit for Gbif Speces is
const GBIFSPECIESLIMIT int = 1000

//Species file name
const FILENAME_S string = "Species.json"

// Timm used for when the program start up
var startTime time.Time
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
//GbifARRAY n nono no
type GbifARRAY struct{

}
*/
//GbifJSON get the spices result
type GbifJSON struct {
	Results []SpeciesJSON `json:"results"`
}

/*
Species represents species by country
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
Diag represent the diagnostic tool to see API status code for
external APIs, and our version and uptime
	{
		"gbif": <value>, e.g "200"
		"restcountries": <value>, e.g "200"
		"version": <value>, e.g "v1"
		"uptime": <value>, e.g "1200 seconds"
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
			file, err := os.Open(FILENAME_S)
			if err != nil {
				log.Fatal(err)
			}
			
			data, err := ioutil.ReadAll(file)
			file.Close()

			var speciesArray []GbifJSON
			var species GbifJSON;

			json.Unmarshal(data, &speciesArray)
			for index := 0; index < len(speciesArray); index++ {
				species.Results = append(species.Results, speciesArray[index].Results...)
			}
			//Sort the data by key
			sort.SliceStable(species.Results, func (i, j int ) bool {
				return species.Results[i].Key < species.Results[j].Key
			})
			json.NewEncoder(w).Encode(species)
		} else if key, err := strconv.Atoi(parts[4]); err == nil {
			
			file, err := os.Open(FILENAME_S)
			if err != nil {
				log.Fatal(err)
			}
			
			data, err := ioutil.ReadAll(file)
			file.Close()

			var speciesArray []GbifJSON
			var species GbifJSON;

			json.Unmarshal(data, &speciesArray)
			for index := 0; index < len(speciesArray); index++ {
				species.Results = append(species.Results, speciesArray[index].Results...)
			}
			//Sort the data by key
			for index := 0; index < len(species.Results); index++ {
				if species.Results[index].Key == key{
					json.NewEncoder(w).Encode(species.Results[index])
					return
					
				}
			}
			http.Error(w, "could not find the species key", http.StatusBadRequest)
			
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
	var speciesArray []GbifJSON 
	file, err := os.Open(FILENAME_S)
	if err != nil {
		//Create the OS
		file, err = os.Create(FILENAME_S)
		if err != nil {
			fmt.Println(err)
			return
		}

	} else {
		//it exsist and we need to see how old it is
		info, err := os.Stat(FILENAME_S)
		if err != nil {
			fmt.Println(err)
			return
		}
		mtime := info.ModTime()
		fmt.Println("Last changed:", mtime)
		timenow := time.Now()
		fmt.Println("Time now:", timenow)
		fmt.Println("Is", timenow.Sub(mtime).Hours(), "larger than", 24, "?")
		//Does not care for timezones btw
		if timenow.Sub(mtime).Hours() > 24 {
			fmt.Println("Cache is old, Run Update")
			err := os.Remove(FILENAME_S)
			if err != nil {
				fmt.Println(err)
				return
			}
			file, err = os.Create(FILENAME_S)
			if err != nil {
				fmt.Println(err)
				return
			}

		} else {
			//We have the file and it's pretty new, no need to update
			fmt.Println("Found recent cahced data")
			return
		}
	}
	datachan := make(chan []byte)
	var wg sync.WaitGroup
	for i := 0; i < MAXCALL; i++ {
		wg.Add(1)
		go func(i int) {
			limit := strconv.Itoa(GBIFSPECIESLIMIT)
			offset := strconv.Itoa(GBIFSPECIESLIMIT * i)
			query := "http://api.gbif.org/v1/species/search?limit=" + limit + "&offset=" + offset
			resp, err := http.Get(query)
			if err != nil {
				fmt.Println(err)
				return
			}
			data, _ := ioutil.ReadAll(resp.Body)
			datachan <- data
		}(i)
		go func() {
			getdata := <-datachan
			var gbifJSON GbifJSON
			json.Unmarshal(getdata, &gbifJSON)
			speciesArray = append(speciesArray, gbifJSON)

			fmt.Println("We have recived data")
			wg.Done()
		}()

	}
	wg.Wait()
	speciesBytes,_ := json.Marshal(speciesArray)
	file.Write(speciesBytes)
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
