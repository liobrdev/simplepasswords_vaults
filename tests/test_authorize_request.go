package tests

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"

	"github.com/liobrdev/simplepasswords_vaults/config"
	"github.com/liobrdev/simplepasswords_vaults/tests/helpers"
	"github.com/liobrdev/simplepasswords_vaults/utils"
)

func testAuthorizeRequest(t *testing.T, app *fiber.App, conf *config.AppConfig) {
	var dummyToken string
	var err error

	if dummyToken, err = utils.GenerateSlug(80); err != nil {
		t.Fatalf("Generate dummyToken failed: %s", err.Error())
	}

	t.Run("null_empty_auth_header_401_unauthorized", func(t *testing.T) {
		testAuthorizeRequestClientError(t, app, 401, utils.ErrorToken, "", "")
		testAuthorizeRequestClientError(t, app, 401, utils.ErrorToken, "null", "null")
	})

	t.Run("null_empty_token_401_unauthorized", func(t *testing.T) {
		testAuthorizeRequestClientError(t, app, 401, utils.ErrorToken, "Token null", "Token null")
		testAuthorizeRequestClientError(t, app, 401, utils.ErrorToken, "token null", "token null")
		testAuthorizeRequestClientError(t, app, 401, utils.ErrorToken, "Token", "Token ")
		testAuthorizeRequestClientError(t, app, 401, utils.ErrorToken, "token", "token ")
	})

	t.Run("invalid_token_regexp_401_unauthorized", func(t *testing.T) {
		testAuthorizeRequestClientError(
			t, app, 401, utils.ErrorToken, "T0ken " + dummyToken, "T0ken " + dummyToken,
		)

		testAuthorizeRequestClientError(
			t, app, 401, utils.ErrorToken, "Bearer " + dummyToken, "Bearer " + dummyToken,
		)

		testAuthorizeRequestClientError(t, app, 401, utils.ErrorToken, dummyToken, dummyToken)

		testAuthorizeRequestClientError(
			t, app, 401, utils.ErrorToken, "Token " + dummyToken[:79], "Token " + dummyToken[:79],
		)

		testAuthorizeRequestClientError(
			t, app, 401, utils.ErrorToken, "Token a " + dummyToken[2:80], "Token a " + dummyToken[2:80],
		)

		testAuthorizeRequestClientError(
			t, app, 401, utils.ErrorToken,
			"Token " + dummyToken[:79] + "!", "Token " + dummyToken[:79] + "!",
		)
	})

	t.Run("valid_token_no_match_401_unauthorized", func(t *testing.T) {
		testAuthorizeRequestClientError(
			t, app, 401, utils.ErrorToken, dummyToken, "Token " + dummyToken,
		)
	})

	t.Run("valid_token_204_no_content", func(t *testing.T) {
		testAuthorizeRequestSuccess(t, app, "Token " + conf.VAULTS_ACCESS_TOKEN)
	})
}

func testAuthorizeRequestClientError(
	t *testing.T, app *fiber.App, expectedStatus int,
	expectedMessage, expectedDetail, authHeader string,
) {
	resp := newRequestAuthorizeRequest(t, app, authHeader)
	require.Equal(t, expectedStatus, resp.StatusCode)
	helpers.AssertErrorResponseBody(t, resp, utils.ErrorResponseBody{
		ClientOperation: utils.TestAuthReq,
		Message:         expectedMessage,
		Detail:          expectedDetail,
	})
}

func testAuthorizeRequestSuccess(t *testing.T, app *fiber.App, authHeader string) {
	resp := newRequestAuthorizeRequest(t, app, authHeader)
	require.Equal(t, 204, resp.StatusCode)

	if respBody, err := io.ReadAll(resp.Body); err != nil {
		t.Fatalf("Read response body failed: %s", err.Error())
	} else {
		require.Empty(t, respBody)
	}
}

func newRequestAuthorizeRequest(t *testing.T, app *fiber.App, authHeader string) *http.Response {
	req := httptest.NewRequest("GET", "/api/restricted", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Client-Operation", utils.TestAuthReq)
	req.Header.Set("Authorization", authHeader)

	resp, err := app.Test(req, -1)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
