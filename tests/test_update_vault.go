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

func testUpdateVault(t *testing.T, app *fiber.App, db *gorm.DB, conf *config.AppConfig) {
	bodyFmt := `{"vault_title":"%s"}`

	t.Run("empty_body_400_bad_request", func(t *testing.T) {
		testUpdateVaultClientError(
			t, app, conf, 400, utils.ErrorParse,
			"invalid character '\x00' looking for beginning of value", helpers.NewSlug(t), "",
		)
	})

	t.Run("array_body_400_bad_request", func(t *testing.T) {
		testUpdateVaultClientError(
			t, app, conf, 400, utils.ErrorParse, "invalid character '[' looking for beginning of value",
			helpers.NewSlug(t), "[]",
		)

		testUpdateVaultClientError(
			t, app, conf, 400, utils.ErrorParse, "invalid character '[' looking for beginning of value",
			helpers.NewSlug(t), "[{}]",
		)

		testUpdateVaultClientError(
			t, app, conf, 400, utils.ErrorParse, "invalid character '[' looking for beginning of value",
			helpers.NewSlug(t), `[{"vault_title":"updated@0.1.*.*"}]`,
		)
	})

	t.Run("boolean_body_400_bad_request", func(t *testing.T) {
		testUpdateVaultClientError(
			t, app, conf, 400, utils.ErrorParse, "invalid character 't' looking for beginning of value",
			helpers.NewSlug(t), "true",
		)

		testUpdateVaultClientError(
			t, app, conf, 400, utils.ErrorParse, "invalid character 'f' looking for beginning of value",
			helpers.NewSlug(t), "false",
		)
	})

	t.Run("string_body_400_bad_request", func(t *testing.T) {
		testUpdateVaultClientError(
			t, app, conf, 400, utils.ErrorParse, `invalid character '"' looking for beginning of value`,
			helpers.NewSlug(t), `"Valid JSON, but not an object."`,
		)
	})

	t.Run("null_body_400_bad_request", func(t *testing.T) {
		testUpdateVaultClientError(
			t, app, conf, 400, utils.ErrorVaultTitle, "", helpers.NewSlug(t), "null",
		)
	})

	t.Run("empty_object_body_400_bad_request", func(t *testing.T) {
		testUpdateVaultClientError(
			t, app, conf, 400, utils.ErrorVaultTitle, "", helpers.NewSlug(t), "{}",
		)
	})

	t.Run("missing_vault_title_400_bad_request", func(t *testing.T) {
		testUpdateVaultClientError(
			t, app, conf, 400, utils.ErrorVaultTitle, "",
			helpers.NewSlug(t), `{"vualt_title":"Spelled wrong!"}`,
		)
	})

	t.Run("null_vault_title_400_bad_request", func(t *testing.T) {
		testUpdateVaultClientError(
			t, app, conf, 400, utils.ErrorVaultTitle, "", helpers.NewSlug(t), `{"vault_title":null}`,
		)
	})

	t.Run("empty_vault_title_400_bad_request", func(t *testing.T) {
		testUpdateVaultClientError(
			t, app, conf, 400, utils.ErrorVaultTitle, "", helpers.NewSlug(t), `{"vault_title":""}`,
		)
	})

	t.Run("too_long_entry_title_400_bad_request", func(t *testing.T) {
		if title, err := utils.GenerateSlug(256); err != nil {
			t.Fatalf("Generate long string failed: %s", err.Error())
		} else {
			testUpdateVaultClientError(
				t, app, conf, 400, utils.ErrorVaultTitle, "Too long",
				helpers.NewSlug(t), fmt.Sprintf(bodyFmt, title),
			)
		}
	})

	t.Run("valid_body_404_not_found", func(t *testing.T) {
		setup.SetUpWithData(t, db)
		testUpdateVaultClientError(
			t, app, conf, 404, utils.ErrorNoRowsAffected, "Likely that slug was not found.",
			helpers.NewSlug(t), `{"vault_title":"updated@0.1.*.*"}`,
		)
	})

	t.Run("valid_body_vault_title_already_exists_409_conflict", func(t *testing.T) {
		_, vaults, _, _ := setup.SetUpWithData(t, db)

		testUpdateVaultClientError(
			t, app, conf, 500, utils.ErrorFailedDB,
			"UNIQUE constraint failed: vaults.title, vaults.user_slug",
			vaults[0].Slug, fmt.Sprintf(bodyFmt, vaults[1].Title),
		)
	})

	t.Run("valid_body_204_no_content", func(t *testing.T) {
		updatedVaultTitle := "updated@0.1.*.*"

		testUpdateVaultSuccess(
			t, app, db, conf, updatedVaultTitle, fmt.Sprintf(bodyFmt, updatedVaultTitle),
		)
	})

	t.Run("valid_body_irrelevant_data_204_no_content", func(t *testing.T) {
		updatedVaultTitle := "updated@0.1.*.*"

		validBodyIrrelevantData := "{" +
			fmt.Sprintf(`"vault_title":"%s",`, updatedVaultTitle) +
			`"vault_slug":"notEvenARealSlug",` +
			`"vault_created_at":"10/12/22"` +
			"}"

		testUpdateVaultSuccess(t, app, db, conf, updatedVaultTitle, validBodyIrrelevantData)
	})
}

func testUpdateVaultClientError(
	t *testing.T, app *fiber.App, conf *config.AppConfig,
	expectedStatus int, expectedMessage, expectedDetail, slug, body string,
) {
	resp := newRequestUpdateVault(t, app, conf, slug, body)
	require.Equal(t, expectedStatus, resp.StatusCode)
	helpers.AssertErrorResponseBody(t, resp, utils.ErrorResponseBody{
		ClientOperation: utils.UpdateVault,
		Message:         expectedMessage,
		Detail:          expectedDetail,
		RequestBody:     body,
	})
}

func testUpdateVaultSuccess(
	t *testing.T, app *fiber.App, db *gorm.DB, conf *config.AppConfig,
	updatedVaultTitle, body string,
) {
	setup.SetUpWithData(t, db)

	var vaultBeforeUpdate models.Vault
	helpers.QueryTestVault(t, db, &vaultBeforeUpdate, "vault@0.1.*.*")

	resp := newRequestUpdateVault(t, app, conf, vaultBeforeUpdate.Slug, body)
	require.Equal(t, 204, resp.StatusCode)

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

func newRequestUpdateVault(
	t *testing.T, app *fiber.App, conf *config.AppConfig, slug, body string,
) *http.Response {

	reqBody := strings.NewReader(body)
	req := httptest.NewRequest("PATCH", "/api/vaults/" + slug, reqBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Client-Operation", utils.RetrieveVault)
	req.Header.Set("Authorization", "Token " + conf.VAULTS_ACCESS_TOKEN)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
