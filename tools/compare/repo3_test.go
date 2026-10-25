package main

import "testing"

func TestNormalizeVendorTempNames(t *testing.T) {
	cases := map[string]string{
		"x vendor/composer/tmp-0123456789abcdef0123456789abcdef.zip -o/w/vendor/composer/479c6a6e\n": "x vendor/composer/tmp-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx.zip -o/w/vendor/composer/xxxxxxxx\n",
		// wrapped by an exception box
		"  tar -xJf /w/vendor/composer/tmp-f78d3e2f39e87f  \n  a03f1b3f1d9ed9f391.xz -C /w/vendor/composer/445cd40e   \n": "  tar -xJf /w/vendor/composer/tmp-xxxxxxxxxxxxxx  \n  xxxxxxxxxxxxxxxxxx.xz -C /w/vendor/composer/xxxxxxxx   \n",
		// other files and names of another length are kept
		"vendor/composer/autoload_real.php vendor/composer/abc vendor/composer/0123456789": "vendor/composer/autoload_real.php vendor/composer/abc vendor/composer/0123456789",
	}
	for in, want := range cases {
		if got := normalizeVendorTempNames(in); got != want {
			t.Errorf("normalize(%q) = %q, want %q", in, got, want)
		}
	}
}
