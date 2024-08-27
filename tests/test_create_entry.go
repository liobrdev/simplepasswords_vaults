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

func testCreateEntry(t *testing.T, app *fiber.App, db *gorm.DB, conf *config.AppConfig) {
	bodyFmt := `{"user_slug":"%s","vault_slug":"%s","entry_title":"%s","secrets":%s}`

	t.Run("empty_body_400_bad_request", func(t *testing.T) {
		testCreateEntryClientError(
			t, app, conf, 400, utils.ErrorParse,
			"invalid character '\x00' looking for beginning of value", "",
		)
	})

	t.Run("array_body_400_bad_request", func(t *testing.T) {
		testCreateEntryClientError(
			t, app, conf, 400, utils.ErrorParse,
			"invalid character '[' looking for beginning of value", "[]",
		)

		testCreateEntryClientError(
			t, app, conf, 400, utils.ErrorParse,
			"invalid character '[' looking for beginning of value", "[{}]",
		)

		testCreateEntryClientError(
			t, app, conf, 400, utils.ErrorParse,
			"invalid character '[' looking for beginning of value", fmt.Sprintf(
				`[{"user_slug":"%s","vault_slug":"%s","entry_title":"%s","secrets":[]}]`,
				helpers.NewSlug(t), helpers.NewSlug(t), "entry@0.0.2.*",
			),
		)
	})

	t.Run("boolean_body_400_bad_request", func(t *testing.T) {
		testCreateEntryClientError(
			t, app, conf, 400, utils.ErrorParse,
			"invalid character 't' looking for beginning of value", "true",
		)

		testCreateEntryClientError(
			t, app, conf, 400, utils.ErrorParse,
			"invalid character 'f' looking for beginning of value", "false",
		)
	})

	t.Run("string_body_400_bad_request", func(t *testing.T) {
		testCreateEntryClientError(
			t, app, conf, 400, utils.ErrorParse,
			`invalid character '"' looking for beginning of value`, `"Valid JSON, but not an object."`,
		)
	})

	t.Run("null_body_400_bad_request", func(t *testing.T) {
		testCreateEntryClientError(t, app, conf, 400, utils.ErrorUserSlug, "", "null")
	})

	t.Run("empty_object_body_400_bad_request", func(t *testing.T) {
		testCreateEntryClientError(t, app, conf, 400, utils.ErrorUserSlug, "", "{}")
	})

	t.Run("missing_user_slug_400_bad_request", func(t *testing.T) {
		testCreateEntryClientError(
			t, app, conf, 400, utils.ErrorUserSlug, "", fmt.Sprintf(
				`{"userr_slug":"%s","vault_slug":"%s","entry_title":"%s","secrets":[]}`,
				"Spelled wrong!", helpers.NewSlug(t), "entry@0.0.2.*",
			),
		)
	})

	t.Run("missing_vault_slug_400_bad_request", func(t *testing.T) {
		testCreateEntryClientError(
			t, app, conf, 400, utils.ErrorVaultSlug, "", fmt.Sprintf(
				`{"user_slug":"%s","vualt_slug":"%s","entry_title":"%s","secrets":[]}`,
				helpers.NewSlug(t), "Spelled wrong!", "entry@0.0.2.*",
			),
		)
	})

	t.Run("missing_entry_title_400_bad_request", func(t *testing.T) {
		testCreateEntryClientError(
			t, app, conf, 400, utils.ErrorEntryTitle, "", fmt.Sprintf(
				`{"user_slug":"%s","vault_slug":"%s","enrty_title":"%s","secrets":[]}`,
				helpers.NewSlug(t), helpers.NewSlug(t), "Spelled wrong!",
			),
		)
	})

	t.Run("missing_secrets_400_bad_request", func(t *testing.T) {
		testCreateEntryClientError(
			t, app, conf, 400, utils.ErrorSecrets, "", fmt.Sprintf(
				`{"user_slug":"%s","vault_slug":"%s","entry_title":"%s","secerts":[]}`,
				helpers.NewSlug(t), helpers.NewSlug(t), "entry@0.0.2.*",
			),
		)
	})

	t.Run("null_user_slug_400_bad_request", func(t *testing.T) {
		testCreateEntryClientError(
			t, app, conf, 400, utils.ErrorUserSlug, "", fmt.Sprintf(
				`{"user_slug":%s,"vault_slug":"%s","entry_title":"%s","secrets":[]}`,
				"null", helpers.NewSlug(t), "entry@0.0.2.*",
			),
		)
	})

	t.Run("null_vault_slug_400_bad_request", func(t *testing.T) {
		testCreateEntryClientError(
			t, app, conf, 400, utils.ErrorVaultSlug, "", fmt.Sprintf(
				`{"user_slug":"%s","vault_slug":%s,"entry_title":"%s","secrets":[]}`,
				helpers.NewSlug(t), "null", "entry@0.0.2.*",
			),
		)
	})

	t.Run("null_entry_title_400_bad_request", func(t *testing.T) {
		testCreateEntryClientError(
			t, app, conf, 400, utils.ErrorEntryTitle, "", fmt.Sprintf(
				`{"user_slug":"%s","vault_slug":"%s","entry_title":%s,"secrets":[]}`,
				helpers.NewSlug(t), helpers.NewSlug(t), "null",
			),
		)
	})

	t.Run("null_secrets_400_bad_request", func(t *testing.T) {
		testCreateEntryClientError(
			t, app, conf, 400, utils.ErrorSecrets, "",
			fmt.Sprintf(bodyFmt, helpers.NewSlug(t), helpers.NewSlug(t), "entry@0.0.2.*", "null"),
		)
	})

	t.Run("empty_user_slug_400_bad_request", func(t *testing.T) {
		testCreateEntryClientError(
			t, app, conf, 400, utils.ErrorUserSlug, "",
			fmt.Sprintf(bodyFmt, "", helpers.NewSlug(t), "entry@0.0.2.*", "[]"),
		)
	})

	t.Run("empty_vault_slug_400_bad_request", func(t *testing.T) {
		testCreateEntryClientError(
			t, app, conf, 400, utils.ErrorVaultSlug, "",
			fmt.Sprintf(bodyFmt, helpers.NewSlug(t), "", "entry@0.0.2.*", "[]"),
		)
	})

	t.Run("empty_entry_title_400_bad_request", func(t *testing.T) {
		testCreateEntryClientError(
			t, app, conf, 400, utils.ErrorEntryTitle, "",
			fmt.Sprintf(bodyFmt, helpers.NewSlug(t), helpers.NewSlug(t), "", "[]"),
		)
	})

	t.Run("too_long_entry_title_400_bad_request", func(t *testing.T) {
		if title, err := utils.GenerateSlug(256); err != nil {
			t.Fatalf("Generate long string failed: %s", err.Error())
		} else {
			testCreateEntryClientError(
				t, app, conf, 400, utils.ErrorEntryTitle, "Too long",
				fmt.Sprintf(bodyFmt, helpers.NewSlug(t), helpers.NewSlug(t), title, "[]"),
			)
		}
	})

	t.Run("empty_item_in_secrets_body_400_bad_request", func(t *testing.T) {
		secretsStr := `[{` +
			`"secret_label":"secret[_label='username']@0.0.2.0",` +
			`"secret_string":"secret[_string='foodeater1234']@0.0.2.0",` +
			`"secret_priority":0` +
			`},{}]`

		testCreateEntryClientError(
			t, app, conf, 400, utils.ErrorItemSecrets, "secrets[1].Label; len(secrets) == 2",
			fmt.Sprintf(bodyFmt, helpers.NewSlug(t), helpers.NewSlug(t), "entry@0.0.2.*", secretsStr),
		)
	})

	t.Run("missing_secret_label_item_in_secrets_body_400_bad_request", func(t *testing.T) {
		secretsStr := `[{` +
			`"secret_label":"secret[_label='username']@0.0.2.0",` +
			`"secret_string":"secret[_string='foodeater1234']@0.0.2.0",` +
			`"secret_priority": 0` +
			`},{` +
			`"secret_string":"secret[_string='3a7!ng40oD']@0.0.2.1",` +
			`"secret_priority":1` +
			`}]`

		testCreateEntryClientError(
			t, app, conf, 400, utils.ErrorItemSecrets, "secrets[1].Label; len(secrets) == 2",
			fmt.Sprintf(bodyFmt, helpers.NewSlug(t), helpers.NewSlug(t), "entry@0.0.2.*", secretsStr),
		)
	})

	t.Run("missing_secret_string_item_in_secrets_body_400_bad_request", func(t *testing.T) {
		secretsStr := "[{" +
			`"secret_label":"secret[_label='username']@0.0.2.0",` +
			`"secret_string":"secret[_string='foodeater1234']@0.0.2.0",` +
			`"secret_priority":0` +
			`},{` +
			`"secret_label":"secret[_label='password']@0.0.2.1",` +
			`"secret_priority":1` +
			`}]`

		testCreateEntryClientError(
			t, app, conf, 400, utils.ErrorItemSecrets, "secrets[1].String; len(secrets) == 2",
			fmt.Sprintf(bodyFmt, helpers.NewSlug(t), helpers.NewSlug(t), "entry@0.0.2.*", secretsStr),
		)
	})

	t.Run("empty_secret_label_item_in_secrets_body_400_bad_request", func(t *testing.T) {
		secretsStr := `[{` +
			`"secret_label":"secret[_label='username']@0.0.2.0",` +
			`"secret_string":"secret[_string='foodeater1234']@0.0.2.0",` +
			`"secret_priority":0` +
			`},{` +
			`"secret_label":"",` +
			`"secret_string":"secret[_string='3a7!ng40oD']@0.0.2.1",` +
			`"secret_priority":1` +
			`}]`

		testCreateEntryClientError(
			t, app, conf, 400, utils.ErrorItemSecrets, "secrets[1].Label; len(secrets) == 2",
			fmt.Sprintf(bodyFmt, helpers.NewSlug(t), helpers.NewSlug(t), "entry@0.0.2.*", secretsStr),
		)
	})

	t.Run("empty_secret_string_item_in_secrets_body_400_bad_request", func(t *testing.T) {
		secretsStr := `[{` +
			`"secret_label":"secret[_label='username']@0.0.2.0",` +
			`"secret_string":"secret[_string='foodeater1234']@0.0.2.0",` +
			`"secret_priority":0` +
			`},{` +
			`"secret_label":"secret[_label='password']@0.0.2.1",` +
			`"secret_string":"",` +
			`"secret_priority":1` +
			`}]`

		testCreateEntryClientError(
			t, app, conf, 400, utils.ErrorItemSecrets, "secrets[1].String; len(secrets) == 2",
			fmt.Sprintf(bodyFmt, helpers.NewSlug(t), helpers.NewSlug(t), "entry@0.0.2.*", secretsStr),
		)
	})

	t.Run("too_long_secret_label_item_in_secrets_body_400_bad_request", func(t *testing.T) {
		if label, err := utils.GenerateSlug(256); err != nil {
			t.Fatalf("Generate long string failed: %s", err.Error())
		} else {
			secretsStr := `[{` +
				`"secret_label":"secret[_label='username']@0.0.2.0",` +
				`"secret_string":"secret[_string='foodeater1234']@0.0.2.0",` +
				`"secret_priority":0` +
				`},{` +
				fmt.Sprintf(`"secret_label":"%s",`, label) +
				`"secret_string":"secret[_string='3a7!ng40oD']@0.0.2.1",` +
				`"secret_priority":1` +
				`}]`

			testCreateEntryClientError(
				t, app, conf, 400, utils.ErrorItemSecrets, "secrets[1].Label; len(secrets) == 2",
				fmt.Sprintf(bodyFmt, helpers.NewSlug(t), helpers.NewSlug(t), "entry@0.0.2.*", secretsStr),
			)
		}
	})

	t.Run("too_long_secret_string_item_in_secrets_body_400_bad_request", func(t *testing.T) {
		if str, err := utils.GenerateSlug(1001); err != nil {
			t.Fatalf("Generate long string failed: %s", err.Error())
		} else {
			secretsStr := `[{` +
				`"secret_label":"secret[_label='username']@0.0.2.0",` +
				`"secret_string":"secret[_string='foodeater1234']@0.0.2.0",` +
				`"secret_priority":0` +
				`},{` +
				`"secret_label":"secret[_label='password']@0.0.2.1",` +
				fmt.Sprintf(`"secret_string":"%s",`, str) +
				`"secret_priority":1` +
				`}]`

			testCreateEntryClientError(
				t, app, conf, 400, utils.ErrorItemSecrets, "secrets[1].String; len(secrets) == 2",
				fmt.Sprintf(bodyFmt, helpers.NewSlug(t), helpers.NewSlug(t), "entry@0.0.2.*", secretsStr),
			)
		}
	})

	t.Run("valid_body_entry_title_already_exists_500_error", func(t *testing.T) {
		users, vaults, _, _ := setup.SetUpWithData(t, db)
		userSlug := users[0].Slug
		vaultSlug := vaults[0].Slug
		entryTitle := "entry@0.0.1.*"

		secretsStr := `[{` +
			`"secret_label":"secret[_label='username']@0.0.2.0",` +
			`"secret_string":"secret[_string='foodeater1234']@0.0.2.0",` +
			`"secret_priority":0` +
			`},{` +
			`"secret_label":"secret[_label='password']@0.0.2.1",` +
			`"secret_string":"secret[_string='3a7!ng40oD']@0.0.2.1",` +
			`"secret_priority":1` +
			`}]`

		body := fmt.Sprintf(bodyFmt, userSlug, vaultSlug, entryTitle, secretsStr)

		testCreateEntryClientError(
			t, app, conf, 500, utils.ErrorFailedDB,
			"UNIQUE constraint failed: entries.title, entries.vault_slug", body,
		)
	})

	t.Run("valid_body_entry_secret_label_duplicate_400_bad_request", func(t *testing.T) {
		secretLabel := "secret[_label='misc']@0.0.2.*"

		secretsStr := `[{` +
			fmt.Sprintf(`"secret_label":"%s",`, secretLabel) +
			`"secret_string":"secret[_string='foodeater1234']@0.0.2.0",` +
			`"secret_priority":0` +
			`},{` +
			fmt.Sprintf(`"secret_label":"%s",`, secretLabel) +
			`"secret_string":"secret[_string='3a7!ng40oD']@0.0.2.1",` +
			`"secret_priority":1` +
			`}]`

		body := fmt.Sprintf(
			bodyFmt, helpers.NewSlug(t), helpers.NewSlug(t), "entry@0.0.2.*", secretsStr,
		)

		testCreateEntryClientError(
			t, app, conf, 400, utils.ErrorDuplicateSecretsLabel, secretLabel, body,
		)
	})

	t.Run("valid_body_entry_secret_priority_duplicate_400_bad_request", func(t *testing.T) {
		secretsStr := `[{` +
			`"secret_label":"secret[_label='username']@0.0.2.0",` +
			`"secret_string":"secret[_string='foodeater1234']@0.0.2.0",` +
			`"secret_priority":1` +
			`},{` +
			`"secret_label":"secret[_label='password']@0.0.2.1",` +
			`"secret_string":"secret[_string='3a7!ng40oD']@0.0.2.1",` +
			`"secret_priority":1` +
			`}]`

		body := fmt.Sprintf(
			bodyFmt, helpers.NewSlug(t), helpers.NewSlug(t), "entry@0.0.2.*", secretsStr,
		)

		testCreateEntryClientError(
			t, app, conf, 400, utils.ErrorDuplicateSecretsPriority, "", body,
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
			`"secret_string":"secret[_string='foodeater1234']@0.0.2.0",` +
			`"secret_priority":0` +
			`},{` +
			`"secret_label":"secret[_label='password']@0.0.2.1",` +
			`"secret_string":"secret[_string='3a7!ng40oD']@0.0.2.1",` +
			`"secret_priority":1` +
			`}]`

		body := fmt.Sprintf(bodyFmt, userSlug, vault.Slug, entryTitle, secretsStr)

		testCreateEntrySuccess(
			t, app, db, conf, entryCount, secretCount, userSlug, vault.Slug, entryTitle, body,
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
			`"secret_priority":0,` +
			`"secret_slug":"notEvenARealSlug_0",` +
			`"secret_created_at":"10/12/22"` +
			`},{` +
			`"secret_label":"secret[_label='password']@0.0.2.1",` +
			`"secret_string":"secret[_string='3a7!ng40oD']@0.0.2.1",` +
			`"secret_priority":1,` +
			`"secret_slug":"notEvenARealSlug_1",` +
			`"secret_created_at":"10/12/22"` +
			`}]`

		body := fmt.Sprintf(
			`{"user_slug":"%s","vault_slug":"%s","entry_title":"%s","secrets":%s,"abc":123}`,
			userSlug, vault.Slug, entryTitle, secretsStr,
		)

		testCreateEntrySuccess(
			t, app, db, conf, entryCount, secretCount, userSlug, vault.Slug, entryTitle, body,
		)
	})
}

func testCreateEntryClientError(
	t *testing.T, app *fiber.App, conf *config.AppConfig, expectedStatus int,
	expectedMessage, expectedDetail, body string,
) {
	resp := newRequestCreateEntry(t, app, conf, body)
	require.Equal(t, expectedStatus, resp.StatusCode)
	helpers.AssertErrorResponseBody(t, resp, utils.ErrorResponseBody{
		ClientOperation: utils.CreateEntry,
		Message:         expectedMessage,
		Detail:          expectedDetail,
		RequestBody:     body,
	})
}

func testCreateEntrySuccess(
	t *testing.T, app *fiber.App, db *gorm.DB, conf *config.AppConfig,
	entryCount, secretCount int64,
	userSlug, vaultSlug, entryTitle, body string,
) {
	require.EqualValues(t, 8, entryCount)
	require.EqualValues(t, 16, secretCount)

	resp := newRequestCreateEntry(t, app, conf, body)
	require.Equal(t, 204, resp.StatusCode)

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
	require.EqualValues(t, 0, secrets[0].Priority)
	require.Equal(t, "secret[_label='password']@0.0.2.1", secrets[1].Label)
	require.Equal(t, "secret[_string='3a7!ng40oD']@0.0.2.1", secrets[1].String)
	require.EqualValues(t, 1, secrets[1].Priority)

	helpers.CountEntries(t, db, &entryCount)
	require.EqualValues(t, 9, entryCount)
	helpers.CountSecrets(t, db, &secretCount)
	require.EqualValues(t, 18, secretCount)
}

func newRequestCreateEntry(
	t *testing.T, app *fiber.App, conf *config.AppConfig, body string,
) *http.Response {

	reqBody := strings.NewReader(body)
	req := httptest.NewRequest("POST", "/api/entries", reqBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Client-Operation", utils.CreateEntry)
	req.Header.Set("Authorization", "Token " + conf.VAULTS_ACCESS_TOKEN)
	resp, err := app.Test(req, -1)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
