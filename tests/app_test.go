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

	conf.GO_TESTING_CONTEXT = t
	app, dbs := app.CreateApp(&conf)

	runTests(t, app, dbs, &conf)
}

func runTests(t *testing.T, app *fiber.App, db *gorm.DB, conf *config.AppConfig) {
	t.Run("test_authorize_request", func(t *testing.T) {
		testAuthorizeRequest(t, app, conf)
	})

	t.Run("test_create_user", func(t *testing.T) {
		testCreateUser(t, app, db, conf)
	})

	t.Run("test_retrieve_user", func(t *testing.T) {
		testRetrieveUser(t, app, db, conf)
	})

	t.Run("test_create_vault", func(t *testing.T) {
		testCreateVault(t, app, db, conf)
	})

	t.Run("test_retrieve_vault", func(t *testing.T) {
		testRetrieveVault(t, app, db, conf)
	})

	t.Run("test_update_vault", func(t *testing.T) {
		testUpdateVault(t, app, db, conf)
	})

	t.Run("test_delete_vault", func(t *testing.T) {
		testDeleteVault(t, app, db, conf)
	})

	t.Run("test_create_entry", func(t *testing.T) {
		testCreateEntry(t, app, db, conf)
	})

	// t.Run("test_retrieve_entry", func(t *testing.T) {
	// 	testRetrieveEntry(t, app, db)
	// })

	// t.Run("test_update_entry", func(t *testing.T) {
	// 	testUpdateEntry(t, app, db)
	// })

	// t.Run("test_delete_entry", func(t *testing.T) {
	// 	testDeleteEntry(t, app, db)
	// })

	// t.Run("test_create_secret", func(t *testing.T) {
	// 	testCreateSecret(t, app, db)
	// })

	// t.Run("test_update_secret", func(t *testing.T) {
	// 	testUpdateSecret(t, app, db)
	// })

	// t.Run("test_delete_secret", func(t *testing.T) {
	// 	testDeleteSecret(t, app, db)
	// })
}
