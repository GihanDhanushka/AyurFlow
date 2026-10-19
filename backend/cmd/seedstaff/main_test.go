package main

import (
	"strings"
	"testing"
)

func TestReadInputRequiresExplicitDevelopmentAndRoles(t *testing.T) {
	base := map[string]string{
		"APP_ENV":             "development",
		"SEED_STAFF_LOGIN":    " Demo.Admin ",
		"SEED_STAFF_PASSWORD": "synthetic passphrase",
		"SEED_STAFF_ROLES":    "administrator, clinician",
	}
	read := func(values map[string]string) (input, error) {
		return readInput(func(key string) string { return values[key] })
	}
	values, err := read(base)
	if err != nil || values.login != "demo.admin" || strings.Join(values.roles, ",") != "administrator,clinician" {
		t.Fatalf("readInput did not normalize the login and roles: %v", err)
	}
	for _, tc := range []struct {
		key, value string
	}{
		{"APP_ENV", "production"},
		{"APP_ENV", ""},
		{"SEED_STAFF_LOGIN", "bad login"},
		{"SEED_STAFF_PASSWORD", "short"},
		{"SEED_STAFF_ROLES", ""},
		{"SEED_STAFF_ROLES", "administrator,administrator"},
		{"SEED_STAFF_ROLES", "administrator,unknown"},
	} {
		changed := make(map[string]string, len(base))
		for key, value := range base {
			changed[key] = value
		}
		changed[tc.key] = tc.value
		if _, err := read(changed); err == nil {
			t.Errorf("readInput accepted invalid %s", tc.key)
		}
	}
}
