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

func testMoveSecret(t *testing.T, app *fiber.App, db *gorm.DB, conf *config.AppConfig) {
	dummySlug := helpers.NewSlug(t)

	t.Run("empty_body_400_bad_request", func(t *testing.T) {
		testMoveSecretClientError(
			t, app, conf, 400, utils.ErrorParse,
			"invalid character '\x00' looking for beginning of value", dummySlug, "",
		)
	})

	t.Run("array_body_400_bad_request", func(t *testing.T) {
		testMoveSecretClientError(
			t, app, conf, 400, utils.ErrorParse, "invalid character '[' looking for beginning of value",
			dummySlug, "[]",
		)

		testMoveSecretClientError(
			t, app, conf, 400, utils.ErrorParse, "invalid character '[' looking for beginning of value",
			dummySlug, "[{}]",
		)

		testMoveSecretClientError(
			t, app, conf, 400, utils.ErrorParse, "invalid character '[' looking for beginning of value",
			dummySlug, fmt.Sprintf(`[{"secret_priority":"0","entry_slug":"%s"}]`, dummySlug),
		)
	})

	t.Run("boolean_body_400_bad_request", func(t *testing.T) {
		testMoveSecretClientError(
			t, app, conf, 400, utils.ErrorParse, "invalid character 't' looking for beginning of value",
			dummySlug, "true",
		)

		testMoveSecretClientError(
			t, app, conf, 400, utils.ErrorParse, "invalid character 'f' looking for beginning of value",
			dummySlug, "false",
		)
	})

	t.Run("string_body_400_bad_request", func(t *testing.T) {
		testMoveSecretClientError(
			t, app, conf, 400, utils.ErrorParse, `invalid character '"' looking for beginning of value`,
			dummySlug, `"Valid JSON, but not an object."`,
		)
	})

	t.Run("null_body_400_bad_request", func(t *testing.T) {
		testMoveSecretClientError(t, app, conf, 400, utils.ErrorSecretPriority, "", dummySlug, "null")
	})

	t.Run("empty_object_body_400_bad_request", func(t *testing.T) {
		testMoveSecretClientError(t, app, conf, 400, utils.ErrorSecretPriority, "", dummySlug, "{}")
	})

	t.Run("missing_secret_priority_400_bad_request", func(t *testing.T) {
		testMoveSecretClientError(
			t, app, conf, 400, utils.ErrorSecretPriority, "", dummySlug,
			fmt.Sprintf(`{"secert_priority":"0","entry_slug":"%s"}`, dummySlug),
		)
	})

	t.Run("null_secret_priority_400_bad_request", func(t *testing.T) {
		testMoveSecretClientError(
			t, app, conf, 400, utils.ErrorSecretPriority, "", dummySlug,
			fmt.Sprintf(`{"secret_priority":null,"entry_slug":"%s"}`, dummySlug),
		)
	})

	t.Run("empty_secret_priority_400_bad_request", func(t *testing.T) {
		testMoveSecretClientError(
			t, app, conf, 400, utils.ErrorSecretPriority, "", dummySlug,
			fmt.Sprintf(`{"secret_priority":"","entry_slug":"%s"}`, dummySlug),
		)
	})

	t.Run("invalid_secret_priority_400_bad_request", func(t *testing.T) {
		testMoveSecretClientError(
			t, app, conf, 400, utils.ErrorSecretPriority, `strconv.Atoi: parsing "abc": invalid syntax`,
			dummySlug, fmt.Sprintf(`{"secret_priority":"abc","entry_slug":"%s"}`, dummySlug),
		)
	})

	t.Run("missing_entry_slug_400_bad_request", func(t *testing.T) {
		testMoveSecretClientError(
			t, app, conf, 400, utils.ErrorEntrySlug, "", dummySlug,
			fmt.Sprintf(`{"secret_priority":"0","enrty_slug":"%s"}`, dummySlug),
		)
	})

	t.Run("null_entry_slug_400_bad_request", func(t *testing.T) {
		testMoveSecretClientError(
			t, app, conf, 400, utils.ErrorEntrySlug, "", dummySlug,
			`{"secret_priority":"0","entry_slug":null}`,
		)
	})

	t.Run("empty_entry_slug_400_bad_request", func(t *testing.T) {
		testMoveSecretClientError(
			t, app, conf, 400, utils.ErrorEntrySlug, "", dummySlug,
			`{"secret_priority":"0","entry_slug":""}`,
		)
	})

	t.Run("invalid_entry_slug_400_bad_request", func(t *testing.T) {
		testMoveSecretClientError(
			t, app, conf, 400, utils.ErrorEntrySlug, dummySlug + "a", dummySlug,
			fmt.Sprintf(`{"secret_priority":"0","entry_slug":"%s"}`, dummySlug + "a"),
		)

		testMoveSecretClientError(
			t, app, conf, 400, utils.ErrorEntrySlug, dummySlug[:31], dummySlug,
			fmt.Sprintf(`{"secret_priority":"0","entry_slug":"%s"}`, dummySlug[:31]),
		)

		testMoveSecretClientError(
			t, app, conf, 400, utils.ErrorEntrySlug, dummySlug[:31] + "!", dummySlug,
			fmt.Sprintf(`{"secret_priority":"0","entry_slug":"%s"}`, dummySlug[:31] + "!"),
		)
	})

	t.Run("entry_slug_404_not_found", func(t *testing.T) {
		testMoveSecretClientError(
			t, app, conf, 404, utils.ErrorNotFound, "No secrets found", dummySlug,
			fmt.Sprintf(`{"secret_priority":"0","entry_slug":"%s"}`, dummySlug),
		)
	})

	t.Run("secret_slug_404_not_found", func(t *testing.T) {
		_, _, entries, _ := setup.SetUpWithData(t, db)

		testMoveSecretClientError(
			t, app, conf, 404, utils.ErrorSecretSlug, "Not found", dummySlug,
			fmt.Sprintf(`{"secret_priority":"0","entry_slug":"%s"}`, entries[7].Slug),
		)
	})

	t.Run("valid_body_one_secret_204_no_content", func(t *testing.T) {
		_, _, _, secrets := setup.SetUpWithData(t, db)

		testMoveSecretSuccess(
			t, app, db, conf, 0, 3, secrets[12].EntrySlug, secrets[12].Slug,
			fmt.Sprintf(`{"secret_priority":"3","entry_slug":"%s"}`, secrets[12].EntrySlug),
		)
	})

	t.Run("valid_body_lesser_new_priority_204_no_content", func(t *testing.T) {
		_, _, _, secrets := setup.SetUpWithData(t, db)

		testMoveSecretSuccess(
			t, app, db, conf, 4, 1, secrets[17].EntrySlug, secrets[17].Slug,
			fmt.Sprintf(`{"secret_priority":"1","entry_slug":"%s"}`, secrets[17].EntrySlug),
		)
	})

	t.Run("valid_body_much_lesser_new_priority_204_no_content", func(t *testing.T) {
		_, _, _, secrets := setup.SetUpWithData(t, db)

		testMoveSecretSuccess(
			t, app, db, conf, 4, 0, secrets[17].EntrySlug, secrets[17].Slug,
			fmt.Sprintf(`{"secret_priority":"-5","entry_slug":"%s"}`, secrets[17].EntrySlug),
		)
	})

	t.Run("valid_body_greater_new_priority_204_no_content", func(t *testing.T) {
		_, _, _, secrets := setup.SetUpWithData(t, db)

		testMoveSecretSuccess(
			t, app, db, conf, 1, 4, secrets[14].EntrySlug, secrets[14].Slug,
			fmt.Sprintf(`{"secret_priority":"4","entry_slug":"%s"}`, secrets[14].EntrySlug),
		)
	})

	t.Run("valid_body_much_greater_new_priority_204_no_content", func(t *testing.T) {
		_, _, _, secrets := setup.SetUpWithData(t, db)

		testMoveSecretSuccess(
			t, app, db, conf, 1, 6, secrets[14].EntrySlug, secrets[14].Slug,
			fmt.Sprintf(`{"secret_priority":"10","entry_slug":"%s"}`, secrets[14].EntrySlug),
		)
	})

	t.Run("valid_body_same_priority_204_no_content", func(t *testing.T) {
		_, _, _, secrets := setup.SetUpWithData(t, db)

		testMoveSecretSuccess(
			t, app, db, conf, 3, 3, secrets[16].EntrySlug, secrets[16].Slug,
			fmt.Sprintf(`{"secret_priority":"3","entry_slug":"%s"}`, secrets[16].EntrySlug),
		)
	})

	// t.Run("null_or_empty_object_400_bad_request", func(t *testing.T) {
	// 	_, _, _, secrets := setup.SetUpWithData(t, db)

	// 	testMoveSecretClientError(
	// 		t, app, conf, 400, utils.ErrorEmptyMoveSecret, "Null or empty object or fields.",
	// 		secrets[0].Slug, "null",
	// 	)

	// 	testMoveSecretClientError(
	// 		t, app, conf, 400, utils.ErrorEmptyMoveSecret, "Null or empty object or fields.",
	// 		secrets[0].Slug, "{}",
	// 	)
	// })

	// t.Run("both_fields_null_or_empty_or_missing_400_bad_request", func(t *testing.T) {
	// 	_, _, _, secrets := setup.SetUpWithData(t, db)

	// 	testMoveSecretClientError(
	// 		t, app, conf, 400, utils.ErrorEmptyMoveSecret, "Null or empty object or fields.",
	// 		secrets[0].Slug, `{"secret_label":"","secret_string":""}`,
	// 	)

	// 	testMoveSecretClientError(
	// 		t, app, conf, 400, utils.ErrorEmptyMoveSecret, "Null or empty object or fields.",
	// 		secrets[0].Slug, `{"secret_label":null,"secret_string":null}`,
	// 	)

	// 	testMoveSecretClientError(
	// 		t, app, conf, 400, utils.ErrorEmptyMoveSecret, "Null or empty object or fields.",
	// 		secrets[0].Slug, `{"weird_label":"abc","weird_string":"123"}`,
	// 	)
	// })

	// t.Run("too_long_secret_label_400_bad_request", func(t *testing.T) {
	// 	if label, err := utils.GenerateSlug(256); err != nil {
	// 		t.Fatalf("Generate long string failed: %s", err.Error())
	// 	} else {
	// 		testMoveSecretClientError(
	// 			t, app, conf, 400, utils.ErrorSecretLabel, "Too long", dummySlug,
	// 			fmt.Sprintf(`{"secret_label":"%s"}`, label),
	// 		)
	// 	}
	// })

	// t.Run("too_long_secret_string_400_bad_request", func(t *testing.T) {
	// 	if str, err := utils.GenerateSlug(1001); err != nil {
	// 		t.Fatalf("Generate long string failed: %s", err.Error())
	// 	} else {
	// 		testMoveSecretClientError(
	// 			t, app, conf, 400, utils.ErrorSecretString, "Too long", dummySlug,
	// 			fmt.Sprintf(`{"secret_string":"%s"}`, str),
	// 		)
	// 	}
	// })

	// t.Run("valid_body_secret_label_already_exists_500_error", func(t *testing.T) {
	// 	_, _, _, secrets := setup.SetUpWithData(t, db)

	// 	testMoveSecretClientError(
	// 		t, app, conf, 500, utils.ErrorFailedDB,
	// 		"UNIQUE constraint failed: secrets.label, secrets.entry_slug", secrets[0].Slug,
	// 		fmt.Sprintf(`{"secret_label":"%s"}`, secrets[1].Label),
	// 	)
	// })

	// t.Run("valid_body_404_not_found", func(t *testing.T) {
	// 	setup.SetUpWithData(t, db)
	// 	testMoveSecretClientError(
	// 		t, app, conf, 404, utils.ErrorNoRowsAffected, "Likely that slug was not found.",
	// 		dummySlug, `{"secret_label":"updated_secret"}`,
	// 	)
	// })

	// t.Run("valid_body_204_no_content", func(t *testing.T) {
	// 	_, _, _, secrets := setup.SetUpWithData(t, db)
	// 	slug := secrets[0].Slug

	// 	updatedSecretLabel := "secret[_label='updated_label']@0.0.0.0"

	// 	testMoveSecretSuccess(
	// 		t, app, db, conf, slug, updatedSecretLabel, "",
	// 		fmt.Sprintf(`{"secret_label":"%s"}`, updatedSecretLabel),
	// 	)

	// 	updatedSecretString := "secret[_string='updated_string']@0.0.0.0"

	// 	testMoveSecretSuccess(
	// 		t, app, db, conf, slug, "", updatedSecretString,
	// 		fmt.Sprintf(`{"secret_string":"%s"}`, updatedSecretString),
	// 	)

	// 	updatedSecretLabel = "secret[_label='updated_again_label']@0.0.0.0"
	// 	updatedSecretString = "secret[_string='updated_again_string']@0.0.0.0"

	// 	testMoveSecretSuccess(
	// 		t, app, db, conf, slug, updatedSecretLabel, updatedSecretString, fmt.Sprintf(
	// 			`{"secret_label":"%s","secret_string":"%s"}`, updatedSecretLabel, updatedSecretString,
	// 		),
	// 	)
	// })

	// t.Run("valid_body_irrelevant_data_204_no_content", func(t *testing.T) {
	// 	_, _, _, secrets := setup.SetUpWithData(t, db)
	// 	slug := secrets[0].Slug

	// 	updatedSecretLabel := "secret[_label='updated_label']@0.0.0.0"
	// 	updatedSecretString := "secret[_string='updated_string']@0.0.0.0"

	// 	testMoveSecretSuccess(
	// 		t, app, db, conf, slug, updatedSecretLabel, updatedSecretString, fmt.Sprintf(
	// 			`{"secret_label":"%s","secret_string":"%s","is_real_field":false,"a":1}`,
	// 			updatedSecretLabel, updatedSecretString,
	// 		),
	// 	)
	// })
}

