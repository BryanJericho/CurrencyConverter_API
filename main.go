package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const apiKey = "fb97d048b434372ce5b3c65f"
const apiURL = "https://v6.exchangerate-api.com/v6/"

type ExchangeRateResponse struct {
	Result          string             `json:"result"`
	BaseCode        string             `json:"base_code"`
	ConversionRates map[string]float64 `json:"conversion_rates"`
}

func getExchangeRates(baseCurrency string) (*ExchangeRateResponse, error) {
	url := fmt.Sprintf("%s%s/latest/%s", apiURL, apiKey, baseCurrency)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not fetch exchange rates: %v", err)
	}
	defer resp.Body.Close()

	var exchangeRateResponse ExchangeRateResponse
	if err := json.NewDecoder(resp.Body).Decode(&exchangeRateResponse); err != nil {
		return nil, fmt.Errorf("could not decode response: %v", err)
	}

	if exchangeRateResponse.Result != "success" {
		return nil, fmt.Errorf("API request failed with result: %s", exchangeRateResponse.Result)
	}

	return &exchangeRateResponse, nil
}

func convertCurrency(amount float64, fromCurrency string, toCurrency string) {
	rates, err := getExchangeRates(fromCurrency)
	if err != nil {
		log.Fatalf("Error fetching exchange rates: %v", err)
	}

	toRate, exists := rates.ConversionRates[toCurrency]
	if !exists {
		log.Fatalf("Currency %s not found in exchange rates", toCurrency)
	}

	convertedAmount := amount * toRate
	fmt.Printf("%.2f %s = %.2f %s\n", amount, fromCurrency, convertedAmount, toCurrency)
}

func main() {
	// Input dari pengguna
	var amount float64
	var fromCurrency, toCurrency string

	fmt.Print("Enter amount: ")
	fmt.Scanln(&amount)
	fmt.Print("Enter base currency (e.g. USD): ")
	fmt.Scanln(&fromCurrency)
	fmt.Print("Enter target currency (e.g. EUR): ")
	fmt.Scanln(&toCurrency)

	convertCurrency(amount, fromCurrency, toCurrency)
}
