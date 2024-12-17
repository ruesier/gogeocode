package gogeocode

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	ApiKey string
}

type Address struct {
	HouseNumber    string `json:"house_number"`
	Road           string `json:"road"`
	Neighbourhood  string `json:"neighbourhood"`
	Suburb         string `json:"suburb"`
	County         string `json:"county"`
	City           string `json:"city"`
	State          string `json:"state"`
	ISO3166_2_lvl4 string `json:"ISO3166-2-lvl4"`
	PostCode       string `json:"postcode"`
	Country        string `json:"country"`
	CountryCode    string `json:"country_code"`
}

type Response struct {
	PlaceID     uint64   `json:"place_id"`
	Licence     string   `json:"licence"`
	OSMType     string   `json:"osm_type"`
	OSMID       uint64   `json:"osm_id"`
	BoundingBox []string `json:"boundingbox"`
	Latitude    string   `json:"lat"`
	Longitude   string   `json:"lon"`
	DisplayName string   `json:"display_name"`
	Class       string   `json:"class"`
	Type        string   `json:"type"`
	Importance  float64  `json:"importance"`
	Address     Address  `json:"address"`
}

func buildGeocodeURL(query, key string) string {
	safequery := url.QueryEscape(strings.ReplaceAll(query, " ", "+"))
	return "https://geocode.maps.co/search?q=" + safequery + "&api_key=" + key
}

var (
	ErrThrottle = errors.New("Geocode failed due to exceeding rquest limit")
	ErrTraffic  = errors.New("Geocode failed due to high traffic on geocode server")
	ErrFlooding = errors.New("Geocode has detected api key abuse, contact: https://maps.co/contact/ to resolve")
)

func (c Client) Geocode(query string) ([]*Response, error) {
	var results []*Response

	// Build the request
	req, err := http.NewRequest("GET", buildGeocodeURL(query, c.ApiKey), nil)
	if err != nil {
		return results, err
	}

	// For control over HTTP client headers, redirect policy, and other settings, create a Client
	// A Client is an HTTP client
	client := &http.Client{}

	// Send the request via a client
	// Do sends an HTTP request and returns an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		return results, err
	}

	switch resp.StatusCode {
	case 429:
		return results, ErrThrottle
	case 503:
		return results, ErrTraffic
	case 403:
		return results, ErrFlooding
	}

	if resp.StatusCode >= 400 {
		return results, fmt.Errorf("unrecognized error {StatusCode: %d}", resp.StatusCode)
	}

	// Callers should close resp.Body when done reading from it
	// Defer the closing of the body
	defer resp.Body.Close()

	// Use json.Decode for reading streams of JSON data
	err = json.NewDecoder(resp.Body).Decode(&results)
	return results, err
}

func buildReverseURL(lat, long float64, key string) string {
	return fmt.Sprintf("https://geocode.maps.co/reverse?lat=%f&lon=%f&api_key=%s", lat, long, key)
}

func (c Client) Reverse(lat, long float64) (*Response, error) {
	var result *Response

	// Build the request
	req, err := http.NewRequest("GET", buildReverseURL(lat, long, c.ApiKey), nil)
	if err != nil {
		return result, err
	}

	// For control over HTTP client headers, redirect policy, and other settings, create a Client
	// A Client is an HTTP client
	client := &http.Client{}

	// Send the request via a client
	// Do sends an HTTP request and returns an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}

	switch resp.StatusCode {
	case 429:
		return result, ErrThrottle
	case 503:
		return result, ErrTraffic
	case 403:
		return result, ErrFlooding
	}

	if resp.StatusCode >= 400 {
		return result, fmt.Errorf("unrecognized error {StatusCode: %d}", resp.StatusCode)
	}

	// Callers should close resp.Body when done reading from it
	// Defer the closing of the body
	defer resp.Body.Close()

	result = &Response{}
	// Use json.Decode for reading streams of JSON data
	err = json.NewDecoder(resp.Body).Decode(result)
	return result, err
}
