package helpers

import (
	"io"
	"net/http"
	"testing"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/require"

	"github.com/liobrdev/simplepasswords_vaults/utils"
)

func AssertErrorResponseBody(
	t *testing.T,
	resp *http.Response,
	expected utils.ErrorResponseBody,
) {
	if respBody, err := io.ReadAll(resp.Body); err != nil {
		t.Fatalf("Read response body failed: %s", err.Error())
	} else {
		errRespBody := utils.ErrorResponseBody{}

		if err := json.Unmarshal(respBody, &errRespBody); err != nil {
			t.Fatalf("JSON unmarshal error response body failed: %s", err.Error())
		}

		require.Equal(t, expected.Message, errRespBody.Message)
		require.Equal(t, expected.Detail, errRespBody.Detail)
		require.Equal(t, expected.RequestBody, errRespBody.RequestBody)
		require.Equal(t, expected.ClientOperation, errRespBody.ClientOperation)
	}
}
