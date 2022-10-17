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

func testCreateUser(t *testing.T, app *fiber.App, db *gorm.DB) {
	t.Run("empty_body_400_bad_request", func(t *testing.T) {
		testCreateUserClientError(
			t, app, db, "", http.StatusBadRequest, utils.ErrorParse,
			"invalid character '\x00' looking for beginning of value",
		)
	})

	t.Run("array_body_400_bad_request", func(t *testing.T) {
		testCreateUserClientError(
			t, app, db, "[]", http.StatusBadRequest, utils.ErrorParse,
			"invalid character '[' looking for beginning of value",
		)

		testCreateUserClientError(
			t, app, db, "[{}]", http.StatusBadRequest, utils.ErrorParse,
			"invalid character '[' looking for beginning of value",
		)

		testCreateUserClientError(
			t, app, db, fmt.Sprintf(`[{"user_slug":"%s"}]`, helpers.NewSlug(t)),
			http.StatusBadRequest, utils.ErrorParse,
			"invalid character '[' looking for beginning of value",
		)
	})

	t.Run("null_body_400_bad_request", func(t *testing.T) {
		testCreateUserClientError(
			t, app, db, "null", http.StatusBadRequest, utils.ErrorUserSlug, "",
		)
	})

	t.Run("boolean_body_400_bad_request", func(t *testing.T) {
		testCreateUserClientError(
			t, app, db, "true", http.StatusBadRequest, utils.ErrorParse,
			"invalid character 't' looking for beginning of value",
		)

		testCreateUserClientError(
			t, app, db, "false", http.StatusBadRequest, utils.ErrorParse,
			"invalid character 'f' looking for beginning of value",
		)
	})

	t.Run("string_body_400_bad_request", func(t *testing.T) {
		testCreateUserClientError(
			t, app, db, "\"Valid JSON, but not an object.\"", http.StatusBadRequest,
			utils.ErrorParse, "invalid character '\"' looking for beginning of value",
		)
	})

	t.Run("empty_object_body_400_bad_request", func(t *testing.T) {
		testCreateUserClientError(
			t, app, db, "{}", http.StatusBadRequest, utils.ErrorUserSlug, "",
		)
	})

	t.Run("missing_user_slug_400_bad_request", func(t *testing.T) {
		testCreateUserClientError(
			t, app, db, `{"userr_slug":"Spelled wrong!"}`, http.StatusBadRequest,
			utils.ErrorUserSlug, "",
		)
	})

	t.Run("null_user_slug_400_bad_request", func(t *testing.T) {
		testCreateUserClientError(
			t, app, db, `{"user_slug":null}`, http.StatusBadRequest, utils.ErrorUserSlug, "",
		)
	})

	t.Run("empty_user_slug_400_bad_request", func(t *testing.T) {
		testCreateUserClientError(
			t, app, db, `{"user_slug":""}`, http.StatusBadRequest, utils.ErrorUserSlug, "",
		)
	})

	t.Run("too_long_user_slug_400_bad_request", func(t *testing.T) {
		// `slug` is a random string greater than 32 characters in length
		if slug, err := utils.GenerateSlug(33); err != nil {
			t.Fatalf("Generate long string failed: %s", err.Error())
		} else {
			testCreateUserClientError(
				t, app, db, fmt.Sprintf(`{"user_slug":"%s"}`, slug), http.StatusBadRequest,
				utils.ErrorUserSlug, slug,
			)
		}
	})

	t.Run("valid_body_user_slug_already_exists_409_conflict", func(t *testing.T) {
		users, _, _, _ := setup.SetUpWithData(t, db)
		slug := (*users)[0].Slug

		testCreateUserClientError(
			t, app, db, fmt.Sprintf(`{"user_slug":"%s"}`, slug),
			http.StatusConflict, utils.ErrorFailedDB, "UNIQUE constraint failed: users.slug",
		)
	})

	t.Run("valid_body_204_no_content", func(t *testing.T) {
		slug := helpers.NewSlug(t)
		testCreateUserSuccess(t, app, db, slug, fmt.Sprintf(`{"user_slug":"%s"}`, slug))
	})

	t.Run("valid_body_irrelevant_data_204_no_content", func(t *testing.T) {
		slug := helpers.NewSlug(t)

		validBodyIrrelevantData := `{` +
			fmt.Sprintf(`"user_slug":"%s",`, slug) +
			`"user_email":"test@email.co",` +
			`"user_created_at":"10/12/22"` +
			`}`

		testCreateUserSuccess(t, app, db, slug, validBodyIrrelevantData)
	})
}

func testCreateUserClientError(
	t *testing.T,
	app *fiber.App,
	db *gorm.DB,
	body string,
	expectedStatus int,
	expectedMessage utils.ErrorMessage,
	expectedDetail string,
) {
	resp := newRequestCreateUser(t, app, body)
	require.Equal(t, expectedStatus, resp.StatusCode)
	helpers.AssertErrorResponseBody(t, resp, utils.ErrorResponseBody{
		ClientOperation: utils.CreateUser,
		Message:         string(expectedMessage),
		Detail:          expectedDetail,
		RequestBody:     body,
	})
}

func testCreateUserSuccess(
	t *testing.T,
	app *fiber.App,
	db *gorm.DB,
	slug string,
	body string,
) {
	setup.SetUp(t, db)
	resp := newRequestCreateUser(t, app, body)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)

	if respBody, err := io.ReadAll(resp.Body); err != nil {
		t.Fatalf("Read response body failed: %s", err.Error())
	} else {
		require.Empty(t, respBody)
	}

	var user models.User
	helpers.QueryTestUser(t, db, &user, slug)
	require.Equal(t, slug, user.Slug)
	require.Equal(t, user.CreatedAt, user.UpdatedAt)
	require.Empty(t, user.Vaults)
	require.Empty(t, user.Entries)
	require.Empty(t, user.Secrets)

	var userCount int64
	helpers.CountUsers(t, db, &userCount)
	require.EqualValues(t, 1, userCount)
}

func newRequestCreateUser(t *testing.T, app *fiber.App, body string) *http.Response {
	reqBody := strings.NewReader(body)
	req := httptest.NewRequest(http.MethodPost, "/api/users", reqBody)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
