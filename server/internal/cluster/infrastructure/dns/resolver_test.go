package dns

import "testing"

func TestIncludesIP(t *testing.T) {
	t.Run("returns true when expected IP is present", func(t *testing.T) {
		if !includesIP([]string{"203.0.113.1", "203.0.113.2"}, "203.0.113.1") {
			t.Fatal("expected IP to be found")
		}
	})

	t.Run("returns false when expected IP is missing", func(t *testing.T) {
		if includesIP([]string{"203.0.113.2"}, "203.0.113.1") {
			t.Fatal("expected IP not to be found")
		}
	})
}

func TestResolver_AllHostsResolve(t *testing.T) {
	resolver := NewResolver()

	t.Run("returns true when all hosts resolve to expected IP", func(t *testing.T) {
		hosts := []string{"one.one.one.one"}
		if !resolver.AllHostsResolve(hosts, "1.1.1.1") {
			t.Fatalf("expected one.one.one.one to resolve to 1.1.1.1 via public resolver")
		}
	})

	t.Run("returns false for unknown host", func(t *testing.T) {
		hosts := []string{"this-host-should-not-exist.starliner.invalid"}
		if resolver.AllHostsResolve(hosts, "1.1.1.1") {
			t.Fatalf("expected unknown host lookup to fail")
		}
	})

	t.Run("returns false when IP does not match", func(t *testing.T) {
		hosts := []string{"one.one.one.one"}
		if resolver.AllHostsResolve(hosts, "203.0.113.1") {
			t.Fatalf("expected one.one.one.one not to match unrelated IP")
		}
	})
}
