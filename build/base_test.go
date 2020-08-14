package build

import (
	. "github.com/aandryashin/matchers"
	"testing"
)

func TestVersionFromPackageName(t *testing.T) {
	AssertThat(t, extractVersion("firefox_45.0.2+build1-0ubuntu0.14.04.1+aerokube0_amd64.deb"), EqualTo{"45.0.2"})
	AssertThat(t, extractVersion("google-chrome-stable_48.0.2564.109-1+aerokube0_amd64.deb"), EqualTo{"48.0.2564.109"})
	AssertThat(t, extractVersion("opera-stable_38.0.2220.31_amd64.deb"), EqualTo{"38.0.2220.31"})
}

func TestVersionToTag(t *testing.T) {
	AssertThat(t, extractVersion("45.0.2+build1-0ubuntu0.14.04.1+aerokube0"), EqualTo{"45.0.2"})
	AssertThat(t, extractVersion("48.0.2564.109-1+aerokube0"), EqualTo{"48.0.2564.109"})
	AssertThat(t, extractVersion("38.0.2220.31"), EqualTo{"38.0.2220.31"})
}
