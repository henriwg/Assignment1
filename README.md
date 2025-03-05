# Countries API

## Overview
This API gives information about countries, their population data, and API status. It is also possible to get more detailed information about a specific timeperiod. It fetches data from external services and gives a structured JSON file.

## Features
- Retrieve your desired countries information using an ISO country code at the end of a link.
- Fetch population data for a country you want with optional limits for the year of the data.
- Check the status of the APIs being used in case something is wrong.

## Technologies Used
- Golang is used for backend development
- HTTP server mux is used for routing
- Net/HTTP is used forHandling API requests

### Prerequisites
- Go 1.XX or later has to be installed on your system.
- A terminal/command prompt.
- An API key (if required by external services).

## Usage
The service has three resource root paths:
- /countryinfo/v1/info/
- /countryinfo/v1/population/
- /countryinfo/v1/status/

The web service should run on port 8080, and in that case the resource root paths should look like this:

- http://localhost:8080/countryinfo/v1/info/
- http://localhost:8080/countryinfo/v1/population/
- http://localhost:8080/countryinfo/v1/status/

You can enter your desired country to get specific data on info and population, and for population you can add an interval for years

- http://localhost:8080/countryinfo/v1/info/no
- http://localhost:8080/countryinfo/v1/population/no
- http://localhost:8080/countryinfo/v1/population/no?limit=2010-2015

The respons will be in json format

### Clone the Repository
```sh
git clone https://github.com/yourusername/country-info-api.git
cd country-info-api