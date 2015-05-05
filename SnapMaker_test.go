/*
** SnapMaker_test.go
** Author: Marin Alcaraz
** Mail   <mailto@alcarazmar.in>
 */

package main

import (
	"testing"

	"github.com/marinhero/url2snap/urlbox"
)

func TestUploadToS3(t *testing.T) {
	cases := []struct {
		url           string
		widht, height uint
		expected      string
	}{
		{"marinhero.com", 1024, 768,
			"https://s3.amazonaws.com/" +
				"snapshotswootric/marinhero.com1024x768.png"},
		{"wootric.com", 1024, 768,
			"https://s3.amazonaws.com/" +
				"snapshotswootric/wootric.com1024x768.png"},
		{"", 0, 0, ""},
	}
	var shotinfo urlbox.ShotData
	for _, tCase := range cases {
		shotinfo.URL = tCase.url
		shotinfo.Width = tCase.widht
		shotinfo.Height = tCase.height
		urlbox.CreateShot(shotinfo)
		output := uploadToS3(shotinfo)
		if output != tCase.expected {
			t.Errorf("[!CreateShot]IN: %q EXPECTED: %q GOT: %q",
				shotinfo,
				tCase.expected,
				output)
		}
	}
}
