package utils

const (
	ErrorParse             string = "Failed to parse request body."
	ErrorUserSlug          string = "Invalid `user_slug`."
	ErrorVaultSlug         string = "Invalid `vault_slug`."
	ErrorVaultTitle        string = "Invalid `vault_title`."
	ErrorEntrySlug         string = "Invalid `entry_slug`."
	ErrorEntryTitle        string = "Invalid `entry_title`."
	ErrorSecretSlug        string = "Invalid `secret_slug`."
	ErrorSecretLabel       string = "Invalid `secret_label`."
	ErrorSecretString      string = "Invalid `secret_string`."
	ErrorEmptyUpdateSecret string = "Empty 'update_secret' body."
	ErrorSecrets           string = "Invalid `secrets`."
	ErrorItemSecrets       string = "Invalid item in `secrets`."
	ErrorDuplicateSecrets  string = "Duplicate `entry.secrets.secret_label`."
	ErrorDuplicateUser		 string = "User already exists."
	ErrorFailedDB          string = "Failed DB operation."
	ErrorNotFound          string = "Record not found."
	ErrorNoRowsAffected    string = "result.RowsAffected == 0"
)
