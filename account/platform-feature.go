package accountmodel

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

// Platform features are the keys of Profile.PlatformFeatures: the platform
// areas a profile shows, and the parts of Configuration it shows. The web
// client mirrors this list (clients/web/src/lib/platform-features.ts); keep
// the two in sync.
//
// Naming is <feature>_<part>. An area key (`config`) alone shows the whole
// area; adding part keys (`config_users`) narrows it to those parts. The
// Configuration part keys are the sidebar entry ids (`config-users`) in
// snake_case. An empty list shows every feature; `none`, alone, shows none
// (people with the profile see only their apps).
const (
	// PlatformFeatureNone hides every platform feature. It cannot be
	// combined with other keys.
	PlatformFeatureNone = "none"

	PlatformFeatureConfig    = "config"
	PlatformFeatureReports   = "reports"
	PlatformFeatureKnowledge = "knowledge"
	PlatformFeatureCalendar  = "calendar"
	// PlatformFeatureConversationsOwn shows the Conversations hub listing
	// Mine only; PlatformFeatureConversationsAll adds the Mine / All switch
	// and implies Own.
	PlatformFeatureConversationsOwn = "conversations_own"
	PlatformFeatureConversationsAll = "conversations_all"

	// PlatformFeatureConfigPartPrefix starts every Configuration part key.
	PlatformFeatureConfigPartPrefix = "config_"
)

// PlatformFeatures lists every known key, areas first, then the
// Configuration parts in sidebar order.
var PlatformFeatures = []string{
	PlatformFeatureNone,
	PlatformFeatureConfig,
	PlatformFeatureReports,
	PlatformFeatureKnowledge,
	PlatformFeatureCalendar,
	PlatformFeatureConversationsOwn,
	PlatformFeatureConversationsAll,
	"config_canvas",
	"config_entities",
	"config_pages",
	"config_lists",
	"config_platform_designer",
	"config_design_references",
	"config_agent_studio",
	"config_workflows",
	"config_conversations",
	"config_integrations",
	"config_scheduling",
	"config_users",
	"config_roles",
	"config_profiles",
	"config_teams",
	"config_usage",
	"config_events",
	"config_discovery",
}

// UnknownPlatformFeatureError names a key outside PlatformFeatures.
type UnknownPlatformFeatureError struct{ Feature string }

func (err UnknownPlatformFeatureError) Error() string {
	return fmt.Sprintf("unknown platform feature %q", err.Feature)
}

// ErrPlatformFeatureNoneCombined rejects `none` listed beside other keys,
// which would both hide and show features.
var ErrPlatformFeatureNoneCombined = errors.New(`platform feature "none" cannot be combined with other features`)

// CanonicalPlatformFeatures is the stored form of an allowlist: known keys
// only, deduplicated, in PlatformFeatures order, with the implications made
// explicit — conversations_all ⇒ conversations_own, any config_<part> ⇒
// config — so no reader has to know the rules. nil and empty both mean
// "every feature" and canonicalize to an empty slice; `none` (alone) means no
// feature and is kept as `[none]`. Pure; shared by
// account-service (which stores it) and the CLI (which diffs against it).
func CanonicalPlatformFeatures(features []string) ([]string, error) {
	present := make(map[string]bool, len(features))
	for _, feature := range features {
		if !slices.Contains(PlatformFeatures, feature) {
			return nil, UnknownPlatformFeatureError{Feature: feature}
		}
		present[feature] = true
		if feature == PlatformFeatureConversationsAll {
			present[PlatformFeatureConversationsOwn] = true
		}
		if strings.HasPrefix(feature, PlatformFeatureConfigPartPrefix) {
			present[PlatformFeatureConfig] = true
		}
	}
	if present[PlatformFeatureNone] && len(present) > 1 {
		return nil, ErrPlatformFeatureNoneCombined
	}
	canonical := make([]string, 0, len(present))
	for _, feature := range PlatformFeatures {
		if present[feature] {
			canonical = append(canonical, feature)
		}
	}
	return canonical, nil
}
