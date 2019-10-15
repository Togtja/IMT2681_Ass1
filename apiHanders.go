package imt2681ass1

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
)

//NilHandler throws a Bad Request
func NilHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Default Handler: Invalid request received.")
	http.Error(w, "Invalid request", http.StatusBadRequest)
}

//CountryHandler ...
func CountryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.Header.Add(w.Header(), "content-type", "application/json")
		fmt.Println(r.URL.Path)
		parts := strings.Split(r.URL.Path, "/")
		fmt.Println(parts)
		if len(parts) != 5 {
			http.Error(w, "Expecting format .../country_ios2", http.StatusBadRequest)
			return
		}
		if parts[4] == "" {
			//Get all
			//TODO: Cache all country species and then itterate through it
			// The caching is done when we add a new country
			//File name should be Countries/CountriesSpecies.json
			// (aka dir = CountyFileFolder ,finename =  "Countries" + filename)

		} else {

			limit, _ := strconv.ParseInt(r.FormValue("limit"), 10, 32)
			if limit <= 0 {
				//Default limit
				limit = 300
			}
			offset, _ := strconv.ParseInt(r.FormValue("offset"), 10, 32)
			fmt.Println("Limit:",limit)
			fmt.Println("Offet:",offset)
			if len(parts[4]) != 2 {
				fmt.Print(parts[4])
				http.Error(w, "Expecting country ios2 code", http.StatusBadRequest)
				return
			}
			//See if we have cahced data from before
			should, file := shouldFileCache(strings.ToUpper(parts[4])+FilenameS, CountyFileFolder)
			var finalResult Country
			if should == Error || should == DirFail {
				http.Error(w, "Could not find or create file", http.StatusInternalServerError)
				return
			} else if should == Created || should == OldRenew {
				var found bool
				finalResult, found = CachingCounry(parts[4], file)
				if !found {
					http.Error(w, "could not find country ios2 code", http.StatusBadRequest)
				}
			} else if should == Exist { //Could be an else, but just to make sure
				data, err := ioutil.ReadAll(file)
				if err != nil {
					fmt.Println(err)
					return
				}
				json.Unmarshal(data, &finalResult)
				file.Close()
			}

			if limit+offset < int64(len(finalResult.SpeciesKey)){
				finalResult.SpeciesKey = finalResult.SpeciesKey[offset:limit+offset];
				finalResult.Species = finalResult.Species[offset:limit+offset]
				//Else will show everything
			}
			json.NewEncoder(w).Encode(finalResult)

		}
		/*
			else {
				fmt.Print(parts[4])
				http.Error(w, "Expecting country ios2 code", http.StatusBadRequest)
				return
			}
		*/
		return
	}
	http.Error(w, "only get method allowed", http.StatusNotImplemented)
	return
}

/*SpeciesHandler ..

If could I would also do on demand caching here, but I can't search the API for a specific key
So I cache asmuch as I can :)
*/
func SpeciesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.Header.Add(w.Header(), "content-type", "application/json")
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) != 5 {
			http.Error(w, "Expecting format .../specieskey", http.StatusBadRequest)
			return
		}
		if parts[4] == "" {
			file, err := os.Open(FilenameS)
			if err != nil {
				log.Fatal(err)
			}

			data, err := ioutil.ReadAll(file)
			file.Close()

			var speciesArray []GbifJSON
			var species GbifJSON

			json.Unmarshal(data, &speciesArray)
			for index := 0; index < len(speciesArray); index++ {
				species.Results = append(species.Results, speciesArray[index].Results...)
			}
			json.NewEncoder(w).Encode(species)
		} else if key, err := strconv.Atoi(parts[4]); err == nil {

			file, err := os.Open(FilenameS)
			if err != nil {
				log.Fatal(err)
			}

			data, err := ioutil.ReadAll(file)
			file.Close()

			var speciesArray []GbifJSON
			var species GbifJSON

			json.Unmarshal(data, &speciesArray)
			for index := 0; index < len(speciesArray); index++ {
				species.Results = append(species.Results, speciesArray[index].Results...)
			}
			//Sort the data by key
			for index := 0; index < len(species.Results); index++ {
				if species.Results[index].Key == key {
					//Example thing that's wierd
					//Key 1007 exist in species/1007/name
					//However in the species/search?limit=10&offset=750
					query := "http://api.gbif.org/v1/species/" + strconv.Itoa(key) + "/name"
					resp, err := http.Get(query)
					if err != nil {
						fmt.Println(err)
						return
					}
					data, _ := ioutil.ReadAll(resp.Body)
					var year SpeciesYear;
					json.Unmarshal(data, &year)
					species.Results[index].Year = year.Year;
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

//DiagHandler Gives Diagnostic tool in JSON format
func DiagHandler(w http.ResponseWriter, r *http.Request) {
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
		diag := Diag{respneGbif.StatusCode, respneCon.StatusCode, Version, uptimeString}
		json.NewEncoder(w).Encode(diag)
		return
	}
	http.Error(w, "only get method allowed", http.StatusNotImplemented)
	return
}

//Finds current uptime
func uptime() time.Duration {
	return time.Since(startTime)
}

//Start the timer to figure out current uptime
func init() {
	startTime = time.Now()
}
