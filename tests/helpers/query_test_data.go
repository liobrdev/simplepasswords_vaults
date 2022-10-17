package helpers

import (
	"testing"

	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_vaults/models"
)

func QueryTestUser(t *testing.T, db *gorm.DB, user *models.User, slug string) {
	if result := db.First(&user, "slug = ?", slug); result.Error != nil {
		t.Fatalf("User query failed: %s", result.Error.Error())
	}
}

func QueryTestVault(t *testing.T, db *gorm.DB, vault *models.Vault, title string) {
	if result := db.First(&vault, "title = ?", title); result.Error != nil {
		t.Fatalf("Vault query failed: %s", result.Error.Error())
	}
}

func QueryTestVaultEager(t *testing.T, db *gorm.DB, vault *models.Vault, title string) {
	if result := db.Preload(
		"Entries.Secrets",
	).First(
		&vault, "title = ?", title,
	); result.Error != nil {
		t.Fatalf("Vault eager query failed: %s", result.Error.Error())
	}
}

func QueryTestEntry(t *testing.T, db *gorm.DB, entry *models.Entry, title string) {
	if result := db.First(&entry, "title = ?", title); result.Error != nil {
		t.Fatalf("Entry query failed: %s", result.Error.Error())
	}
}

func QueryTestEntryEager(t *testing.T, db *gorm.DB, entry *models.Entry, title string) {
	if result := db.Preload("Secrets").First(&entry, "title = ?", title); result.Error != nil {
		t.Fatalf("Entry eager query failed: %s", result.Error.Error())
	}
}

func QueryTestSecretByLabel(
	t *testing.T,
	db *gorm.DB,
	secret *models.Secret,
	secretLabel string,
) {
	if result := db.First(&secret, "label = ?", secretLabel); result.Error != nil {
		t.Fatalf("Secret query by label failed: %s", result.Error.Error())
	}
}

func QueryTestSecretBySlug(
	t *testing.T,
	db *gorm.DB,
	secret *models.Secret,
	secretSlug string,
) {
	if result := db.First(&secret, "slug = ?", secretSlug); result.Error != nil {
		t.Fatalf("Secret query by slug failed: %s", result.Error.Error())
	}
}

func QueryTestSecretsByEntry(
	t *testing.T,
	db *gorm.DB,
	secrets *[]models.Secret,
	entrySlug string,
) {
	if result := db.Where("entry_slug = ?", entrySlug).Find(&secrets); result.Error != nil {
		t.Fatalf("Secrets by entry query failed: %s", result.Error.Error())
	}
}
