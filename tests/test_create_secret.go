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

func testCreateSecret(t *testing.T, app *fiber.App, db *gorm.DB) {
	bodyFmt := `{` +
		`"user_slug":"%s",` +
		`"vault_slug":"%s",` +
		`"entry_slug":"%s",` +
		`"secret_label":"%s",` +
		`"secret_string":"%s"` +
		`}`

	t.Run("empty_body_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, db, "", http.StatusBadRequest, utils.ErrorParse,
			"invalid character '\x00' looking for beginning of value",
		)
	})

	t.Run("array_body_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, db, "[]", http.StatusBadRequest, utils.ErrorParse,
			"invalid character '[' looking for beginning of value",
		)

		testCreateSecretClientError(
			t, app, db, "[{}]", http.StatusBadRequest, utils.ErrorParse,
			"invalid character '[' looking for beginning of value",
		)

		testCreateSecretClientError(
			t, app, db, fmt.Sprintf(
				`[`+bodyFmt+`]`,
				helpers.NewSlug(t), helpers.NewSlug(t), helpers.NewSlug(t), "abc", "123",
			), http.StatusBadRequest, utils.ErrorParse,
			"invalid character '[' looking for beginning of value",
		)
	})

	t.Run("null_body_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, db, "null", http.StatusBadRequest, utils.ErrorUserSlug, "",
		)
	})

	t.Run("boolean_body_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, db, "true", http.StatusBadRequest, utils.ErrorParse,
			"invalid character 't' looking for beginning of value",
		)

		testCreateSecretClientError(
			t, app, db, "false", http.StatusBadRequest, utils.ErrorParse,
			"invalid character 'f' looking for beginning of value",
		)
	})

	t.Run("string_body_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, db, "\"Valid JSON, but not an object.\"", http.StatusBadRequest,
			utils.ErrorParse, "invalid character '\"' looking for beginning of value",
		)
	})

	t.Run("empty_object_body_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, db, "{}", http.StatusBadRequest, utils.ErrorUserSlug, "",
		)
	})

	t.Run("missing_user_slug_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, db, fmt.Sprintf(
				`{`+
					`"usr_slug":"Spelled wrong!",`+
					`"vault_slug":"%s",`+
					`"entry_slug":"%s",`+
					`"secret_label":"%s",`+
					`"secret_string":"%s"`+
					`}`,
				helpers.NewSlug(t), helpers.NewSlug(t), "abc", "123",
			), http.StatusBadRequest, utils.ErrorUserSlug, "",
		)
	})

	t.Run("missing_vault_slug_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, db, fmt.Sprintf(
				`{`+
					`"user_slug":"%s",`+
					`"vualt_slug":"Spelled wrong!",`+
					`"entry_slug":"%s",`+
					`"secret_label":"%s",`+
					`"secret_string":"%s"`+
					`}`,
				helpers.NewSlug(t), helpers.NewSlug(t), "abc", "123",
			), http.StatusBadRequest, utils.ErrorVaultSlug, "",
		)
	})

	t.Run("missing_entry_slug_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, db, fmt.Sprintf(
				`{`+
					`"user_slug":"%s",`+
					`"vault_slug":"%s",`+
					`"enry_slug":"Spelled wrong!",`+
					`"secret_label":"%s",`+
					`"secret_string":"%s"`+
					`}`,
				helpers.NewSlug(t), helpers.NewSlug(t), "abc", "123",
			), http.StatusBadRequest, utils.ErrorEntrySlug, "",
		)
	})

	t.Run("missing_secret_label_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, db, fmt.Sprintf(
				`{`+
					`"user_slug":"%s",`+
					`"vault_slug":"%s",`+
					`"entry_slug":"%s",`+
					`"secret_labl":"Spelled wrong!",`+
					`"secret_string":"%s"`+
					`}`,
				helpers.NewSlug(t), helpers.NewSlug(t), helpers.NewSlug(t), "123",
			), http.StatusBadRequest, utils.ErrorSecretLabel, "",
		)
	})

	t.Run("missing_secret_string_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, db, fmt.Sprintf(
				`{`+
					`"user_slug":"%s",`+
					`"vault_slug":"%s",`+
					`"entry_slug":"%s",`+
					`"secret_label":"%s",`+
					`"secret_strng":"Spelled wrong!"`+
					`}`,
				helpers.NewSlug(t), helpers.NewSlug(t), helpers.NewSlug(t), "abc",
			), http.StatusBadRequest, utils.ErrorSecretString, "",
		)
	})

	t.Run("null_user_slug_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, db, fmt.Sprintf(
				`{`+
					`"user_slug":null,`+
					`"vault_slug":"%s",`+
					`"entry_slug":"%s",`+
					`"secret_label":"%s",`+
					`"secret_string":"%s"`+
					`}`,
				helpers.NewSlug(t), helpers.NewSlug(t), "abc", "123",
			), http.StatusBadRequest, utils.ErrorUserSlug, "",
		)
	})

	t.Run("null_vault_slug_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, db, fmt.Sprintf(
				`{`+
					`"user_slug":"%s",`+
					`"vault_slug":null,`+
					`"entry_slug":"%s",`+
					`"secret_label":"%s",`+
					`"secret_string":"%s"`+
					`}`,
				helpers.NewSlug(t), helpers.NewSlug(t), "abc", "123",
			), http.StatusBadRequest, utils.ErrorVaultSlug, "",
		)
	})

	t.Run("null_entry_slug_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, db, fmt.Sprintf(
				`{`+
					`"user_slug":"%s",`+
					`"vault_slug":"%s",`+
					`"entry_slug":null,`+
					`"secret_label":"%s",`+
					`"secret_string":"%s"`+
					`}`,
				helpers.NewSlug(t), helpers.NewSlug(t), "abc", "123",
			), http.StatusBadRequest, utils.ErrorEntrySlug, "",
		)
	})

	t.Run("null_secret_label_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, db, fmt.Sprintf(
				`{`+
					`"user_slug":"%s",`+
					`"vault_slug":"%s",`+
					`"entry_slug":"%s",`+
					`"secret_label":null,`+
					`"secret_string":"%s"`+
					`}`,
				helpers.NewSlug(t), helpers.NewSlug(t), helpers.NewSlug(t), "123",
			), http.StatusBadRequest, utils.ErrorSecretLabel, "",
		)
	})

	t.Run("null_secret_string_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, db, fmt.Sprintf(
				`{`+
					`"user_slug":"%s",`+
					`"vault_slug":"%s",`+
					`"entry_slug":"%s",`+
					`"secret_label":"%s",`+
					`"secret_string":null`+
					`}`,
				helpers.NewSlug(t), helpers.NewSlug(t), helpers.NewSlug(t), "abc",
			), http.StatusBadRequest, utils.ErrorSecretString, "",
		)
	})

	t.Run("empty_user_slug_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, db, fmt.Sprintf(
				bodyFmt, "", helpers.NewSlug(t), helpers.NewSlug(t), "abc", "123",
			), http.StatusBadRequest, utils.ErrorUserSlug, "",
		)
	})

	t.Run("empty_vault_slug_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, db, fmt.Sprintf(
				bodyFmt, helpers.NewSlug(t), "", helpers.NewSlug(t), "abc", "123",
			), http.StatusBadRequest, utils.ErrorVaultSlug, "",
		)
	})

	t.Run("empty_entry_slug_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, db, fmt.Sprintf(
				bodyFmt,
				helpers.NewSlug(t), helpers.NewSlug(t), "", "abc", "123",
			), http.StatusBadRequest, utils.ErrorEntrySlug, "",
		)
	})

	t.Run("empty_secret_label_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, db, fmt.Sprintf(
				bodyFmt,
				helpers.NewSlug(t), helpers.NewSlug(t), helpers.NewSlug(t), "", "123",
			), http.StatusBadRequest, utils.ErrorSecretLabel, "",
		)
	})

	t.Run("empty_secret_string_400_bad_request", func(t *testing.T) {
		testCreateSecretClientError(
			t, app, db, fmt.Sprintf(
				bodyFmt,
				helpers.NewSlug(t), helpers.NewSlug(t), helpers.NewSlug(t), "abc", "",
			), http.StatusBadRequest, utils.ErrorSecretString, "",
		)
	})

	t.Run("too_long_secret_label_400_bad_request", func(t *testing.T) {
		// `label` is a random string greater than 255 characters in length
		if label, err := utils.GenerateSlug(256); err != nil {
			t.Fatalf("Generate long string failed: %s", err.Error())
		} else {
			testCreateSecretClientError(
				t, app, db, fmt.Sprintf(
					bodyFmt,
					helpers.NewSlug(t), helpers.NewSlug(t), helpers.NewSlug(t), label, "123",
				), http.StatusBadRequest, utils.ErrorSecretLabel, "Too long (256 > 255)",
			)
		}
	})

	t.Run("too_long_secret_string_400_bad_request", func(t *testing.T) {
		// `str` is a random string greater than 1000 characters in length
		if str, err := utils.GenerateSlug(1001); err != nil {
			t.Fatalf("Generate long string failed: %s", err.Error())
		} else {
			testCreateSecretClientError(
				t, app, db, fmt.Sprintf(
					bodyFmt,
					helpers.NewSlug(t), helpers.NewSlug(t), helpers.NewSlug(t), "abc", str,
				), http.StatusBadRequest, utils.ErrorSecretString, "Too long (1001 > 1000)",
			)
		}
	})

	t.Run("valid_body_secret_label_already_exists_409_conflict", func(t *testing.T) {
		users, vaults, entries, secrets := setup.SetUpWithData(t, db)
		userSlug := (*users)[0].Slug
		vaultSlug := (*vaults)[1].Slug
		entrySlug := (*entries)[3].Slug

		testCreateSecretClientError(
			t, app, db,
			fmt.Sprintf(bodyFmt, userSlug, vaultSlug, entrySlug, (*secrets)[7].Label, "123"),
			http.StatusConflict, utils.ErrorFailedDB,
			"UNIQUE constraint failed: secrets.label, secrets.entry_slug",
		)
	})

	t.Run("valid_body_204_no_content", func(t *testing.T) {
		users, vaults, entries, _ := setup.SetUpWithData(t, db)
		userSlug := (*users)[0].Slug
		vaultSlug := (*vaults)[1].Slug
		entrySlug := (*entries)[3].Slug
		secretLabel := "secret[_label='email']@0.1.1.2"
		secretString := "secret[_string='food.eater@email.dev']@0.1.1.2"

		testCreateSecretSuccess(
			t, app, db, userSlug, vaultSlug, entrySlug, secretLabel, secretString,
			fmt.Sprintf(bodyFmt, userSlug, vaultSlug, entrySlug, secretLabel, secretString),
		)
	})

	t.Run("valid_body_irrelevant_data_204_no_content", func(t *testing.T) {
		users, vaults, entries, _ := setup.SetUpWithData(t, db)
		userSlug := (*users)[0].Slug
		vaultSlug := (*vaults)[1].Slug
		entrySlug := (*entries)[3].Slug
		secretLabel := "secret[_label='email']@0.1.1.2"
		secretString := "secret[_string='food.eater@email.dev']@0.1.1.2"

		validBodyIrrelevantData := fmt.Sprintf(
			`{`+
				`"user_slug":"%s",`+
				`"vault_slug":"%s",`+
				`"entry_slug":"%s",`+
				`"secret_label":"%s",`+
				`"secret_string":"%s",`+
				`"secret_slug":"notEvenARealSlug",`+
				`"secret_created_at":"10/12/22"`+
				`}`,
			userSlug, vaultSlug, entrySlug, secretLabel, secretString,
		)

		testCreateSecretSuccess(
			t, app, db, userSlug, vaultSlug, entrySlug, secretLabel, secretString,
			validBodyIrrelevantData,
		)
	})
}

