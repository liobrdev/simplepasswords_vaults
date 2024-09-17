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

func testCreateVault(t *testing.T, app *fiber.App, db *gorm.DB, conf *config.AppConfig) {
	bodyFmt := `{"user_slug":"%s","vault_title":"%s"}`

	t.Run("empty_body_400_bad_request", func(t *testing.T) {
		testCreateVaultClientError(
			t, app, conf, 400, utils.ErrorParse,
			"invalid character '\x00' looking for beginning of value", "",
		)
	})

	t.Run("array_body_400_bad_request", func(t *testing.T) {
		testCreateVaultClientError(
			t, app, conf, 400, utils.ErrorParse,
			"invalid character '[' looking for beginning of value", "[]",
		)

		testCreateVaultClientError(
			t, app, conf, 400, utils.ErrorParse,
			"invalid character '[' looking for beginning of value", "[{}]",
		)

		testCreateVaultClientError(
			t, app, conf, 400, utils.ErrorParse, "invalid character '[' looking for beginning of value",
			fmt.Sprintf("[" + bodyFmt + "]", helpers.NewSlug(t), "vault@0.2.*.*"),
		)
	})

	t.Run("boolean_body_400_bad_request", func(t *testing.T) {
		testCreateVaultClientError(
			t, app, conf, 400, utils.ErrorParse, "invalid character 't' looking for beginning of value",
			"true",
		)

		testCreateVaultClientError(
			t, app, conf, 400, utils.ErrorParse, "invalid character 'f' looking for beginning of value",
			"false",
		)
	})

	t.Run("string_body_400_bad_request", func(t *testing.T) {
		testCreateVaultClientError(
			t, app, conf, 400, utils.ErrorParse, "invalid character '\"' looking for beginning of value",
			"\"Valid JSON, but not an object.\"",
		)
	})

	t.Run("null_body_400_bad_request", func(t *testing.T) {
		testCreateVaultClientError(t, app, conf, 400, utils.ErrorUserSlug, "", "null")
	})

	t.Run("empty_object_body_400_bad_request", func(t *testing.T) {
		testCreateVaultClientError(t, app, conf, 400, utils.ErrorUserSlug, "", "{}")
	})

	t.Run("missing_user_slug_400_bad_request", func(t *testing.T) {
		testCreateVaultClientError(
			t, app, conf, 400, utils.ErrorUserSlug, "",
			`{"userr_slug":"Spelled wrong!","vault_title":"vault@0.2.*.*"}`,
		)
	})

	t.Run("missing_vault_title_400_bad_request", func(t *testing.T) {
		testCreateVaultClientError(
			t, app, conf, 400, utils.ErrorVaultTitle, "",
			fmt.Sprintf(`{"user_slug":"%s","vault_titel":"Spelled wrong!"}`, helpers.NewSlug(t)),
		)
	})

	t.Run("null_user_slug_400_bad_request", func(t *testing.T) {
		testCreateVaultClientError(
			t, app, conf, 400, utils.ErrorUserSlug, "", `{"user_slug":null,"vault_title":"vault@0.2.*.*"}`,
		)
	})

	t.Run("null_vault_title_400_bad_request", func(t *testing.T) {
		testCreateVaultClientError(
			t, app, conf, 400, utils.ErrorVaultTitle, "",
			fmt.Sprintf(`{"user_slug":"%s","vault_title":null}`, helpers.NewSlug(t)),
		)
	})

	t.Run("empty_user_slug_400_bad_request", func(t *testing.T) {
		testCreateVaultClientError(
			t, app, conf, 400, utils.ErrorUserSlug, "", fmt.Sprintf(bodyFmt, "", "vault@0.2.*.*"),
		)
	})

	t.Run("empty_vault_title_400_bad_request", func(t *testing.T) {
		testCreateVaultClientError(
			t, app, conf, 400, utils.ErrorVaultTitle, "", fmt.Sprintf(bodyFmt, helpers.NewSlug(t), ""),
		)
	})

	t.Run("too_long_vault_title_400_bad_request", func(t *testing.T) {
		if title, err := utils.GenerateSlug(256); err != nil {
			t.Fatalf("Generate long string failed: %s", err.Error())
		} else {
			testCreateVaultClientError(
				t, app, conf, 400, utils.ErrorVaultTitle, "Too long",
				fmt.Sprintf(bodyFmt, helpers.NewSlug(t), title),
			)
		}
	})

	t.Run("valid_body_vault_title_already_exists_409_conflict", func(t *testing.T) {
		users, _, _, _ := setup.SetUpWithData(t, db)
		slug := users[0].Slug

		testCreateVaultClientError(
			t, app, conf, 409, utils.ErrorDuplicateVault,
			"UNIQUE constraint failed: vaults.title, vaults.user_slug",
			fmt.Sprintf(bodyFmt, slug, "vault@0.1.*.*"),
		)
	})

	t.Run("valid_body_204_no_content", func(t *testing.T) {
		users, _, _, _ := setup.SetUpWithData(t, db)
		userSlug := users[0].Slug
		vaultTitle := "vault@0.2.*.*"

		testCreateVaultSuccess(
			t, app, db, conf, userSlug, vaultTitle, fmt.Sprintf(bodyFmt, userSlug, vaultTitle),
		)
	})

	t.Run("valid_body_irrelevant_data_204_no_content", func(t *testing.T) {
		users, _, _, _ := setup.SetUpWithData(t, db)
		userSlug := users[0].Slug
		vaultTitle := "vault@0.2.*.*"

		validBodyIrrelevantData := `{` +
			fmt.Sprintf(`"user_slug":"%s","vault_title":"%s",`, userSlug, vaultTitle) +
			`"vault_slug":"notARealSlug",` +
			`"vault_created_at":"10/12/22"` +
			`}`

		testCreateVaultSuccess(t, app, db, conf, userSlug, vaultTitle, validBodyIrrelevantData)
	})
}

