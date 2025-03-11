package gogeocode

import (
	"encoding/json"
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var APIKEY string

func TestMain(m *testing.M) {
	APIKEY = os.Getenv("GEOCODE_APIKEY")
	flag.StringVar(&APIKEY, "api", APIKEY, "provide an API key for testing")
	flag.Parse()

	os.Exit(m.Run())
}

func TestGeocode(t *testing.T) {
	testClient := Client{
		ApiKey: APIKEY,
	}
	// https://geocode.maps.co/search?q=555+5th+Ave+New+York+NY+10017+US&api_key=
	got, err := testClient.Geocode("555 5th Ave New York NY 10017 US")
	if err != nil {
		t.Fatalf("Geocode failed: %s", err)
	}

	var want []*Response
	json.NewDecoder(strings.NewReader(`[{"place_id":319634989,"licence":"Data © OpenStreetMap contributors, ODbL 1.0. https://osm.org/copyright","osm_type":"node","osm_id":1000793154,"boundingbox":["40.7557728","40.7558728","-73.9788465","-73.9787465"],"lat":"40.7558228","lon":"-73.9787965","display_name":"Barnes & Noble, 555, 5th Avenue, Midtown East, Manhattan, New York County, New York, 10017, United States","class":"shop","type":"books","importance":0.62001},{"place_id":319634907,"licence":"Data © OpenStreetMap contributors, ODbL 1.0. https://osm.org/copyright","osm_type":"node","osm_id":2716012085,"boundingbox":["40.7557517","40.7558517","-73.9787914","-73.9786914"],"lat":"40.7558017","lon":"-73.9787414","display_name":"555, 5th Avenue, Midtown East, Manhattan, New York County, New York, 10017, United States","class":"place","type":"house","importance":0.62001}]"`)).Decode(&want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Fatal(diff)
	}
}

func TestAddressGeocode(t *testing.T) {
	testClient := Client{
		ApiKey: APIKEY,
	}
	// https://geocode.maps.co/search?street=555+5th+Ave&city=New+York&state=NY&postalcode=10017&country=US&api_key=
	got, err := testClient.AddressGeocode("555 5th Ave", "New York", "", "NY", "US", "10017")
	if err != nil {
		t.Fatal("Geocode Failed:", err)
	}

	var want []*Response
	json.NewDecoder(strings.NewReader(`[{"place_id":319634989,"licence":"Data © OpenStreetMap contributors, ODbL 1.0. https://osm.org/copyright","osm_type":"node","osm_id":1000793154,"boundingbox":["40.7557728","40.7558728","-73.9788465","-73.9787465"],"lat":"40.7558228","lon":"-73.9787965","display_name":"Barnes & Noble, 555, 5th Avenue, Midtown East, Manhattan, New York County, New York, 10017, United States","class":"shop","type":"books","importance":0.62001},{"place_id":319634907,"licence":"Data © OpenStreetMap contributors, ODbL 1.0. https://osm.org/copyright","osm_type":"node","osm_id":2716012085,"boundingbox":["40.7557517","40.7558517","-73.9787914","-73.9786914"],"lat":"40.7558017","lon":"-73.9787414","display_name":"555, 5th Avenue, Midtown East, Manhattan, New York County, New York, 10017, United States","class":"place","type":"house","importance":0.62001}]"`)).Decode(&want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Fatal(diff)
	}
}

func TestReverse(t *testing.T) {
	testClient := Client{
		ApiKey: APIKEY,
	}
	// https://geocode.maps.co/reverse?lat=40.7558017&lon=-73.9787414&api_key=
	got, err := testClient.Reverse(40.7558017, -73.9787414)
	if err != nil {
		t.Fatalf("Reverse failed: %s", err)
	}

	var want = &Response{}
	json.NewDecoder(strings.NewReader(`{"place_id":319634907,"licence":"Data © OpenStreetMap contributors, ODbL 1.0. https://osm.org/copyright","osm_type":"node","osm_id":2716012085,"lat":"40.7558017","lon":"-73.9787414","display_name":"555, 5th Avenue, Midtown East, Manhattan, New York County, New York, 10017, United States","address":{"house_number":"555","road":"5th Avenue","neighbourhood":"Midtown East","suburb":"Manhattan","county":"New York County","city":"New York","state":"New York","ISO3166-2-lvl4":"US-NY","postcode":"10017","country":"United States","country_code":"us"},"boundingbox":["40.7557517","40.7558517","-73.9787914","-73.9786914"]}`)).Decode(want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Fatal(diff)
	}
}
