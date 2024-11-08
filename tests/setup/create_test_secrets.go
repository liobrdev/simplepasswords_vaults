package setup

import (
	"testing"

	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_vaults/models"
	"github.com/liobrdev/simplepasswords_vaults/tests/helpers"
	"github.com/liobrdev/simplepasswords_vaults/utils"
)

func createTestSecrets(
	users *[]models.User, vaults *[]models.Vault, entries *[]models.Entry, t *testing.T, db *gorm.DB,
) (secrets []models.Secret) {
	encryptedStrings := encryptTestStrings(t)

	secrets = []models.Secret{
		{
			Slug:      helpers.NewSlug(t),
			Label:     "secret[_label='username']@0.0.0.0",
			String:    encryptedStrings[0],
			EntrySlug: (*entries)[0].Slug,
			VaultSlug: (*vaults)[0].Slug,
			UserSlug:  (*users)[0].Slug,
			Priority:	 0,
		},
		{
			Slug:      helpers.NewSlug(t),
			Label:     "secret[_label='password']@0.0.0.1",
			String:    encryptedStrings[1],
			EntrySlug: (*entries)[0].Slug,
			VaultSlug: (*vaults)[0].Slug,
			UserSlug:  (*users)[0].Slug,
			Priority:	 1,
		},
		{
			Slug:      helpers.NewSlug(t),
			Label:     "secret[_label='username']@0.0.1.0",
			String:    encryptedStrings[2],
			EntrySlug: (*entries)[1].Slug,
			VaultSlug: (*vaults)[0].Slug,
			UserSlug:  (*users)[0].Slug,
			Priority:	 0,
		},
		{
			Slug:      helpers.NewSlug(t),
			Label:     "secret[_label='password']@0.0.1.1",
			String:    encryptedStrings[3],
			EntrySlug: (*entries)[1].Slug,
			VaultSlug: (*vaults)[0].Slug,
			UserSlug:  (*users)[0].Slug,
			Priority:	 1,
		},
		{
			Slug:      helpers.NewSlug(t),
			Label:     "secret[_label='username']@0.1.0.0",
			String:    encryptedStrings[4],
			EntrySlug: (*entries)[2].Slug,
			VaultSlug: (*vaults)[1].Slug,
			UserSlug:  (*users)[0].Slug,
			Priority:	 0,
		},
		{
			Slug:      helpers.NewSlug(t),
			Label:     "secret[_label='password']@0.1.0.1",
			String:    encryptedStrings[5],
			EntrySlug: (*entries)[2].Slug,
			VaultSlug: (*vaults)[1].Slug,
			UserSlug:  (*users)[0].Slug,
			Priority:	 1,
		},
		{
			Slug:      helpers.NewSlug(t),
			Label:     "secret[_label='username']@0.1.1.0",
			String:    encryptedStrings[6],
			EntrySlug: (*entries)[3].Slug,
			VaultSlug: (*vaults)[1].Slug,
			UserSlug:  (*users)[0].Slug,
			Priority:	 0,
		},
		{
			Slug:      helpers.NewSlug(t),
			Label:     "secret[_label='password']@0.1.1.1",
			String:    encryptedStrings[7],
			EntrySlug: (*entries)[3].Slug,
			VaultSlug: (*vaults)[1].Slug,
			UserSlug:  (*users)[0].Slug,
			Priority:	 1,
		},
		{
			Slug:      helpers.NewSlug(t),
			Label:     "secret[_label='username']@1.0.0.0",
			String:    encryptedStrings[8],
			EntrySlug: (*entries)[4].Slug,
			VaultSlug: (*vaults)[2].Slug,
			UserSlug:  (*users)[1].Slug,
			Priority:	 0,
		},
		{
			Slug:      helpers.NewSlug(t),
			Label:     "secret[_label='password']@1.0.0.1",
			String:    encryptedStrings[9],
			EntrySlug: (*entries)[4].Slug,
			VaultSlug: (*vaults)[2].Slug,
			UserSlug:  (*users)[1].Slug,
			Priority:	 1,
		},
		{
			Slug:      helpers.NewSlug(t),
			Label:     "secret[_label='username']@1.0.1.0",
			String:    encryptedStrings[10],
			EntrySlug: (*entries)[5].Slug,
			VaultSlug: (*vaults)[2].Slug,
			UserSlug:  (*users)[1].Slug,
			Priority:	 0,
		},
		{
			Slug:      helpers.NewSlug(t),
			Label:     "secret[_label='password']@1.0.1.1",
			String:    encryptedStrings[11],
			EntrySlug: (*entries)[5].Slug,
			VaultSlug: (*vaults)[2].Slug,
			UserSlug:  (*users)[1].Slug,
			Priority:	 1,
		},
		{
			Slug:      helpers.NewSlug(t),
			Label:     "secret[_label='username']@1.1.0.0",
			String:    encryptedStrings[12],
			EntrySlug: (*entries)[6].Slug,
			VaultSlug: (*vaults)[3].Slug,
			UserSlug:  (*users)[1].Slug,
			Priority:	 0,
		},
		{
			Slug:      helpers.NewSlug(t),
			Label:     "secret[_label='username']@1.1.1.0",
			String:    encryptedStrings[13],
			EntrySlug: (*entries)[7].Slug,
			VaultSlug: (*vaults)[3].Slug,
			UserSlug:  (*users)[1].Slug,
			Priority:	 0,
		},
		{
			Slug:      helpers.NewSlug(t),
			Label:     "secret[_label='password']@1.1.1.1",
			String:    encryptedStrings[14],
			EntrySlug: (*entries)[7].Slug,
			VaultSlug: (*vaults)[3].Slug,
			UserSlug:  (*users)[1].Slug,
			Priority:	 1,
		},
		{
			Slug:      helpers.NewSlug(t),
			Label:     "secret[_label='email']@1.1.1.2",
			String:    encryptedStrings[15],
			EntrySlug: (*entries)[7].Slug,
			VaultSlug: (*vaults)[3].Slug,
			UserSlug:  (*users)[1].Slug,
			Priority:	 2,
		},
		{
			Slug:      helpers.NewSlug(t),
			Label:     "secret[_label='foo']@1.1.1.3",
			String:    encryptedStrings[16],
			EntrySlug: (*entries)[7].Slug,
			VaultSlug: (*vaults)[3].Slug,
			UserSlug:  (*users)[1].Slug,
			Priority:	 3,
		},
		{
			Slug:      helpers.NewSlug(t),
			Label:     "secret[_label='bar']@1.1.1.4",
			String:    encryptedStrings[17],
			EntrySlug: (*entries)[7].Slug,
			VaultSlug: (*vaults)[3].Slug,
			UserSlug:  (*users)[1].Slug,
			Priority:	 4,
		},
		{
			Slug:      helpers.NewSlug(t),
			Label:     "secret[_label='this']@1.1.1.5",
			String:    encryptedStrings[18],
			EntrySlug: (*entries)[7].Slug,
			VaultSlug: (*vaults)[3].Slug,
			UserSlug:  (*users)[1].Slug,
			Priority:	 5,
		},
		{
			Slug:      helpers.NewSlug(t),
			Label:     "secret[_label='that']@1.1.1.6",
			String:    encryptedStrings[19],
			EntrySlug: (*entries)[7].Slug,
			VaultSlug: (*vaults)[3].Slug,
			UserSlug:  (*users)[1].Slug,
			Priority:	 6,
		},
	}

	if result := db.Create(&secrets); result.Error != nil {
		t.Fatalf("Create test secrets failed: %s", result.Error.Error())
	}

	return
}

