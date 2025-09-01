package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/santosjordi/fc_challenge_weather_cloudrun/config"
	"github.com/stretchr/testify/require"
)

// These global variables will be overridden in the tests to use the mock server URLs.
// This is a simple form of dependency injection for testing purposes.
var weatherAPIKey = "YOUR_MOCK_API_KEY"

// TestHandler tests the main handler function with different scenarios.
func TestHandler(t *testing.T) {
	// Start mock servers for external APIs
	viaCepMock := mockViaCEPServer()
	defer viaCepMock.Close()

	weatherAPIMock := mockWeatherAPIServer()
	defer weatherAPIMock.Close()

	// Override the base URLs to point to our mock servers for the test
	viaCepURLBase = viaCepMock.URL
	weatherAPIURLBase = weatherAPIMock.URL

	appInstance := &app{
		cfg: &config.Config{
			WeatherAPIKey: weatherAPIKey,
			// Fill other fields as needed for your test
		},
	}

	// Define test cases using a table-driven approach
	testCases := []struct {
		name           string
		cep            string
		expectedStatus int
		expectedBody   string
		expect         func(*testing.T, *http.Response)
	}{
		{
			name:           "Valid CEP - Success",
			cep:            "89068210",
			expectedStatus: http.StatusOK,
			expect: func(t *testing.T, resp *http.Response) {
				var result TempResponse
				err := json.NewDecoder(resp.Body).Decode(&result)
				require.NoError(t, err)

				require.Equal(t, 25.5, result.TempC)
				require.Equal(t, 77.9, result.TempF)
				require.Equal(t, 298.65, result.TempK)
			},
		},
		{
			name:           "Invalid Format CEP",
			cep:            "12345",
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   "invalid zipcode\n",
		},
		{
			name:           "Non-Existent CEP",
			cep:            "99999999",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "can not find zipcode\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a request to the handler with a test URL
			req := httptest.NewRequest(http.MethodGet, "/"+tc.cep, nil)
			rr := httptest.NewRecorder()

			// Call the handler function
			appInstance.handler(rr, req)

			// Assert HTTP status code
			require.Equal(t, tc.expectedStatus, rr.Code)

			// Assert body content if a specific string is expected
			if tc.expectedBody != "" {
				body, _ := io.ReadAll(rr.Result().Body)
				require.Equal(t, tc.expectedBody, string(body))
			}

			// Run custom assertions for specific test cases
			if tc.expect != nil {
				tc.expect(t, rr.Result())
			}
		})
	}
}

// mockViaCEPServer creates a test server that mimics the ViaCEP API.
func mockViaCEPServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ws/89068210/json/" {
			resp := map[string]string{"localidade": "Indaial"}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	}))
}

// mockWeatherAPIServer creates a test server that mimics the WeatherAPI.
func mockWeatherAPIServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("q") == "Indaial" {
			resp := map[string]interface{}{"current": map[string]float64{"temp_c": 25.5}}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}
		http.Error(w, "{}", http.StatusNotFound)
	}))
}