func testCreateVaultClientError(
	t *testing.T, app *fiber.App, conf *config.AppConfig, expectedStatus int,
	expectedMessage, expectedDetail, body string,
) {
	resp := newRequestCreateVault(t, app, conf, body)
	require.Equal(t, expectedStatus, resp.StatusCode)
	helpers.AssertErrorResponseBody(t, resp, utils.ErrorResponseBody{
		ClientOperation: utils.CreateVault,
		Message:         expectedMessage,
		Detail:          expectedDetail,
		RequestBody:     body,
	})
}

func testCreateVaultSuccess(
	t *testing.T, app *fiber.App, db *gorm.DB, conf *config.AppConfig,
	userSlug, vaultTitle, body string,
) {
	var vaultCount int64
	helpers.CountVaults(t, db, &vaultCount)
	require.EqualValues(t, 4, vaultCount)

	resp := newRequestCreateVault(t, app, conf, body)

	require.Equal(t, 204, resp.StatusCode)

	if respBody, err := io.ReadAll(resp.Body); err != nil {
		t.Fatalf("Read response body failed: %s", err.Error())
	} else {
		require.Empty(t, respBody)
	}

	var vault models.Vault
	helpers.QueryTestVault(t, db, &vault, vaultTitle)
	require.Equal(t, vault.CreatedAt, vault.UpdatedAt)
	require.Equal(t, vaultTitle, vault.Title)
	require.Equal(t, userSlug, vault.UserSlug)
	require.Empty(t, vault.Entries)

	helpers.CountVaults(t, db, &vaultCount)
	require.EqualValues(t, 5, vaultCount)
}

func newRequestCreateVault(
	t *testing.T, app *fiber.App, conf *config.AppConfig, body string,
) *http.Response {

	reqBody := strings.NewReader(body)
	req := httptest.NewRequest(http.MethodPost, "/api/vaults", reqBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Client-Operation", utils.CreateVault)
	req.Header.Set("Authorization", "Token " + conf.VAULTS_ACCESS_TOKEN)

	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
