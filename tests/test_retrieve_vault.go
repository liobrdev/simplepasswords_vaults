package tests

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_vaults/models"
	"github.com/liobrdev/simplepasswords_vaults/tests/helpers"
	"github.com/liobrdev/simplepasswords_vaults/tests/setup"
	"github.com/liobrdev/simplepasswords_vaults/utils"
)

func testRetrieveVault(t *testing.T, app *fiber.App, db *gorm.DB) {
	t.Run("invalid_slug_400_bad_request", func(t *testing.T) {
		slug := "notEvenARealSlug"
		testRetrieveVaultClientError(
			t, app, db, slug, http.StatusBadRequest, utils.ErrorVaultSlug, slug,
		)
	})

	t.Run("valid_slug_404_not_found", func(t *testing.T) {
		slug := helpers.NewSlug(t)
		testRetrieveVaultClientError(
			t, app, db, slug, http.StatusNotFound, utils.ErrorNotFound, slug,
		)
	})

	t.Run("valid_slug_200_ok", func(t *testing.T) {
		testRetrieveVaultSuccess(t, app, db)
	})
}

func testRetrieveVaultClientError(
	t *testing.T, app *fiber.App, db *gorm.DB, slug string, expectedStatus int,
	expectedMessage string, expectedDetail string,
) {
	resp := newRequestRetrieveVault(t, app, slug)
	require.Equal(t, expectedStatus, resp.StatusCode)
	helpers.AssertErrorResponseBody(t, resp, utils.ErrorResponseBody{
		ClientOperation: utils.RetrieveVault,
		Message:         expectedMessage,
		Detail:          expectedDetail,
	})
}

func testRetrieveVaultSuccess(t *testing.T, app *fiber.App, db *gorm.DB) {
	setup.SetUpWithData(t, db)

	var expectedVault models.Vault
	helpers.QueryTestVault(t, db, &expectedVault, "vault@0.1.*.*")

	resp := newRequestRetrieveVault(t, app, expectedVault.Slug)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	if respBody, err := io.ReadAll(resp.Body); err != nil {
		t.Fatalf("Read response body failed: %s", err.Error())
	} else {
		var actualVault models.Vault

		if err := json.Unmarshal(respBody, &actualVault); err != nil {
			t.Fatalf("JSON unmarshal failed: %s", err.Error())
		}

		require.Equal(t, expectedVault.Slug, actualVault.Slug)
		require.Equal(t, expectedVault.UserSlug, actualVault.UserSlug)
		require.Equal(t, expectedVault.Title, actualVault.Title)
		require.Equal(t, "vault@0.1.*.*", actualVault.Title)
		require.Len(t, actualVault.Entries, 2)
	}
}

func newRequestRetrieveVault(t *testing.T, app *fiber.App, slug string) *http.Response {
	req := httptest.NewRequest(http.MethodGet, "/api/vaults/"+slug, nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
