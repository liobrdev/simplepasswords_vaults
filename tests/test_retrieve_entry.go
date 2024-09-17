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
	setup.SetUpWithData(t, db)

	var entryFromDB models.Entry
	helpers.QueryTestEntryEager(t, db, &entryFromDB, "entry@0.1.1.*")

	resp := newRequestRetrieveEntry(t, app, conf, entryFromDB.Slug)
	require.Equal(t, 200, resp.StatusCode)

	if respBody, err := io.ReadAll(resp.Body); err != nil {
		t.Fatalf("Read response body failed: %s", err.Error())
	} else {
		var respEntry models.Entry

		if err := json.Unmarshal(respBody, &respEntry); err != nil {
			t.Fatalf("JSON unmarshal failed: %s", err.Error())
		}

		require.Equal(t, entryFromDB.Slug, respEntry.Slug)
		require.Equal(t, entryFromDB.Title, respEntry.Title)
		require.Equal(t, "entry@0.1.1.*", respEntry.Title)

		if plaintext, err := utils.Decrypt(entryFromDB.Secrets[0].String, helpers.HexHash[:64]);
		err != nil {
			t.Fatalf("Password decryption failed: %s", err.Error())
		} else {
			require.Equal(t, "secret[_string='foodeater1234']@0.1.1.0", plaintext)
			require.Equal(t, plaintext, respEntry.Secrets[0].String)
			require.Equal(t, "secret[_label='username']@0.1.1.0", respEntry.Secrets[0].Label)
			require.EqualValues(t, 0, respEntry.Secrets[0].Priority)
		}

		if plaintext, err := utils.Decrypt(entryFromDB.Secrets[1].String, helpers.HexHash[:64]);
		err != nil {
			t.Fatalf("Password decryption failed: %s", err.Error())
		} else {
			require.Equal(t, "secret[_string='3a7!ng40oD']@0.1.1.1", plaintext)
			require.Equal(t, plaintext, respEntry.Secrets[1].String)
			require.Equal(t, "secret[_label='password']@0.1.1.1", respEntry.Secrets[1].Label)
			require.EqualValues(t, 1, respEntry.Secrets[1].Priority)
		}
	}
}

func newRequestRetrieveEntry(
	t *testing.T, app *fiber.App, conf *config.AppConfig, slug string,
) *http.Response {

	req := httptest.NewRequest("GET", "/api/entries/" + slug, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Client-Operation", utils.RetrieveEntry)
	req.Header.Set("Authorization", "Token " + conf.VAULTS_ACCESS_TOKEN)
	req.Header.Set(conf.PASSWORD_HEADER_KEY, helpers.HexHash[:64])

	resp, err := app.Test(req, -1)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
