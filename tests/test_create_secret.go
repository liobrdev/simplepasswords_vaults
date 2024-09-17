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

func testCreateSecret(t *testing.T, app *fiber.App, db *gorm.DB, conf *config.AppConfig) {
	bodyFmt := `{` +
		`"user_slug":"%s",` +
		`"vault_slug":"%s",` +
		`"entry_slug":"%s",` +
		`"secret_label":"%s",` +
		`"secret_string":"%s"` +
		`}`

	dummySlug := helpers.NewSlug(t)

	t.Run("empty_body_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, conf, 400, utils.ErrorParse,
			"invalid character '\x00' looking for beginning of value", "",
		)
	})

	t.Run("array_body_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, conf, 400, utils.ErrorParse, "invalid character '[' looking for beginning of value",
			"[]",
		)

		testCreateSecretClientError(
			t, app, conf, 400, utils.ErrorParse, "invalid character '[' looking for beginning of value",
			"[{}]",
		)

		testCreateSecretClientError(
			t, app, conf, 400, utils.ErrorParse, "invalid character '[' looking for beginning of value",
			fmt.Sprintf("[" + bodyFmt + "]", dummySlug, dummySlug, dummySlug, "abc", "123"),
		)
	})

	t.Run("boolean_body_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, conf, 400, utils.ErrorParse, "invalid character 't' looking for beginning of value",
			"true",
		)

		testCreateSecretClientError(
			t, app, conf, 400, utils.ErrorParse, "invalid character 'f' looking for beginning of value",
			"false",
		)
	})

	t.Run("string_body_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, conf, 400, utils.ErrorParse, `invalid character '"' looking for beginning of value`,
			`"Valid JSON, but not an object."`,
		)
	})

	t.Run("null_body_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(t, app, conf, 400, utils.ErrorUserSlug, "", "null")
	})

	t.Run("empty_object_body_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(t, app, conf, 400, utils.ErrorUserSlug, "", "{}")
	})

	t.Run("missing_user_slug_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, conf, 400, utils.ErrorUserSlug, "", fmt.Sprintf(
				`{` +
					`"usr_slug":"Spelled wrong!",` +
					`"vault_slug":"%s",` +
					`"entry_slug":"%s",` +
					`"secret_label":"%s",` +
					`"secret_string":"%s"` +
					`}`,
				dummySlug, dummySlug, "abc", "123",
			),
		)
	})

	t.Run("missing_vault_slug_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, conf, 400, utils.ErrorVaultSlug, "", fmt.Sprintf(
				`{` +
					`"user_slug":"%s",` +
					`"vualt_slug":"Spelled wrong!",` +
					`"entry_slug":"%s",` +
					`"secret_label":"%s",` +
					`"secret_string":"%s"` +
					`}`,
				dummySlug, dummySlug, "abc", "123",
			),
		)
	})

	t.Run("missing_entry_slug_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, conf, 400, utils.ErrorEntrySlug, "", fmt.Sprintf(
				`{` +
					`"user_slug":"%s",` +
					`"vault_slug":"%s",` +
					`"enry_slug":"Spelled wrong!",` +
					`"secret_label":"%s",` +
					`"secret_string":"%s"` +
					`}`,
				dummySlug, dummySlug, "abc", "123",
			),
		)
	})

	t.Run("missing_secret_label_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, conf, 400, utils.ErrorSecretLabel, "", fmt.Sprintf(
				`{` +
					`"user_slug":"%s",` +
					`"vault_slug":"%s",` +
					`"entry_slug":"%s",` +
					`"secret_labl":"Spelled wrong!",` +
					`"secret_string":"%s"` +
					`}`,
				dummySlug, dummySlug, dummySlug, "123",
			),
		)
	})

	t.Run("missing_secret_string_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, conf, 400, utils.ErrorSecretString, "", fmt.Sprintf(
				`{` +
					`"user_slug":"%s",` +
					`"vault_slug":"%s",` +
					`"entry_slug":"%s",` +
					`"secret_label":"%s",` +
					`"secret_strng":"Spelled wrong!"` +
					`}`,
				dummySlug, dummySlug, dummySlug, "abc",
			),
		)
	})

	t.Run("null_user_slug_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, conf, 400, utils.ErrorUserSlug, "", fmt.Sprintf(
				`{` +
					`"user_slug":null,` +
					`"vault_slug":"%s",` +
					`"entry_slug":"%s",` +
					`"secret_label":"%s",` +
					`"secret_string":"%s"` +
					`}`,
				dummySlug, dummySlug, "abc", "123",
			),
		)
	})

	t.Run("null_vault_slug_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, conf, 400, utils.ErrorVaultSlug, "", fmt.Sprintf(
				`{` +
					`"user_slug":"%s",` +
					`"vault_slug":null,` +
					`"entry_slug":"%s",` +
					`"secret_label":"%s",` +
					`"secret_string":"%s"` +
					`}`,
				dummySlug, dummySlug, "abc", "123",
			),
		)
	})

	t.Run("null_entry_slug_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, conf, 400, utils.ErrorEntrySlug, "", fmt.Sprintf(
				`{` +
					`"user_slug":"%s",` +
					`"vault_slug":"%s",` +
					`"entry_slug":null,` +
					`"secret_label":"%s",` +
					`"secret_string":"%s"` +
					`}`,
				dummySlug, dummySlug, "abc", "123",
			),
		)
	})

	t.Run("null_secret_label_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, conf, 400, utils.ErrorSecretLabel, "", fmt.Sprintf(
				`{` +
					`"user_slug":"%s",` +
					`"vault_slug":"%s",` +
					`"entry_slug":"%s",` +
					`"secret_label":null,` +
					`"secret_string":"%s"` +
					`}`,
				dummySlug, dummySlug, dummySlug, "123",
			),
		)
	})

	t.Run("null_secret_string_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, conf, 400, utils.ErrorSecretString, "", fmt.Sprintf(
				`{` +
					`"user_slug":"%s",` +
					`"vault_slug":"%s",` +
					`"entry_slug":"%s",` +
					`"secret_label":"%s",` +
					`"secret_string":null` +
					`}`,
				dummySlug, dummySlug, dummySlug, "abc",
			),
		)
	})

	t.Run("empty_user_slug_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, conf, 400, utils.ErrorUserSlug, "",
			fmt.Sprintf(bodyFmt, "", dummySlug, dummySlug, "abc", "123"),
		)
	})

	t.Run("empty_vault_slug_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, conf, 400, utils.ErrorVaultSlug, "",
			fmt.Sprintf(bodyFmt, dummySlug, "", dummySlug, "abc", "123"),
		)
	})

	t.Run("empty_entry_slug_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, conf, 400, utils.ErrorEntrySlug, "",
			fmt.Sprintf(bodyFmt, dummySlug, dummySlug, "", "abc", "123"),
		)
	})

	t.Run("empty_secret_label_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, conf, 400, utils.ErrorSecretLabel, "",
			fmt.Sprintf(bodyFmt, dummySlug, dummySlug, dummySlug, "", "123"),
		)
	})

	t.Run("empty_secret_string_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, conf, 400, utils.ErrorSecretString, "",
			fmt.Sprintf(bodyFmt, dummySlug, dummySlug, dummySlug, "abc", ""),
		)
	})

	t.Run("too_long_secret_label_400_bad_request", func(t *testing.T) {
		if label, err := utils.GenerateSlug(256); err != nil {
			t.Fatalf("Generate long string failed: %s", err.Error())
		} else {
			testCreateSecretClientError(
				t, app, conf, 400, utils.ErrorSecretLabel, "Too long",
				fmt.Sprintf(bodyFmt, dummySlug, dummySlug, dummySlug, label, "123"),
			)
		}
	})

	t.Run("too_long_secret_string_400_bad_request", func(t *testing.T) {
		if str, err := utils.GenerateSlug(1001); err != nil {
			t.Fatalf("Generate long string failed: %s", err.Error())
		} else {
			testCreateSecretClientError(
				t, app, conf, 400, utils.ErrorSecretString, "Too long",
				fmt.Sprintf(bodyFmt, dummySlug, dummySlug, dummySlug, "abc", str),
			)
		}
	})

	t.Run("valid_body_secret_label_already_exists_500_error", func(t *testing.T) {
		users, vaults, entries, secrets := setup.SetUpWithData(t, db)
		userSlug := users[0].Slug
		vaultSlug := vaults[1].Slug
		entrySlug := entries[3].Slug

		testCreateSecretClientError(
			t, app, conf, 500, utils.ErrorFailedDB,
			"UNIQUE constraint failed: secrets.label, secrets.entry_slug",
			fmt.Sprintf(bodyFmt, userSlug, vaultSlug, entrySlug, secrets[7].Label, "123"),
		)
	})

	t.Run("valid_body_204_no_content", func(t *testing.T) {
		users, vaults, entries, _ := setup.SetUpWithData(t, db)
		userSlug := users[0].Slug
		vaultSlug := vaults[1].Slug
		entrySlug := entries[3].Slug
		secretLabel := "secret[_label='email']@0.1.1.2"
		secretString := "secret[_string='food.eater@email.dev']@0.1.1.2"

		testCreateSecretSuccess(
			t, app, db, conf, 2, secretLabel, secretString,
			fmt.Sprintf(bodyFmt, userSlug, vaultSlug, entrySlug, secretLabel, secretString),
		)
	})

	t.Run("valid_body_irrelevant_data_204_no_content", func(t *testing.T) {
		users, vaults, entries, _ := setup.SetUpWithData(t, db)
		userSlug := users[0].Slug
		vaultSlug := vaults[1].Slug
		entrySlug := entries[3].Slug
		secretLabel := "secret[_label='email']@0.1.1.2"
		secretString := "secret[_string='food.eater@email.dev']@0.1.1.2"

		validBodyIrrelevantData := fmt.Sprintf(
			`{`+
				`"user_slug":"%s",`+
				`"vault_slug":"%s",`+
				`"entry_slug":"%s",`+
				`"secret_label":"%s",`+
				`"secret_string":"%s",`+
				`"secret_slug":"notARealSlug",`+
				`"secret_created_at":"10/12/22"`+
				`}`,
			userSlug, vaultSlug, entrySlug, secretLabel, secretString,
		)

		testCreateSecretSuccess(t, app, db, conf, 2, secretLabel, secretString, validBodyIrrelevantData)
	})
}

