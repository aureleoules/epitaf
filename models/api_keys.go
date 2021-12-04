package models

const (
	apiKeySchema = `
		CREATE TABLE api_keys (
			key VARCHAR(64) NOT NULL UNIQUE,

			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP(),

			PRIMARY KEY (key)
		);
	`
)

// Members string
type ApiKey []string

// Includes checks if s is included in slice
func (m ApiKey) Includes(s string) bool {
	for _, a := range m {
		if a == s {
			return true
		}
	}
	return false
}
