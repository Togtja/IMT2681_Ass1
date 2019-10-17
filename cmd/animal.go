package main

import (
	"fmt"
	"imt2681ass1"
	"log"
	"net/http"
	"os"
)

/*
Question?
What do do with empty Species Key from Occurances?

What do do with Species without names or key?
http://api.gbif.org/v1/occurrence/search?country=se&limit=300&offset=0 (184)

*/
func main() {
	fmt.Println("Caching API's ")
	imt2681ass1.CachingSpecies()
	fmt.Println("Starting application:")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", imt2681ass1.NilHandler)
	http.HandleFunc("/conservation/v1/country/", imt2681ass1.CountryHandler)
	http.HandleFunc("/conservation/v1/species/", imt2681ass1.SpeciesHandler)
	http.HandleFunc("/conservation/v1/diag/", imt2681ass1.DiagHandler)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
