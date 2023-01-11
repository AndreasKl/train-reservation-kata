package api

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReferenceDoesSerialize(t *testing.T) {
	response, err := json.Marshal(referenceResponse{})
	require.NoError(t, err)

	fmt.Printf("JSON: %s", string(response))
}
