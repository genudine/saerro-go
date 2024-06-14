package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventNames(t *testing.T) {
	result := getEventNames()
	assert.Contains(t, result, "GainExperience_experience_id_55")
	assert.Contains(t, result, "Death")
	assert.Contains(t, result, "VehicleDestroy")
}
