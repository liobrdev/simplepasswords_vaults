package tests

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_vaults/models"
	"github.com/liobrdev/simplepasswords_vaults/tests/helpers"
	"github.com/liobrdev/simplepasswords_vaults/tests/setup"
	"github.com/liobrdev/simplepasswords_vaults/utils"
)

func testUpdateEntry(t *testing.T, app *fiber.App, db *gorm.DB) {
	t.Run("empty_body_400_bad_request", func(t *testing.T) {
		testUpdateEntryClientError(
			t, app, db, helpers.NewSlug(t), "", http.StatusBadRequest, utils.ErrorParse,
			"invalid character '\x00' looking for beginning of value",
		)
	})

	t.Run("array_body_400_bad_request", func(t *testing.T) {
		testUpdateEntryClientError(
			t, app, db, helpers.NewSlug(t), "[]", http.StatusBadRequest, utils.ErrorParse,
			"invalid character '[' looking for beginning of value",
		)

		testUpdateEntryClientError(
			t, app, db, helpers.NewSlug(t), "[{}]", http.StatusBadRequest, utils.ErrorParse,
			"invalid character '[' looking for beginning of value",
		)

		testUpdateEntryClientError(
			t, app, db, helpers.NewSlug(t), "[{\"entry_title\":\"updated@0.1.1.*\"}]",
			http.StatusBadRequest, utils.ErrorParse,
			"invalid character '[' looking for beginning of value",
		)
	})

	t.Run("null_body_400_bad_request", func(t *testing.T) {
		testUpdateEntryClientError(
			t, app, db, helpers.NewSlug(t), "null",
			http.StatusBadRequest, utils.ErrorEntryTitle, "",
		)
	})

	t.Run("boolean_body_400_bad_request", func(t *testing.T) {
		testUpdateEntryClientError(
			t, app, db, helpers.NewSlug(t), "true", http.StatusBadRequest, utils.ErrorParse,
			"invalid character 't' looking for beginning of value",
		)

		testUpdateEntryClientError(
			t, app, db, helpers.NewSlug(t), "false", http.StatusBadRequest, utils.ErrorParse,
			"invalid character 'f' looking for beginning of value",
		)
	})

	t.Run("string_body_400_bad_request", func(t *testing.T) {
		testUpdateEntryClientError(
			t, app, db, helpers.NewSlug(t), "\"Valid JSON, but not an object.\"",
			http.StatusBadRequest, utils.ErrorParse,
			"invalid character '\"' looking for beginning of value",
		)
	})

	t.Run("empty_object_body_400_bad_request", func(t *testing.T) {
		testUpdateEntryClientError(
			t, app, db, helpers.NewSlug(t), "{}",
			http.StatusBadRequest, utils.ErrorEntryTitle, "",
		)
	})

	t.Run("missing_entry_title_400_bad_request", func(t *testing.T) {
		testUpdateEntryClientError(
			t, app, db, helpers.NewSlug(t), `{"enrty_title":"Spelled wrong!"}`,
			http.StatusBadRequest, utils.ErrorEntryTitle, "",
		)
	})

	t.Run("null_entry_title_400_bad_request", func(t *testing.T) {
		testUpdateEntryClientError(
			t, app, db, helpers.NewSlug(t), `{"entry_title":null}`,
			http.StatusBadRequest, utils.ErrorEntryTitle, "",
		)
	})

	t.Run("empty_entry_title_400_bad_request", func(t *testing.T) {
		testUpdateEntryClientError(
			t, app, db, helpers.NewSlug(t), `{"entry_title":""}`,
			http.StatusBadRequest, utils.ErrorEntryTitle, "",
		)
	})

	t.Run("too_long_entry_title_400_bad_request", func(t *testing.T) {
		// `title` is a random string greater than 255 characters in length
		if title, err := utils.GenerateSlug(256); err != nil {
			t.Fatalf("Generate long string failed: %s", err.Error())
		} else {
			testUpdateEntryClientError(
				t, app, db, helpers.NewSlug(t), fmt.Sprintf(`{"entry_title":"%s"}`, title),
				http.StatusBadRequest, utils.ErrorEntryTitle, "Too long (256 > 255)",
			)
		}
	})

	t.Run("valid_body_entry_title_already_exists_409_conflict", func(t *testing.T) {
		_, _, entries, _ := setup.SetUpWithData(t, db)

		testUpdateEntryClientError(
			t, app, db, (*entries)[0].Slug,
			fmt.Sprintf(`{"entry_title":"%s"}`, (*entries)[1].Title),
			http.StatusConflict, utils.ErrorFailedDB,
			"UNIQUE constraint failed: entries.title, entries.vault_slug",
		)
	})

	t.Run("valid_body_404_not_found", func(t *testing.T) {
		setup.SetUpWithData(t, db)
		testUpdateEntryClientError(
			t, app, db, helpers.NewSlug(t), `{"entry_title":"updated@0.1.1.*"}`,
			http.StatusNotFound, utils.ErrorNoRowsAffected, "Likely that slug was not found.",
		)
	})

	t.Run("valid_body_204_no_content", func(t *testing.T) {
		updatedEntryTitle := "updated@0.1.1.*"

		testUpdateEntrySuccess(t, app, db, updatedEntryTitle, fmt.Sprintf(
			`{"entry_title":"%s"}`,
			updatedEntryTitle,
		))
	})

	t.Run("valid_body_irrelevant_data_204_no_content", func(t *testing.T) {
		updatedEntryTitle := "updated@0.1.1.*"

		validBodyIrrelevantData := "{" +
			fmt.Sprintf(`"entry_title":"%s",`, updatedEntryTitle) +
			`"entry_slug":"notEvenARealSlug",` +
			`"entry_created_at":"10/12/22"` +
			"}"

		testUpdateEntrySuccess(t, app, db, updatedEntryTitle, validBodyIrrelevantData)
	})
}

