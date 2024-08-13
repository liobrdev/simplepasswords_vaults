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

func testUpdateSecret(t *testing.T, app *fiber.App, db *gorm.DB, conf *config.AppConfig) {
	dummySlug := helpers.NewSlug(t)

	t.Run("empty_body_400_bad_request", func(t *testing.T) {
		testUpdateSecretClientError(
			t, app, conf, 400, utils.ErrorParse,
			"invalid character '\x00' looking for beginning of value", dummySlug, "",
		)
	})

	t.Run("array_body_400_bad_request", func(t *testing.T) {
		testUpdateSecretClientError(
			t, app, conf, 400, utils.ErrorParse, "invalid character '[' looking for beginning of value",
			dummySlug, "[]",
		)

		testUpdateSecretClientError(
			t, app, conf, 400, utils.ErrorParse, "invalid character '[' looking for beginning of value",
			dummySlug, "[{}]",
		)

		testUpdateSecretClientError(
			t, app, conf, 400, utils.ErrorParse, "invalid character '[' looking for beginning of value",
			dummySlug, `[{"secret_label":"abc","secret_string":"123"}]`,
		)
	})

	t.Run("boolean_body_400_bad_request", func(t *testing.T) {
		testUpdateSecretClientError(
			t, app, conf, 400, utils.ErrorParse, "invalid character 't' looking for beginning of value",
			dummySlug, "true",
		)

		testUpdateSecretClientError(
			t, app, conf, 400, utils.ErrorParse, "invalid character 'f' looking for beginning of value",
			dummySlug, "false",
		)
	})

	t.Run("string_body_400_bad_request", func(t *testing.T) {
		testUpdateSecretClientError(
			t, app, conf, 400, utils.ErrorParse, `invalid character '"' looking for beginning of value`,
			dummySlug, `"Valid JSON, but not an object."`,
		)
	})

	t.Run("null_or_empty_object_400_bad_request", func(t *testing.T) {
		_, _, _, secrets := setup.SetUpWithData(t, db)

		testUpdateSecretClientError(
			t, app, conf, 400, utils.ErrorEmptyUpdateSecret, "Null or empty object or fields.",
			secrets[0].Slug, "null",
		)

		testUpdateSecretClientError(
			t, app, conf, 400, utils.ErrorEmptyUpdateSecret, "Null or empty object or fields.",
			secrets[0].Slug, "{}",
		)
	})

	t.Run("both_fields_null_or_empty_or_missing_400_bad_request", func(t *testing.T) {
		_, _, _, secrets := setup.SetUpWithData(t, db)

		testUpdateSecretClientError(
			t, app, conf, 400, utils.ErrorEmptyUpdateSecret, "Null or empty object or fields.",
			secrets[0].Slug, `{"secret_label":"","secret_string":""}`,
		)

		testUpdateSecretClientError(
			t, app, conf, 400, utils.ErrorEmptyUpdateSecret, "Null or empty object or fields.",
			secrets[0].Slug, `{"secret_label":null,"secret_string":null}`,
		)

		testUpdateSecretClientError(
			t, app, conf, 400, utils.ErrorEmptyUpdateSecret, "Null or empty object or fields.",
			secrets[0].Slug, `{"weird_label":"abc","weird_string":"123"}`,
		)
	})

	t.Run("too_long_secret_label_400_bad_request", func(t *testing.T) {
		if label, err := utils.GenerateSlug(256); err != nil {
			t.Fatalf("Generate long string failed: %s", err.Error())
		} else {
			testUpdateSecretClientError(
				t, app, conf, 400, utils.ErrorSecretLabel, "Too long", dummySlug,
				fmt.Sprintf(`{"secret_label":"%s"}`, label),
			)
		}
	})

	t.Run("too_long_secret_string_400_bad_request", func(t *testing.T) {
		if str, err := utils.GenerateSlug(1001); err != nil {
			t.Fatalf("Generate long string failed: %s", err.Error())
		} else {
			testUpdateSecretClientError(
				t, app, conf, 400, utils.ErrorSecretString, "Too long", dummySlug,
				fmt.Sprintf(`{"secret_string":"%s"}`, str),
			)
		}
	})

	t.Run("valid_body_secret_label_already_exists_500_error", func(t *testing.T) {
		_, _, _, secrets := setup.SetUpWithData(t, db)

		testUpdateSecretClientError(
			t, app, conf, 500, utils.ErrorFailedDB,
			"UNIQUE constraint failed: secrets.label, secrets.entry_slug", secrets[0].Slug,
			fmt.Sprintf(`{"secret_label":"%s"}`, secrets[1].Label),
		)
	})

	t.Run("valid_body_404_not_found", func(t *testing.T) {
		setup.SetUpWithData(t, db)
		testUpdateSecretClientError(
			t, app, conf, 404, utils.ErrorNoRowsAffected, "Likely that slug was not found.",
			dummySlug, `{"secret_label":"updated_secret"}`,
		)
	})

	t.Run("valid_body_204_no_content", func(t *testing.T) {
		_, _, _, secrets := setup.SetUpWithData(t, db)
		slug := secrets[0].Slug

		updatedSecretLabel := "secret[_label='updated_label']@0.0.0.0"

		testUpdateSecretSuccess(
			t, app, db, conf, slug, updatedSecretLabel, "",
			fmt.Sprintf(`{"secret_label":"%s"}`, updatedSecretLabel),
		)

		updatedSecretString := "secret[_string='updated_string']@0.0.0.0"

		testUpdateSecretSuccess(
			t, app, db, conf, slug, "", updatedSecretString,
			fmt.Sprintf(`{"secret_string":"%s"}`, updatedSecretString),
		)

		updatedSecretLabel = "secret[_label='updated_again_label']@0.0.0.0"
		updatedSecretString = "secret[_string='updated_again_string']@0.0.0.0"

		testUpdateSecretSuccess(
			t, app, db, conf, slug, updatedSecretLabel, updatedSecretString, fmt.Sprintf(
				`{"secret_label":"%s","secret_string":"%s"}`, updatedSecretLabel, updatedSecretString,
			),
		)
	})

	t.Run("valid_body_irrelevant_data_204_no_content", func(t *testing.T) {
		_, _, _, secrets := setup.SetUpWithData(t, db)
		slug := secrets[0].Slug

		updatedSecretLabel := "secret[_label='updated_label']@0.0.0.0"
		updatedSecretString := "secret[_string='updated_string']@0.0.0.0"

		testUpdateSecretSuccess(
			t, app, db, conf, slug, updatedSecretLabel, updatedSecretString, fmt.Sprintf(
				`{"secret_label":"%s","secret_string":"%s","is_real_field":false,"a":1}`,
				updatedSecretLabel, updatedSecretString,
			),
		)
	})
}

