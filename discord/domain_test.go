package discord

import (
	"testing"
)

func TestDomainInfo(t *testing.T) {
	info := Domain{}.Info()
	if info.Scheme != "discord" {
		t.Errorf("Scheme = %q, want discord", info.Scheme)
	}
	if info.Identity.Binary != "discord" {
		t.Errorf("Identity.Binary = %q, want discord", info.Identity.Binary)
	}
}

func TestClassify(t *testing.T) {
	cases := []struct{ in, typ, id string }{
		{"invite/abc123", "page", "invite/abc123"},
		{"https://discord.com/invite/abc", "page", "invite/abc"},
	}
	for _, tc := range cases {
		typ, id, err := Domain{}.Classify(tc.in)
		if err != nil || typ != tc.typ || id != tc.id {
			t.Errorf("Classify(%q) = (%q, %q, %v), want (%q, %q, nil)",
				tc.in, typ, id, err, tc.typ, tc.id)
		}
	}
}

func TestLocate(t *testing.T) {
	got, err := Domain{}.Locate("page", "invite/abc")
	want := "https://discord.com/invite/abc"
	if err != nil || got != want {
		t.Errorf("Locate = (%q, %v), want (%q, nil)", got, err, want)
	}
}
