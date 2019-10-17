package imt2681ass1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
)

/*CachingSpecies ...
First see if the FilenameS exisit
if it dosen't it creates it and then fill it with API calls to GBIF

if the file exist
Checks if it is needed to cache species (if the FilenameS is older than 24 hours)

*/
func CachingSpecies() {

	//See if we should run caching
	should, file := shouldFileCache(FilenameS, "")
	//it either threw an error or it already existed
	if should == Error || should == DirFail {
		fmt.Println("Failed to find or create file")
		return
	}
	if should == Exist {
		//No need to cache the files
		return
	}
	//The array we store all the date in and will write to file
	var speciesArray []GbifJSON
	species := make(chan GbifJSON)
	var wg sync.WaitGroup
	for i := 0; i < MaxCall; i++ {
		wg.Add(1)
		go func(i int) {
			limit := strconv.Itoa(GBIFSpeciesLimit)
			offset := strconv.Itoa(GBIFSpeciesLimit * i)
			query := "http://api.gbif.org/v1/species/search?limit=" + limit + "&offset=" + offset
			resp, err := http.Get(query)
			if err != nil {
				fmt.Println(err)
				return
			}
			data, _ := ioutil.ReadAll(resp.Body)
			var gbifJSON GbifJSON
			json.Unmarshal(data, &gbifJSON)
			gbifJSON.index = i
			species <- gbifJSON
		}(i)
		go func() {
			getspecies := <-species

			speciesArray = append(speciesArray, getspecies)

			fmt.Println("We have recived data")
			wg.Done()
		}()

	}
	wg.Wait()
	//Sort the data by when we called it
	//Quicker to sort it here than when displaying
	sort.SliceStable(speciesArray, func(i, j int) bool {
		return speciesArray[i].index < speciesArray[j].index
	})
	speciesBytes, _ := json.Marshal(speciesArray)
	file.Write(speciesBytes)
	fmt.Println("We are done")
	file.Close()
}

//CachingCounry Input the IOS2 of the country and
//file to where to want to cache country information
func CachingCounry(ios2 string, file *os.File) (country Country, exsist bool) {
	query := "https://restcountries.eu/rest/v2/alpha/" + ios2 + "?fields=name;flag;alpha2Code"
	fmt.Println(query)
	resp, err := http.Get(query)
	if err != nil {
		//We fucked up
		fmt.Println("Query fucked")
		return country, false
	}
	data, _ := ioutil.ReadAll(resp.Body)
	var countryJSON CountryJSON
	json.Unmarshal(data, &countryJSON)
	if countryJSON.Code == "" {
		return country, false
	}
	//We do one initsal query so we can figure out how many we need to complete it
	//Query for country
	query = "http://api.gbif.org/v1/occurrence/search?limit=" + strconv.Itoa(GBIFOccurancesLimit) + "&country=" + ios2
	resp, err = http.Get(query)
	if err != nil {
		//We fucked up
		fmt.Println("Query fucked")
		return country, false
	}
	//CounrtySpeciesArray
	var CSA CountryGbifJSON
	data, _ = ioutil.ReadAll(resp.Body)
	json.Unmarshal(data, &CSA)
	//Do several API call till we have all the results or we run out of calls
	resultSum := CSA.Count
	fmt.Println("We got one")
	datachan := make(chan CountryGbifJSON)
	var wg sync.WaitGroup
	for calls := 1; GBIFOccurancesLimit*calls < resultSum && calls < MaxCall; calls++ {
		fmt.Println("We go")
		wg.Add(1)
		go func(calls int) {
			limit := strconv.Itoa(GBIFOccurancesLimit)
			offset := strconv.Itoa(GBIFOccurancesLimit * calls)
			query = "http://api.gbif.org/v1/occurrence/search?limit=" + limit + "&offset=" + offset + "&country=" + ios2
			fmt.Println(query)
			resp, err = http.Get(query)
			if err != nil {
				fmt.Println("Query fucked")
				return
			}
			data, _ = ioutil.ReadAll(resp.Body)
			var temp CountryGbifJSON
			json.Unmarshal(data, &temp)
			datachan <- temp
		}(calls)
		go func() {
			countrydata := <-datachan
			CSA.Results = append(CSA.Results, countrydata.Results...)
			fmt.Println("We have recived the data")
			wg.Done()
		}()

	}
	wg.Wait()
	fmt.Println("We done")
	//Sort the result based on index
	sort.SliceStable(CSA.Results, func(i, j int) bool {
		return CSA.Results[i].Key < CSA.Results[j].Key
	})
/*
	for i := 1; i < len(CSA.Results); i++ {
		if CSA.Results[i-1].Species == CSA.Results[i].Species {
			CSA.Results = append(CSA.Results[:i], CSA.Results[i+1:]...)
		}
	}
*/
	countySpecies := unique(CSA);
	CSA.Results = countySpecies;

	country.Code = countryJSON.Code
	country.CountryFlag = countryJSON.CountryFlag
	country.CountryName = countryJSON.CountryName
	var Species []string
	var SpeciesKey []int
	for index := 0; index < len(CSA.Results); index++ {
		Species = append(Species, CSA.Results[index].Species)
		SpeciesKey = append(SpeciesKey, CSA.Results[index].Key)
	}
	country.Species = Species
	country.SpeciesKey = SpeciesKey
	countryBytes, _ := json.Marshal(country)
	file.Write(countryBytes)
	file.Close()
	return country, true
}

//Given a filename return true and the file if the gile need to be cahced
//The directory is optional
//And false and a nil if it dosen't
func shouldFileCache(filename string, dir string) (Msg, *os.File) {
	file, err := os.Open(dir + filename)
	if err != nil {
		fmt.Println("Can not open file")
		if dir != "" {
			fmt.Println("Trying to create directory")
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				fmt.Println("Failed to create directory")
				return DirFail, nil
			}
		}
		//Create the file
		fmt.Println("trying to create file")
		file, err = os.Create(dir + filename)
		if err != nil {

			return Error, nil
		}
		fmt.Println("file created")
		return Created, file
	}

	//The file exist and we need to see how old it is
	info, err := os.Stat(dir + filename)
	if err != nil {
		fmt.Println(err)
		return Error, nil
	}
	mtime := info.ModTime()
	fmt.Println("Last changed:", mtime)
	timenow := time.Now()
	fmt.Println("Time now:", timenow)
	fmt.Println("Is", timenow.Sub(mtime).Hours(), "larger than", 24, "?")
	//Does not care for timezones btw
	if timenow.Sub(mtime).Hours() > 24 {
		fmt.Println("Cache is old, Run Update")
		err := os.Remove(dir + filename)
		if err != nil {
			fmt.Println(err)
			return Error, nil
		}
		file, err = os.Create(dir + filename)
		if err != nil {
			fmt.Println(err)
			return Error, nil
		}
		return OldRenew, file
	}
	fmt.Println("Cache is recent, No need to update")
	return Exist, file
}
func unique(country CountryGbifJSON) ([]CountySpecies) {
    keys := make(map[CountySpecies]bool)
	list := []CountySpecies{}

	 for index := 0; index < len(country.Results); index++ {
		 entry := country.Results[index]
		if _, value := keys[entry]; !value {
            keys[entry] = true
			list = append(list, entry)
        }
	 }   
    return list
}
