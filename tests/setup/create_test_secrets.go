package setup

import (
	"testing"

	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_vaults/models"
	"github.com/liobrdev/simplepasswords_vaults/tests/helpers"
)

func createTestSecrets(
	users *[]models.User, vaults *[]models.Vault, entries *[]models.Entry, t *testing.T, db *gorm.DB,
) (secrets []models.Secret) {
	secrets = []models.Secret{
		{
			Slug:      helpers.NewSlug(t),
			Label:     "secret[_label='username']@0.0.0.0",
			String:    "secret[_string='foodeater1234']@0.0.0.0",
			EntrySlug: (*entries)[0].Slug,
			VaultSlug: (*vaults)[0].Slug,
			UserSlug:  (*users)[0].Slug,
			Priority:	 0,
		},
		{
			Slug:      helpers.NewSlug(t),
			Label:     "secret[_label='password']@0.0.0.1",
			String:    "secret[_string='3a7!ng40oD']@0.0.0.1",
			EntrySlug: (*entries)[0].Slug,
			VaultSlug: (*vaults)[0].Slug,
			UserSlug:  (*users)[0].Slug,
			Priority:	 1,
		},
		{
			Slug:      helpers.NewSlug(t),
			Label:     "secret[_label='username']@0.0.1.0",
			String:    "secret[_string='foodeater1234']@0.0.1.0",
			EntrySlug: (*entries)[1].Slug,
			VaultSlug: (*vaults)[0].Slug,
			UserSlug:  (*users)[0].Slug,
			Priority:	 0,
		},
		{
			Slug:      helpers.NewSlug(t),
			Label:     "secret[_label='password']@0.0.1.1",
			String:    "secret[_string='3a7!ng40oD']@0.0.1.1",
			EntrySlug: (*entries)[1].Slug,
			VaultSlug: (*vaults)[0].Slug,
			UserSlug:  (*users)[0].Slug,
			Priority:	 1,
		},
		{
			Slug:      helpers.NewSlug(t),
			Label:     "secret[_label='username']@0.1.0.0",
			String:    "secret[_string='foodeater1234']@0.1.0.0",
			EntrySlug: (*entries)[2].Slug,
			VaultSlug: (*vaults)[1].Slug,
			UserSlug:  (*users)[0].Slug,
			Priority:	 0,
		},
		{
			Slug:      helpers.NewSlug(t),
			Label:     "secret[_label='password']@0.1.0.1",
			String:    "secret[_string='3a7!ng40oD']@0.1.0.1",
			EntrySlug: (*entries)[2].Slug,
			VaultSlug: (*vaults)[1].Slug,
			UserSlug:  (*users)[0].Slug,
			Priority:	 1,
		},
		{
			Slug:      helpers.NewSlug(t),
			Label:     "secret[_label='username']@0.1.1.0",
			String:    "secret[_string='foodeater1234']@0.1.1.0",
			EntrySlug: (*entries)[3].Slug,
			VaultSlug: (*vaults)[1].Slug,
			UserSlug:  (*users)[0].Slug,
			Priority:	 0,
		},
		{
			Slug:      helpers.NewSlug(t),
			Label:     "secret[_label='password']@0.1.1.1",
			String:    "secret[_string='3a7!ng40oD']@0.1.1.1",
			EntrySlug: (*entries)[3].Slug,
			VaultSlug: (*vaults)[1].Slug,
			UserSlug:  (*users)[0].Slug,
			Priority:	 1,
		},
		{
			Slug:      helpers.NewSlug(t),
			Label:     "secret[_label='username']@1.0.0.0",
			String:    "secret[_string='foodeater1234']@1.0.0.0",
			EntrySlug: (*entries)[4].Slug,
			VaultSlug: (*vaults)[2].Slug,
			UserSlug:  (*users)[1].Slug,
			Priority:	 0,
		},
		{
			Slug:      helpers.NewSlug(t),
			Label:     "secret[_label='password']@1.0.0.1",
			String:    "secret[_string='3a7!ng40oD']@1.0.0.1",
			EntrySlug: (*entries)[4].Slug,
			VaultSlug: (*vaults)[2].Slug,
			UserSlug:  (*users)[1].Slug,
			Priority:	 1,
		},
		{
			Slug:      helpers.NewSlug(t),
			Label:     "secret[_label='username']@1.0.1.0",
			String:    "secret[_string='foodeater1234']@1.0.1.0",
			EntrySlug: (*entries)[5].Slug,
			VaultSlug: (*vaults)[2].Slug,
			UserSlug:  (*users)[1].Slug,
			Priority:	 0,
		},
		{
			Slug:      helpers.NewSlug(t),
			Label:     "secret[_label='password']@1.0.1.1",
			String:    "secret[_string='3a7!ng40oD']@1.0.1.1",
			EntrySlug: (*entries)[5].Slug,
			VaultSlug: (*vaults)[2].Slug,
			UserSlug:  (*users)[1].Slug,
			Priority:	 1,
		},
		{
			Slug:      helpers.NewSlug(t),
			Label:     "secret[_label='username']@1.1.0.0",
			String:    "secret[_string='foodeater1234']@1.1.0.0",
			EntrySlug: (*entries)[6].Slug,
			VaultSlug: (*vaults)[3].Slug,
			UserSlug:  (*users)[1].Slug,
			Priority:	 0,
		},
		{
			Slug:      helpers.NewSlug(t),
			Label:     "secret[_label='password']@1.1.0.1",
			String:    "secret[_string='3a7!ng40oD']@1.1.0.1",
			EntrySlug: (*entries)[6].Slug,
			VaultSlug: (*vaults)[3].Slug,
			UserSlug:  (*users)[1].Slug,
			Priority:	 1,
		},
		{
			Slug:      helpers.NewSlug(t),
			Label:     "secret[_label='username']@1.1.1.0",
			String:    "secret[_string='foodeater1234']@1.1.1.0",
			EntrySlug: (*entries)[7].Slug,
			VaultSlug: (*vaults)[3].Slug,
			UserSlug:  (*users)[1].Slug,
			Priority:	 0,
		},
		{
			Slug:      helpers.NewSlug(t),
			Label:     "secret[_label='password']@1.1.1.1",
			String:    "secret[_string='3a7!ng40oD']@1.1.1.1",
			EntrySlug: (*entries)[7].Slug,
			VaultSlug: (*vaults)[3].Slug,
			UserSlug:  (*users)[1].Slug,
			Priority:	 1,
		},
	}

	for _, secret := range secrets {
		if result := db.Create(&secret); result.Error != nil {
			t.Fatalf("Create test secret failed: %s", result.Error.Error())
		}
	}

	if result := db.Find(&secrets); result.Error != nil {
		t.Fatalf("Find test secrets failed: %s", result.Error.Error())
	}

	return
}
