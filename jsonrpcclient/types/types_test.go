package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarshalArgUint64(t *testing.T) {
	var value ArgUint64
	value = 1
	a, err := json.Marshal(value)
	require.NoError(t, err)
	require.Equal(t, "\"0x1\"", string(a))
}

func TestUnMarshalArgUint64(t *testing.T) {
	var value ArgUint64
	value = 1
	err := json.Unmarshal([]byte("\"0x1\""), &value)
	require.NoError(t, err)

}
