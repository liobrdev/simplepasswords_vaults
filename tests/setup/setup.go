package setup

import (
	"testing"

	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_vaults/models"
)

func SetUp(t *testing.T, db *gorm.DB) {
	TearDown(t, db)

	if err := db.AutoMigrate(
		&models.User{},
		&models.Vault{},
		&models.Entry{},
		&models.Secret{},
	); err != nil {
		t.Fatalf("Failed database auto-migrate: %s", err)
	}
}

func SetUpWithData(t *testing.T, db *gorm.DB) (
	*[]models.User,
	*[]models.Vault,
	*[]models.Entry,
	*[]models.Secret,
) {
	SetUp(t, db)
	return populateTestDB(t, db)
}

func TearDown(t *testing.T, db *gorm.DB) {
	if result := db.Exec("DROP TABLE IF EXISTS secrets"); result.Error != nil {
		t.Fatalf("Test database tearDown failed: %s", result.Error.Error())
	}

	if result := db.Exec("DROP TABLE IF EXISTS entries"); result.Error != nil {
		t.Fatalf("Test database tearDown failed: %s", result.Error.Error())
	}

	if result := db.Exec("DROP TABLE IF EXISTS vaults"); result.Error != nil {
		t.Fatalf("Test database tearDown failed: %s", result.Error.Error())
	}

	if result := db.Exec("DROP TABLE IF EXISTS users"); result.Error != nil {
		t.Fatalf("Test database tearDown failed: %s", result.Error.Error())
	}
}
