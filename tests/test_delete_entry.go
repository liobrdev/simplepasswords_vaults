package tests

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_vaults/models"
	"github.com/liobrdev/simplepasswords_vaults/tests/helpers"
	"github.com/liobrdev/simplepasswords_vaults/tests/setup"
	"github.com/liobrdev/simplepasswords_vaults/utils"
)

func testDeleteEntry(t *testing.T, app *fiber.App, db *gorm.DB) {
	t.Run("valid_slug_404_not_found", func(t *testing.T) {
		testDeleteEntryClientError(
			t, app, db, helpers.NewSlug(t), http.StatusNotFound, utils.ErrorNoRowsAffected,
			"Likely that slug was not found.",
		)
	})

	t.Run("invalid_slug_400_bad_request", func(t *testing.T) {
		slug := "notEvenARealSlug"
		testDeleteEntryClientError(
			t, app, db, slug, http.StatusBadRequest, utils.ErrorEntrySlug, slug,
		)
	})

	t.Run("valid_slug_204_no_content", func(t *testing.T) {
		testDeleteEntrySuccess(t, app, db)
	})
}

func testDeleteEntryClientError(
	t *testing.T, app *fiber.App, db *gorm.DB, slug string, expectedStatus int,
	expectedMessage string, expectedDetail string,
) {
	resp := newRequestDeleteEntry(t, app, slug)
	require.Equal(t, expectedStatus, resp.StatusCode)
	helpers.AssertErrorResponseBody(t, resp, utils.ErrorResponseBody{
		ClientOperation: utils.DeleteEntry,
		Message:         expectedMessage,
		Detail:          expectedDetail,
	})
}

func testDeleteEntrySuccess(t *testing.T, app *fiber.App, db *gorm.DB) {
	setup.SetUpWithData(t, db)

	var entry models.Entry
	helpers.QueryTestEntryEager(t, db, &entry, "entry@0.1.1.*")
	require.Len(t, entry.Secrets, 2)

	secret1 := entry.Secrets[0]
	secret2 := entry.Secrets[1]
	require.Equal(t, "secret[_label='username']@0.1.1.0", secret1.Label)
	require.Equal(t, "secret[_string='foodeater1234']@0.1.1.0", secret1.String)
	require.Equal(t, "secret[_label='password']@0.1.1.1", secret2.Label)
	require.Equal(t, "secret[_string='3a7!ng40oD']@0.1.1.1", secret2.String)

	var entryCount int64
	helpers.CountEntries(t, db, &entryCount)
	require.EqualValues(t, 8, entryCount)

	var secretCount int64
	helpers.CountSecrets(t, db, &secretCount)
	require.EqualValues(t, 16, secretCount)

	resp := newRequestDeleteEntry(t, app, entry.Slug)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)

	if respBody, err := io.ReadAll(resp.Body); err != nil {
		t.Fatalf("Read response body failed: %s", err.Error())
	} else {
		require.Empty(t, respBody)
	}

	if result := db.First(&entry, "slug = ?", entry.Slug); result.Error != nil {
		require.ErrorIs(t, result.Error, gorm.ErrRecordNotFound)
	} else {
		t.Fatalf("Deleted entry query failed: %s", result.Error.Error())
	}

	if result := db.First(&secret1, "slug = ?", secret1.Slug); result.Error != nil {
		require.ErrorIs(t, result.Error, gorm.ErrRecordNotFound)
	} else {
		t.Fatalf("Deleted secret1 query failed: %s", result.Error.Error())
	}

	if result := db.First(&secret2, "slug = ?", secret2.Slug); result.Error != nil {
		require.ErrorIs(t, result.Error, gorm.ErrRecordNotFound)
	} else {
		t.Fatalf("Deleted secret2 query failed: %s", result.Error.Error())
	}

	helpers.CountEntries(t, db, &entryCount)
	require.EqualValues(t, 7, entryCount)

	helpers.CountSecrets(t, db, &secretCount)
	require.EqualValues(t, 14, secretCount)
}

func newRequestDeleteEntry(t *testing.T, app *fiber.App, slug string) *http.Response {
	req := httptest.NewRequest(http.MethodDelete, "/api/entries/"+slug, nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
