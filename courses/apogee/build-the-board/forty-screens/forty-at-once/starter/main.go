package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

// main serves the Board on 127.0.0.1, on a port the system picks, and prints its address.
// In the building it would listen on the office network, with passplan and the weather
// service at the addresses its flags give. No lesson runs it.
func main() {
	passplan := flag.String("passplan", "", "passplan's address, such as http://127.0.0.1:8070")
	weather := flag.String("weather", "", "the weather service's address")
	flag.Parse()
	client := &http.Client{Timeout: 2 * time.Second}
	src := Sources{
		Passes:  schedule(client, *passplan),
		Weather: (&WeatherClient{BaseURL: *weather, HTTP: client}).Reading,
		Rota:    rota(),
	}
	b, err := NewBoard(src, stationTable(), time.Now)
	if err != nil {
		log.Fatal(err)
	}
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("The Board is at http://%s/\n", ln.Addr())
	go func() {
		for {
			if err := b.Refresh(context.Background()); err != nil {
				log.Print(err)
			}
			time.Sleep(15 * time.Second)
		}
	}()
	log.Fatal(http.Serve(ln, b.Handler()))
}

// WeatherClient asks the weather service for station readings.
type WeatherClient struct {
	BaseURL string       // the service's address
	HTTP    *http.Client // with a Timeout
}

// Reading returns the latest reading of the station with this code.
func (c *WeatherClient) Reading(ctx context.Context, code string) (Weather, error) {
	var w Weather
	if err := getJSON(ctx, c.HTTP, c.BaseURL+"/weather/"+code, &w); err != nil {
		return Weather{}, fmt.Errorf("weather for %s: %w", code, err)
	}
	return w, nil
}

// schedule returns a source that asks passplan at base for every booked pass.
func schedule(client *http.Client, base string) func(ctx context.Context) ([]Booking, error) {
	return func(ctx context.Context) ([]Booking, error) {
		var all []Booking
		if err := getJSON(ctx, client, base+"/schedule", &all); err != nil {
			return nil, fmt.Errorf("schedule: %w", err)
		}
		return all, nil
	}
}

// getJSON asks url for a JSON answer and decodes it into v.
func getJSON(ctx context.Context, client *http.Client, url string, v any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return errors.New(res.Status)
	}
	return json.NewDecoder(res.Body).Decode(v)
}