func encryptTestStrings(t *testing.T) (encryptedStrings []string) {
	testStrings := []string{
		"secret[_string='foodeater1234']@0.0.0.0",
		"secret[_string='3a7!ng40oD']@0.0.0.1",
		"secret[_string='foodeater1234']@0.0.1.0",
		"secret[_string='3a7!ng40oD']@0.0.1.1",
		"secret[_string='foodeater1234']@0.1.0.0",
		"secret[_string='3a7!ng40oD']@0.1.0.1",
		"secret[_string='foodeater1234']@0.1.1.0",
		"secret[_string='3a7!ng40oD']@0.1.1.1",
		"secret[_string='foodeater1234']@1.0.0.0",
		"secret[_string='3a7!ng40oD']@1.0.0.1",
		"secret[_string='foodeater1234']@1.0.1.0",
		"secret[_string='3a7!ng40oD']@1.0.1.1",
		"secret[_string='foodeater1234']@1.1.0.0",
		"secret[_string='foodeater1234']@1.1.1.0",
		"secret[_string='3a7!ng40oD']@1.1.1.1",
		"secret[_string='foodeater@email.co']@1.1.1.2",
		"secret[_string='foodeater']@1.1.1.3",
		"secret[_string='bardeater']@1.1.1.4",
		"secret[_string='thiseater']@1.1.1.5",
		"secret[_string='thateater']@1.1.1.6",
	}

	for _, str := range testStrings {
		if encryptedString, err := utils.Encrypt(str, helpers.HexHash[:64]); err != nil {
			t.Fatal("Failed encryption: ", err.Error())
		} else {
			encryptedStrings = append(encryptedStrings, encryptedString)
		}
	}

	return
}
