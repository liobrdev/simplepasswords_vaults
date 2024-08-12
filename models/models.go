package models

import "time"

type User struct {
	Slug      string    `json:"user_slug" gorm:"primaryKey;not null"`
	CreatedAt time.Time `json:"-" gorm:"autoCreateTime:nano;not null"`
	UpdatedAt time.Time `json:"-" gorm:"autoUpdateTime:nano;not null"`
	Vaults    []Vault   `json:"vaults" gorm:"foreignKey:UserSlug;references:Slug;not null"`
	Entries   []Entry   `json:"-" gorm:"foreignKey:UserSlug;references:Slug;not null"`
	Secrets   []Secret  `json:"-" gorm:"foreignKey:UserSlug;references:Slug;not null"`
}

type Vault struct {
	Slug      string    `json:"vault_slug" gorm:"primaryKey;not null"`
	CreatedAt time.Time `json:"vault_created_at" gorm:"autoCreateTime:nano;not null"`
	UpdatedAt time.Time `json:"vault_updated_at" gorm:"autoUpdateTime:nano;not null"`
	Title     string    `json:"vault_title" gorm:"uniqueIndex:unique_title_user_slug;not null"`
	UserSlug  string    `json:"-" gorm:"uniqueIndex:unique_title_user_slug;index;not null"`
	User      User      `json:"-" gorm:"foreignKey:UserSlug;constraint:OnDelete:CASCADE"`
	Entries   []Entry   `json:"entries" gorm:"foreignKey:VaultSlug;references:Slug;not null"`
}

type Entry struct {
	Slug      string    `json:"entry_slug" gorm:"primaryKey;not null"`
	CreatedAt time.Time `json:"entry_created_at" gorm:"autoCreateTime:nano;not null"`
	UpdatedAt time.Time `json:"entry_updated_at" gorm:"autoUpdateTime:nano;not null"`
	Title     string    `json:"entry_title" gorm:"uniqueIndex:unique_title_vault_slug;not null"`
	VaultSlug string    `json:"-" gorm:"uniqueIndex:unique_title_vault_slug;index;not null"`
	Vault     Vault     `json:"-" gorm:"foreignKey:VaultSlug;constraint:OnDelete:CASCADE"`
	UserSlug  string    `json:"-" gorm:"not null"`
	Secrets   []Secret  `json:"secrets" gorm:"foreignKey:EntrySlug;references:Slug;not null"`
}

func (Entry) TableName() string {
	return "entries"
}

type Secret struct {
	Slug      string    `json:"secret_slug" gorm:"primaryKey;not null"`
	CreatedAt time.Time `json:"secret_created_at" gorm:"autoCreateTime:nano;not null"`
	UpdatedAt time.Time `json:"secret_updated_at" gorm:"autoUpdateTime:nano;not null"`
	Label     string    `json:"secret_label" gorm:"uniqueIndex:unique_label_entry_slug;not null"`
	String    string    `json:"secret_string" gorm:"not null"`
	EntrySlug string    `json:"-" gorm:"uniqueIndex:unique_label_entry_slug;index;not null"`
	Entry     Entry     `json:"-" gorm:"foreignKey:EntrySlug;constraint:OnDelete:CASCADE"`
	VaultSlug string    `json:"-" gorm:"not null"`
	UserSlug  string    `json:"-" gorm:"not null"`
}
