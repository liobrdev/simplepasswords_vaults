package setup

import (
	"testing"

	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_vaults/models"
	"github.com/liobrdev/simplepasswords_vaults/tests/helpers"
)

func createTestVaults(users *[]models.User, t *testing.T, db *gorm.DB) (vaults []models.Vault) {
	vaults = []models.Vault{
		{
			Slug:     helpers.NewSlug(t),
			Title:    "vault@0.0.*.*",
			UserSlug: (*users)[0].Slug,
		},
		{
			Slug:     helpers.NewSlug(t),
			Title:    "vault@0.1.*.*",
			UserSlug: (*users)[0].Slug,
		},
		{
			Slug:     helpers.NewSlug(t),
			Title:    "vault@1.0.*.*",
			UserSlug: (*users)[1].Slug,
		},
		{
			Slug:     helpers.NewSlug(t),
			Title:    "vault@1.1.*.*",
			UserSlug: (*users)[1].Slug,
		},
	}

	for _, vault := range vaults {
		if result := db.Create(&vault); result.Error != nil {
			t.Fatalf("Create test vault failed: %s", result.Error.Error())
		}
	}

	if result := db.Find(&vaults); result.Error != nil {
		t.Fatalf("Find test vaults failed: %s", result.Error.Error())
	}

	return
}
