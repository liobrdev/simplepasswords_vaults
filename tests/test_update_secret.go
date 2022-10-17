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

func testUpdateSecret(t *testing.T, app *fiber.App, db *gorm.DB) {
	t.Run("empty_body_400_bad_request", func(t *testing.T) {
		testUpdateSecretClientError(
			t, app, db, helpers.NewSlug(t), "", http.StatusBadRequest, utils.ErrorParse,
			"invalid character '\x00' looking for beginning of value",
		)
	})

	t.Run("array_body_400_bad_request", func(t *testing.T) {
		testUpdateSecretClientError(
			t, app, db, helpers.NewSlug(t), "[]", http.StatusBadRequest, utils.ErrorParse,
			"invalid character '[' looking for beginning of value",
		)

		testUpdateSecretClientError(
			t, app, db, helpers.NewSlug(t), "[{}]", http.StatusBadRequest, utils.ErrorParse,
			"invalid character '[' looking for beginning of value",
		)

		testUpdateSecretClientError(
			t, app, db, helpers.NewSlug(t), `[{"secret_label":"abc","secret_string":"123"}]`,
			http.StatusBadRequest, utils.ErrorParse,
			"invalid character '[' looking for beginning of value",
		)
	})

	t.Run("boolean_body_400_bad_request", func(t *testing.T) {
		testUpdateSecretClientError(
			t, app, db, helpers.NewSlug(t), "true", http.StatusBadRequest, utils.ErrorParse,
			"invalid character 't' looking for beginning of value",
		)

		testUpdateSecretClientError(
			t, app, db, helpers.NewSlug(t), "false", http.StatusBadRequest, utils.ErrorParse,
			"invalid character 'f' looking for beginning of value",
		)
	})

	t.Run("string_body_400_bad_request", func(t *testing.T) {
		testUpdateSecretClientError(
			t, app, db, helpers.NewSlug(t), "\"Valid JSON, but not an object.\"",
			http.StatusBadRequest, utils.ErrorParse,
			"invalid character '\"' looking for beginning of value",
		)
	})

	t.Run("null_or_empty_object_400_bad_request", func(t *testing.T) {
		_, _, _, secrets := setup.SetUpWithData(t, db)

		testUpdateSecretClientError(
			t, app, db, (*secrets)[0].Slug, "null",
			http.StatusBadRequest, utils.ErrorEmptyUpdateSecret,
			"Likely (null|empty) (object|fields).",
		)

		testUpdateSecretClientError(
			t, app, db, (*secrets)[0].Slug, "{}",
			http.StatusBadRequest, utils.ErrorEmptyUpdateSecret,
			"Likely (null|empty) (object|fields).",
		)
	})

	t.Run("both_fields_null_or_empty_or_missing_400_bad_request", func(t *testing.T) {
		_, _, _, secrets := setup.SetUpWithData(t, db)

		testUpdateSecretClientError(
			t, app, db, (*secrets)[0].Slug, `{"secret_label":"","secret_string":""}`,
			http.StatusBadRequest, utils.ErrorEmptyUpdateSecret,
			"Likely (null|empty) (object|fields).",
		)

		testUpdateSecretClientError(
			t, app, db, (*secrets)[0].Slug, `{"secret_label":null,"secret_string":null}`,
			http.StatusBadRequest, utils.ErrorEmptyUpdateSecret,
			"Likely (null|empty) (object|fields).",
		)

		testUpdateSecretClientError(
			t, app, db, (*secrets)[0].Slug, `{"weird_label":"abc","weird_string":"123"}`,
			http.StatusBadRequest, utils.ErrorEmptyUpdateSecret,
			"Likely (null|empty) (object|fields).",
		)
	})

	t.Run("too_long_secret_label_400_bad_request", func(t *testing.T) {
		// `label` is a random string greater than 255 characters in length
		if label, err := utils.GenerateSlug(256); err != nil {
			t.Fatalf("Generate long string failed: %s", err.Error())
		} else {
			testUpdateSecretClientError(
				t, app, db, helpers.NewSlug(t), fmt.Sprintf(`{"secret_label":"%s"}`, label),
				http.StatusBadRequest, utils.ErrorSecretLabel, "Too long (256 > 255)",
			)
		}
	})

	t.Run("too_long_secret_string_400_bad_request", func(t *testing.T) {
		// `str` is a random string greater than 1000 characters in length
		if str, err := utils.GenerateSlug(1001); err != nil {
			t.Fatalf("Generate long string failed: %s", err.Error())
		} else {
			testUpdateSecretClientError(
				t, app, db, helpers.NewSlug(t), fmt.Sprintf(`{"secret_string":"%s"}`, str),
				http.StatusBadRequest, utils.ErrorSecretString, "Too long (1001 > 1000)",
			)
		}
	})

	t.Run("valid_body_secret_label_already_exists_409_conflict", func(t *testing.T) {
		_, _, _, secrets := setup.SetUpWithData(t, db)

		testUpdateSecretClientError(
			t, app, db, (*secrets)[0].Slug,
			fmt.Sprintf(`{"secret_label":"%s"}`, (*secrets)[1].Label),
			http.StatusConflict, utils.ErrorFailedDB,
			"UNIQUE constraint failed: secrets.label, secrets.entry_slug",
		)
	})

	t.Run("valid_body_404_not_found", func(t *testing.T) {
		setup.SetUpWithData(t, db)
		testUpdateSecretClientError(
			t, app, db, helpers.NewSlug(t), `{"secret_label":"updated_secret"}`,
			http.StatusNotFound, utils.ErrorNoRowsAffected, "Likely that slug was not found.",
		)
	})

	t.Run("valid_body_204_no_content", func(t *testing.T) {
		_, _, _, secrets := setup.SetUpWithData(t, db)
		slug := (*secrets)[0].Slug

		updatedSecretLabel := "secret[_label='updated_label']@0.0.0.0"

		testUpdateSecretSuccess(
			t, app, db, slug, updatedSecretLabel, "",
			fmt.Sprintf(`{"secret_label":"%s"}`, updatedSecretLabel),
		)

		updatedSecretString := "secret[_string='updated_string']@0.0.0.0"

		testUpdateSecretSuccess(
			t, app, db, slug, "", updatedSecretString,
			fmt.Sprintf(`{"secret_string":"%s"}`, updatedSecretString),
		)

		updatedSecretLabel = "secret[_label='updated_again_label']@0.0.0.0"
		updatedSecretString = "secret[_string='updated_again_string']@0.0.0.0"

		testUpdateSecretSuccess(
			t, app, db, slug, updatedSecretLabel, updatedSecretString, fmt.Sprintf(
				`{"secret_label":"%s","secret_string":"%s"}`,
				updatedSecretLabel, updatedSecretString,
			),
		)
	})

	t.Run("valid_body_irrelevant_data_204_no_content", func(t *testing.T) {
		_, _, _, secrets := setup.SetUpWithData(t, db)
		slug := (*secrets)[0].Slug

		updatedSecretLabel := "secret[_label='updated_label']@0.0.0.0"
		updatedSecretString := "secret[_string='updated_string']@0.0.0.0"

		testUpdateSecretSuccess(
			t, app, db, slug, updatedSecretLabel, updatedSecretString, fmt.Sprintf(
				`{"secret_label":"%s","secret_string":"%s","is_real_field":false,"a":1}`,
				updatedSecretLabel, updatedSecretString,
			),
		)
	})
}

func testUpdateSecretClientError(
	t *testing.T,
	app *fiber.App,
	db *gorm.DB,
	slug string,
	body string,
	expectedStatus int,
	expectedMessage utils.ErrorMessage,
	expectedDetail string,
) {
	resp := newRequestUpdateSecret(t, app, slug, body)
	require.Equal(t, expectedStatus, resp.StatusCode)
	helpers.AssertErrorResponseBody(t, resp, utils.ErrorResponseBody{
		ClientOperation: utils.UpdateSecret,
		Message:         string(expectedMessage),
		Detail:          expectedDetail,
		RequestBody:     body,
	})
}

func testUpdateSecretSuccess(
	t *testing.T,
	app *fiber.App,
	db *gorm.DB,
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

	resp := newRequestUpdateSecret(t, app, slug, body)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)

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
	t *testing.T,
	app *fiber.App,
	slug string,
	body string,
) *http.Response {
	reqBody := strings.NewReader(body)
	req := httptest.NewRequest(http.MethodPatch, "/api/secrets/"+slug, reqBody)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
