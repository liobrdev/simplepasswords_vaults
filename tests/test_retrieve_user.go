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

func testRetrieveUser(t *testing.T, app *fiber.App, db *gorm.DB) {
	t.Run("invalid_slug_400_bad_request", func(t *testing.T) {
		slug := "notEvenARealSlug"
		testRetrieveUserClientError(
			t, app, db, slug, http.StatusBadRequest, utils.ErrorUserSlug, slug,
		)
	})

	t.Run("valid_slug_404_not_found", func(t *testing.T) {
		slug := helpers.NewSlug(t)
		testRetrieveUserClientError(
			t, app, db, slug, http.StatusNotFound, utils.ErrorNotFound, slug,
		)
	})

	t.Run("valid_slug_200_ok", func(t *testing.T) {
		testRetrieveUserSuccess(t, app, db)
	})
}

func testRetrieveUserClientError(
	t *testing.T,
	app *fiber.App,
	db *gorm.DB,
	slug string,
	expectedStatus int,
	expectedMessage utils.ErrorMessage,
	expectedDetail string,
) {
	resp := newRequestRetrieveUser(t, app, slug)
	require.Equal(t, expectedStatus, resp.StatusCode)
	helpers.AssertErrorResponseBody(t, resp, utils.ErrorResponseBody{
		ClientOperation: utils.RetrieveUser,
		Message:         string(expectedMessage),
		Detail:          expectedDetail,
	})
}

func testRetrieveUserSuccess(t *testing.T, app *fiber.App, db *gorm.DB) {
	users, _, _, _ := setup.SetUpWithData(t, db)
	slug := (*users)[0].Slug
	resp := newRequestRetrieveUser(t, app, slug)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	if respBody, err := io.ReadAll(resp.Body); err != nil {
		t.Fatalf("Read response body failed: %s", err.Error())
	} else {
		var user models.User

		if err := json.Unmarshal(respBody, &user); err != nil {
			t.Fatalf("JSON unmarshal failed: %s", err.Error())
		}

		require.Equal(t, slug, user.Slug)
		require.Len(t, user.Vaults, 2)
	}
}

func newRequestRetrieveUser(t *testing.T, app *fiber.App, slug string) *http.Response {
	req := httptest.NewRequest(http.MethodGet, "/api/users/"+slug, nil)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
