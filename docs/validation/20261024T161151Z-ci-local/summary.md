# CI local run 20261024T161151Z

Workflow: .github/workflows/ci.yml (sha256 1c45e55e145d770080fdff3f660f9c3d0fab90c02687151f6b14acbae7d21dcb)
Host: Darwin 27.0.0 arm64
Linux image: docker.io/library/golang:1.23 (sha256:ac2b6b754dbe064bd10a97c2f46f9253a5e8220565fc3ae4edf8b362b87bbe2f)
GitHub Actions: not run (requires push; out of scope)

| job | result | log |
|---|---|---|
| unit (macos, go1.23.12) | PASS | unit-macos.log |
| unit (linux, docker.io/library/golang:1.23) | PASS | unit-linux.log |
| oracle: prepare-8.4 | PASS | oracle-prepare-8.4.log |
| oracle: compare-8.4 | PASS | oracle-compare-8.4.log |
| oracle: prepare-8.3 | PASS | oracle-prepare-8.3.log |
| oracle: compare-8.3 | PASS | oracle-compare-8.3.log |
| oracle: prepare-8.5 | PASS | oracle-prepare-8.5.log |
| oracle: compare-8.5 | PASS | oracle-compare-8.5.log |
| oracle: integration-prepare | PASS | oracle-integration-prepare.log |
| oracle: integration | PASS | oracle-integration.log |
