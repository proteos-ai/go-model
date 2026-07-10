package metamodel

import "golang.org/x/text/currency"

// IsValidCurrencyCode reports whether code is a valid ISO-4217 alphabetic
// currency code in canonical uppercase form (e.g. "USD", "JPY", "BHD"). It is
// the single source of truth for currency-code validity across the platform —
// functions-codegen meta validation, metadata-service attribute-config validation,
// and data-service record-value validation all call it, so the
// x/text/currency dependency lives in exactly one place.
//
// currency.ParseISO is lenient about case, but the platform stores codes in
// canonical uppercase so stored values, allow-lists, and defaults all compare
// consistently — we therefore reject anything that isn't three uppercase
// letters before consulting ParseISO.
func IsValidCurrencyCode(code string) bool {
	if len(code) != 3 {
		return false
	}
	for _, r := range code {
		if r < 'A' || r > 'Z' {
			return false
		}
	}
	_, err := currency.ParseISO(code)
	return err == nil
}
