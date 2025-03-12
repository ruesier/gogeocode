# GoGeocode
simple go package for using the geocoding api at [geocode.maps.co]

## API Key
You can sign up for an api key at: [https://geocode.maps.co/join/]

## Geocode
Geocoding takes an address or similar location description, then returns a precise location description

```go
package main
import "github.com/ruesier/gogeocode"

func main() {
    client := gogeocode.Client{
        ApiKey: "GEOCODE_API_KEY"
    }

    response, err := client.Geocode("The Statue of Liberty")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("The Statue of Liberty is at Lat: %s, Long: %s", response[0].Latitude, response[0].Longitude)
}
```
Output: `The Statue of Liberty is at Lat: 40.689253199999996, Long: -74.04454817144321`

## Reverse Geocode
Reverse Geocoding takes a latitude and longitude then returns Addresses that are at that point.

```go
package main

func main() {
    client := gogeocode.Client{
        ApiKey: "GEOCODE_API_KEY"
    }

    response, err := client.Reverse(40.689253199999996, -74.04454817144321)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(respose.Address.Tourism)
}
```
Output: `Statue of Liberty`