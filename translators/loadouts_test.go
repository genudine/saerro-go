package translators_test

import (
	"testing"

	"github.com/genudine/saerro-go/translators"
	"github.com/stretchr/testify/assert"
)

func TestLoadouts(t *testing.T) {
	assert.Equal(t, translators.ClassFromLoadout(1), translators.Infiltrator)
	assert.Equal(t, translators.ClassFromLoadout(0), translators.Class("unknown"))
}
