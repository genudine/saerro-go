package util_test

import (
	"strconv"
	"testing"

	"github.com/genudine/saerro-go/util"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	dolls := []int64{44203, 41666, 79579, 63741, 57213}

	result := util.Map(dolls, func(doll int64) string {
		return strconv.FormatInt(doll, 16)
	})

	assert.Contains(t, result, "acab")
	assert.Len(t, result, len(dolls))
}
