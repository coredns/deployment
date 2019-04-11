package migration

import (
	"testing"
)

func TestMigrate(t *testing.T) {
	startCorefile := `.:53 {
    errors
    health
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        endpoint thing1 thing2
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    proxy . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
}
`

	expected := `.:53 {
    errors
    health
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        endpoint thing1
        pods insecure
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    forward . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
}
`
	result, err := Migrate("1.3.1", "1.5.0", startCorefile, true)

	if err != nil {
		t.Errorf("%v", err)
	}

	if result != expected {
		t.Errorf("expected %v; got %v", expected, result)
	}
}

func TestDeprecated(t *testing.T) {
	startCorefile := `.:53 {
    errors
    health
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    proxy . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
}
`

	expected := []Notice{
		{Plugin: "kubernetes", Option: "upstream", Severity: deprecated, Version: "1.4.0"},
		{Plugin: "proxy", Severity: deprecated, ReplacedBy: "forward", Version: "1.4.0"},
	}

	result, err := Deprecated("1.3.1", "1.5.0", startCorefile)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != len(expected) {
		t.Fatalf("expected to find %v deprecations; got %v", len(expected), len(result))
	}

	for i, dep := range expected {
		if result[i].ToString() != dep.ToString() {
			t.Errorf("expected to get '%v'; got '%v'", dep.ToString(), result[i].ToString())
		}
	}
}

func TestRemoved(t *testing.T) {
	startCorefile := `.:53 {
    errors
    health
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    proxy . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
}
`

	expected := []Notice{
		{Plugin: "proxy", Severity: removed, ReplacedBy: "forward", Version: "1.5.0"},
	}

	result, err := Removed("1.3.1", "1.5.0", startCorefile)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != len(expected) {
		t.Fatalf("expected to find %v deprecations; got %v", len(expected), len(result))
	}

	for i, dep := range expected {
		if result[i].ToString() != dep.ToString() {
			t.Errorf("expected to get '%v'; got '%v'", dep.ToString(), result[i].ToString())
		}
	}
}

func TestUnsupported(t *testing.T) {
	startCorefile := `.:53 {
    errors
    health
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
    }
    route53 example.org.:Z1Z2Z3Z4DZ5Z6Z7
    prometheus :9153
    proxy . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
}
`

	expected := []Notice{
		{Plugin: "route53", Severity: unsupported, Version: "1.4.0"},
		{Plugin: "route53", Severity: unsupported, Version: "1.5.0"},
	}

	result, err := Unsupported("1.3.1", "1.5.0", startCorefile)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != len(expected) {
		t.Fatalf("expected to find %v deprecations; got %v", len(expected), len(result))
	}

	for i, dep := range expected {
		if result[i].ToString() != dep.ToString() {
			t.Errorf("expected to get '%v'; got '%v'", dep.ToString(), result[i].ToString())
		}
	}
}

func TestDefault(t *testing.T) {
	defaultCorefiles := []string{`.:53 {
    errors
    health
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    forward . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
}
`,
		`.:53 {
    errors
    health
    kubernetes myzone.org in-addr.arpa ip6.arpa {
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    forward . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
}
`}

	nonDefaultCorefiles := []string{`.:53 {
    errors
    health
    rewrite name suffix myzone.org cluster.local
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    forward . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
}
`,
		`.:53 {
    errors
    health
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    forward . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
}
stubzone.org:53 {
    forward . 1.2.3.4
}
`}

	for _, d := range defaultCorefiles {
		if !Default("", d) {
			t.Errorf("expected config to be identified as a default: %v", d)
		}
	}
	for _, d := range nonDefaultCorefiles {
		if Default("", d) {
			t.Errorf("expected config to NOT be identified as a default: %v", d)
		}
	}
}

func TestValidateVersions(t *testing.T) {
	testCases := []struct {
		from   string
		to     string
		shouldErr    bool
	}{
		{"1.3.1", "1.5.0", false},
		{"banana", "1.5.0", true},
		{"1.3.1", "apple", true},
		{"banana", "apple", true},
	}

	for _, tc := range testCases {
		err := validateVersions(tc.from, tc.to)

		if !tc.shouldErr && err != nil {
			t.Errorf("expected to '%v' to '%v' to be valid versions.", tc.from, tc.to)
		}
		if tc.shouldErr && err == nil {
			t.Errorf("expected to '%v' to '%v' to be invalid versions.", tc.from, tc.to)
		}
	}
}
