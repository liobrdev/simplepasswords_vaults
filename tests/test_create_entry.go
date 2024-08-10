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

func testCreateEntry(t *testing.T, app *fiber.App, db *gorm.DB) {
	t.Run("empty_body_400_bad_request", func(t *testing.T) {
		testCreateEntryClientError(
			t, app, db, "", http.StatusBadRequest, utils.ErrorParse,
			"invalid character '\x00' looking for beginning of value",
		)
	})

	t.Run("array_body_400_bad_request", func(t *testing.T) {
		testCreateEntryClientError(
			t, app, db, "[]", http.StatusBadRequest, utils.ErrorParse,
			"invalid character '[' looking for beginning of value",
		)

		testCreateEntryClientError(
			t, app, db, "[{}]", http.StatusBadRequest, utils.ErrorParse,
			"invalid character '[' looking for beginning of value",
		)

		testCreateEntryClientError(
			t, app, db, fmt.Sprintf(
				`[{"user_slug":"%s","vault_slug":"%s","entry_title":"%s","secrets":[]}]`,
				helpers.NewSlug(t), helpers.NewSlug(t), "entry@0.0.2.*",
			), http.StatusBadRequest, utils.ErrorParse,
			"invalid character '[' looking for beginning of value",
		)
	})

	t.Run("null_body_400_bad_request", func(t *testing.T) {
		testCreateEntryClientError(
			t, app, db, "null", http.StatusBadRequest, utils.ErrorUserSlug, "",
		)
	})

	t.Run("boolean_body_400_bad_request", func(t *testing.T) {
		testCreateEntryClientError(
			t, app, db, "true", http.StatusBadRequest, utils.ErrorParse,
			"invalid character 't' looking for beginning of value",
		)

		testCreateEntryClientError(
			t, app, db, "false", http.StatusBadRequest, utils.ErrorParse,
			"invalid character 'f' looking for beginning of value",
		)
	})

	t.Run("string_body_400_bad_request", func(t *testing.T) {
		testCreateEntryClientError(
			t, app, db, "\"Valid JSON, but not an object.\"", http.StatusBadRequest,
			utils.ErrorParse, "invalid character '\"' looking for beginning of value",
		)
	})

	t.Run("empty_object_body_400_bad_request", func(t *testing.T) {
		testCreateEntryClientError(
			t, app, db, "{}", http.StatusBadRequest, utils.ErrorUserSlug, "",
		)
	})

	t.Run("missing_user_slug_400_bad_request", func(t *testing.T) {
		testCreateEntryClientError(
			t, app, db, fmt.Sprintf(
				`{"userr_slug":"%s","vault_slug":"%s","entry_title":"%s","secrets":[]}`,
				"Spelled wrong!", helpers.NewSlug(t), "entry@0.0.2.*",
			), http.StatusBadRequest, utils.ErrorUserSlug, "",
		)
	})

	t.Run("missing_vault_slug_400_bad_request", func(t *testing.T) {
		testCreateEntryClientError(
			t, app, db, fmt.Sprintf(
				`{"user_slug":"%s","vualt_slug":"%s","entry_title":"%s","secrets":[]}`,
				helpers.NewSlug(t), "Spelled wrong!", "entry@0.0.2.*",
			), http.StatusBadRequest, utils.ErrorVaultSlug, "",
		)
	})

	t.Run("missing_entry_title_400_bad_request", func(t *testing.T) {
		testCreateEntryClientError(
			t, app, db, fmt.Sprintf(
				`{"user_slug":"%s","vault_slug":"%s","enrty_title":"%s","secrets":[]}`,
				helpers.NewSlug(t), helpers.NewSlug(t), "Spelled wrong!",
			), http.StatusBadRequest, utils.ErrorEntryTitle, "",
		)
	})

	t.Run("missing_secrets_400_bad_request", func(t *testing.T) {
		testCreateEntryClientError(
			t, app, db, fmt.Sprintf(
				`{"user_slug":"%s","vault_slug":"%s","entry_title":"%s","secerts":[]}`,
				helpers.NewSlug(t), helpers.NewSlug(t), "entry@0.0.2.*",
			), http.StatusBadRequest, utils.ErrorSecrets, "",
		)
	})

	t.Run("null_user_slug_400_bad_request", func(t *testing.T) {
		testCreateEntryClientError(
			t, app, db, fmt.Sprintf(
				`{"user_slug":%s,"vault_slug":"%s","entry_title":"%s","secrets":[]}`,
				"null", helpers.NewSlug(t), "entry@0.0.2.*",
			), http.StatusBadRequest, utils.ErrorUserSlug, "",
		)
	})

	t.Run("null_vault_slug_400_bad_request", func(t *testing.T) {
		testCreateEntryClientError(
			t, app, db, fmt.Sprintf(
				`{"user_slug":"%s","vault_slug":%s,"entry_title":"%s","secrets":[]}`,
				helpers.NewSlug(t), "null", "entry@0.0.2.*",
			), http.StatusBadRequest, utils.ErrorVaultSlug, "",
		)
	})

	t.Run("null_entry_title_400_bad_request", func(t *testing.T) {
		testCreateEntryClientError(
			t, app, db, fmt.Sprintf(
				`{"user_slug":"%s","vault_slug":"%s","entry_title":%s,"secrets":[]}`,
				helpers.NewSlug(t), helpers.NewSlug(t), "null",
			), http.StatusBadRequest, utils.ErrorEntryTitle, "",
		)
	})

	t.Run("null_secrets_400_bad_request", func(t *testing.T) {
		testCreateEntryClientError(
			t, app, db, fmt.Sprintf(
				`{"user_slug":"%s","vault_slug":"%s","entry_title":"%s","secrets":null}`,
				helpers.NewSlug(t), helpers.NewSlug(t), "entry@0.0.2.*",
			), http.StatusBadRequest, utils.ErrorSecrets, "",
		)
	})

	t.Run("empty_user_slug_400_bad_request", func(t *testing.T) {
		testCreateEntryClientError(
			t, app, db, fmt.Sprintf(
				`{"user_slug":"","vault_slug":"%s","entry_title":"%s","secrets":[]}`,
				helpers.NewSlug(t), "entry@0.0.2.*",
			), http.StatusBadRequest, utils.ErrorUserSlug, "",
		)
	})

	t.Run("empty_vault_slug_400_bad_request", func(t *testing.T) {
		testCreateEntryClientError(
			t, app, db, fmt.Sprintf(
				`{"user_slug":"%s","vault_slug":"","entry_title":"%s","secrets":[]}`,
				helpers.NewSlug(t), "entry@0.0.2.*",
			), http.StatusBadRequest, utils.ErrorVaultSlug, "",
		)
	})

	t.Run("empty_entry_title_400_bad_request", func(t *testing.T) {
		testCreateEntryClientError(
			t, app, db, fmt.Sprintf(
				`{"user_slug":"%s","vault_slug":"%s","entry_title":"","secrets":[]}`,
				helpers.NewSlug(t), helpers.NewSlug(t),
			), http.StatusBadRequest, utils.ErrorEntryTitle, "",
		)
	})

	t.Run("too_long_entry_title_400_bad_request", func(t *testing.T) {
		// `title` is a random string greater than 255 characters in length
		if title, err := utils.GenerateSlug(256); err != nil {
			t.Fatalf("Generate long string failed: %s", err.Error())
		} else {
			testCreateEntryClientError(
				t, app, db, fmt.Sprintf(
					`{"user_slug":"%s","vault_slug":"%s","entry_title":"%s","secrets":[]}`,
					helpers.NewSlug(t), helpers.NewSlug(t), title,
				), http.StatusBadRequest, utils.ErrorEntryTitle, "Too long (256 > 255)",
			)
		}
	})

	t.Run("empty_item_in_secrets_body_400_bad_request", func(t *testing.T) {
		secretsStr := `[{` +
			`"secret_label":"secret[_label='username']@0.0.2.0",` +
			`"secret_string":"secret[_string='foodeater1234']@0.0.2.0"` +
			`},{}]`

		testCreateEntryClientError(
			t, app, db, fmt.Sprintf(
				`{"user_slug":"%s","vault_slug":"%s","entry_title":"%s","secrets":%s}`,
				helpers.NewSlug(t), helpers.NewSlug(t), "entry@0.0.2.*", secretsStr,
			), http.StatusBadRequest, utils.ErrorItemSecrets,
			"secrets[1].Label; len(secrets) == 2",
		)
	})

	t.Run("missing_secret_label_item_in_secrets_body_400_bad_request", func(t *testing.T) {
		secretsStr := `[{` +
			`"secret_label":"secret[_label='username']@0.0.2.0",` +
			`"secret_string":"secret[_string='foodeater1234']@0.0.2.0"` +
			`},{` +
			`"secret_string":"secret[_string='3a7!ng40oD']@0.0.2.1"` +
			`}]`

		testCreateEntryClientError(
			t, app, db, fmt.Sprintf(
				`{"user_slug":"%s","vault_slug":"%s","entry_title":"%s","secrets":%s}`,
				helpers.NewSlug(t), helpers.NewSlug(t), "entry@0.0.2.*", secretsStr,
			), http.StatusBadRequest, utils.ErrorItemSecrets,
			"secrets[1].Label; len(secrets) == 2",
		)
	})

	t.Run("missing_secret_string_item_in_secrets_body_400_bad_request", func(t *testing.T) {
		secretsStr := "[{" +
			`"secret_label":"secret[_label='username']@0.0.2.0",` +
			`"secret_string":"secret[_string='foodeater1234']@0.0.2.0"` +
			`},{` +
			`"secret_label":"secret[_label='password']@0.0.2.1"` +
			`}]`

		testCreateEntryClientError(
			t, app, db, fmt.Sprintf(
				`{"user_slug":"%s","vault_slug":"%s","entry_title":"%s","secrets":%s}`,
				helpers.NewSlug(t), helpers.NewSlug(t), "entry@0.0.2.*", secretsStr,
			), http.StatusBadRequest, utils.ErrorItemSecrets,
			"secrets[1].String; len(secrets) == 2",
		)
	})

	t.Run("empty_secret_label_item_in_secrets_body_400_bad_request", func(t *testing.T) {
		secretsStr := `[{` +
			`"secret_label":"secret[_label='username']@0.0.2.0",` +
			`"secret_string":"secret[_string='foodeater1234']@0.0.2.0"` +
			`},{` +
			`"secret_label":"",` +
			`"secret_string":"secret[_string='3a7!ng40oD']@0.0.2.1"` +
			`}]`

		testCreateEntryClientError(
			t, app, db, fmt.Sprintf(
				`{"user_slug":"%s","vault_slug":"%s","entry_title":"%s","secrets":%s}`,
				helpers.NewSlug(t), helpers.NewSlug(t), "entry@0.0.2.*", secretsStr,
			), http.StatusBadRequest, utils.ErrorItemSecrets,
			"secrets[1].Label; len(secrets) == 2",
		)
	})

	t.Run("empty_secret_string_item_in_secrets_body_400_bad_request", func(t *testing.T) {
		secretsStr := `[{` +
			`"secret_label":"secret[_label='username']@0.0.2.0",` +
			`"secret_string":"secret[_string='foodeater1234']@0.0.2.0"` +
			`},{` +
			`"secret_label":"secret[_label='password']@0.0.2.1",` +
			`"secret_string":""` +
			`}]`

		testCreateEntryClientError(
			t, app, db, fmt.Sprintf(
				`{"user_slug":"%s","vault_slug":"%s","entry_title":"%s","secrets":%s}`,
				helpers.NewSlug(t), helpers.NewSlug(t), "entry@0.0.2.*", secretsStr,
			), http.StatusBadRequest, utils.ErrorItemSecrets,
			"secrets[1].String; len(secrets) == 2",
		)
	})

	t.Run("too_long_secret_label_item_in_secrets_body_400_bad_request", func(t *testing.T) {
		// `label` is a random string greater than 255 characters in length
		if label, err := utils.GenerateSlug(256); err != nil {
			t.Fatalf("Generate long string failed: %s", err.Error())
		} else {
			secretsStr := `[{` +
				`"secret_label":"secret[_label='username']@0.0.2.0",` +
				`"secret_string":"secret[_string='foodeater1234']@0.0.2.0"` +
				`},{` +
				fmt.Sprintf(`"secret_label":"%s",`, label) +
				`"secret_string":"secret[_string='3a7!ng40oD']@0.0.2.1"` +
				`}]`

			testCreateEntryClientError(
				t, app, db, fmt.Sprintf(
					`{"user_slug":"%s","vault_slug":"%s","entry_title":"%s","secrets":%s}`,
					helpers.NewSlug(t), helpers.NewSlug(t), "entry@0.0.2.*", secretsStr,
				), http.StatusBadRequest, utils.ErrorItemSecrets,
				"secrets[1].Label; len(secrets) == 2",
			)
		}
	})

	t.Run("too_long_secret_string_item_in_secrets_body_400_bad_request", func(t *testing.T) {
		// `str` is a random string greater than 1000 characters in length
		if str, err := utils.GenerateSlug(1001); err != nil {
			t.Fatalf("Generate long string failed: %s", err.Error())
		} else {
			secretsStr := `[{` +
				`"secret_label":"secret[_label='username']@0.0.2.0",` +
				`"secret_string":"secret[_string='foodeater1234']@0.0.2.0"` +
				`},{` +
				`"secret_label":"secret[_label='password']@0.0.2.1",` +
				fmt.Sprintf(`"secret_string":"%s"`, str) +
				`}]`

			testCreateEntryClientError(
				t, app, db, fmt.Sprintf(
					`{"user_slug":"%s","vault_slug":"%s","entry_title":"%s","secrets":%s}`,
					helpers.NewSlug(t), helpers.NewSlug(t), "entry@0.0.2.*", secretsStr,
				), http.StatusBadRequest, utils.ErrorItemSecrets,
				"secrets[1].String; len(secrets) == 2",
			)
		}
	})

	t.Run("valid_body_entry_title_already_exists_409_conflict", func(t *testing.T) {
		users, vaults, _, _ := setup.SetUpWithData(t, db)
		userSlug := users[0].Slug
		vaultSlug := vaults[0].Slug
		entryTitle := "entry@0.0.1.*"

		secretsStr := `[{` +
			`"secret_label":"secret[_label='username']@0.0.2.0",` +
			`"secret_string":"secret[_string='foodeater1234']@0.0.2.0"` +
			`},{` +
			`"secret_label":"secret[_label='password']@0.0.2.1",` +
			`"secret_string":"secret[_string='3a7!ng40oD']@0.0.2.1"` +
			`}]`

		body := fmt.Sprintf(
			`{"user_slug":"%s","vault_slug":"%s","entry_title":"%s","secrets":%s}`,
			userSlug, vaultSlug, entryTitle, secretsStr,
		)

		testCreateEntryClientError(
			t, app, db, body,
			http.StatusConflict, utils.ErrorFailedDB,
			"UNIQUE constraint failed: entries.title, entries.vault_slug",
		)
	})

	t.Run("valid_body_entry_secret_label_duplicate_400_bad_request", func(t *testing.T) {
		secretLabel := "secret[_label='misc']@0.0.2.*"

		secretsStr := `[{` +
			fmt.Sprintf(`"secret_label":"%s",`, secretLabel) +
			`"secret_string":"secret[_string='foodeater1234']@0.0.2.0"` +
			`},{` +
			fmt.Sprintf(`"secret_label":"%s",`, secretLabel) +
			`"secret_string":"secret[_string='3a7!ng40oD']@0.0.2.1"` +
			`}]`

		body := fmt.Sprintf(
			`{"user_slug":"%s","vault_slug":"%s","entry_title":"%s","secrets":%s}`,
			helpers.NewSlug(t), helpers.NewSlug(t), "entry@0.0.2.*", secretsStr,
		)

		testCreateEntryClientError(
			t, app, db, body, http.StatusBadRequest, utils.ErrorDuplicateSecrets,
			secretLabel,
		)
	})

	t.Run("valid_body_204_no_content", func(t *testing.T) {
		users, _, _, _ := setup.SetUpWithData(t, db)
		userSlug := users[0].Slug

		var entryCount int64
		helpers.CountEntries(t, db, &entryCount)

		var secretCount int64
		helpers.CountSecrets(t, db, &secretCount)

		var vault models.Vault
		helpers.QueryTestVault(t, db, &vault, "vault@0.0.*.*")

		entryTitle := "entry@0.0.2.*"

		secretsStr := `[{` +
			`"secret_label":"secret[_label='username']@0.0.2.0",` +
			`"secret_string":"secret[_string='foodeater1234']@0.0.2.0"` +
			`},{` +
			`"secret_label":"secret[_label='password']@0.0.2.1",` +
			`"secret_string":"secret[_string='3a7!ng40oD']@0.0.2.1"` +
			`}]`

		body := fmt.Sprintf(
			`{"user_slug":"%s","vault_slug":"%s","entry_title":"%s","secrets":%s}`,
			userSlug, vault.Slug, entryTitle, secretsStr,
		)

		testCreateEntrySuccess(
			t, app, db, entryCount, secretCount, userSlug, vault.Slug, entryTitle, body,
		)
	})

	t.Run("valid_body_irrelevant_data_204_no_content", func(t *testing.T) {
		users, _, _, _ := setup.SetUpWithData(t, db)
		userSlug := users[0].Slug

		var entryCount int64
		helpers.CountEntries(t, db, &entryCount)

		var secretCount int64
		helpers.CountSecrets(t, db, &secretCount)

		var vault models.Vault
		helpers.QueryTestVault(t, db, &vault, "vault@0.0.*.*")

		entryTitle := "entry@0.0.2.*"

		secretsStr := `[{` +
			`"secret_label":"secret[_label='username']@0.0.2.0",` +
			`"secret_string":"secret[_string='foodeater1234']@0.0.2.0",` +
			`"secret_slug":"notEvenARealSlug_0",` +
			`"secret_created_at":"10/12/22"` +
			`},{` +
			`"secret_label":"secret[_label='password']@0.0.2.1",` +
			`"secret_string":"secret[_string='3a7!ng40oD']@0.0.2.1",` +
			`"secret_slug":"notEvenARealSlug_1",` +
			`"secret_created_at":"10/12/22"` +
			`}]`

		body := fmt.Sprintf(
			`{"user_slug":"%s","vault_slug":"%s","entry_title":"%s","secrets":%s,"abc":123}`,
			userSlug, vault.Slug, entryTitle, secretsStr,
		)

		testCreateEntrySuccess(
			t, app, db, entryCount, secretCount, userSlug, vault.Slug, entryTitle, body,
		)
	})
}

