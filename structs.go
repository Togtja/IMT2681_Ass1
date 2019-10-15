package imt2681ass1

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

//CountryGbifJSON get the result from occurances of country
type CountryGbifJSON struct {
	Results []CountySpecies `json:"results"`
	Count   int             `json:"count"`
}

/*
CountySpecies represents species by country
*/
type CountySpecies struct {
	Key     int    `json:"speciesKey"`
	Species string `json:"species"`
}

//GbifJSON get the species result
type GbifJSON struct {
	Results []SpeciesJSON `json:"results"`
	index   int
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
	Year 	  string `json:"year"`
}
//SpeciesYear is to grap the year and add it to SpeicesJSON's year
type SpeciesYear struct {
	Year 	  string `json: "year"`
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
