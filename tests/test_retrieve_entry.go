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

func testRetrieveEntry(t *testing.T, app *fiber.App, db *gorm.DB, conf *config.AppConfig) {
	t.Run("invalid_slug_400_bad_request", func(t *testing.T) {
		slug := "notEvenARealSlug"
		testRetrieveEntryClientError(t, app, conf, 400, utils.ErrorEntrySlug, slug, slug)
	})

	t.Run("valid_slug_404_not_found", func(t *testing.T) {
		slug := helpers.NewSlug(t)
		testRetrieveEntryClientError(t, app, conf, 404, utils.ErrorNotFound, slug, slug)
	})

	t.Run("valid_slug_200_ok", func(t *testing.T) {
		testRetrieveEntrySuccess(t, app, db, conf)
	})
}

func testRetrieveEntryClientError(
	t *testing.T, app *fiber.App, conf *config.AppConfig, expectedStatus int,
	expectedMessage, expectedDetail, slug string,
) {
	resp := newRequestRetrieveEntry(t, app, conf, slug)
	require.Equal(t, expectedStatus, resp.StatusCode)
	helpers.AssertErrorResponseBody(t, resp, utils.ErrorResponseBody{
		ClientOperation: utils.RetrieveEntry,
		Message:         expectedMessage,
		Detail:          expectedDetail,
	})
}

func testRetrieveEntrySuccess(t *testing.T, app *fiber.App, db *gorm.DB, conf *config.AppConfig) {
	_, _, _, secrets := setup.SetUpWithData(t, db)

	var expectedEntry models.Entry
	helpers.QueryTestEntry(t, db, &expectedEntry, "entry@0.1.1.*")

	resp := newRequestRetrieveEntry(t, app, conf, expectedEntry.Slug)
	require.Equal(t, 200, resp.StatusCode)

	if respBody, err := io.ReadAll(resp.Body); err != nil {
		t.Fatalf("Read response body failed: %s", err.Error())
	} else {
		var actualEntry models.Entry

		if err := json.Unmarshal(respBody, &actualEntry); err != nil {
			t.Fatalf("JSON unmarshal failed: %s", err.Error())
		}

		require.Equal(t, expectedEntry.Slug, actualEntry.Slug)
		require.Equal(t, expectedEntry.Title, actualEntry.Title)
		require.Equal(t, "entry@0.1.1.*", actualEntry.Title)

		var secretsJSON []models.Secret

		if secretsBytes, err := json.Marshal(secrets[6:8]); err != nil {
			t.Fatalf("JSON marshal failed: %s", err.Error())
		} else if err := json.Unmarshal(secretsBytes, &secretsJSON); err != nil {
			t.Fatalf("JSON unmarshal failed: %s", err.Error())
		}

		require.ElementsMatch(t, secretsJSON, actualEntry.Secrets)
		require.Less(t, actualEntry.Secrets[0].Priority, actualEntry.Secrets[1].Priority)
		require.Equal(t, "secret[_label='username']@0.1.1.0", actualEntry.Secrets[0].Label)
		require.Equal(t, "secret[_string='foodeater1234']@0.1.1.0", actualEntry.Secrets[0].String)
		require.Equal(t, "secret[_label='password']@0.1.1.1", actualEntry.Secrets[1].Label)
		require.Equal(t, "secret[_string='3a7!ng40oD']@0.1.1.1", actualEntry.Secrets[1].String)
	}
}

func newRequestRetrieveEntry(
	t *testing.T, app *fiber.App, conf *config.AppConfig, slug string,
) *http.Response {

	req := httptest.NewRequest("GET", "/api/entries/" + slug, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Client-Operation", utils.RetrieveEntry)
	req.Header.Set("Authorization", "Token " + conf.VAULTS_ACCESS_TOKEN)

	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