func testCreateEntryClientError(
	t *testing.T, app *fiber.App, db *gorm.DB, body string, expectedStatus int,
	expectedMessage string, expectedDetail string,
) {
	resp := newRequestCreateEntry(t, app, body)
	require.Equal(t, expectedStatus, resp.StatusCode)
	helpers.AssertErrorResponseBody(t, resp, utils.ErrorResponseBody{
		ClientOperation: utils.CreateEntry,
		Message:         expectedMessage,
		Detail:          expectedDetail,
		RequestBody:     body,
	})
}

func testCreateEntrySuccess(
	t *testing.T,
	app *fiber.App,
	db *gorm.DB,
	entryCount int64,
	secretCount int64,
	userSlug, vaultSlug, entryTitle,
	body string,
) {
	require.EqualValues(t, 8, entryCount)
	require.EqualValues(t, 16, secretCount)

	resp := newRequestCreateEntry(t, app, body)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)

	if respBody, err := io.ReadAll(resp.Body); err != nil {
		t.Fatalf("Read response body failed: %s", err.Error())
	} else {
		require.Empty(t, respBody)
	}

	var entry models.Entry
	helpers.QueryTestEntry(t, db, &entry, entryTitle)
	require.Equal(t, entry.CreatedAt, entry.UpdatedAt)
	require.Equal(t, entryTitle, entry.Title)
	require.Equal(t, userSlug, entry.UserSlug)
	require.Equal(t, vaultSlug, entry.VaultSlug)

	var secrets []models.Secret
	helpers.QueryTestSecretsByEntry(t, db, &secrets, entry.Slug)
	require.Equal(t, "secret[_label='username']@0.0.2.0", secrets[0].Label)
	require.Equal(t, "secret[_string='foodeater1234']@0.0.2.0", secrets[0].String)
	require.Equal(t, "secret[_label='password']@0.0.2.1", secrets[1].Label)
	require.Equal(t, "secret[_string='3a7!ng40oD']@0.0.2.1", secrets[1].String)

	helpers.CountEntries(t, db, &entryCount)
	require.EqualValues(t, 9, entryCount)
	helpers.CountSecrets(t, db, &secretCount)
	require.EqualValues(t, 18, secretCount)
}

func newRequestCreateEntry(t *testing.T, app *fiber.App, body string) *http.Response {
	reqBody := strings.NewReader(body)
	req := httptest.NewRequest(http.MethodPost, "/api/entries", reqBody)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
