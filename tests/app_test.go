package tests

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_vaults/app"
	"github.com/liobrdev/simplepasswords_vaults/config"
)

func TestApp(t *testing.T) {
	var conf config.AppConfig

	if err := config.LoadConfigFromEnv(&conf); err != nil {
		t.Fatal("Failed to load config from environment:", err)
	}

	t.Run("is_behind_proxy", func(t *testing.T) {
		conf.GO_FIBER_BEHIND_PROXY = true
		app, dbs := app.CreateApp(&conf)
		runTests(t, app, dbs, &conf)
	})

	t.Run("is_not_behind_proxy", func(t *testing.T) {
		conf.GO_FIBER_BEHIND_PROXY = false
		app, db := app.CreateApp(&conf)
		runTests(t, app, db, &conf)
	})
}

func runTests(
	t *testing.T,
	app *fiber.App,
	db *gorm.DB,
	conf *config.AppConfig,
) {
	t.Run("test_create_user", func(t *testing.T) {
		testCreateUser(t, app, db)
	})

	t.Run("test_retrieve_user", func(t *testing.T) {
		testRetrieveUser(t, app, db)
	})

	t.Run("test_create_vault", func(t *testing.T) {
		testCreateVault(t, app, db)
	})

	t.Run("test_retrieve_vault", func(t *testing.T) {
		testRetrieveVault(t, app, db)
	})

	t.Run("test_update_vault", func(t *testing.T) {
		testUpdateVault(t, app, db)
	})

	t.Run("test_delete_vault", func(t *testing.T) {
		testDeleteVault(t, app, db)
	})

	t.Run("test_create_entry", func(t *testing.T) {
		testCreateEntry(t, app, db)
	})

	t.Run("test_retrieve_entry", func(t *testing.T) {
		testRetrieveEntry(t, app, db)
	})

	t.Run("test_update_entry", func(t *testing.T) {
		testUpdateEntry(t, app, db)
	})

	t.Run("test_delete_entry", func(t *testing.T) {
		testDeleteEntry(t, app, db)
	})

	t.Run("test_create_secret", func(t *testing.T) {
		testCreateSecret(t, app, db)
	})

	t.Run("test_update_secret", func(t *testing.T) {
		testUpdateSecret(t, app, db)
	})

	t.Run("test_delete_secret", func(t *testing.T) {
		testDeleteSecret(t, app, db)
	})
}
