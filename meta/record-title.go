package metamodel

import "fmt"

// RecordTitleField is the virtual column that stands for a record's rendered
// title — the entity's title_template evaluated against the row. Lists show
// it, and the records query sorts and filters on it when the template is
// plain enough to compile to SQL (see titletemplate.ParseSimple). No attribute
// ever backs it, which is exactly why the name is reserved: an attribute
// called record_title would shadow the column in every path grammar that
// reads it (list columns, sort keys, filter fields).
const RecordTitleField = "record_title"

// IsReservedAttributeName reports whether a name is taken by a virtual
// column and therefore not available to user-defined attributes.
func IsReservedAttributeName(name string) bool {
	return name == RecordTitleField
}

// ValidateReservedAttributeNames rejects a top-level attribute that uses a
// reserved name. Only the top level matters: the virtual columns are
// addressed as top-level paths (`record_title`, `company_id.record_title`),
// so a leaf inside an object attribute cannot collide with them.
func ValidateReservedAttributeNames(attributes []Attribute) error {
	for _, attribute := range attributes {
		if IsReservedAttributeName(attribute.Name) {
			return fmt.Errorf("attribute name %q is reserved for the record title column", attribute.Name)
		}
	}
	return nil
}
