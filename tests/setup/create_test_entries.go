package setup

import (
	"testing"

	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_vaults/models"
	"github.com/liobrdev/simplepasswords_vaults/tests/helpers"
)

func createTestEntries(
	users *[]models.User, vaults *[]models.Vault, t *testing.T, db *gorm.DB,
) (entries []models.Entry) {
	entries = []models.Entry{
		{
			Slug:      helpers.NewSlug(t),
			Title:     "entry@0.0.0.*",
			VaultSlug: (*vaults)[0].Slug,
			UserSlug:  (*users)[0].Slug,
		},
		{
			Slug:      helpers.NewSlug(t),
			Title:     "entry@0.0.1.*",
			VaultSlug: (*vaults)[0].Slug,
			UserSlug:  (*users)[0].Slug,
		},
		{
			Slug:      helpers.NewSlug(t),
			Title:     "entry@0.1.0.*",
			VaultSlug: (*vaults)[1].Slug,
			UserSlug:  (*users)[0].Slug,
		},
		{
			Slug:      helpers.NewSlug(t),
			Title:     "entry@0.1.1.*",
			VaultSlug: (*vaults)[1].Slug,
			UserSlug:  (*users)[0].Slug,
		},
		{
			Slug:      helpers.NewSlug(t),
			Title:     "entry@1.0.0.*",
			VaultSlug: (*vaults)[2].Slug,
			UserSlug:  (*users)[1].Slug,
		},
		{
			Slug:      helpers.NewSlug(t),
			Title:     "entry@1.0.1.*",
			VaultSlug: (*vaults)[2].Slug,
			UserSlug:  (*users)[1].Slug,
		},
		{
			Slug:      helpers.NewSlug(t),
			Title:     "entry@1.1.0.*",
			VaultSlug: (*vaults)[3].Slug,
			UserSlug:  (*users)[1].Slug,
		},
		{
			Slug:      helpers.NewSlug(t),
			Title:     "entry@1.1.1.*",
			VaultSlug: (*vaults)[3].Slug,
			UserSlug:  (*users)[1].Slug,
		},
	}

	for _, entry := range entries {
		if result := db.Create(&entry); result.Error != nil {
			t.Fatalf("Create test entry failed: %s", result.Error.Error())
		}
	}

	if result := db.Find(&entries); result.Error != nil {
		t.Fatalf("Find test entries failed: %s", result.Error.Error())
	}

	return
}
