// / Generate pkgs/translators/vehicles_map.gen.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type VehicleResponse struct {
	VehicleList []Vehicle `json:"vehicle_list"`
}

type Vehicle struct {
	VehicleID string       `json:"vehicle_id"`
	Name      LocaleString `json:"name"`
}

type LocaleString struct {
	En string `json:"en"`
}

func fetchCensusVehicles() (VehicleResponse, error) {
	var vehicleResponse VehicleResponse

	client := http.Client{
		Timeout: time.Second * 30,
	}

	resp, err := client.Get("https://census.lithafalcon.cc/get/ps2/vehicle")
	if err != nil {
		return vehicleResponse, fmt.Errorf("census request failed: %w", err)
	}

	err = json.NewDecoder(resp.Body).Decode(&vehicleResponse)
	if err != nil {
		return vehicleResponse, fmt.Errorf("census response decode failed: %w", err)
	}

	return vehicleResponse, nil
}

func generateRegexp(vehicles []string) *regexp.Regexp {
	pipes := strings.Join(vehicles, "|")
	expr := fmt.Sprintf("(%s)", pipes)

	log.Println(expr)

	return regexp.MustCompile(expr)
}

func main() {
	filterRegexp := generateRegexp(AllVehicles)

	censusVehicles, err := fetchCensusVehicles()
	if err != nil {
		log.Fatalln("fetch census failed", err)
	}

	vehicles := []VehicleItem{}
	for _, vehicle := range censusVehicles.VehicleList {
		if vehicle.Name.En == "" || strings.Contains(vehicle.Name.En, "Turret") {
			continue
		}

		match := filterRegexp.FindString(strings.ToLower(vehicle.Name.En))
		if match == "" {
			continue
		}

		switch match {
		case "wasp":
			match = "valkyrie"
		case "deliverer":
			match = "ant"
		case "lodestar":
			match = "galaxy"
		}

		enumName := fmt.Sprintf("%s%s", strings.ToUpper(match[0:1]), match[1:])

		vehicles = append(vehicles, VehicleItem{
			VehicleID:       vehicle.VehicleID,
			VehicleEnumName: enumName,
		})
	}

	output, err := renderTemplate(TemplateData{
		Vehicles: vehicles,
	})
	if err != nil {
		log.Fatalln("render failed", err)
	}

	fmt.Println(output)
}
