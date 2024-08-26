package tests

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_vaults/config"
	"github.com/liobrdev/simplepasswords_vaults/controllers"
	"github.com/liobrdev/simplepasswords_vaults/models"
	"github.com/liobrdev/simplepasswords_vaults/tests/helpers"
	"github.com/liobrdev/simplepasswords_vaults/tests/setup"
	"github.com/liobrdev/simplepasswords_vaults/utils"
)

func testListVaults(t *testing.T, app *fiber.App, db *gorm.DB, conf *config.AppConfig) {
	t.Run("invalid_slug_400_bad_request", func(t *testing.T) {
		slug := "notEvenARealSlug"
		testListVaultsClientError(t, app, conf, 400, utils.ErrorUserSlug, slug, slug)
	})

	t.Run("valid_slug_200_ok", func(t *testing.T) {
		testListVaultsSuccess(t, app, db, conf)
	})
}

func testListVaultsClientError(
	t *testing.T, app *fiber.App, conf *config.AppConfig, expectedStatus int,
	expectedMessage, expectedDetail, slug string,
) {
	resp := newRequestListVaults(t, app, conf, slug)
	require.Equal(t, expectedStatus, resp.StatusCode)
	helpers.AssertErrorResponseBody(t, resp, utils.ErrorResponseBody{
		ClientOperation: utils.ListVaults,
		Message:         expectedMessage,
		Detail:          expectedDetail,
	})
}

func testListVaultsSuccess(t *testing.T, app *fiber.App, db *gorm.DB, conf *config.AppConfig) {
	users, vaults, _, _ := setup.SetUpWithData(t, db)
	slug := users[0].Slug
	resp := newRequestListVaults(t, app, conf, slug)
	require.Equal(t, 200, resp.StatusCode)

	if respBody, err := io.ReadAll(resp.Body); err != nil {
		t.Fatalf("Read response body failed: %s", err.Error())
	} else {
		var vaultsJSON []models.Vault

		if vaultsBytes, err := json.Marshal(vaults[:2]); err != nil {
			t.Fatalf("JSON marshal failed: %s", err.Error())
		} else if err := json.Unmarshal(vaultsBytes, &vaultsJSON); err != nil {
			t.Fatalf("JSON unmarshal failed: %s", err.Error())
		}

		var listVaultsRespBody controllers.ListVaultsResponseBody

		if err := json.Unmarshal(respBody, &listVaultsRespBody); err != nil {
			t.Fatalf("JSON unmarshal failed: %s", err.Error())
		}

		require.ElementsMatch(t, vaultsJSON, listVaultsRespBody.Vaults)
	}
}

func newRequestListVaults(
	t *testing.T, app *fiber.App, conf *config.AppConfig, slug string,
) *http.Response {

	req := httptest.NewRequest("GET", "/api/vaults", nil)
	req.Header.Set("Authorization", "Token " + conf.VAULTS_ACCESS_TOKEN)
	req.Header.Set("Client-Operation", utils.ListVaults)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Slug", slug)

	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
