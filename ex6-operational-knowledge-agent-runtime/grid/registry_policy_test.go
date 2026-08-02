package grid

import "testing"

func TestRegistryAllowListPersistsCanonicalHosts(t *testing.T) {
	store, err := OpenPolicyStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if err := store.AllowRegistry("REGISTRY.example:5000"); err != nil {
		t.Fatal(err)
	}
	if err := store.AllowRegistry("registry.example:5000"); err != nil {
		t.Fatal(err)
	}
	if got := store.RegistryAllowList(); len(got) != 1 || got[0] != "registry.example:5000" {
		t.Fatalf("allow-list = %#v", got)
	}
	if err := store.RemoveRegistry("registry.example:5000"); err != nil {
		t.Fatal(err)
	}
	if got := store.RegistryAllowList(); len(got) != 0 {
		t.Fatalf("allow-list = %#v", got)
	}
	for _, host := range []string{"registry.example:bad", "registry.example:5000:extra", "registry .example", "*.example"} {
		if err := store.AllowRegistry(host); err == nil {
			t.Fatalf("accepted invalid host %q", host)
		}
	}
}
