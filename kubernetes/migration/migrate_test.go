package migration

import (
	"testing"
)

func TestMigrate(t *testing.T) {
	testCases := []struct {
		name             string
		fromVersion      string
		toVersion        string
		deprecations     bool
		startCorefile    string
		expectedCorefile string
	}{
		{
			name:         "Remove invalid proxy option",
			fromVersion:  "1.1.3",
			toVersion:    "1.2.6",
			deprecations: true,
			startCorefile: `.:53 {
    errors
    health
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        endpoint thing1 thing2
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    proxy example.org 1.2.3.4:53 {
        protocol https_google
    }
    cache 30
    loop
    reload
    loadbalance
}
`,
			expectedCorefile: `.:53 {
    errors
    health
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        endpoint thing1 thing2
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    proxy example.org 1.2.3.4:53
    cache 30
    loop
    reload
    loadbalance
}
`,
		},
		{
			name:         "Migrate from proxy to forward and handle Kubernetes deprecations",
			fromVersion:  "1.3.1",
			toVersion:    "1.5.0",
			deprecations: true,
			startCorefile: `.:53 {
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
`,
			expectedCorefile: `.:53 {
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
    ready
}
`,
		},
		{
			name:         "add missing loop and ready plugins",
			fromVersion:  "1.1.3",
			toVersion:    "1.5.0",
			deprecations: true,
			startCorefile: `.:53 {
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
    reload
    loadbalance
}
`,
			expectedCorefile: `.:53 {
    errors
    health
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    forward . /etc/resolv.conf
    cache 30
    reload
    loadbalance
    loop
    ready
}
`,
		},
		{
			name:         "handle multiple proxy plugins",
			fromVersion:  "1.1.3",
			toVersion:    "1.5.0",
			deprecations: true,
			startCorefile: `.:53 {
    errors
    health
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    proxy mystub-1.example.org 1.2.3.4
    proxy mystub-2.example.org 5.6.7.8
    proxy . /etc/resolv.conf
    cache 30
    reload
    loadbalance
}
`,
			expectedCorefile: `.:53 {
    errors
    health
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    forward . /etc/resolv.conf
    cache 30
    reload
    loadbalance
    loop
    ready
}

mystub-1.example.org {
    forward . 1.2.3.4
    loop
    errors
    cache 30
}

mystub-2.example.org {
    forward . 5.6.7.8
    loop
    errors
    cache 30
}
`,
		},
		{
			name:         "no-op same version migration",
			fromVersion:  "1.3.1",
			toVersion:    "1.3.1",
			deprecations: true,
			startCorefile: `.:53 {
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
`,
			expectedCorefile: `.:53 {
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
`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			result, err := Migrate(testCase.fromVersion, testCase.toVersion, testCase.startCorefile, testCase.deprecations)

			if err != nil {
				t.Errorf("%v", err)
			}

			if result != testCase.expectedCorefile {
				t.Errorf("expected != result\n%v\n%v", testCase.expectedCorefile, result)
			}
		})
	}
}

func TestMigrateDown(t *testing.T) {
	testCases := []struct {
		name             string
		fromVersion      string
		toVersion        string
		deprecations     bool
		startCorefile    string
		expectedCorefile string
	}{
		{
			name:        "from 1.5.0 to 1.1.3",
			fromVersion: "1.5.0",
			toVersion:   "1.1.3",
			startCorefile: `.:53 {
    errors
    health
    ready
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
			expectedCorefile: `.:53 {
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
    reload
    loadbalance
}
`,
		},
		{
			name:         "no-op same version migration",
			fromVersion:  "1.3.1",
			toVersion:    "1.3.1",
			deprecations: true,
			startCorefile: `.:53 {
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
`,
			expectedCorefile: `.:53 {
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
`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			result, err := MigrateDown(testCase.fromVersion, testCase.toVersion, testCase.startCorefile)

			if err != nil {
				t.Errorf("%v", err)
			}

			if result != testCase.expectedCorefile {
				t.Errorf("expected != result:\n%v\n%v", testCase.expectedCorefile, result)
			}
		})
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
		{Option: "upstream", Plugin: "kubernetes", Severity: ignored, Version: "1.5.0"},
		{Plugin: "proxy", Severity: removed, ReplacedBy: "forward", Version: "1.5.0"},
		{Plugin: "ready", Severity: newdefault, Version: "1.5.0"},
	}

	result, err := Deprecated("1.2.0", "1.5.0", startCorefile)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != len(expected) {
		t.Fatalf("expected to find %v notifications; got %v", len(expected), len(result))
	}

	for i, dep := range expected {
		if result[i].ToString() != dep.ToString() {
			t.Errorf("expected to get '%v'; got '%v'", dep.ToString(), result[i].ToString())
		}
	}

	result, err = Deprecated("1.3.1", "1.3.1", startCorefile)
	if err != nil {
		t.Fatal(err)
	}
	expected = []Notice{}
	if len(result) != len(expected) {
		t.Fatalf("expected to find %v notifications in no-op upgrade; got %v", len(expected), len(result))
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

func TestValidUpMigration(t *testing.T) {
	testCases := []struct {
		from      string
		to        string
		shouldErr bool
	}{
		{"1.3.1", "1.3.1", true},
		{"1.3.1", "1.5.0", false},
		{"1.5.0", "1.3.1", true},
		{"banana", "1.5.0", true},
		{"1.3.1", "apple", true},
		{"banana", "apple", true},
	}

	for _, tc := range testCases {
		err := validUpMigration(tc.from, tc.to)

		if !tc.shouldErr && err != nil {
			t.Errorf("expected '%v' to '%v' to be valid versions.", tc.from, tc.to)
		}
		if tc.shouldErr && err == nil {
			t.Errorf("expected '%v' to '%v' to be invalid versions.", tc.from, tc.to)
		}
	}
}

func TestValidDownMigration(t *testing.T) {
	testCases := []struct {
		from      string
		to        string
		shouldErr bool
	}{
		{"1.3.1", "1.3.1", true},
		{"1.3.1", "1.5.0", true},
		{"1.5.0", "1.3.1", false},
		{"banana", "1.5.0", true},
		{"1.3.1", "apple", true},
		{"banana", "apple", true},
	}

	for _, tc := range testCases {
		err := validDownMigration(tc.from, tc.to)

		if !tc.shouldErr && err != nil {
			t.Errorf("expected '%v' to '%v' to be valid versions.", tc.from, tc.to)
		}
		if tc.shouldErr && err == nil {
			t.Errorf("expected '%v' to '%v' to be invalid versions.", tc.from, tc.to)
		}
	}
}
