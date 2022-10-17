package utils

type ErrorMessage string

const (
	ErrorParse             ErrorMessage = "Failed to parse request body."
	ErrorUserSlug          ErrorMessage = "Invalid `user_slug`."
	ErrorVaultSlug         ErrorMessage = "Invalid `vault_slug`."
	ErrorVaultTitle        ErrorMessage = "Invalid `vault_title`."
	ErrorEntrySlug         ErrorMessage = "Invalid `entry_slug`."
	ErrorEntryTitle        ErrorMessage = "Invalid `entry_title`."
	ErrorSecretSlug        ErrorMessage = "Invalid `secret_slug`."
	ErrorSecretLabel       ErrorMessage = "Invalid `secret_label`."
	ErrorSecretString      ErrorMessage = "Invalid `secret_string`."
	ErrorEmptyUpdateSecret ErrorMessage = "Empty 'update_secret' body."
	ErrorSecrets           ErrorMessage = "Invalid `secrets`."
	ErrorItemSecrets       ErrorMessage = "Invalid item in `secrets`."
	ErrorDuplicateSecrets  ErrorMessage = "Duplicate `entry.secrets.secret_label`."
	ErrorFailedDB          ErrorMessage = "Failed DB operation."
	ErrorNotFound          ErrorMessage = "Record not found."
	ErrorNoRowsAffected    ErrorMessage = "result.RowsAffected == 0"
)
