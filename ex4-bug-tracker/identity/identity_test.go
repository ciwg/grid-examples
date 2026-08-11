package identity_test

import (
	"testing"

	"github.com/computerscienceiscool/grid-examples/ex4-bug-tracker/identity"
)

func TestEnrollmentProofVerifiesForMatchingKey(t *testing.T) {
	key, err := identity.NewAgentKey()
	if err != nil {
		t.Fatalf("new key: %v", err)
	}
	enrollment := identity.NewEnrollment(key, "reporter")
	proof, err := key.SignEnrollment(enrollment)
	if err != nil {
		t.Fatalf("sign enrollment: %v", err)
	}
	if err := identity.VerifyEnrollment(enrollment, proof); err != nil {
		t.Fatalf("verify enrollment: %v", err)
	}
}

func TestEnrollmentProofRejectsWrongKey(t *testing.T) {
	key, err := identity.NewAgentKey()
	if err != nil {
		t.Fatalf("new key: %v", err)
	}
	other, err := identity.NewAgentKey()
	if err != nil {
		t.Fatalf("new other key: %v", err)
	}
	enrollment := identity.NewEnrollment(key, "reporter")
	proof, err := other.SignEnrollment(enrollment)
	if err != nil {
		t.Fatalf("sign enrollment: %v", err)
	}
	if err := identity.VerifyEnrollment(enrollment, proof); err == nil {
		t.Fatal("expected wrong-key proof rejection")
	}
}
