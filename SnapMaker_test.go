/*
** SnapMaker_test.go
** Author: Marin Alcaraz
** Mail   <mailto@alcarazmar.in>
 */

package main

import "testing"

func TestuploadToS3(t *testing.T) {
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
	}
	var shotinfo ShotData
	for _, tCase := range cases {
		shotinfo.URL = tCase.url
		shotinfo.Width = tCase.widht
		shotinfo.Height = tCase.height
		output := uploadToS3(shotinfo)
		if output != tCase.expected {
			t.Errorf("[!CreateShot]IN: %q EXPECTED: %q GOT: %q",
				shotinfo,
				tCase.expected,
				output)
		}
	}
}
