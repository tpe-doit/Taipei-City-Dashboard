package util

import (
	"strings"
	"testing"
)

func TestPasswordHashRoundTrip(t *testing.T) {
	for _, password := range []string{"normal password", "長密碼密碼密碼", strings.Repeat("x", 100)} {
		hash, err := HashPassword(password)
		if err != nil {
			t.Fatal(err)
		}
		second, err := HashPassword(password)
		if err != nil {
			t.Fatal(err)
		}
		if hash == second || !strings.HasPrefix(hash, passwordHashPrefix) {
			t.Fatal("password hashes must use independent random salts and the current format")
		}
		if valid, upgrade := VerifyPassword(password, hash); !valid || upgrade {
			t.Fatal("current hash did not verify")
		}
		if valid, _ := VerifyPassword(password+"wrong", hash); valid {
			t.Fatal("wrong password accepted")
		}
	}
}

func TestLegacyPasswordUpgrade(t *testing.T) {
	password := "legacy password"
	legacy := HashString(password)
	if valid, upgrade := VerifyPassword(password, legacy); !valid || !upgrade {
		t.Fatal("legacy password must verify and request an upgrade")
	}
	if valid, upgrade := VerifyPassword("wrong", legacy); valid || upgrade {
		t.Fatal("wrong legacy password accepted")
	}
	replacement, err := HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}
	if valid, upgrade := VerifyPassword(password, replacement); !valid || upgrade {
		t.Fatal("upgraded password did not verify")
	}
	// Preserve the deterministic identity hash used by SSO.
	if got := HashString("abc"); got != "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad" {
		t.Fatal("identity hashing changed")
	}
}

func TestPasswordHashRejectsMalformedFormats(t *testing.T) {
	hash, err := HashPassword("password")
	if err != nil {
		t.Fatal(err)
	}
	for _, value := range []string{
		"", "not-a-hash", strings.Repeat("z", 64), strings.ToUpper(HashString("password")),
		strings.Replace(hash, "m=19456", "m=4294967295", 1),
		strings.Replace(hash, "t=2", "t=9999999", 1),
		strings.Replace(hash, "v=19", "v=16", 1),
		hash + "$extra", hash[:len(hash)-1], passwordHashPrefix + strings.Repeat("A", 22) + "$" + strings.Repeat("!", 43),
	} {
		if valid, upgrade := VerifyPassword("password", value); valid || upgrade {
			t.Fatal("malformed hash accepted")
		}
	}
}
