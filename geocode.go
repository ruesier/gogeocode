// Package gogeocode provides utilities for utilizing the geocoding API provided through https://geocode.maps.co.
package gogeocode

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Client is used to call geocoding api
type Client struct {
	ApiKey string
}

// Address is apart of the Revese response
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

// Response contains all possible fields returned by the API
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
	// Address is only included when calling Reverse
	Address Address `json:"address"`
}

func buildGeocodeURL(query, key string) string {
	safequery := url.QueryEscape(strings.ReplaceAll(query, " ", "+"))
	return "https://geocode.maps.co/search?q=" + safequery + "&api_key=" + key
}

var (
	ErrAuthorization = errors.New("Geocode invalid API Key")
	ErrThrottle      = errors.New("Geocode failed due to exceeding rquest limit")
	ErrTraffic       = errors.New("Geocode failed due to high traffic on geocode server")
	ErrFlooding      = errors.New("Geocode has detected api key abuse, contact: https://maps.co/contact/ to resolve")
)

// Geocode takes a string description of a location and returns precise location data.
// Possible queries include addresses or famous place names.
func (c Client) Geocode(query string) ([]*Response, error) {
	return c.GeocodeWithContext(context.Background(), query)
}

// GeocodeWithContext performs the same request as Geocode using the given context.
func (c Client) GeocodeWithContext(ctx context.Context, query string) ([]*Response, error) {
	var results []*Response

	resp, err := callAPI(ctx, buildGeocodeURL(query, c.ApiKey))
	if err != nil {
		return results, err
	}

	// Callers should close resp.Body when done reading from it
	// Defer the closing of the body
	defer resp.Body.Close()

	// Use json.Decode for reading streams of JSON data
	err = json.NewDecoder(resp.Body).Decode(&results)
	return results, err
}

// AddressGeocode geocodes a specific address.
func (c Client) AddressGeocode(street, city, county, state, country, postalcode string) ([]*Response, error) {
	return c.AddressGeocodeWithContext(context.Background(), street, city, county, state, country, postalcode)
}

// AddressGeocodeWithContext performs same action as AddressGeocode with provided context
func (c Client) AddressGeocodeWithContext(ctx context.Context, street, city, county, state, country, postalcode string) ([]*Response, error) {
	var fields []string
	if len(street) > 0 {
		fields = append(fields, "street="+strings.ReplaceAll(street, " ", "+"))
	}
	if len(city) > 0 {
		fields = append(fields, "city="+strings.ReplaceAll(city, " ", "+"))
	}
	if len(county) > 0 {
		fields = append(fields, "county="+strings.ReplaceAll(county, " ", "+"))
	}
	if len(state) > 0 {
		fields = append(fields, "state="+strings.ReplaceAll(state, " ", "+"))
	}
	if len(country) > 0 {
		fields = append(fields, "country="+strings.ReplaceAll(country, " ", "+"))
	}
	if len(postalcode) > 0 {
		fields = append(fields, "postalcode="+strings.ReplaceAll(postalcode, " ", "+"))
	}
	fields = append(fields, "api_key="+c.ApiKey)
	url := "https://geocode.maps.co/search?" + strings.Join(fields, "&")

	var results []*Response

	resp, err := callAPI(ctx, url)
	if err != nil {
		return results, err
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

// Reverse takes a latitude and longitude and returns nearest address
func (c Client) Reverse(lat, long float64) (*Response, error) {
	return c.ReverseWithContext(context.Background(), lat, long)
}

func (c Client) ReverseWithContext(ctx context.Context, lat, long float64) (*Response, error) {
	var result *Response

	resp, err := callAPI(ctx, buildReverseURL(lat, long, c.ApiKey))
	if err != nil {
		return result, err
	}

	// Callers should close resp.Body when done reading from it
	// Defer the closing of the body
	defer resp.Body.Close()

	result = &Response{}
	// Use json.Decode for reading streams of JSON data
	err = json.NewDecoder(resp.Body).Decode(result)
	return result, err
}

func callAPI(ctx context.Context, url string) (*http.Response, error) {
	// Build the request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	// For control over HTTP client headers, redirect policy, and other settings, create a Client
	// A Client is an HTTP client
	client := &http.Client{}

	// Send the request via a client
	// Do sends an HTTP request and returns an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		return resp, err
	}

	switch resp.StatusCode {
	case 401:
		return resp, ErrAuthorization
	case 429:
		return resp, ErrThrottle
	case 503:
		return resp, ErrTraffic
	case 403:
		return resp, ErrFlooding
	}

	if resp.StatusCode >= 400 {
		return resp, fmt.Errorf("unrecognized error {StatusCode: %d}", resp.StatusCode)
	}

	return resp, nil
}
