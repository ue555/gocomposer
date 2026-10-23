Linux image: docker.io/library/golang:1.23 (sha256:ac2b6b754dbe064bd10a97c2f46f9253a5e8220565fc3ae4edf8b362b87bbe2f)

# CI local run 20261023T023955Z

Workflow: .github/workflows/ci.yml (sha256 bde9d6300b56ee9326321cc466478750c24fa30f7d7c5d1edd02a6eeaa97d79f)
Host: Darwin 27.0.0 arm64
GitHub Actions: not run (requires push; out of scope)

| job | result | log |
|---|---|---|
| unit (macos, go1.23.12) | PASS | unit-macos.log |
| unit (linux, docker.io/library/golang:1.23) | PASS | unit-linux.log |
| oracle: prepare | PASS | oracle-prepare.log |
| oracle: compare-8.4 | FAIL (exit 1) | oracle-compare-8.4.log |
| oracle: compare-8.3 | PASS | oracle-compare-8.3.log |
| oracle: compare-8.5 | PASS | oracle-compare-8.5.log |
| oracle: integration | PASS | oracle-integration.log |
