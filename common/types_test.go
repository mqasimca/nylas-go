package common

import "testing"

func TestParticipant(t *testing.T) {
	p := Participant{
		Name:  "Test User",
		Email: "test@example.com",
	}

	if p.Name != "Test User" {
		t.Errorf("Name = %s, want Test User", p.Name)
	}
	if p.Email != "test@example.com" {
		t.Errorf("Email = %s, want test@example.com", p.Email)
	}
}

func TestParticipant_EmailOnly(t *testing.T) {
	p := Participant{
		Email: "noreply@example.com",
	}

	if p.Name != "" {
		t.Errorf("Name = %s, want empty", p.Name)
	}
	if p.Email != "noreply@example.com" {
		t.Errorf("Email = %s, want noreply@example.com", p.Email)
	}
}
