// Tests for FormatFloat, the json_encode() float formatting used by
// JsonFile::encode in src/Composer/Json/JsonManipulator.php (composer/composer
// f83588392df43d4fd027d84d26de459946ef4e80).
// Composer is MIT licensed, Copyright (c) Nils Adermann, Jordi Boggiano.

package jsonmanip

import "testing"

// Expected values are PHP 8.4 json_encode() output (serialize_precision=-1).
func TestFormatFloat(t *testing.T) {
	for in, want := range map[float64]string{
		1.0: "1", 100.0: "100", 1e15: "1000000000000000", 1e20: "1.0e+20", 0.1: "0.1",
		1.5e-7: "1.5e-7", 123456789.0: "123456789", 1e16: "10000000000000000", 1e17: "1.0e+17",
		0.0001: "0.0001", 0.00001: "1.0e-5", -2.5: "-2.5", 1.25e100: "1.25e+100",
	} {
		if got := FormatFloat(in); got != want {
			t.Errorf("FormatFloat(%v) = %q, want %q", in, got, want)
		}
	}
}
