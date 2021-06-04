package main

import(
	"net/http"
	"log"
	"crypto/tls"
    "bytes"
	"errors"
	

	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/go-spiffe/v2/bundle/spiffebundle"
	"github.com/spiffe/go-spiffe/v2/federation"
    "github.com/spiffe/go-spiffe/v2/logger"
	// "github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
	// "github.com/spiffe/go-spiffe/v2/workloadapi"
)

const jwks = `{
    "keys": [
        {
            "use": "x509-svid",
            "kty": "EC",
            "crv": "P-384",
            "x": "WjB-nSGSxIYiznb84xu5WGDZj80nL7W1c3zf48Why0ma7Y7mCBKzfQkrgDguI4j0",
            "y": "Z-0_tDH_r8gtOtLLrIpuMwWHoe4vbVBFte1vj6Xt6WeE8lXwcCvLs_mcmvPqVK9j",
            "x5c": [
                "MIIBzDCCAVOgAwIBAgIJAJM4DhRH0vmuMAoGCCqGSM49BAMEMB4xCzAJBgNVBAYTAlVTMQ8wDQYDVQQKDAZTUElGRkUwHhcNMTgwNTEzMTkzMzQ3WhcNMjMwNTEyMTkzMzQ3WjAeMQswCQYDVQQGEwJVUzEPMA0GA1UECgwGU1BJRkZFMHYwEAYHKoZIzj0CAQYFK4EEACIDYgAEWjB+nSGSxIYiznb84xu5WGDZj80nL7W1c3zf48Why0ma7Y7mCBKzfQkrgDguI4j0Z+0/tDH/r8gtOtLLrIpuMwWHoe4vbVBFte1vj6Xt6WeE8lXwcCvLs/mcmvPqVK9jo10wWzAdBgNVHQ4EFgQUh6XzV6LwNazA+GTEVOdu07o5yOgwDwYDVR0TAQH/BAUwAwEB/zAOBgNVHQ8BAf8EBAMCAQYwGQYDVR0RBBIwEIYOc3BpZmZlOi8vbG9jYWwwCgYIKoZIzj0EAwQDZwAwZAIwE4Me13qMC9i6Fkx0h26y09QZIbuRqA9puLg9AeeAAyo5tBzRl1YL0KNEp02VKSYJAjBdeJvqjJ9wW55OGj1JQwDFD7kWeEB6oMlwPbI/5hEY3azJi16I0uN1JSYTSWGSqWc="
            ]
        },
        {
            "use": "jwt-svid",
            "kty": "EC",
            "kid": "IRsID4VIM3T11TsK43Ny1DgCD5UNWhva",
            "crv": "P-256",
            "x": "64Mm92hnvqSdLxl6XQVu32-3rydodal1S8JgGYg5AGk",
            "y": "x5hVyA9Z4OgQxpvkhTXwNsZACk1jz1xyuTDr6JdH3R0"
        },
        {
            "use": "jwt-svid",
            "kty": "EC",
            "kid": "qjwWkiMpkHzIxsSrAsLxSZ2WZ8AyMESx",
            "crv": "P-256",
            "x": "mLo0vBg7xWrcSOEhWpSmrcoVpZRBGoDDxwNJQugFzR4",
            "y": "7Kb3afZVERcOBWHOTWwTJLTwWX4913TSxeoTU9A0hYQ"
        },
        {
            "use": "jwt-svid",
            "kty": "EC",
            "kid": "uNhqAaPI7NDn7IHOsa2ac1BF4O5qGxjZ",
            "crv": "P-256",
            "x": "pm3pJKQjBVx7x1h_dbNVCoHoXZuTwD5EAS2DjcUNsTY",
            "y": "jAe4nHFq0Jtr4ugIz500GfjjMfCfupIoWcPnwVWnpJc"
        },
        {
            "use": "jwt-svid",
            "kty": "EC",
            "kid": "y3UHKFp0WqPpG7gVr3FKieiEzwH8fTMm",
            "crv": "P-256",
            "x": "SttM6EWWCPBRYDqGKIAqVbCcelCIJE9VBqj-uX4sgwE",
            "y": "4jVOkUVM-lA9GYNU8_GX5Us5fjNR_f9Hcaj7PGkgZDo"
        },
        {
            "use": "jwt-svid",
            "kty": "EC",
            "kid": "mbrcuIaIUUapdCCmhQon4xJSicDmAVfK",
            "crv": "P-256",
            "x": "HvINid075yH0ssKsbmalal1LDceTuz9dN_RScnmLiaw",
            "y": "VBFCHK_RS4oKamHY3UdGkwbd03cC6fqEm5PI4V1crSg"
        },
        {
            "use": "jwt-svid",
            "kty": "EC",
            "kid": "cHPeHMMEtvTeSMBc20DzPPhkF41BN2WJ",
            "crv": "P-256",
            "x": "hjhyHcq6nNph9QbcSIf3VzpkVWtfOT3HbRx4aZ3j-D8",
            "y": "KrdsDLurVCbnnYmYS2Gm0rjKkXOJ9x1-SVOjqdAQEBk"
        }
    ],
    "spiffe_refresh_hint": 60,
    "spiffe_sequence": 1
}`

type fakeSource struct {
	bundles map[spiffeid.TrustDomain]*spiffebundle.Bundle
}

func (s *fakeSource) GetBundleForTrustDomain(trustDomain spiffeid.TrustDomain) (*spiffebundle.Bundle, error) {
	b, ok := s.bundles[trustDomain]
	if !ok {
		return nil, errors.New("bundle not found")
	}
	return b, nil
}

func main() {

	// generate a `Certificate` struct
	cert, _ := tls.LoadX509KeyPair( "localhost.crt", "localhost.key" )
	trustDomain, _ := spiffeid.TrustDomainFromString("localhost")
	bundle, _ := spiffebundle.Parse(trustDomain, []byte(jwks))
	source := &fakeSource{}
	source.bundles = map[spiffeid.TrustDomain]*spiffebundle.Bundle{
		trustDomain: bundle,
	}
	writer := new(bytes.Buffer)
	handler := federation.Handler(trustDomain, source, logger.Writer(writer))
	// create a custom server with `TLSConfig`
	s := &http.Server{
	  Addr: "localhost:443",
	  Handler: handler, // use `http.DefaultServeMux`
	  TLSConfig: &tls.Config{
		Certificates: []tls.Certificate{ cert },
	  },
	}

	// run server on port "443"
	log.Fatal( s.ListenAndServeTLS("", "") )

}
