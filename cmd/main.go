package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/santosjordi/fc_challenge_weather_cloudrun/config"
)

// ViaCEPResponse represents the structure of the response from ViaCEP
type ViaCEPResponse struct {
	Localidade string `json:"localidade"`
	Erro       bool   `json:"erro"`
}

// WeatherAPIResponse represents the structure of the response from WeatherAPI
type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

// TempResponse represents the final JSON structure for the client
type TempResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type app struct {
	cfg *config.Config
}

// Function to handle the main logic
func (a *app) handler(w http.ResponseWriter, r *http.Request) {
	// Get the ZIP code from the URL path
	cep := r.URL.Path[1:]

	// Validate ZIP code format
	if !regexp.MustCompile(`^\d{8}$`).MatchString(cep) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	// Fetch city from ViaCEP
	viaCepURL := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	resp, err := http.Get(viaCepURL)
	if err != nil {
		http.Error(w, "can not find zipcode", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var viaCepData ViaCEPResponse
	if err := json.Unmarshal(body, &viaCepData); err != nil {
		http.Error(w, "can not find zipcode", http.StatusInternalServerError)
		return
	}
	if viaCepData.Erro || viaCepData.Localidade == "" {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	// Fetch temperature from WeatherAPI
	weatherAPIKey := a.cfg.WeatherAPIKey
	weatherURL := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", weatherAPIKey, viaCepData.Localidade)
	resp, err = http.Get(weatherURL)
	if err != nil {
		http.Error(w, "can not find temperature", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, _ = io.ReadAll(resp.Body)
	var weatherData WeatherAPIResponse
	if err := json.Unmarshal(body, &weatherData); err != nil {
		http.Error(w, "can not find temperature", http.StatusInternalServerError)
		return
	}

	tempC := weatherData.Current.TempC
	tempF := tempC*1.8 + 32
	tempK := tempC + 273.15 // Use 273.15 for more precision

	// Construct and send the final response
	response := TempResponse{
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func main() {
	// 1. Load the configuration file at startup
	cfg, err := config.LoadConfig(".env") // or pass it as an argument
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	// 2. Create an instance of the app struct with the loaded config
	application := &app{cfg: cfg}

	// 3. Register the handler method
	http.HandleFunc("/", application.handler)

	port := cfg.ServerPort // Get port from the config
	fmt.Printf("Server is running on port %s\n", port)
	http.ListenAndServe(port, nil)
}
