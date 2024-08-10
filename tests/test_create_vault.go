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

func testCreateVault(t *testing.T, app *fiber.App, db *gorm.DB) {
	t.Run("empty_body_400_bad_request", func(t *testing.T) {
		testCreateVaultClientError(
			t, app, db, "", http.StatusBadRequest, utils.ErrorParse,
			"invalid character '\x00' looking for beginning of value",
		)
	})

	t.Run("array_body_400_bad_request", func(t *testing.T) {
		testCreateVaultClientError(
			t, app, db, "[]", http.StatusBadRequest, utils.ErrorParse,
			"invalid character '[' looking for beginning of value",
		)

		testCreateVaultClientError(
			t, app, db, "[{}]", http.StatusBadRequest, utils.ErrorParse,
			"invalid character '[' looking for beginning of value",
		)

		testCreateVaultClientError(
			t, app, db, fmt.Sprintf(
				`[{"user_slug":"%s","vault_title":"%s"}]`, helpers.NewSlug(t), "vault@0.2.*.*",
			), http.StatusBadRequest, utils.ErrorParse,
			"invalid character '[' looking for beginning of value",
		)
	})

	t.Run("null_body_400_bad_request", func(t *testing.T) {
		testCreateVaultClientError(
			t, app, db, "null", http.StatusBadRequest, utils.ErrorUserSlug, "",
		)
	})

	t.Run("boolean_body_400_bad_request", func(t *testing.T) {
		testCreateVaultClientError(
			t, app, db, "true", http.StatusBadRequest, utils.ErrorParse,
			"invalid character 't' looking for beginning of value",
		)

		testCreateVaultClientError(
			t, app, db, "false", http.StatusBadRequest, utils.ErrorParse,
			"invalid character 'f' looking for beginning of value",
		)
	})

	t.Run("string_body_400_bad_request", func(t *testing.T) {
		testCreateVaultClientError(
			t, app, db, "\"Valid JSON, but not an object.\"", http.StatusBadRequest,
			utils.ErrorParse, "invalid character '\"' looking for beginning of value",
		)
	})

	t.Run("empty_object_body_400_bad_request", func(t *testing.T) {
		testCreateVaultClientError(
			t, app, db, "{}", http.StatusBadRequest, utils.ErrorUserSlug, "",
		)
	})

	t.Run("missing_user_slug_400_bad_request", func(t *testing.T) {
		testCreateVaultClientError(
			t, app, db, `{"userr_slug":"Spelled wrong!","vault_title":"vault@0.2.*.*"}`,
			http.StatusBadRequest, utils.ErrorUserSlug, "",
		)
	})

	t.Run("missing_vault_title_400_bad_request", func(t *testing.T) {
		testCreateVaultClientError(
			t, app, db, fmt.Sprintf(
				`{"user_slug":"%s","vault_titel":"Spelled wrong!"}`, helpers.NewSlug(t),
			), http.StatusBadRequest, utils.ErrorVaultTitle, "",
		)
	})

	t.Run("null_user_slug_400_bad_request", func(t *testing.T) {
		testCreateVaultClientError(
			t, app, db, `{"user_slug":null,"vault_title":"vault@0.2.*.*"}`,
			http.StatusBadRequest, utils.ErrorUserSlug, "",
		)
	})

	t.Run("null_vault_title_400_bad_request", func(t *testing.T) {
		testCreateVaultClientError(
			t, app, db,
			fmt.Sprintf(`{"user_slug":"%s","vault_title":null}`, helpers.NewSlug(t)),
			http.StatusBadRequest, utils.ErrorVaultTitle, "",
		)
	})

	t.Run("empty_user_slug_400_bad_request", func(t *testing.T) {
		testCreateVaultClientError(
			t, app, db, `{"user_slug":"","vault_title":"vault@0.2.*.*"}`,
			http.StatusBadRequest, utils.ErrorUserSlug, "",
		)
	})

	t.Run("empty_vault_title_400_bad_request", func(t *testing.T) {
		testCreateVaultClientError(
			t, app, db,
			fmt.Sprintf(`{"user_slug":"%s","vault_title":""}`, helpers.NewSlug(t)),
			http.StatusBadRequest, utils.ErrorVaultTitle, "",
		)
	})

	t.Run("too_long_vault_title_400_bad_request", func(t *testing.T) {
		// `title` is a random string greater than 255 characters in length
		if title, err := utils.GenerateSlug(256); err != nil {
			t.Fatalf("Generate long string failed: %s", err.Error())
		} else {
			testCreateVaultClientError(
				t, app, db, fmt.Sprintf(
					`{"user_slug":"%s","vault_title":"%s"}`, helpers.NewSlug(t), title,
				), http.StatusBadRequest, utils.ErrorVaultTitle, "Too long (256 > 255)",
			)
		}
	})

	t.Run("valid_body_vault_title_already_exists_409_conflict", func(t *testing.T) {
		users, _, _, _ := setup.SetUpWithData(t, db)
		slug := users[0].Slug

		testCreateVaultClientError(
			t, app, db, fmt.Sprintf(`{"user_slug":"%s","vault_title":"vault@0.1.*.*"}`, slug),
			http.StatusConflict, utils.ErrorFailedDB,
			"UNIQUE constraint failed: vaults.title, vaults.user_slug",
		)
	})

	t.Run("valid_body_204_no_content", func(t *testing.T) {
		users, _, _, _ := setup.SetUpWithData(t, db)
		userSlug := users[0].Slug
		vaultTitle := "vault@0.2.*.*"

		testCreateVaultSuccess(t, app, db, userSlug, vaultTitle, fmt.Sprintf(
			`{"user_slug":"%s","vault_title":"%s"}`,
			userSlug,
			vaultTitle,
		))
	})

	t.Run("valid_body_irrelevant_data_204_no_content", func(t *testing.T) {
		users, _, _, _ := setup.SetUpWithData(t, db)
		userSlug := users[0].Slug
		vaultTitle := "vault@0.2.*.*"

		validBodyIrrelevantData := `{` +
			fmt.Sprintf(`"user_slug":"%s","vault_title":"%s",`, userSlug, vaultTitle) +
			`"vault_slug":"notEvenARealSlug",` +
			`"vault_created_at":"10/12/22"` +
			`}`

		testCreateVaultSuccess(t, app, db, userSlug, vaultTitle, validBodyIrrelevantData)
	})
}

func testCreateVaultClientError(
	t *testing.T, app *fiber.App, db *gorm.DB, body string, expectedStatus int,
	expectedMessage string, expectedDetail string,
) {
	resp := newRequestCreateVault(t, app, body)
	require.Equal(t, expectedStatus, resp.StatusCode)
	helpers.AssertErrorResponseBody(t, resp, utils.ErrorResponseBody{
		ClientOperation: utils.CreateVault,
		Message:         expectedMessage,
		Detail:          expectedDetail,
		RequestBody:     body,
	})
}

func testCreateVaultSuccess(
	t *testing.T, app *fiber.App, db *gorm.DB, userSlug, vaultTitle, body string,
) {
	var vaultCount int64
	helpers.CountVaults(t, db, &vaultCount)
	require.EqualValues(t, 4, vaultCount)

	resp := newRequestCreateVault(t, app, body)

	require.Equal(t, http.StatusNoContent, resp.StatusCode)

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

func newRequestCreateVault(t *testing.T, app *fiber.App, body string) *http.Response {
	reqBody := strings.NewReader(body)
	req := httptest.NewRequest(http.MethodPost, "/api/vaults", reqBody)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
