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

func testUpdateVault(t *testing.T, app *fiber.App, db *gorm.DB) {
	t.Run("empty_body_400_bad_request", func(t *testing.T) {
		testUpdateVaultClientError(
			t, app, db, helpers.NewSlug(t), "", http.StatusBadRequest, utils.ErrorParse,
			"invalid character '\x00' looking for beginning of value",
		)
	})

	t.Run("array_body_400_bad_request", func(t *testing.T) {
		testUpdateVaultClientError(
			t, app, db, helpers.NewSlug(t), "[]", http.StatusBadRequest, utils.ErrorParse,
			"invalid character '[' looking for beginning of value",
		)

		testUpdateVaultClientError(
			t, app, db, helpers.NewSlug(t), "[{}]", http.StatusBadRequest, utils.ErrorParse,
			"invalid character '[' looking for beginning of value",
		)

		testUpdateVaultClientError(
			t, app, db, helpers.NewSlug(t), "[{\"vault_title\":\"updated@0.1.*.*\"}]",
			http.StatusBadRequest, utils.ErrorParse,
			"invalid character '[' looking for beginning of value",
		)
	})

	t.Run("null_body_400_bad_request", func(t *testing.T) {
		testUpdateVaultClientError(
			t, app, db, helpers.NewSlug(t), "null", http.StatusBadRequest, utils.ErrorVaultTitle, "",
		)
	})

	t.Run("boolean_body_400_bad_request", func(t *testing.T) {
		testUpdateVaultClientError(
			t, app, db, helpers.NewSlug(t), "true", http.StatusBadRequest, utils.ErrorParse,
			"invalid character 't' looking for beginning of value",
		)

		testUpdateVaultClientError(
			t, app, db, helpers.NewSlug(t), "false", http.StatusBadRequest, utils.ErrorParse,
			"invalid character 'f' looking for beginning of value",
		)
	})

	t.Run("string_body_400_bad_request", func(t *testing.T) {
		testUpdateVaultClientError(
			t, app, db, helpers.NewSlug(t), "\"Valid JSON, but not an object.\"",
			http.StatusBadRequest, utils.ErrorParse,
			"invalid character '\"' looking for beginning of value",
		)
	})

	t.Run("empty_object_body_400_bad_request", func(t *testing.T) {
		testUpdateVaultClientError(
			t, app, db, helpers.NewSlug(t), "{}", http.StatusBadRequest, utils.ErrorVaultTitle, "",
		)
	})

	t.Run("missing_vault_title_400_bad_request", func(t *testing.T) {
		testUpdateVaultClientError(
			t, app, db, helpers.NewSlug(t), `{"vualt_title":"Spelled wrong!"}`,
			http.StatusBadRequest, utils.ErrorVaultTitle, "",
		)
	})

	t.Run("null_vault_title_400_bad_request", func(t *testing.T) {
		testUpdateVaultClientError(
			t, app, db, helpers.NewSlug(t), `{"vault_title":null}`,
			http.StatusBadRequest, utils.ErrorVaultTitle, "",
		)
	})

	t.Run("empty_vault_title_400_bad_request", func(t *testing.T) {
		testUpdateVaultClientError(
			t, app, db, helpers.NewSlug(t), `{"vault_title":""}`,
			http.StatusBadRequest, utils.ErrorVaultTitle, "",
		)
	})

	t.Run("too_long_entry_title_400_bad_request", func(t *testing.T) {
		// `title` is a random string greater than 255 characters in length
		if title, err := utils.GenerateSlug(256); err != nil {
			t.Fatalf("Generate long string failed: %s", err.Error())
		} else {
			testUpdateVaultClientError(
				t, app, db, helpers.NewSlug(t), fmt.Sprintf(`{"vault_title":"%s"}`, title),
				http.StatusBadRequest, utils.ErrorVaultTitle, "Too long (256 > 255)",
			)
		}
	})

	t.Run("valid_body_404_not_found", func(t *testing.T) {
		setup.SetUpWithData(t, db)
		testUpdateVaultClientError(
			t, app, db, helpers.NewSlug(t), `{"vault_title":"updated@0.1.*.*"}`,
			http.StatusNotFound, utils.ErrorNoRowsAffected, "Likely that slug was not found.",
		)
	})

	t.Run("valid_body_vault_title_already_exists_409_conflict", func(t *testing.T) {
		_, vaults, _, _ := setup.SetUpWithData(t, db)

		testUpdateVaultClientError(
			t, app, db, vaults[0].Slug,
			fmt.Sprintf(`{"vault_title":"%s"}`, vaults[1].Title),
			http.StatusConflict, utils.ErrorFailedDB,
			"UNIQUE constraint failed: vaults.title, vaults.user_slug",
		)
	})

	t.Run("valid_body_204_no_content", func(t *testing.T) {
		updatedVaultTitle := "updated@0.1.*.*"

		testUpdateVaultSuccess(t, app, db, updatedVaultTitle, fmt.Sprintf(
			`{"vault_title":"%s"}`,
			updatedVaultTitle,
		))
	})

	t.Run("valid_body_irrelevant_data_204_no_content", func(t *testing.T) {
		updatedVaultTitle := "updated@0.1.*.*"

		validBodyIrrelevantData := "{" +
			fmt.Sprintf(`"vault_title":"%s",`, updatedVaultTitle) +
			`"vault_slug":"notEvenARealSlug",` +
			`"vault_created_at":"10/12/22"` +
			"}"

		testUpdateVaultSuccess(t, app, db, updatedVaultTitle, validBodyIrrelevantData)
	})
}

