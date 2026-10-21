package config

import (
	"testing"
	"time"
)

func TestLoadAuthPolicy(t *testing.T) {
	tests := []struct {
		name        string
		values      map[string]string
		wantEnabled bool
		wantError   bool
	}{
		{name: "disabled when unset", values: map[string]string{}},
		{name: "enabled with explicit values", values: map[string]string{
			"AUTH_ENABLED": "true", "AUTH_IDLE_TIMEOUT": "30m", "AUTH_ABSOLUTE_TIMEOUT": "8h",
			"AUTH_ATTEMPT_LIMIT": "5", "AUTH_ATTEMPT_WINDOW": "15m",
			"AUTH_BLOCK_DURATION": "15m", "AUTH_ALLOWED_ORIGIN": "https://clinic.example",
		}, wantEnabled: true},
		{name: "missing expiry", values: map[string]string{
			"AUTH_ENABLED": "true", "AUTH_ALLOWED_ORIGIN": "https://clinic.example",
		}, wantError: true},
		{name: "idle longer than absolute", values: map[string]string{
			"AUTH_ENABLED": "true", "AUTH_IDLE_TIMEOUT": "9h", "AUTH_ABSOLUTE_TIMEOUT": "8h",
			"AUTH_ATTEMPT_LIMIT": "5", "AUTH_ATTEMPT_WINDOW": "15m",
			"AUTH_BLOCK_DURATION": "15m", "AUTH_ALLOWED_ORIGIN": "https://clinic.example",
		}, wantError: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			policy, err := loadAuthPolicy(func(key string) string { return tc.values[key] })
			if (err != nil) != tc.wantError {
				t.Fatalf("error = %v, wantError %v", err, tc.wantError)
			}
			if policy.Enabled != tc.wantEnabled {
				t.Fatalf("enabled = %v, want %v", policy.Enabled, tc.wantEnabled)
			}
			if tc.wantEnabled && policy.IdleTimeout != 30*time.Minute {
				t.Fatalf("idle timeout = %v", policy.IdleTimeout)
			}
		})
	}
}
