package metamodel

// CurrencyAttributeMeta holds the metadata for an attribute of type `currency`.
// A currency attribute stores a composite value `{ amount, currency_code }`,
// where `amount` is an exact decimal STRING (never float64, to preserve
// financial-system-level precision) and `currency_code` is an ISO-4217 code.
//
// Both config fields are optional. When `AllowedCurrencyCodes` is empty, any
// valid ISO-4217 code is accepted. `DefaultCurrencyCode`, when set, seeds the
// currency picker for new values and must be a member of `AllowedCurrencyCodes`
// when that allow-list is also set.
type CurrencyAttributeMeta struct {
	DefaultCurrencyCode  *string  `json:"default_currency_code,omitempty"`
	AllowedCurrencyCodes []string `json:"allowed_currency_codes,omitempty"`
}
