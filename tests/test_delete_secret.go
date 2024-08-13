package tests

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_vaults/config"
	"github.com/liobrdev/simplepasswords_vaults/tests/helpers"
	"github.com/liobrdev/simplepasswords_vaults/tests/setup"
	"github.com/liobrdev/simplepasswords_vaults/utils"
)

func testDeleteSecret(t *testing.T, app *fiber.App, db *gorm.DB, conf *config.AppConfig) {
	t.Run("valid_slug_404_not_found", func(t *testing.T) {
		testDeleteSecretClientError(
			t, app, conf, 404, utils.ErrorNoRowsAffected, "Likely that slug was not found.",
			helpers.NewSlug(t),
		)
	})

	t.Run("invalid_slug_400_bad_request", func(t *testing.T) {
		slug := "notEvenARealSlug"
		testDeleteSecretClientError(t, app, conf, 400, utils.ErrorSecretSlug, slug, slug)
	})

	t.Run("valid_slug_204_no_content", func(t *testing.T) {
		testDeleteSecretSuccess(t, app, db, conf)
	})
}

func testDeleteSecretClientError(
	t *testing.T, app *fiber.App, conf *config.AppConfig, expectedStatus int,
	expectedMessage, expectedDetail, slug string,
) {
	resp := newRequestDeleteSecret(t, app, conf, slug)
	require.Equal(t, expectedStatus, resp.StatusCode)
	helpers.AssertErrorResponseBody(t, resp, utils.ErrorResponseBody{
		ClientOperation: utils.DeleteSecret,
		Message:         expectedMessage,
		Detail:          expectedDetail,
	})
}

func testDeleteSecretSuccess(t *testing.T, app *fiber.App, db *gorm.DB, conf *config.AppConfig) {
	_, _, _, secrets := setup.SetUpWithData(t, db)
	secret := secrets[0]

	var secretCount int64
	helpers.CountSecrets(t, db, &secretCount)
	require.EqualValues(t, 16, secretCount)

	resp := newRequestDeleteSecret(t, app, conf, secret.Slug)
	require.Equal(t, 204, resp.StatusCode)

	if respBody, err := io.ReadAll(resp.Body); err != nil {
		t.Fatalf("Read response body failed: %s", err.Error())
	} else {
		require.Empty(t, respBody)
	}

	if result := db.First(&secret, "slug = ?", secret.Slug); result.Error != nil {
		require.ErrorIs(t, result.Error, gorm.ErrRecordNotFound)
	} else {
		t.Fatalf("Deleted secret query failed: %s", result.Error.Error())
	}

	helpers.CountSecrets(t, db, &secretCount)
	require.EqualValues(t, 15, secretCount)
}

func newRequestDeleteSecret(
	t *testing.T, app *fiber.App, conf *config.AppConfig, slug string,
) *http.Response {

	req := httptest.NewRequest("DELETE", "/api/secrets/" + slug, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Client-Operation", utils.DeleteSecret)
	req.Header.Set("Authorization", "Token " + conf.VAULTS_ACCESS_TOKEN)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
