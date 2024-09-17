package utils

import "regexp"

var (
	SlugRegexp						 = regexp.MustCompile(`^[\w-]{16}$`)
	RowsRegexp						 = regexp.MustCompile(`^result.RowsAffected \([0-9]+\) > 1$`)
	FailedSecretSlugRegexp = regexp.MustCompile("^Failed to generate `secret.Slug`:")
	AuthHeaderRegexp			 = regexp.MustCompile(`^[Tt]oken [\w-]{80}$`)
	TokenNullRegexp				 = regexp.MustCompile(`^[Tt]oken (null)?$`)
)
