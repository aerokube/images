package build

import (
	. "github.com/aandryashin/matchers"
	"testing"
)

func TestVersionFromPackageName(t *testing.T) {
	AssertThat(t, versionFromPackageName("firefox_45.0.2+build1-0ubuntu0.14.04.1+aerokube0_amd64.deb"), EqualTo{"45.0.2"})
	AssertThat(t, versionFromPackageName("google-chrome-stable_48.0.2564.109-1+aerokube0_amd64.deb"), EqualTo{"48.0.2564.109"})
	AssertThat(t, versionFromPackageName("opera-stable_38.0.2220.31_amd64.deb"), EqualTo{"38.0.2220.31"})
}