func testCreateSecretClientError(
	t *testing.T, app *fiber.App, conf *config.AppConfig, expectedStatus int,
	expectedMessage, expectedDetail, body string,
) {
	resp := newRequestCreateSecret(t, app, conf, body)
	require.Equal(t, expectedStatus, resp.StatusCode)
	helpers.AssertErrorResponseBody(t, resp, utils.ErrorResponseBody{
		ClientOperation: utils.CreateSecret,
		Message:         expectedMessage,
		Detail:          expectedDetail,
		RequestBody:     body,
	})
}

func testCreateSecretSuccess(
	t *testing.T, app *fiber.App, db *gorm.DB, conf *config.AppConfig,
	secretPriority uint8, secretLabel, secretString, body string,
) {
	var secretCount int64
	helpers.CountSecrets(t, db, &secretCount)
	require.EqualValues(t, 20, secretCount)

	resp := newRequestCreateSecret(t, app, conf, body)
	require.Equal(t, 204, resp.StatusCode)

	if respBody, err := io.ReadAll(resp.Body); err != nil {
		t.Fatalf("Read response body failed: %s", err.Error())
	} else {
		require.Empty(t, respBody)
	}

	var secret models.Secret
	helpers.QueryTestSecretByLabel(t, db, &secret, secretLabel)

	if plaintext, err := utils.Decrypt(secret.String, helpers.HexHash[:64]);
	err != nil {
		t.Fatalf("Password decryption failed: %s", err.Error())
	} else {
		require.Equal(t, secretString, plaintext)
	}

	require.Equal(t, secretLabel, secret.Label)
	require.Equal(t, secretPriority, secret.Priority)
	require.Equal(t, secret.CreatedAt, secret.UpdatedAt)
	helpers.CountSecrets(t, db, &secretCount)
	require.EqualValues(t, 21, secretCount)
}

func newRequestCreateSecret(
	t *testing.T, app *fiber.App, conf *config.AppConfig, body string,
) *http.Response {

	reqBody := strings.NewReader(body)
	req := httptest.NewRequest("POST", "/api/secrets", reqBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Client-Operation", utils.CreateSecret)
	req.Header.Set("Authorization", "Token " + conf.VAULTS_ACCESS_TOKEN)
	req.Header.Set(conf.PASSWORD_HEADER_KEY, helpers.HexHash[:64])

	resp, err := app.Test(req, -1)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