func testUpdateSecretClientError(
	t *testing.T, app *fiber.App, conf *config.AppConfig, expectedStatus int,
	expectedMessage, expectedDetail, slug, body string,
) {
	resp := newRequestUpdateSecret(t, app, conf, slug, body)
	require.Equal(t, expectedStatus, resp.StatusCode)
	helpers.AssertErrorResponseBody(t, resp, utils.ErrorResponseBody{
		ClientOperation: utils.UpdateSecret,
		Message:         expectedMessage,
		Detail:          expectedDetail,
		RequestBody:     body,
	})
}

func testUpdateSecretSuccess(
	t *testing.T, app *fiber.App, db *gorm.DB, conf *config.AppConfig,
	slug, updatedSecretLabel, updatedSecretString, body string,
) {
	var secretBeforeUpdate models.Secret
	helpers.QueryTestSecretBySlug(t, db, &secretBeforeUpdate, slug)

	if updatedSecretLabel == "" {
		updatedSecretLabel = secretBeforeUpdate.Label
	}

	if updatedSecretString == "" {
		updatedSecretString = secretBeforeUpdate.String
	}

	resp := newRequestUpdateSecret(t, app, conf, slug, body)
	require.Equal(t, 204, resp.StatusCode)

	if respBody, err := io.ReadAll(resp.Body); err != nil {
		t.Fatalf("Read response body failed: %s", err.Error())
	} else {
		require.Empty(t, respBody)
	}

	var secretAfterUpdate models.Secret
	helpers.QueryTestSecretBySlug(t, db, &secretAfterUpdate, slug)
	require.Equal(t, updatedSecretLabel, secretAfterUpdate.Label)
	require.Equal(t, updatedSecretString, secretAfterUpdate.String)
	require.Equal(t, secretBeforeUpdate.CreatedAt, secretAfterUpdate.CreatedAt)
	require.True(t, secretAfterUpdate.UpdatedAt.After(secretAfterUpdate.CreatedAt))
	require.True(t, secretAfterUpdate.UpdatedAt.After(secretBeforeUpdate.UpdatedAt))
}

func newRequestUpdateSecret(
	t *testing.T, app *fiber.App, conf *config.AppConfig, slug, body string,
) *http.Response {

	reqBody := strings.NewReader(body)
	req := httptest.NewRequest("PATCH", "/api/secrets/" + slug, reqBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Client-Operation", utils.UpdateSecret)
	req.Header.Set("Authorization", "Token " + conf.VAULTS_ACCESS_TOKEN)

	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
