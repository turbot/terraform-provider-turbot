package turbot

import (
	"regexp"
	"slices"
	"sort"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
)

// credentialAttribute matches attribute names that hold a credential. Terraform prints an
// attribute's value in plan and apply output unless its schema sets Sensitive, so every match
// must set it. See GHSA-7jmh-m2pr-47hj.
var credentialAttribute = regexp.MustCompile(`secret|password|private_key|token`)

// notCredentials match credentialAttribute but hold no credential.
var notCredentials = map[string]bool{
	"turbot_local_directory_user.password_timestamp": true,
}

func TestCredentialAttributesAreSensitive(t *testing.T) {
	p := Provider().(*schema.Provider)

	var found, missing []string
	check := func(path string, s *schema.Schema) {
		found = append(found, path)
		if !s.Sensitive {
			missing = append(missing, path)
		}
	}
	walkCredentialAttributes("provider", p.Schema, check)
	for name, r := range p.ResourcesMap {
		walkCredentialAttributes(name, r.Schema, check)
	}
	for name, r := range p.DataSourcesMap {
		walkCredentialAttributes("data."+name, r.Schema, check)
	}

	sort.Strings(missing)
	for _, path := range missing {
		t.Errorf("%s holds a credential but is not Sensitive", path)
	}

	// Guard the scan itself: if these stop matching, the walk is broken and the test proves nothing.
	for _, want := range []string{
		"provider.secret_key",
		"turbot_google_directory.client_secret",
		"turbot_ldap_directory.password",
		"turbot_saml_directory.signature_private_key",
	} {
		if !slices.Contains(found, want) {
			t.Errorf("scan did not reach %s", want)
		}
	}
}

func walkCredentialAttributes(prefix string, attrs map[string]*schema.Schema, check func(string, *schema.Schema)) {
	for name, s := range attrs {
		path := prefix + "." + name
		if credentialAttribute.MatchString(name) && !notCredentials[path] {
			check(path, s)
		}
		if nested, ok := s.Elem.(*schema.Resource); ok {
			walkCredentialAttributes(path, nested.Schema, check)
		}
	}
}
