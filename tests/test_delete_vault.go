package tests

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_vaults/models"
	"github.com/liobrdev/simplepasswords_vaults/tests/helpers"
	"github.com/liobrdev/simplepasswords_vaults/tests/setup"
	"github.com/liobrdev/simplepasswords_vaults/utils"
)

func testDeleteVault(t *testing.T, app *fiber.App, db *gorm.DB) {
	t.Run("valid_slug_404_not_found", func(t *testing.T) {
		testDeleteVaultClientError(
			t, app, db, helpers.NewSlug(t), http.StatusNotFound, utils.ErrorNoRowsAffected,
			"Likely that slug was not found.",
		)
	})

	t.Run("invalid_slug_400_bad_request", func(t *testing.T) {
		slug := "notEvenARealSlug"
		testDeleteVaultClientError(
			t, app, db, slug, http.StatusBadRequest, utils.ErrorVaultSlug, slug,
		)
	})

	t.Run("valid_slug_204_no_content", func(t *testing.T) {
		testDeleteVaultSuccess(t, app, db)
	})
}

func testDeleteVaultClientError(
	t *testing.T,
	app *fiber.App,
	db *gorm.DB,
	slug string,
	expectedStatus int,
	expectedMessage utils.ErrorMessage,
	expectedDetail string,
) {
	resp := newRequestDeleteVault(t, app, slug)
	require.Equal(t, expectedStatus, resp.StatusCode)
	helpers.AssertErrorResponseBody(t, resp, utils.ErrorResponseBody{
		ClientOperation: utils.DeleteVault,
		Message:         string(expectedMessage),
		Detail:          expectedDetail,
	})
}

func testDeleteVaultSuccess(t *testing.T, app *fiber.App, db *gorm.DB) {
	setup.SetUpWithData(t, db)

	var vault models.Vault
	helpers.QueryTestVaultEager(t, db, &vault, "vault@0.1.*.*")
	require.Len(t, vault.Entries, 2)

	entry1 := vault.Entries[0]
	require.Equal(t, "entry@0.1.0.*", entry1.Title)
	require.Len(t, entry1.Secrets, 2)

	secret1 := entry1.Secrets[0]
	secret2 := entry1.Secrets[1]
	require.Equal(t, "secret[_label='username']@0.1.0.0", secret1.Label)
	require.Equal(t, "secret[_string='foodeater1234']@0.1.0.0", secret1.String)
	require.Equal(t, "secret[_label='password']@0.1.0.1", secret2.Label)
	require.Equal(t, "secret[_string='3a7!ng40oD']@0.1.0.1", secret2.String)

	entry2 := vault.Entries[1]
	require.Equal(t, "entry@0.1.1.*", entry2.Title)
	require.Len(t, entry2.Secrets, 2)

	secret3 := entry2.Secrets[0]
	secret4 := entry2.Secrets[1]
	require.Equal(t, "secret[_label='username']@0.1.1.0", secret3.Label)
	require.Equal(t, "secret[_string='foodeater1234']@0.1.1.0", secret3.String)
	require.Equal(t, "secret[_label='password']@0.1.1.1", secret4.Label)
	require.Equal(t, "secret[_string='3a7!ng40oD']@0.1.1.1", secret4.String)

	var vaultCount int64
	helpers.CountVaults(t, db, &vaultCount)
	require.EqualValues(t, 4, vaultCount)

	var entryCount int64
	helpers.CountEntries(t, db, &entryCount)
	require.EqualValues(t, 8, entryCount)

	var secretCount int64
	helpers.CountSecrets(t, db, &secretCount)
	require.EqualValues(t, 16, secretCount)

	resp := newRequestDeleteVault(t, app, vault.Slug)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)

	if respBody, err := io.ReadAll(resp.Body); err != nil {
		t.Fatalf("Read response body failed: %s", err.Error())
	} else {
		require.Empty(t, respBody)
	}

	if result := db.First(&vault, "slug = ?", vault.Slug); result.Error != nil {
		require.ErrorIs(t, result.Error, gorm.ErrRecordNotFound)
	} else {
		t.Fatalf("Deleted vault query failed: %s", result.Error.Error())
	}

	if result := db.First(&entry1, "slug = ?", entry1.Slug); result.Error != nil {
		require.ErrorIs(t, result.Error, gorm.ErrRecordNotFound)
	} else {
		t.Fatalf("Deleted entry1 query failed: %s", result.Error.Error())
	}

	if result := db.First(&entry2, "slug = ?", entry2.Slug); result.Error != nil {
		require.ErrorIs(t, result.Error, gorm.ErrRecordNotFound)
	} else {
		t.Fatalf("Deleted entry2 query failed: %s", result.Error.Error())
	}

	if result := db.First(&secret1, "slug = ?", secret1.Slug); result.Error != nil {
		require.ErrorIs(t, result.Error, gorm.ErrRecordNotFound)
	} else {
		t.Fatalf("Deleted secret1 query failed: %s", result.Error.Error())
	}

	if result := db.First(&secret2, "slug = ?", secret2.Slug); result.Error != nil {
		require.ErrorIs(t, result.Error, gorm.ErrRecordNotFound)
	} else {
		t.Fatalf("Deleted secret2 query failed: %s", result.Error.Error())
	}

	if result := db.First(&secret3, "slug = ?", secret3.Slug); result.Error != nil {
		require.ErrorIs(t, result.Error, gorm.ErrRecordNotFound)
	} else {
		t.Fatalf("Deleted secret3 query failed: %s", result.Error.Error())
	}

	if result := db.First(&secret4, "slug = ?", secret4.Slug); result.Error != nil {
		require.ErrorIs(t, result.Error, gorm.ErrRecordNotFound)
	} else {
		t.Fatalf("Deleted secret4 query failed: %s", result.Error.Error())
	}

	helpers.CountVaults(t, db, &vaultCount)
	require.EqualValues(t, 3, vaultCount)

	helpers.CountEntries(t, db, &entryCount)
	require.EqualValues(t, 6, entryCount)

	helpers.CountSecrets(t, db, &secretCount)
	require.EqualValues(t, 12, secretCount)
}

func newRequestDeleteVault(
	t *testing.T,
	app *fiber.App,
	slug string,
) *http.Response {
	req := httptest.NewRequest(http.MethodDelete, "/api/vaults/"+slug, nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
