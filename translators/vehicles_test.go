package translators_test

import (
	"testing"

	"github.com/genudine/saerro-go/pkg/translators"
	"github.com/stretchr/testify/assert"
)

func TestVehicles(t *testing.T) {
	assert.Equal(t, translators.VehicleNameFromID("12"), translators.Harasser)
	assert.Equal(t, translators.VehicleNameFromID("0"), translators.Vehicle("unknown"))
}
