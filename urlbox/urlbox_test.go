/*
** urlbox_test.go
** Author: Marin Alcaraz
** Mail:   <marin.alcaraz@gmail.com>
 */

package urlbox

import (
	"os"
	"testing"
)

type testCase struct {
	input    string
	expected string
	fileName string
}

func TestGetScreenshot(t *testing.T) {
	cases := []testCase{
		{"https://api.urlbox.io/v1/db411da0-2525-40a1-add5-ad2cf85a3a51/" +
			"3810736cc9d2bb5a0dd0b4bb7ce8de9eb57e27d6/" +
			"png?url=marinhero.com", "200 OK", "TestGetScreenshot-1.png"},
		{"https://api.urlbox.io/v1/db411da0-2525-40a1-add5-ad2cf85a3a51/" +
			"95828f59b086bc19c353ebf9da616712ae1097ee/" +
			"png?url=wootric.com", "200 OK", "TestGetScreenshot-2.png"},
		{"https://api.urlbox.io/v1/db411da0-2525-40a1-add5-ad2cf85a3a51/" +
			"336e6da581fbbc688bc58857456fdb72a117dbaf/" +
			"png?url=marin", "400 Bad Request", "test-3.png"},
		{"https://api.urlbox.io/v1/db411da0-2525-40a1-add5-ad2cf85a3a51/" +
			"invalidToken/png?url=wootric.com", "401 Unauthorized", "test-4.png"},
	}
	for _, tCase := range cases {
		output := GetScreenshot(tCase.input, tCase.fileName)
		if output != tCase.expected {
			t.Errorf("[!GetScreenshot] IN: %q EXPECTED: %q GOT: %q",
				tCase.input,
				tCase.expected,
				output)
		}
		os.Remove(tCase.fileName)
	}
}

func TestGenerateToken(t *testing.T) {
	cases := []testCase{
		{"url=marinhero.com",
			"3810736cc9d2bb5a0dd0b4bb7ce8de9eb57e27d6", ""},
		{"url=wootric.com",
			"95828f59b086bc19c353ebf9da616712ae1097ee", ""},
		{"url=marinhero.com&width=100&height=100",
			"6268e30b9f27b6a9ff9dd49ddabe46c46163c714", ""},
		{"url=wootric.com&width=1024&height=768",
			"f23f6656a857b02853276eb47c44c7cc6adff566", ""},
	}
	for _, tCase := range cases {
		output := GenerateToken(tCase.input)
		if output != tCase.expected {
			t.Errorf("[!GenerateToken] IN: %q EXPECTED: %q GOT: %q",
				tCase.input,
				tCase.expected,
				output)
		}
	}
}

func TestCreateRequestString(t *testing.T) {
	cases := []testCase{
		{"url=marinhero.com",
			"https://api.urlbox.io/v1/db411da0-2525-40a1-add5-ad2cf85a3a51/" +
				"3810736cc9d2bb5a0dd0b4bb7ce8de9eb57e27d6/" +
				"png?url=marinhero.com", ""},
		{"url=wootric.com",
			"https://api.urlbox.io/v1/db411da0-2525-40a1-add5-ad2cf85a3a51/" +
				"95828f59b086bc19c353ebf9da616712ae1097ee/" +
				"png?url=wootric.com", ""},
		{"url=marinhero.com&width=100&height=100",
			"https://api.urlbox.io/v1/db411da0-2525-40a1-add5-ad2cf85a3a51/" +
				"6268e30b9f27b6a9ff9dd49ddabe46c46163c714/" +
				"png?url=marinhero.com&width=100&height=100", ""},
		{"url=wootric.com&width=1024&height=768",
			"https://api.urlbox.io/v1/db411da0-2525-40a1-add5-ad2cf85a3a51/" +
				"f23f6656a857b02853276eb47c44c7cc6adff566/" +
				"png?url=wootric.com&width=1024&height=768", ""},
	}
	for _, tCase := range cases {
		output := CreateRequestString(tCase.input)
		if output != tCase.expected {
			t.Errorf("[!CreateRequestString]IN: %q EXPECTED: %q GOT: %q",
				tCase.input,
				tCase.expected,
				output)
		}
	}
}

func TestCreateShot(t *testing.T) {
	cases := []struct {
		url           string
		widht, height uint
		expected      string
	}{
		{"marinhero.com", 1024, 768, "OK"},
		{"wootric.com", 1024, 768, "OK"},
		{"", 0, 0, "KO"},
		{"", 1024, 768, "KO"},
		{"wootric.com", 1024, 0, "KO"},
		{"keatonrow.com", 0, 768, "KO"},
	}
	var shotinfo ShotData
	for _, tCase := range cases {
		shotinfo.URL = tCase.url
		shotinfo.Width = tCase.widht
		shotinfo.Height = tCase.height
		output := CreateShot(shotinfo)
		if output != tCase.expected {
			t.Errorf("[!CreateShot]IN: %q EXPECTED: %q GOT: %q",
				shotinfo,
				tCase.expected,
				output)
		}
		os.Remove(GetFileName(shotinfo))
	}
}
