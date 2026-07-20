package conversationmodel

import "testing"

func TestFilterConfigRoundTrip(t *testing.T) {
	cases := []struct {
		name string
		typ  ConversationFilterType
		cfg  ConversationFilterConfig
	}{
		{"address", FilterTypeAddress, AddressFilterConfig{Kind: ContactAddressKindEmail, Value: "spam@x.com", MatchOn: FilterMatchOnSender}},
		{"domain", FilterTypeDomain, DomainFilterConfig{Domain: "x.com", MatchOn: FilterMatchOnAnyParticipant}},
		{"role_based", FilterTypeRoleBased, RoleBasedFilterConfig{Prefixes: []string{"noreply", "info"}}},
		{"automated", FilterTypeAutomated, AutomatedFilterConfig{Signals: []AutomatedSignal{AutomatedSignalBounce, AutomatedSignalMailingList}}},
		{"internal", FilterTypeInternalConversations, InternalConversationsFilterConfig{Domains: []string{"proteos.ai", "onefabric.io"}}},
		{"all", FilterTypeAll, AllFilterConfig{}},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			raw, err := MarshalFilterConfig(testCase.cfg)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			decoded, err := DecodeFilterConfig(testCase.typ, raw)
			if err != nil {
				t.Fatalf("decode: %v", err)
			}
			if decoded.FilterType() != testCase.typ {
				t.Fatalf("type = %q, want %q", decoded.FilterType(), testCase.typ)
			}
			switch want := testCase.cfg.(type) {
			case AddressFilterConfig:
				if got, ok := decoded.(AddressFilterConfig); !ok || got != want {
					t.Fatalf("address round trip lost data: %#v", decoded)
				}
			case DomainFilterConfig:
				if got, ok := decoded.(DomainFilterConfig); !ok || got != want {
					t.Fatalf("domain round trip lost data: %#v", decoded)
				}
			case RoleBasedFilterConfig:
				got, ok := decoded.(RoleBasedFilterConfig)
				if !ok || len(got.Prefixes) != len(want.Prefixes) {
					t.Fatalf("role_based round trip lost data: %#v", decoded)
				}
			case AutomatedFilterConfig:
				got, ok := decoded.(AutomatedFilterConfig)
				if !ok || len(got.Signals) != len(want.Signals) {
					t.Fatalf("automated round trip lost data: %#v", decoded)
				}
			case InternalConversationsFilterConfig:
				got, ok := decoded.(InternalConversationsFilterConfig)
				if !ok || len(got.Domains) != len(want.Domains) {
					t.Fatalf("internal round trip lost data: %#v", decoded)
				}
			}
		})
	}
}

func TestFilterConfigEmptyNilAndUnknown(t *testing.T) {
	// nil (all with no params) marshals to '{}'.
	raw, err := MarshalFilterConfig(nil)
	if err != nil || string(raw) != "{}" {
		t.Fatalf("nil config must marshal to '{}': %q %v", raw, err)
	}
	// Empty/null raw tolerated for parameterized types (decodes zero variant).
	for _, raw := range [][]byte{nil, []byte("null"), []byte("{}")} {
		decoded, err := DecodeFilterConfig(FilterTypeDomain, raw)
		if err != nil {
			t.Fatalf("empty raw must decode: %v", err)
		}
		if _, ok := decoded.(DomainFilterConfig); !ok {
			t.Fatalf("empty raw must decode to the typed zero variant: %#v", decoded)
		}
	}
	if _, err := DecodeFilterConfig("carrier_pigeon", []byte("{}")); err == nil {
		t.Fatalf("an unknown filter type must error, not silently drop")
	}
}
