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

	"github.com/liobrdev/simplepasswords_vaults/models"
	"github.com/liobrdev/simplepasswords_vaults/tests/helpers"
	"github.com/liobrdev/simplepasswords_vaults/tests/setup"
	"github.com/liobrdev/simplepasswords_vaults/utils"
)

func testRetrieveEntry(t *testing.T, app *fiber.App, db *gorm.DB) {
	t.Run("invalid_slug_400_bad_request", func(t *testing.T) {
		slug := "notEvenARealSlug"
		testRetrieveEntryClientError(
			t, app, db, slug, http.StatusBadRequest, utils.ErrorEntrySlug, slug,
		)
	})

	t.Run("valid_slug_404_not_found", func(t *testing.T) {
		slug := helpers.NewSlug(t)
		testRetrieveEntryClientError(
			t, app, db, slug, http.StatusNotFound, utils.ErrorNotFound, slug,
		)
	})

	t.Run("valid_slug_200_ok", func(t *testing.T) {
		testRetrieveEntrySuccess(t, app, db)
	})
}

func testRetrieveEntryClientError(
	t *testing.T,
	app *fiber.App,
	db *gorm.DB,
	slug string,
	expectedStatus int,
	expectedMessage utils.ErrorMessage,
	expectedDetail string,
) {
	resp := newRequestRetrieveEntry(t, app, slug)
	require.Equal(t, expectedStatus, resp.StatusCode)
	helpers.AssertErrorResponseBody(t, resp, utils.ErrorResponseBody{
		ClientOperation: utils.RetrieveEntry,
		Message:         string(expectedMessage),
		Detail:          expectedDetail,
	})
}

func testRetrieveEntrySuccess(t *testing.T, app *fiber.App, db *gorm.DB) {
	setup.SetUpWithData(t, db)

	var expectedEntry models.Entry
	helpers.QueryTestEntry(t, db, &expectedEntry, "entry@0.1.1.*")

	resp := newRequestRetrieveEntry(t, app, expectedEntry.Slug)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	if respBody, err := io.ReadAll(resp.Body); err != nil {
		t.Fatalf("Read response body failed: %s", err.Error())
	} else {
		var actualEntry models.Entry

		if err := json.Unmarshal(respBody, &actualEntry); err != nil {
			t.Fatalf("JSON unmarshal failed: %s", err.Error())
		}

		require.Equal(t, expectedEntry.Slug, actualEntry.Slug)
		require.Equal(t, expectedEntry.UserSlug, actualEntry.UserSlug)
		require.Equal(t, expectedEntry.Title, actualEntry.Title)
		require.Equal(t, "entry@0.1.1.*", actualEntry.Title)
		require.Len(t, actualEntry.Secrets, 2)
		require.Equal(t, "secret[_label='username']@0.1.1.0", actualEntry.Secrets[0].Label)
		require.Equal(t, "secret[_string='foodeater1234']@0.1.1.0", actualEntry.Secrets[0].String)
		require.Equal(t, "secret[_label='password']@0.1.1.1", actualEntry.Secrets[1].Label)
		require.Equal(t, "secret[_string='3a7!ng40oD']@0.1.1.1", actualEntry.Secrets[1].String)
	}
}

func newRequestRetrieveEntry(t *testing.T, app *fiber.App, slug string) *http.Response {
	req := httptest.NewRequest(http.MethodGet, "/api/entries/"+slug, nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