func testUpdateVaultClientError(
	t *testing.T, app *fiber.App, db *gorm.DB, slug string, body string, expectedStatus int,
	expectedMessage string, expectedDetail string,
) {
	resp := newRequestUpdateVault(t, app, slug, body)
	require.Equal(t, expectedStatus, resp.StatusCode)
	helpers.AssertErrorResponseBody(t, resp, utils.ErrorResponseBody{
		ClientOperation: utils.UpdateVault,
		Message:         string(expectedMessage),
		Detail:          expectedDetail,
		RequestBody:     body,
	})
}

func testUpdateVaultSuccess(
	t *testing.T, app *fiber.App, db *gorm.DB, updatedVaultTitle, body string,
) {
	setup.SetUpWithData(t, db)
	var vaultBeforeUpdate models.Vault
	helpers.QueryTestVault(t, db, &vaultBeforeUpdate, "vault@0.1.*.*")
	resp := newRequestUpdateVault(t, app, vaultBeforeUpdate.Slug, body)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)

	if respBody, err := io.ReadAll(resp.Body); err != nil {
		t.Fatalf("Read response body failed: %s", err.Error())
	} else {
		require.Empty(t, respBody)
	}

	var vaultAfterUpdate models.Vault
	helpers.QueryTestVault(t, db, &vaultAfterUpdate, updatedVaultTitle)
	require.Equal(t, vaultBeforeUpdate.Slug, vaultAfterUpdate.Slug)
	require.NotEqual(t, vaultBeforeUpdate.Title, vaultAfterUpdate.Title)
	require.Equal(t, updatedVaultTitle, vaultAfterUpdate.Title)
	require.Equal(t, vaultBeforeUpdate.CreatedAt, vaultAfterUpdate.CreatedAt)
	require.True(t, vaultAfterUpdate.UpdatedAt.After(vaultAfterUpdate.CreatedAt))
	require.True(t, vaultAfterUpdate.UpdatedAt.After(vaultBeforeUpdate.UpdatedAt))
}

func newRequestUpdateVault(t *testing.T, app *fiber.App, slug, body string) *http.Response {
	reqBody := strings.NewReader(body)
	req := httptest.NewRequest(http.MethodPatch, "/api/vaults/" + slug, reqBody)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
