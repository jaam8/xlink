package tests

import (
	"testing"
	"xlink/user_service/internal/utils"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	q := utils.GenerateToken(int8(16))
	w := utils.GenerateToken(int8(16))
	e := utils.GenerateToken(int8(16))
	r := utils.GenerateToken(int8(16))
	p := utils.GenerateToken(int8(16))
	y := utils.GenerateToken(int8(16))
	u := utils.GenerateToken(int8(16))
	i := utils.GenerateToken(int8(16))
	assert.NotEqual(t, q, w, e, r, y, u, i, p)
}
