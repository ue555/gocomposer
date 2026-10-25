package main

import "testing"

func TestNormalizeBanner(t *testing.T) {
	logo := "   ______\n  / ____/___  ____ ___  ____  ____  ________  _____\n / /   / __ \\/ __ `__ \\/ __ \\/ __ \\/ ___/ _ \\/ ___/\n/ /___/ /_/ / / / / / / /_/ / /_/ (__  )  __/ /\n\\____/\\____/_/ /_/ /_/ .___/\\____/____/\\___/_/\n                    /_/\n"
	osc := `\e]7501;state=working:app=composer\e\`
	o := osc + logo + `\e[32mComposer\e[39m version \e[33m2.11.999-dev+source\e[39m @release_date@` + "\nx"
	g := osc + `\e[32mgocomposer\e[39m version \e[33m0.1.0-dev\e[39m (Composer-compatible, based on Composer 2.11.999-dev+source f83588392df43d4fd027d84d26de459946ef4e80)` + "\nx"
	if a, b := normalizeBanner(o), normalizeBanner(g); a != b || a != osc+"<APPLICATION BANNER>\nx" {
		t.Fatalf("oracle %q\ngo %q", a, b)
	}
	if a, b := normalizeBanner(logo+"Composer version 2.11.999-dev+source @release_date@"), normalizeBanner("gocomposer version 0.1.0-dev (Composer-compatible, based on Composer 2.11.999-dev+source f83588392df43d4fd027d84d26de459946ef4e80)"); a != b {
		t.Fatalf("plain %q %q", a, b)
	}
}
