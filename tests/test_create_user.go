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

func testCreateUser(t *testing.T, app *fiber.App, db *gorm.DB, conf *config.AppConfig) {
	bodyFmt := `{"user_slug":"%s"}`

	t.Run("empty_body_400_bad_request", func(t *testing.T) {
		testCreateUserClientError(
			t, app, conf, 400, utils.ErrorParse,
			"invalid character '\x00' looking for beginning of value", "",
		)
	})

	t.Run("array_body_400_bad_request", func(t *testing.T) {
		testCreateUserClientError(
			t, app, conf, 400, utils.ErrorParse,
			"invalid character '[' looking for beginning of value", "[]",
		)

		testCreateUserClientError(
			t, app, conf, 400, utils.ErrorParse,
			"invalid character '[' looking for beginning of value", "[{}]",
		)

		testCreateUserClientError(
			t, app, conf, 400, utils.ErrorParse, "invalid character '[' looking for beginning of value",
			fmt.Sprintf(`[{"user_slug":"%s"}]`, helpers.NewSlug(t)),
		)
	})

	t.Run("boolean_body_400_bad_request", func(t *testing.T) {
		testCreateUserClientError(
			t, app, conf, 400, utils.ErrorParse, "invalid character 't' looking for beginning of value",
			"true",
		)

		testCreateUserClientError(
			t, app, conf, 400, utils.ErrorParse, "invalid character 'f' looking for beginning of value",
			"false",
		)
	})

	t.Run("string_body_400_bad_request", func(t *testing.T) {
		testCreateUserClientError(
			t, app, conf, 400, utils.ErrorParse, "invalid character '\"' looking for beginning of value",
			"\"Valid JSON, but not an object.\"",
		)
	})


	t.Run("null_body_400_bad_request", func(t *testing.T) {
		testCreateUserClientError(t, app, conf, 400, utils.ErrorUserSlug, "", "null")
	})

	t.Run("empty_object_body_400_bad_request", func(t *testing.T) {
		testCreateUserClientError(t, app, conf, 400, utils.ErrorUserSlug, "", "{}")
	})

	t.Run("missing_user_slug_400_bad_request", func(t *testing.T) {
		testCreateUserClientError(
			t, app, conf, 400, utils.ErrorUserSlug, "", `{"userr_slug":"Spelled wrong!"}`,
		)
	})

	t.Run("null_user_slug_400_bad_request", func(t *testing.T) {
		testCreateUserClientError(t, app, conf, 400, utils.ErrorUserSlug, "", `{"user_slug":null}`)
	})

	t.Run("empty_user_slug_400_bad_request", func(t *testing.T) {
		testCreateUserClientError(t, app, conf, 400, utils.ErrorUserSlug, "", `{"user_slug":""}`)
	})

	t.Run("too_long_user_slug_400_bad_request", func(t *testing.T) {
		slug := helpers.NewSlug(t) + "aA1!"
		testCreateUserClientError(
			t, app, conf, 400, utils.ErrorUserSlug, slug, fmt.Sprintf(bodyFmt, slug),
		)
	})

	t.Run("too_short_user_slug_400_bad_request", func(t *testing.T) {
		slug := helpers.NewSlug(t)[:15]
		testCreateUserClientError(
			t, app, conf, 400, utils.ErrorUserSlug, slug, fmt.Sprintf(bodyFmt, slug),
		)
	})

	t.Run("valid_body_user_slug_already_exists_409_conflict", func(t *testing.T) {
		users, _, _, _ := setup.SetUpWithData(t, db)
		slug := users[0].Slug

		testCreateUserClientError(
			t, app, conf, 409, utils.ErrorDuplicateUser, "UNIQUE constraint failed: users.slug",
			fmt.Sprintf(bodyFmt, slug),
		)
	})

	t.Run("valid_body_204_no_content", func(t *testing.T) {
		slug := helpers.NewSlug(t)
		testCreateUserSuccess(t, app, db, conf, slug, fmt.Sprintf(bodyFmt, slug))
	})

	t.Run("valid_body_irrelevant_data_204_no_content", func(t *testing.T) {
		slug := helpers.NewSlug(t)

		validBodyIrrelevantData := `{` +
			fmt.Sprintf(`"user_slug":"%s",`, slug) +
			`"user_email":"test@email.co",` +
			`"user_created_at":"10/12/22"` +
			`}`

		testCreateUserSuccess(t, app, db, conf, slug, validBodyIrrelevantData)
	})
}

func testCreateUserClientError(
	t *testing.T, app *fiber.App, conf *config.AppConfig, expectedStatus int,
	expectedMessage, expectedDetail, body string,
) {
	resp := newRequestCreateUser(t, app, conf, body)
	require.Equal(t, expectedStatus, resp.StatusCode)
	helpers.AssertErrorResponseBody(t, resp, utils.ErrorResponseBody{
		ClientOperation: utils.CreateUser,
		Message:         expectedMessage,
		Detail:          expectedDetail,
		RequestBody:     body,
	})
}

func testCreateUserSuccess(
	t *testing.T, app *fiber.App, db *gorm.DB, conf *config.AppConfig, slug, body string,
) {
	setup.SetUp(t, db)

	var userCount int64
	helpers.CountUsers(t, db, &userCount)
	require.EqualValues(t, 0, userCount)

	resp := newRequestCreateUser(t, app, conf, body)
	require.Equal(t, 204, resp.StatusCode)

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

	helpers.CountUsers(t, db, &userCount)
	require.EqualValues(t, 1, userCount)
}

func newRequestCreateUser(
	t *testing.T, app *fiber.App, conf *config.AppConfig, body string,
) *http.Response {

	reqBody := strings.NewReader(body)
	req := httptest.NewRequest(http.MethodPost, "/api/users", reqBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Client-Operation", utils.CreateUser)
	req.Header.Set("Authorization", "Token " + conf.VAULTS_ACCESS_TOKEN)

	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
