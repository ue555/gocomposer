package composer

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"gocomposer/internal/phpjson"
	"gocomposer/internal/solver"
)

func TestMaxParallelHTTP(t *testing.T) {
	for _, tc := range []struct {
		env  string
		want int
	}{{"", 12}, {"1", 1}, {"0", 1}, {"3.7", 3}, {"99", 50}, {"x", 12}, {" 4", 4}} {
		t.Setenv("COMPOSER_MAX_PARALLEL_HTTP", tc.env)
		if got := maxParallelHTTP(); got != tc.want {
			t.Errorf("COMPOSER_MAX_PARALLEL_HTTP=%q: got %d, want %d", tc.env, got, tc.want)
		}
	}
}

func TestFilterSpecRoundTrip(t *testing.T) {
	pm := &PluginManager{}
	if _, ok := pm.filterFromSpec(goFilterSpec(solver.IgnoreAll{})).(solver.IgnoreAll); !ok {
		t.Error("IgnoreAll does not round trip")
	}
	if _, ok := pm.filterFromSpec(goFilterSpec(nil)).(solver.IgnoreNothing); !ok {
		t.Error("nil filter is not IgnoreNothing")
	}
	list := pm.filterFromSpec(goFilterSpec(solver.NewIgnoreList([]string{"ext-foo", "php+"})))
	if !list.IsIgnored("ext-foo") || list.IsIgnored("ext-bar") || list.IsIgnored("php") || !list.IsUpperBoundIgnored("php") {
		t.Error("IgnoreList does not round trip")
	}
	if _, ok := pm.filterFromSpec(true).(solver.IgnoreAll); !ok {
		t.Error("true is not IgnoreAll")
	}
}

func selfSignedCert(t *testing.T) (*x509.Certificate, []byte, *rsa.PrivateKey) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	tpl := &x509.Certificate{SerialNumber: big.NewInt(1), Subject: pkix.Name{CommonName: "plugfix3"}, NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(time.Hour)}
	der, err := x509.CreateCertificate(rand.Reader, tpl, tpl, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}
	return cert, der, key
}

func TestFingerprintMatches(t *testing.T) {
	cert, der, _ := selfSignedCert(t)
	sum := sha1.Sum(der)
	fp := hex.EncodeToString(sum[:])
	if !fingerprintMatches(cert, fp) {
		t.Error("sha1 fingerprint does not match")
	}
	m := phpjson.NewMap()
	m.Set("sha1", fp)
	if !fingerprintMatches(cert, m) {
		t.Error("sha1 fingerprint array does not match")
	}
	if fingerprintMatches(cert, "0000000000000000000000000000000000000000") {
		t.Error("wrong fingerprint matches")
	}
}

func TestLoadClientCertEncryptedKey(t *testing.T) {
	_, der, key := selfSignedCert(t)
	dir := t.TempDir()
	certFile := filepath.Join(dir, "client.pem")
	if err := os.WriteFile(certFile, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}), 0o600); err != nil {
		t.Fatal(err)
	}
	//nolint:staticcheck // legacy PEM encryption, like `openssl rsa -aes128 -traditional`
	block, err := x509.EncryptPEMBlock(rand.Reader, "RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(key), []byte("secret"), x509.PEMCipherAES128)
	if err != nil {
		t.Fatal(err)
	}
	keyFile := filepath.Join(dir, "client-enc.key")
	if err := os.WriteFile(keyFile, pem.EncodeToMemory(block), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := loadClientCert(certFile, keyFile, "secret"); err != nil {
		t.Fatalf("encrypted key: %v", err)
	}
	if _, err := loadClientCert(certFile, keyFile, "wrong"); err == nil {
		t.Fatal("wrong passphrase accepted")
	}
}

func TestValidateWithSchemaFile(t *testing.T) {
	dir := t.TempDir()
	schema := filepath.Join(dir, "schema.json")
	if err := os.WriteFile(schema, []byte(`{"type": "object", "required": ["plug3"], "properties": {"name": {"type": "integer"}}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	data, err := phpjson.Decode([]byte(`{"name": "x"}`))
	if err != nil {
		t.Fatal(err)
	}
	v, err := validateWithSchemaFile(data, schema)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	phpjson.Each(v, func(_ string, x phpjson.Value) {
		s, _ := x.(string)
		got = append(got, s)
	})
	want := []string{"plug3 : The property plug3 is required", "name : String value found, but an integer is required"}
	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("got %q, want %q", got, want)
	}
}
