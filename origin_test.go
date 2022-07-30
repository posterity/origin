package origin

import (
	"testing"
)

func TestPattern(t *testing.T) {
	type testCase struct {
		Origin   string
		Pattern  string
		HasError bool
		IsMatch  bool
	}

	var cases = []*testCase{
		{"example.com", "https://example.com", true, false},
		{"*://example.com", "https://example.com", true, false},
		{"", "*", true, false},
		{"https://example.com", "https://example.com", false, true},
		{"https://a.sub.example.com", "https://*.sub.example.com", false, true},
		{"https://a.sub.example.com", "https://*.*.sub.example.com", false, false},
		{"https://a.458.sub.example.com", "https://*.*.sub.example.com", false, true},
		{"https://a.sub.example.com:443", "https://*.sub.example.com", false, true},
		{"https://sub.example.com", "https://*.example.com", false, true},
		{"http://sub.example.com", "*://sub.example.com:80", false, true},
		{"http://sub.example.com", "*://sub.example.com:8000", false, false},
		{"https://sub.example.com", "https://sub.example.dev", false, false},
		{"ws://sub.example.com", "https://sub.example.dev", false, false},
		{"https://sub.example.dev", "https://sub.*.dev", false, true},
		{"https://example.com", "https://example.dev", false, false},
		{"https://example.example", "https://example.example", false, true},
		{"https://example.example:8080", "https://example.example:*", false, true},
		{"https://example.dev:443", "https://example.dev:*", false, true},
		{"custom://example.com:54232", "custom://example.com:54232", false, true},
		{"custom://example.com:54232", "*://example.com:54232", false, true},
		{"custom://example.com:54232", "*", false, true},
		{"custom://example.com:54232", "*://*:*", false, true},
		{"abcdef", "*://*:*", true, false},
	}

	for _, tc := range cases {
		isMatch, err := Match(tc.Origin, tc.Pattern)
		if hasErr := (err != nil); hasErr != tc.HasError {
			t.Errorf("Origin: %s, Pattern: %s - Error: %v", tc.Origin, tc.Pattern, err)
		}
		if tc.IsMatch != isMatch {
			t.Errorf("Origin: %s, Pattern: %s - Wanted: %v, Got: %v", tc.Origin, tc.Pattern, tc.IsMatch, isMatch)
		}
	}
}
