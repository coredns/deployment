package corefile

import (
	"testing"
)

func TestCorefile(t *testing.T) {

	startCorefile := `.:53 {
    error
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

.:5353 {
    proxy . /etc/resolv.conf
}
`
	c, err := New(startCorefile)

	if err != nil {
		t.Error(err)
	}

	got := c.ToString()

	if got != startCorefile {
		t.Errorf("Corefile did not match expected.\nExpected:\n%v\nGot:\n%v", startCorefile, got)
	}
}

func TestServer_FindMatch(t *testing.T) {
	tests := []struct{
		server *Server
		match bool
	}{
		{server: &Server{DomPorts: []string{".:53"}}, match: true},
		{server: &Server{DomPorts: []string{".:54"}}, match: false},
		{server: &Server{DomPorts: []string{"abc:53"}}, match: false},
		{server: &Server{DomPorts: []string{"abc:53", "blah"}}, match: true},
		{server: &Server{DomPorts: []string{"abc:53", "blah", "blah"}}, match: false},
		{server: &Server{DomPorts: []string{"xyz:53"}}, match: true},
		{server: &Server{DomPorts: []string{"xyz:53", "blah", "blah"}}, match: true},
	}

	def := []*Server{
		{DomPorts: []string{".:53"}},
		{DomPorts: []string{"abc:53", "*"}},
		{DomPorts: []string{"xyz:53", "***"}},
	}
	for i, test := range tests {
		_, match := test.server.FindMatch(def)
		if match != test.match {
			t.Errorf("In test #%v, expected match to be %v but got %v.", i, test.match, match)
		}
	}
}

func TestPlugin_FindMatch(t *testing.T) {
	tests := []struct{
		plugin *Plugin
		match bool
	}{
		{plugin: &Plugin{Name: "plugin1", Args: []string{}}, match: true},
		{plugin: &Plugin{Name: "plugin2", Args: []string{"1","1.5","2"}}, match: true},
		{plugin: &Plugin{Name: "plugin3", Args: []string{"1","2","3","4"}}, match: true},
		{plugin: &Plugin{Name: "plugin1", Args: []string{"a"}}, match: false},
		{plugin: &Plugin{Name: "plugin2", Args: []string{"1","1.5","b"}}, match: false},
		{plugin: &Plugin{Name: "plugin3", Args: []string{"a","2","3","4"}}, match: false},
		{plugin: &Plugin{Name: "plugin4", Args: []string{}}, match: false},
	}

	def := []*Plugin{
		{Name: "plugin1", Args: []string{}},
		{Name: "plugin2", Args: []string{"1", "*", "2"}},
		{Name: "plugin3", Args: []string{"1", "***"}},
	}
	for i, test := range tests {
		_, match := test.plugin.FindMatch(def)
		if match != test.match {
			t.Errorf("In test #%v, expected match to be %v but got %v.", i, test.match, match)
		}
	}
}

func TestOption_FindMatch(t *testing.T) {
	tests := []struct{
		option *Plugin
		match bool
	}{
		{option: &Plugin{Name: "option1", Args: []string{}}, match: true},
		{option: &Plugin{Name: "option2", Args: []string{"1","1.5","2"}}, match: true},
		{option: &Plugin{Name: "option3", Args: []string{"1","2","3","4"}}, match: true},
		{option: &Plugin{Name: "option1", Args: []string{"a"}}, match: false},
		{option: &Plugin{Name: "option2", Args: []string{"1","1.5","b"}}, match: false},
		{option: &Plugin{Name: "option3", Args: []string{"a","2","3","4"}}, match: false},
		{option: &Plugin{Name: "option4", Args: []string{}}, match: false},
	}

	def := []*Plugin{
		{Name: "option1", Args: []string{}},
		{Name: "option2", Args: []string{"1", "*", "2"}},
		{Name: "option3", Args: []string{"1", "***"}},
	}
	for i, test := range tests {
		_, match := test.option.FindMatch(def)
		if match != test.match {
			t.Errorf("In test #%v, expected match to be %v but got %v.", i, test.match, match)
		}
	}
}