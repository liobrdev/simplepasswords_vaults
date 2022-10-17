package utils

import "regexp"

var SlugRegexp = regexp.MustCompile(`^[\w-]{32}$`)
var RowsRegexp = regexp.MustCompile(`^result.RowsAffected \([0-9]+\) > 1$`)
var FailedSecretSlugRegexp = regexp.MustCompile("^Failed to generate `secret.Slug`:")
