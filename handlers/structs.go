package handlers

type RestCountriesData struct {
	Name struct {
		CountryName string `json:"common"`
		Official    string `json:"official"`
	} `json:"name"`
	Continents  []string          `json:"continents"`
	Population  int               `json:"population"`
	Languages   map[string]string `json:"languages"`
	Bordering   []string          `json:"borders"`
	Capital     []string          `json:"capital"`
	CountryFlag struct {
		Png string `json:"png"`
	} `json:"flags"`
	Currencies map[string]struct {
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	} `json:"currencies"`
	Region    string   `json:"region"`
	Subregion string   `json:"subregion"`
	Timezones []string `json:"timezones"`
}

type StatusResponse struct {
	CountriesNowAPI  string `json:"countriesnowapi"`
	RestCountriesAPI string `json:"restcountriesapi"`
	Version          string `json:"version"`
	Uptime           int    `json:"uptime"`
}

type PopulationResponse struct {
	Mean   int `json:"mean"`
	Values []struct {
		Year  int `json:"year"`
		Value int `json:"value"`
	} `json:"values"`
}

type ErrorResponse struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}