func testUpdateEntryClientError(
	t *testing.T,
	app *fiber.App,
	db *gorm.DB,
	slug string,
	body string,
	expectedStatus int,
	expectedMessage utils.ErrorMessage,
	expectedDetail string,
) {
	resp := newRequestUpdateEntry(t, app, slug, body)
	require.Equal(t, expectedStatus, resp.StatusCode)
	helpers.AssertErrorResponseBody(t, resp, utils.ErrorResponseBody{
		ClientOperation: utils.UpdateEntry,
		Message:         string(expectedMessage),
		Detail:          expectedDetail,
		RequestBody:     body,
	})
}

func testUpdateEntrySuccess(
	t *testing.T,
	app *fiber.App,
	db *gorm.DB,
	updatedEntryTitle string,
	body string,
) {
	setup.SetUpWithData(t, db)
	var entryBeforeUpdate models.Entry
	helpers.QueryTestEntry(t, db, &entryBeforeUpdate, "entry@0.1.1.*")
	resp := newRequestUpdateEntry(t, app, entryBeforeUpdate.Slug, body)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)

	if respBody, err := io.ReadAll(resp.Body); err != nil {
		t.Fatalf("Read response body failed: %s", err.Error())
	} else {
		require.Empty(t, respBody)
	}

	var entryAfterUpdate models.Entry
	helpers.QueryTestEntry(t, db, &entryAfterUpdate, updatedEntryTitle)
	require.Equal(t, entryBeforeUpdate.Slug, entryAfterUpdate.Slug)
	require.NotEqual(t, entryBeforeUpdate.Title, entryAfterUpdate.Title)
	require.Equal(t, updatedEntryTitle, entryAfterUpdate.Title)
	require.Equal(t, entryBeforeUpdate.CreatedAt, entryAfterUpdate.CreatedAt)
	require.True(t, entryAfterUpdate.UpdatedAt.After(entryAfterUpdate.CreatedAt))
	require.True(t, entryAfterUpdate.UpdatedAt.After(entryBeforeUpdate.UpdatedAt))
}

func newRequestUpdateEntry(
	t *testing.T,
	app *fiber.App,
	slug string,
	body string,
) *http.Response {
	reqBody := strings.NewReader(body)
	req := httptest.NewRequest(http.MethodPatch, "/api/entries/"+slug, reqBody)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