func testMoveSecretClientError(
	t *testing.T, app *fiber.App, conf *config.AppConfig, expectedStatus int,
	expectedMessage, expectedDetail, slug, body string,
) {
	resp := newRequestMoveSecret(t, app, conf, slug, body)
	require.Equal(t, expectedStatus, resp.StatusCode)
	helpers.AssertErrorResponseBody(t, resp, utils.ErrorResponseBody{
		ClientOperation: utils.MoveSecret,
		Message:         expectedMessage,
		Detail:          expectedDetail,
		RequestBody:     body,
	})
}

func testMoveSecretSuccess(
	t *testing.T, app *fiber.App, db *gorm.DB, conf *config.AppConfig,
	oldPriority, newPriority uint8, entrySlug, secretSlug, body string,
) {
	var secretsBeforeMove []models.Secret
	helpers.QueryTestSecretsByEntry(t, db, &secretsBeforeMove, entrySlug)

	resp := newRequestMoveSecret(t, app, conf, secretSlug, body)
	require.Equal(t, 204, resp.StatusCode)

	if respBody, err := io.ReadAll(resp.Body); err != nil {
		t.Fatalf("Read response body failed: %s", err.Error())
	} else {
		require.Empty(t, respBody)
	}

	var secretsAfterMove []models.Secret
	helpers.QueryTestSecretsByEntry(t, db, &secretsAfterMove, entrySlug)

	if len(secretsBeforeMove) == 1 {
		require.EqualValues(t, 0, secretsAfterMove[0].Priority)
	} else {
		for _, newSecret := range secretsAfterMove {
			for _, oldSecret := range secretsBeforeMove {
				if newSecret.Slug == oldSecret.Slug {
					if newSecret.Slug == secretSlug {
						require.Equal(t, newPriority, newSecret.Priority)

						if newPriority != oldPriority {	
							require.True(t, newSecret.UpdatedAt.After(oldSecret.UpdatedAt))
						} else {
							require.Equal(t, oldSecret.UpdatedAt, newSecret.UpdatedAt)
						}
					} else if oldSecret.Priority < oldPriority && oldSecret.Priority >= newPriority {
						require.Equal(t, oldSecret.Priority + 1, newSecret.Priority)
						require.True(t, newSecret.UpdatedAt.After(oldSecret.UpdatedAt))
					} else if oldSecret.Priority > oldPriority && oldSecret.Priority <= newPriority {
						require.Equal(t, oldSecret.Priority - 1, newSecret.Priority)
						require.True(t, newSecret.UpdatedAt.After(oldSecret.UpdatedAt))
					} else {
						require.Equal(t, oldSecret.Priority, newSecret.Priority)
						require.Equal(t, oldSecret.UpdatedAt, newSecret.UpdatedAt)
					}
				}
			}
		}
	}
}

func newRequestMoveSecret(
	t *testing.T, app *fiber.App, conf *config.AppConfig, slug, body string,
) *http.Response {
	reqBody := strings.NewReader(body)
	req := httptest.NewRequest("PATCH", "/api/secrets/" + slug, reqBody)
	req.Header.Set("Authorization", "Token " + conf.VAULTS_ACCESS_TOKEN)
	req.Header.Set("Client-Operation", utils.MoveSecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
