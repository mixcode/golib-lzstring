package lzstring

import (
	"log"
	"testing"
)

func TestCompress(t *testing.T) {

	testF := func(s, b64 string) {
		sB64 := CompressToBase64(s)
		if sB64 != b64 {
			t.Errorf("encode failed")
			log.Printf("%v", sB64)
		}
		d, err := DecompressBase64(sB64)
		if err != nil {
			t.Error(err)
		}
		if d != s {
			t.Errorf("decode failed")
		}
	}

	// sample from pieroxy's test code
	// https://github.com/pieroxy/lz-string/blob/master/tests/lz-string-spec.js
	testF(
		`During tattooing, ink is injected into the skin, initiating an immune response, and cells called "macrophages" move into the area and "eat up" the ink. The macrophages carry some of the ink to the body's lymph nodes, but some that are filled with ink stay put, embedded in the skin. That's what makes the tattoo visible under the skin. Dalhousie Uiversity's Alec Falkenham is developing a topical cream that works by targeting the macrophages that have remained at the site of the tattoo. New macrophages move in to consume the previously pigment-filled macrophages and then migrate to the lymph nodes, eventually taking all the dye with them. "When comparing it to laser-based tattoo removal, in which you see the burns, the scarring, the blisters, in this case, we've designed a drug that doesn't really have much off-target effect," he said. "We're not targeting any of the normal skin cells, so you won't see a lot of inflammation. In fact, based on the process that we're actually using, we don't think there will be any inflammation at all and it would actually be anti-inflammatory.`,
		`CIVwTglgdg5gBAFwIYIQezdGAaO0DWeAznlAFYCmAxghQCanqIAWFcR+0u0ECEKWOEih4AtqJBQ2YCkQAOaKEQq5hDKhQA2mklSTb6cAESikVMGjnMkMWUbii0ANzbQmCVkJlIhUBkYoUOBA5ew9XKHwAOjgAFU9Tc0trW10kMDAAT3Y0UTY0ADMWCMJ3TwAjNDpMgHISTUzRKzgoKtlccpAEHLyWIPS2AogDBgB3XmZSQiJkbLku3ApRcvo6Q2hi9k4oGPiUOrhR627TfFlN5FQMOCcIIghyzTZJNbBNjmgY4H1mNBB7tgAVQgLjA9wQtRIAEEnlQ4AAxfRnKDWUTEOBrFyaSyCHzoOQQPSaODmQJojxBUZoMD4EjlbLIMC2PiwTaJCxWGznCndawuOAyUzQQxBcLsXj5Ipiy7oNAxAByFFGDjMHJS50c/I2TCoiiIIF6YrkMlufyIDTgBJgeSgCAAtEMRiqkpzUr4GOERKIIDAwCg2GU2A0mpNWmsiIsXLaQPoLchtvBY5tqmxxh5iqIYkYAOqsES6prpQS8RBoOCaJDKMB28qVwwy66C5z6bgiI6EyaZP7sCgBirgJS4MVEPQZLBDiqaO60MGtlh3El13CjCg1fnhn1SBg+OhgEDwHkYtCyKA1brebTZPlsCRUSaFAp2xnMuAUAoFagIbD2TxEJAQOgs2zVcZBaNBumfCgWUTKBskKTZWjAUxiQ+fMtB0XAiDLLsQEORQzx7NgfGxbp4OgAoK3EARFBiABJEQCjML84FrZQGEUTZjTQDQiBIQ8VxqUCmJjS9gnuWBlzYOh8Ig5gCGKUDxm0FiiNg0gKKQKi+A4/plLUPBuipEBNG3GgRItFZfD4O1yMo0x0CyKIgA`,
	)

	// Test code with unicode
	testF(
		`で、適当に日本語入れて正しく出るか見ればいいってことでしょ？`,
		`uYMkAGSWCYyvqNYMinpoNOaDyowpoqBiGQZgyGNrQ6gyHkGQX4qDRDINIMg0ZEqAODICIMlgxgxqDKDIBYMwWg4QyD4/0AA=`,
	)

	// Json test
	testF(
		`{"alwaysDash":true,"commandRemember":false,"bgmVolume":30,"bgsVolume":60,"meVolume":40,"seVolume":40,"messageSpeedKey":7}`,
		"N4IghgNg7mCeDOARM8AWIBcAXATgVwFMAaEAYwHsBbSsAOwBMAlAylgIwJ0wDNJ5iQbAOaUAauQh5WmAMwAGEsPjjJ0jADYFIViqkFMAFi39daoyVbx4YIQQDKABwIF6AaQKxMAdgC+QAA==",
	)

}
