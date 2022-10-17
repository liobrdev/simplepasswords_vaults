package helpers

import (
	"testing"

	"gorm.io/gorm"
)

func CountUsers(t *testing.T, db *gorm.DB, userCount *int64) {
	if result := db.Table("users").Count(userCount); result.Error != nil {
		t.Fatalf("User count failed: %s", result.Error.Error())
	}
}

func CountVaults(t *testing.T, db *gorm.DB, vaultCount *int64) {
	if result := db.Table("vaults").Count(vaultCount); result.Error != nil {
		t.Fatalf("Vault count failed: %s", result.Error.Error())
	}
}

func CountEntries(t *testing.T, db *gorm.DB, entryCount *int64) {
	if result := db.Table("entries").Count(entryCount); result.Error != nil {
		t.Fatalf("Entry count failed: %s", result.Error.Error())
	}
}

func CountSecrets(t *testing.T, db *gorm.DB, secretCount *int64) {
	if result := db.Table("secrets").Count(secretCount); result.Error != nil {
		t.Fatalf("Secret count failed: %s", result.Error.Error())
	}
}
