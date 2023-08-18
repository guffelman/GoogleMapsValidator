package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/tealeg/xlsx"
)

const (
	baseURL = "https://maps.googleapis.com/maps/api/geocode/json"
)

func main() {
	var apiKey, filename string
	fmt.Print("Enter your Google Maps API key: ")
	fmt.Scanln(&apiKey)
	fmt.Print("Enter the file path of the xlsx file: ")
	fmt.Scanln(&filename)

	columnIndex := 0

	xlFile, err := xlsx.OpenFile(filename)
	if err != nil {
		panic(err)
	}

	sheet := xlFile.Sheets[0]
	startCol := columnIndex + 1

	headers := []string{"Address", "City", "State", "Country", "Zip"}

	for i, header := range headers {
		sheet.Cell(0, startCol+i).SetValue(header)
	}

	totalRows := sheet.MaxRow - 1
	ticker := time.NewTicker(time.Second / 10)

	for rowIndex := 1; rowIndex < sheet.MaxRow; rowIndex++ {
		<-ticker.C

		cell := sheet.Cell(rowIndex, columnIndex)
		rawAddress := cell.String()

		if rawAddress == "" {
			fmt.Printf("%d of %d processed (skipped).\n", rowIndex, totalRows)
			continue
		}

		addressComponents, formatted := validateAndFormatAddress(rawAddress, apiKey)

		fmt.Printf("%d of %d processed. API Response: %s\n", rowIndex, totalRows, formatted)

		for i, component := range addressComponents {
			cell := sheet.Cell(rowIndex, startCol+i)
			cell.SetValue(component)
		}
	}

	ticker.Stop()

	err = xlFile.Save(filename)
	if err != nil {
		panic(err)
	}

	fmt.Println("Processing done!")
}

func validateAndFormatAddress(address string, apiKey string) ([]string, string) {
	encodedAddress := url.QueryEscape(address)
	fullURL := fmt.Sprintf("%s?address=%s&key=%s", baseURL, encodedAddress, apiKey)

	resp, err := http.Get(fullURL)
	if err != nil {
		return make([]string, 5), ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return make([]string, 5), ""
	}

	var geocodeResponse GeocodeResponse
	err = json.Unmarshal(body, &geocodeResponse)
	if err != nil {
		return make([]string, 5), ""
	}

	if geocodeResponse.Status != "OK" {
		fmt.Println("API Status Not OK:", geocodeResponse.Status)
		return make([]string, 5), ""
	}

	addressComponents := getAddressComponents(&geocodeResponse)
	return addressComponents, geocodeResponse.Results[0].FormattedAddress
}

func getAddressComponents(geocodeResponse *GeocodeResponse) []string {
	var addressComponents []string
	addressMap := make(map[string]string)

	for _, component := range geocodeResponse.Results[0].AddressComponents {
		switch component.Types[0] {
		case "street_number", "route":
			addressMap["address"] = addressMap["address"] + " " + component.LongName
		case "locality":
			addressMap["city"] = component.LongName
		case "administrative_area_level_1":
			addressMap["state"] = component.LongName
		case "country":
			addressMap["country"] = component.LongName
		case "postal_code":
			addressMap["zip"] = component.LongName
		}
	}

	addressComponents = append(addressComponents, addressMap["address"], addressMap["city"], addressMap["state"], addressMap["country"], addressMap["zip"])

	return addressComponents
}

type GeocodeResponse struct {
	Results []struct {
		AddressComponents []struct {
			LongName  string   `json:"long_name"`
			ShortName string   `json:"short_name"`
			Types     []string `json:"types"`
		} `json:"address_components"`
		FormattedAddress string `json:"formatted_address"`
	} `json:"results"`
	Status string `json:"status"`
}
