package imt2681ass1

import "time"

//Version The current api version
const Version string = "v1"

//MaxCall How may calls we are going to make to the GBIF API on startup
//This number is called when you look for countries Speices
const MaxCall int = 500

//GBIFOccurancesLimit What the limit for Gbif occurances is
const GBIFOccurancesLimit int = 300

//GBIFOccurancesMaxOffset The max offset for occurances is max(200000, offset + limit)
const GBIFOccurancesMaxOffset int = 200000 - GBIFOccurancesLimit

//GBIFSpeciesLimit What the limit for Gbif Speces is
const GBIFSpeciesLimit int = 1000

//GBIFSpeciesMaxOffset The max offset allowed in Species
const GBIFSpeciesMaxOffset int = 100000

//FilenameS Species file name for Species
const FilenameS string = "Species.json"

//CountyFileFolder is the directory all API request for secific countries goes
const CountyFileFolder string = "Countries/"

// Timm used for when the program start up
var startTime time.Time

//Msg A "go enum" for what the shouldChacheFile returned
type Msg int

const (
	//Error an error occured
	Error Msg = 0
	//OldRenew The file did exist but is now recreated due to age
	OldRenew Msg = 1
	//Created we created the file
	Created Msg = 2
	//Exist the file Exist and is recent
	Exist Msg = 3
	//DirFail directory failed to create
	DirFail Msg = 4
)
