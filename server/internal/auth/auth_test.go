package auth

import "testing"

func TestHashAndVerify(t *testing.T) {
	hash, err := HashPassword("s3cret")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	if !VerifyPassword(hash, "s3cret") {
		t.Fatal("VerifyPassword rejected the correct password")
	}
	if VerifyPassword(hash, "wrong") {
		t.Fatal("VerifyPassword accepted a wrong password")
	}
	if VerifyPassword("not-a-hash", "s3cret") {
		t.Fatal("VerifyPassword accepted a malformed hash")
	}
}

func TestHashIsSalted(t *testing.T) {
	a, _ := HashPassword("same")
	b, _ := HashPassword("same")
	if a == b {
		t.Fatal("hashes of the same password should differ (random salt)")
	}
}

func TestNewTokenIsUnique(t *testing.T) {
	if NewToken() == NewToken() {
		t.Fatal("NewToken returned duplicate tokens")
	}
}
