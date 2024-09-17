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

	"github.com/liobrdev/simplepasswords_vaults/config"
	"github.com/liobrdev/simplepasswords_vaults/models"
	"github.com/liobrdev/simplepasswords_vaults/tests/helpers"
	"github.com/liobrdev/simplepasswords_vaults/tests/setup"
	"github.com/liobrdev/simplepasswords_vaults/utils"
)

func testUpdateEntry(t *testing.T, app *fiber.App, db *gorm.DB, conf *config.AppConfig) {
	bodyFmt := `{"entry_title":"%s"}`
	dummySlug := helpers.NewSlug(t)

	t.Run("empty_body_400_bad_request", func(t *testing.T) {
		testUpdateEntryClientError(
			t, app, conf, 400, utils.ErrorParse,
			"invalid character '\x00' looking for beginning of value", dummySlug, "",
		)
	})

	t.Run("array_body_400_bad_request", func(t *testing.T) {
		testUpdateEntryClientError(
			t, app, conf, 400, utils.ErrorParse, "invalid character '[' looking for beginning of value",
			dummySlug, "[]",
		)

		testUpdateEntryClientError(
			t, app, conf, 400, utils.ErrorParse, "invalid character '[' looking for beginning of value",
			dummySlug, "[{}]",
		)

		testUpdateEntryClientError(
			t, app, conf, 400, utils.ErrorParse, "invalid character '[' looking for beginning of value",
			dummySlug, `[{"entry_title":"updated@0.1.1.*"}]`,
		)
	})

	t.Run("boolean_body_400_bad_request", func(t *testing.T) {
		testUpdateEntryClientError(
			t, app, conf, 400, utils.ErrorParse, "invalid character 't' looking for beginning of value",
			dummySlug, "true",
		)

		testUpdateEntryClientError(
			t, app, conf, 400, utils.ErrorParse, "invalid character 'f' looking for beginning of value",
			dummySlug, "false",
		)
	})

	t.Run("string_body_400_bad_request", func(t *testing.T) {
		testUpdateEntryClientError(
			t, app, conf, 400, utils.ErrorParse, `invalid character '"' looking for beginning of value`,
			dummySlug, `"Valid JSON, but not an object."`,
		)
	})

	t.Run("null_body_400_bad_request", func(t *testing.T) {
		testUpdateEntryClientError(t, app, conf, 400, utils.ErrorEntryTitle, "", dummySlug, "null")
	})

	t.Run("empty_object_body_400_bad_request", func(t *testing.T) {
		testUpdateEntryClientError(t, app, conf, 400, utils.ErrorEntryTitle, "", dummySlug, "{}")
	})

	t.Run("missing_entry_title_400_bad_request", func(t *testing.T) {
		testUpdateEntryClientError(
			t, app, conf, 400, utils.ErrorEntryTitle, "", dummySlug, `{"enrty_title":"Spelled wrong!"}`,
		)
	})

	t.Run("null_entry_title_400_bad_request", func(t *testing.T) {
		testUpdateEntryClientError(
			t, app, conf, 400, utils.ErrorEntryTitle, "", dummySlug, `{"entry_title":null}`,
		)
	})

	t.Run("empty_entry_title_400_bad_request", func(t *testing.T) {
		testUpdateEntryClientError(
			t, app, conf, 400, utils.ErrorEntryTitle, "", dummySlug, `{"entry_title":""}`,
		)
	})

	t.Run("too_long_entry_title_400_bad_request", func(t *testing.T) {
		if title, err := utils.GenerateSlug(256); err != nil {
			t.Fatalf("Generate long string failed: %s", err.Error())
		} else {
			testUpdateEntryClientError(
				t, app, conf, 400, utils.ErrorEntryTitle, "Too long", dummySlug,
				fmt.Sprintf(bodyFmt, title),
			)
		}
	})

	t.Run("valid_body_entry_title_already_exists_409_conflict", func(t *testing.T) {
		_, _, entries, _ := setup.SetUpWithData(t, db)

		testUpdateEntryClientError(
			t, app, conf, 500, utils.ErrorFailedDB,
			"UNIQUE constraint failed: entries.title, entries.vault_slug",
			entries[0].Slug, fmt.Sprintf(bodyFmt, entries[1].Title),
		)
	})

	t.Run("valid_body_404_not_found", func(t *testing.T) {
		setup.SetUpWithData(t, db)

		testUpdateEntryClientError(
			t, app, conf, 404, utils.ErrorNoRowsAffected, "Likely that slug was not found.",
			dummySlug, fmt.Sprintf(bodyFmt, "updated@0.1.1.*"),
		)
	})

	t.Run("valid_body_204_no_content", func(t *testing.T) {
		updatedEntryTitle := "updated@0.1.1.*"

		testUpdateEntrySuccess(
			t, app, db, conf, updatedEntryTitle, fmt.Sprintf(bodyFmt, updatedEntryTitle),
		)
	})

	t.Run("valid_body_irrelevant_data_204_no_content", func(t *testing.T) {
		updatedEntryTitle := "updated@0.1.1.*"

		validBodyIrrelevantData := "{" +
			fmt.Sprintf(`"entry_title":"%s",`, updatedEntryTitle) +
			`"entry_slug":"notARealSlug",` +
			`"entry_created_at":"10/12/22"` +
			"}"

		testUpdateEntrySuccess(t, app, db, conf, updatedEntryTitle, validBodyIrrelevantData)
	})
}

func testUpdateEntryClientError(
	t *testing.T, app *fiber.App, conf *config.AppConfig, expectedStatus int,
	expectedMessage, expectedDetail, slug, body string,
) {
	resp := newRequestUpdateEntry(t, app, conf, slug, body)
	require.Equal(t, expectedStatus, resp.StatusCode)
	helpers.AssertErrorResponseBody(t, resp, utils.ErrorResponseBody{
		ClientOperation: utils.UpdateEntry,
		Message:         expectedMessage,
		Detail:          expectedDetail,
		RequestBody:     body,
	})
}

func testUpdateEntrySuccess(
	t *testing.T, app *fiber.App, db *gorm.DB, conf *config.AppConfig,
	updatedEntryTitle, body string,
) {
	setup.SetUpWithData(t, db)

	var entryBeforeUpdate models.Entry
	helpers.QueryTestEntry(t, db, &entryBeforeUpdate, "entry@0.1.1.*")

	resp := newRequestUpdateEntry(t, app, conf, entryBeforeUpdate.Slug, body)
	require.Equal(t, 204, resp.StatusCode)

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
	t *testing.T, app *fiber.App, conf *config.AppConfig, slug, body string,
) *http.Response {

	reqBody := strings.NewReader(body)
	req := httptest.NewRequest("PATCH", "/api/entries/" + slug, reqBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Client-Operation", utils.UpdateEntry)
	req.Header.Set("Authorization", "Token " + conf.VAULTS_ACCESS_TOKEN)

	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
