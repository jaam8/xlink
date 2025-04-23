package tests

import (
	"testing"
	"xlink/shortener/internal/service/utils"

	"github.com/stretchr/testify/assert"
)

func TestGenerateShortURL(t *testing.T) {
	a := utils.GenerateShortURL()
	b := utils.GenerateShortURL()
	c := utils.GenerateShortURL()
	d := utils.GenerateShortURL()
	e := utils.GenerateShortURL()
	assert.NotEqual(t, a, b, c, d, e)
}
