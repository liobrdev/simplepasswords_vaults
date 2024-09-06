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
	"github.com/liobrdev/simplepasswords_vaults/models"
	"github.com/liobrdev/simplepasswords_vaults/tests/helpers"
	"github.com/liobrdev/simplepasswords_vaults/tests/setup"
	"github.com/liobrdev/simplepasswords_vaults/utils"
)

func testDeleteSecret(t *testing.T, app *fiber.App, db *gorm.DB, conf *config.AppConfig) {
	t.Run("valid_slug_404_not_found", func(t *testing.T) {
		slug := helpers.NewSlug(t)
		testDeleteSecretClientError(t, app, conf, 404, utils.ErrorNotFound, slug, slug)
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
	secret := secrets[16]

	var secretCount int64
	helpers.CountSecrets(t, db, &secretCount)
	require.EqualValues(t, 20, secretCount)

	var secretsBeforeDelete []models.Secret
	helpers.QueryTestSecretsByEntry(t, db, &secretsBeforeDelete, secret.EntrySlug)
	require.Equal(t, 7, len(secretsBeforeDelete))

	resp := newRequestDeleteSecret(t, app, conf, secret.Slug)
	require.Equal(t, 204, resp.StatusCode)

	if respBody, err := io.ReadAll(resp.Body); err != nil {
		t.Fatalf("Read response body failed: %s", err.Error())
	} else {
		require.Empty(t, respBody)
	}

	helpers.CountSecrets(t, db, &secretCount)
	require.EqualValues(t, 19, secretCount)

	var secretsAfterDelete []models.Secret
	helpers.QueryTestSecretsByEntry(t, db, &secretsAfterDelete, secret.EntrySlug)
	require.Equal(t, 6, len(secretsAfterDelete))

	for _, newSecret := range secretsAfterDelete {
		require.NotEqual(t, secret.Slug, newSecret.Slug)

		for _, oldSecret := range secretsBeforeDelete {
			if newSecret.Slug == oldSecret.Slug {
				if oldSecret.Priority > secret.Priority {
					require.Equal(t, oldSecret.Priority - 1, newSecret.Priority)
					require.True(t, newSecret.UpdatedAt.After(oldSecret.UpdatedAt))
				} else {
					require.Equal(t, oldSecret.Priority, newSecret.Priority)
					require.Equal(t, oldSecret.UpdatedAt, newSecret.UpdatedAt)
				}
			}
		}
	}

	result := db.First(&secret, "slug = ?", secret.Slug)
	require.NotNil(t, result.Error)
	require.ErrorIs(t, result.Error, gorm.ErrRecordNotFound)
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
