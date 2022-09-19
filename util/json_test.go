package util

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeJson(t *testing.T) {
	type person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	wantP := person{Name: "raymond", Age: 12}

	var p person
	r := strings.NewReader(`{"name":"raymond", "age":12}`)

	err := DecodeJson(r, &p)
	require.NoError(t, err)

	require.Equal(t, wantP, p)
}
