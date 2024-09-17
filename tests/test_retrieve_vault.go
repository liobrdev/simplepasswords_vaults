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
	"github.com/liobrdev/simplepasswords_vaults/models"
	"github.com/liobrdev/simplepasswords_vaults/tests/helpers"
	"github.com/liobrdev/simplepasswords_vaults/tests/setup"
	"github.com/liobrdev/simplepasswords_vaults/utils"
)

func testRetrieveVault(t *testing.T, app *fiber.App, db *gorm.DB, conf *config.AppConfig) {
	t.Run("invalid_slug_400_bad_request", func(t *testing.T) {
		slug := "notARealSlug"
		testRetrieveVaultClientError(t, app, conf, 400, utils.ErrorVaultSlug, slug, slug)
	})

	t.Run("valid_slug_404_not_found", func(t *testing.T) {
		slug := helpers.NewSlug(t)
		testRetrieveVaultClientError(t, app, conf, 404, utils.ErrorNotFound, slug, slug)
	})

	t.Run("valid_slug_200_ok", func(t *testing.T) {
		testRetrieveVaultSuccess(t, app, db, conf)
	})
}

func testRetrieveVaultClientError(
	t *testing.T, app *fiber.App, conf *config.AppConfig,
	expectedStatus int, expectedMessage, expectedDetail, slug string,
) {
	resp := newRequestRetrieveVault(t, app, conf, slug)
	require.Equal(t, expectedStatus, resp.StatusCode)
	helpers.AssertErrorResponseBody(t, resp, utils.ErrorResponseBody{
		ClientOperation: utils.RetrieveVault,
		Message:         expectedMessage,
		Detail:          expectedDetail,
	})
}

func testRetrieveVaultSuccess(t *testing.T, app *fiber.App, db *gorm.DB, conf *config.AppConfig) {
	_, _, entries, _ := setup.SetUpWithData(t, db)

	var expectedVault models.Vault
	helpers.QueryTestVault(t, db, &expectedVault, "vault@0.1.*.*")

	resp := newRequestRetrieveVault(t, app, conf, expectedVault.Slug)
	require.Equal(t, 200, resp.StatusCode)

	if respBody, err := io.ReadAll(resp.Body); err != nil {
		t.Fatalf("Read response body failed: %s", err.Error())
	} else {
		var actualVault models.Vault

		if err := json.Unmarshal(respBody, &actualVault); err != nil {
			t.Fatalf("JSON unmarshal failed: %s", err.Error())
		}

		var entriesJSON []models.Entry

		if entriesBytes, err := json.Marshal(entries[2:4]); err != nil {
			t.Fatalf("JSON marshal failed: %s", err.Error())
		} else if err := json.Unmarshal(entriesBytes, &entriesJSON); err != nil {
			t.Fatalf("JSON unmarshal failed: %s", err.Error())
		}

		require.Equal(t, expectedVault.Slug, actualVault.Slug)
		require.Equal(t, expectedVault.Title, actualVault.Title)
		require.Equal(t, "vault@0.1.*.*", actualVault.Title)
		require.ElementsMatch(t, entriesJSON, actualVault.Entries)
		require.True(t, actualVault.Entries[0].CreatedAt.After(actualVault.Entries[1].CreatedAt))
	}
}

func newRequestRetrieveVault(
	t *testing.T, app *fiber.App, conf *config.AppConfig, slug string,
) *http.Response {

	req := httptest.NewRequest(http.MethodGet, "/api/vaults/" + slug, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Client-Operation", utils.RetrieveVault)
	req.Header.Set("Authorization", "Token " + conf.VAULTS_ACCESS_TOKEN)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