func testCreateSecretClientError(
	t *testing.T,
	app *fiber.App,
	db *gorm.DB,
	body string,
	expectedStatus int,
	expectedMessage utils.ErrorMessage,
	expectedDetail string,
) {
	resp := newRequestCreateSecret(t, app, body)
	require.Equal(t, expectedStatus, resp.StatusCode)
	helpers.AssertErrorResponseBody(t, resp, utils.ErrorResponseBody{
		ClientOperation: utils.CreateSecret,
		Message:         string(expectedMessage),
		Detail:          expectedDetail,
		RequestBody:     body,
	})
}

func testCreateSecretSuccess(
	t *testing.T,
	app *fiber.App,
	db *gorm.DB,
	userSlug, vaultSlug, entrySlug, secretLabel, secretString, body string,
) {
	var secretCount int64
	helpers.CountSecrets(t, db, &secretCount)
	require.EqualValues(t, 16, secretCount)

	resp := newRequestCreateSecret(t, app, body)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)

	if respBody, err := io.ReadAll(resp.Body); err != nil {
		t.Fatalf("Read response body failed: %s", err.Error())
	} else {
		require.Empty(t, respBody)
	}

	var secret models.Secret
	helpers.QueryTestSecretByLabel(t, db, &secret, secretLabel)
	require.Equal(t, secret.CreatedAt, secret.UpdatedAt)
	require.Equal(t, secretLabel, secret.Label)
	require.Equal(t, secretString, secret.String)
	require.Equal(t, userSlug, secret.UserSlug)
	require.Equal(t, vaultSlug, secret.VaultSlug)
	require.Equal(t, entrySlug, secret.EntrySlug)
	helpers.CountSecrets(t, db, &secretCount)
	require.EqualValues(t, 17, secretCount)
}

func newRequestCreateSecret(t *testing.T, app *fiber.App, body string) *http.Response {
	reqBody := strings.NewReader(body)
	req := httptest.NewRequest(http.MethodPost, "/api/secrets", reqBody)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
