package metamodel

// Platform attribute names — the canonical, system-managed fields that every
// entity carries. They map 1:1 to the physical columns on each entity's record
// table in data-service. Keep this set in sync with the data-service table DDL.
const (
	PlatformAttributeID        = "id"
	PlatformAttributeCreatedAt = "created_at"
	PlatformAttributeUpdatedAt = "updated_at"
	PlatformAttributeCreatedBy = "created_by"
	PlatformAttributeUpdatedBy = "updated_by"
)

// retiredPlatformAttributeNames are names that USED to be platform attributes
// and are now stripped from every schema they still appear in.
//
// `owned_by` was the record-ownership platform attribute (a JSONB column on
// every record table). Ownership is now expressed on the business field
// itself — a principal attribute whose `grants` include `share` — so the
// platform attribute is gone. Stored entity schemas persisted it as JSONB;
// without this list EnsurePlatformAttributes would keep such an entry as an
// ordinary USER attribute (principal-typed, no column behind it, writes landing
// in `data`), which is worse than either having or not having it. The
// metadata-service PlatformAttributeBackfill re-runs EnsurePlatformAttributes
// over every stored entity, so retirement propagates the same way addition did.
var retiredPlatformAttributeNames = map[string]struct{}{
	"owned_by": {},
}

// IsRetiredPlatformAttributeName reports whether a name is a former platform
// attribute that must not survive in any schema.
func IsRetiredPlatformAttributeName(name string) bool {
	_, ok := retiredPlatformAttributeNames[name]
	return ok
}

// platformAttributeNames is the membership set behind IsPlatformAttributeName.
var platformAttributeNames = map[string]struct{}{
	PlatformAttributeID:        {},
	PlatformAttributeCreatedAt: {},
	PlatformAttributeUpdatedAt: {},
	PlatformAttributeCreatedBy: {},
	PlatformAttributeUpdatedBy: {},
}

// IsPlatformAttributeName reports whether a name belongs to the canonical
// platform attribute set.
func IsPlatformAttributeName(name string) bool {
	_, ok := platformAttributeNames[name]
	return ok
}

// PlatformAttributes returns the canonical platform attributes, in display order.
// Every entity must carry exactly these (see EnsurePlatformAttributes). They are
// flagged IsPlatformManaged + IsReadOnly so the UI locks them and the backend
// rejects client redefinition.
func PlatformAttributes() []Attribute {
	return []Attribute{
		{
			Name:              PlatformAttributeID,
			Type:              AttributeTypeString,
			Label:             "Id",
			Description:       "Unique identifier of the record.",
			IsRequired:        false,
			IsUnique:          true,
			IsReadOnly:        true,
			IsPlatformManaged: true,
		},
		{
			Name:              PlatformAttributeCreatedAt,
			Type:              AttributeTypeDatetime,
			Label:             "Created At",
			Description:       "Timestamp when the record was created.",
			IsRequired:        false,
			IsReadOnly:        true,
			IsPlatformManaged: true,
			Meta:              &DatetimeAttributeMeta{Format: DatetimeFormatDateTime},
		},
		{
			Name:              PlatformAttributeUpdatedAt,
			Type:              AttributeTypeDatetime,
			Label:             "Updated At",
			Description:       "Timestamp when the record was last updated.",
			IsRequired:        false,
			IsReadOnly:        true,
			IsPlatformManaged: true,
			Meta:              &DatetimeAttributeMeta{Format: DatetimeFormatDateTime},
		},
		{
			Name:              PlatformAttributeCreatedBy,
			Type:              AttributeTypeUser,
			Label:             "Created By",
			Description:       "User who created the record (the \"platform\" sentinel for system writes).",
			IsRequired:        false,
			IsReadOnly:        true,
			IsPlatformManaged: true,
		},
		{
			Name:              PlatformAttributeUpdatedBy,
			Type:              AttributeTypeUser,
			Label:             "Updated By",
			Description:       "User who last updated the record (the \"platform\" sentinel for system writes).",
			IsRequired:        false,
			IsReadOnly:        true,
			IsPlatformManaged: true,
		},
	}
}

// EnsurePlatformAttributes returns the attribute list with the canonical platform
// attributes guaranteed present, correct, and first. Any client-supplied
// attributes whose name collides with a platform attribute are dropped (the
// canonical definition is authoritative), so this doubles as enforcement: callers
// run it on create AND update, making platform attributes impossible to remove or
// redefine. Retired platform names are dropped outright. User attribute order is
// otherwise preserved.
func EnsurePlatformAttributes(attributes []Attribute) []Attribute {
	result := PlatformAttributes()
	for _, attr := range attributes {
		if IsPlatformAttributeName(attr.Name) || IsRetiredPlatformAttributeName(attr.Name) {
			continue
		}
		result = append(result, attr)
	}
	return result
}
