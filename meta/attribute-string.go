package metamodel

// StringFormat represents JSON Schema string formats
type StringFormat string

const (
	StringFormatEmail    StringFormat = "email"
	StringFormatURI      StringFormat = "uri"
	StringFormatUUID     StringFormat = "uuid"
	StringFormatHostname StringFormat = "hostname"
	StringFormatIPv4     StringFormat = "ipv4"
	StringFormatIPv6     StringFormat = "ipv6"
)

// StringAttributeMeta holds string-specific validation and formatting options
type StringAttributeMeta struct {
	MinLength *int         `json:"min_length,omitempty"`
	MaxLength *int         `json:"max_length,omitempty"`
	Pattern   string       `json:"pattern,omitempty"`
	Format    StringFormat `json:"format,omitempty"`
}
