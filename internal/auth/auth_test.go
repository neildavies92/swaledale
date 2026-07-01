package auth

import "testing"

func TestPasswordRoundTrip(t *testing.T) {
	hash, err := HashPassword("correct horse battery staple")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	if !CheckPassword(hash, "correct horse battery staple") {
		t.Fatal("expected matching password to verify")
	}
	if CheckPassword(hash, "wrong password") {
		t.Fatal("expected wrong password to fail")
	}
}

func TestSessionTokens(t *testing.T) {
	first, err := NewSessionToken()
	if err != nil {
		t.Fatalf("new session token: %v", err)
	}
	second, err := NewSessionToken()
	if err != nil {
		t.Fatalf("new session token: %v", err)
	}
	if first == second {
		t.Fatal("expected unique tokens")
	}
	if HashSessionToken(first) == first {
		t.Fatal("expected hash to differ from token")
	}
	if HashSessionToken(first) != HashSessionToken(first) {
		t.Fatal("expected hash to be deterministic")
	}
}
