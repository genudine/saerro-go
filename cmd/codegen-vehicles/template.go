package main

import (
	"bytes"
	"fmt"
	"text/template"
)

var (
	vehicleMapTmpl = `package translators

var (
	VehicleMap = map[string]Vehicle{
		{{ range .Vehicles }}"{{ .VehicleID }}": {{ .VehicleEnumName }},
		{{end}}
	}
)`

	vehicleMapTemplate = template.Must(template.New("vehicle_map").Parse(vehicleMapTmpl))
)

type TemplateData struct {
	Vehicles []VehicleItem
}

type VehicleItem struct {
	VehicleID       string
	VehicleEnumName string
}

func renderTemplate(data TemplateData) (string, error) {
	var buffer bytes.Buffer
	err := vehicleMapTemplate.Execute(&buffer, data)
	if err != nil {
		return "", fmt.Errorf("template render failed, %w", err)
	}

	return buffer.String(), nil
}
