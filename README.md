# IMT2681_Ass1
#Assigment specification:
# Assignment 1: Species Information Service

**NOTE** Due to the problems with the access to tokens to the redlist server, the assignment has been re-written to use a different, open species data source, GBIF. 

## Overview

In this assignment, you are going to develop a REST web application in Golang that provides the client to retrieve information about animal species. For this purpose, you will interrogate an existing web service and return the result in a given output format. 

The REST web service you will be using is the [Global Biodiversity Information Facility](https://www.gbif.org/what-is-gbif), GBIF. It is based on open standards and open access, and it allows the client to query various aspects of Earth's biodiversity. 

[The documentation of GBIF is available online](https://www.gbif.org/developer/summary). To use the GBIF web service through GET requests you do **not** need an API token! Nevertheless, be **mindful of rate limits** - we will talk about mitigation strategies in class.

In addition to the information provided by the GBIF database, you will need to use the [restcountries.eu](https://restcountries.eu) API for obtaining information about countries and their capitals, country codes, currency, etc. 

The final web service should be deployed on Heroku. The initial development should occur on your local machine. For the submission, you will both need to provide a URL to the deployed Heroku service as well as your repository.

In the following, you will find the specification for the REST API exposed to the user for interrogation/testing.



# Specification

Note: the specification may be subject to refinement (e.g., because of typos, technical problems).

## Endpoints

Your web service will have three resource root paths: 

```
/conservation/v1/country/
/conservation/v1/species/
/conservation/v1/diag/
```

Assuming your web service should run on localhost, port 8080, your resource root paths would thus be

```
http://localhost:8080/conservation/v1/country/
http://localhost:8080/conservation/v1/species/
http://localhost:8080/conservation/v1/diag/
````

The supported request/response pairs are specified in the following.



## List Species by Country

The purpose of this endpoint is to list a given number of species entries by country.

### Request

Syntax: `{:value}` indicates input parameter specified by the user - it is mandatory. `{entry}` indicates the optional nature of the entry.

```
Method: GET
Path: country/{:country_identifier}{?limit={:limit}}
```

{:country_identifier} refers to the [2-letter ISO code](https://www.iban.com/country-codes) for the country.

{:limit} indicates the number of entries to be retrieved from the target services. You may report back less entries if the target entries contain duplicates. See this [issue](https://git.gvk.idi.ntnu.no/course/imt2681/imt2681-2019/issues/8#note_6515) for clarification.

### Response

Syntax: `<value>` indicates JSON values populated by the webservice.

* Content type: `application/json`
* Status code: 200 if everything is OK, appropriate error code otherwise. 

Body:
```
{
   "code": "<country code in 2-letter ISO format>",
   "countryname": "<English human-readable country name>",
   "countryflag": "<link to svg version of country flag>",
   "species": [],
   "speciesKey": []
}
```

`species` is the value of the species entry (`species` field) from GBIF for a given country, and `speciesKey` is the equivalent numerical key for the species. Both lists should have the exact same length and should have the matching name-key pairing. The list should not contain duplicate entries (i.e., the number of returned entries can be smaller than the specified limit if duplicates are returned from the GBIF service).


## Information about specific species

This request-response pair provides information about specific species. This relates to ANY species, not only extinct.

### Request

Syntax: `{:value}` indicates input parameter specified by the user.

```
Method: GET
Path: species/{:speciesKey}
```

The entry `speciesKey` refers to the numeric identifier for a specific species.

## Response

Syntax: `<value>` indicates JSON values populated by the webservice.

* Content type: `application/json`
* Status code: 200 if everything is OK, appropriate error code otherwise. 

Body:
```
{
   "key": "<species key>",
   "kingdom": "<kingdom>",
   "phylum": "<phylum>",
   "order": "<order>",
   "family": "<family>",
   "genus": "<genus>",
   "scientificName":"<scientific name>",
   "canonicalName": "<canonical name>",
   "year": "<four-letter year>"
}
```

Hint: For the year, check the */species/{key}/name* path for a given key.

## Diagnostics interface

The diagnostics interface indicates the availability of individual services this service depends on. The reporting occurs based on status codes returned by the dependent services.

### Request

```
Method: GET
Path: diag/
```

### Response

Syntax: `<value>` indicates JSON values populated by the webservice.

* Content type: `application/json`
* Status code: 200 if everything is OK, appropriate error code otherwise. 

Body:
```
{
   "gbif": "<http status code for GBIF API>",
   "restcountries": "<http status code for restcountries API>",
   "version": "v1",
   "uptime": <time in seconds from the last service restart>
}
```



# Submission

The submission deadline is shown on the course [main page](Home). No extensions will be given for late submissions. 

The submission occurs via our submission system that not only facilitates the submission, but also the peer review of the assignment. The system will be made available closer to the deadline.

# Peer Review

After the submission deadline, there will be a second deadline during which you will review other students' submissions. To do this the system provides you with a checklist of aspects to assess. You will need to review at least two submissions to meet mandatory requirements of peer review, but you can review as many submissions as you like, which counts towards your participation mark for the course. The peer-review deadline will be shown in the submission system.
