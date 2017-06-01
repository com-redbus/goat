package goat

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"log"

	"github.com/stretchr/testify/assert"
)

func Test_CSP(t *testing.T) {
	h := &TestNoCacheHandler{}
	csp := NewCSP(CSPOptions{
		DefaultSrc:     []string{"'self'", "s1.rdbuz.com"},
		ScriptSrc:      []string{"'self'"},
		StyleSrc:       []string{"'self'"},
		ImgSrc:         []string{"'self'"},
		ConnectSrc:     []string{"'self'"},
		FontSrc:        []string{"'self'"},
		ObjectSrc:      []string{"'self'"},
		MediaSrc:       []string{"'self'"},
		ChildSrc:       []string{"'self'"},
		FormAction:     []string{"'self'"},
		FrameAncestors: []string{"'none'"},
		PluginTypes:    []string{"application/pdf"},
		Sandbox:        []string{"allow-forms", "allow-scripts"},
		ReportURI:      "/some-dummy-report-api",
		IsReportOnly:   true,
	})
	server := httptest.NewServer(csp.CSP(h))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	if !csp.cspOptions.IsReportOnly {
		log.Println(resp.Header.Get("Content-Security-Policy"))
		assert.Contains(t, resp.Header.Get("Content-Security-Policy"), "default-src 'self' s1.rdbuz.com", "CSP not working")
		assert.Contains(t, resp.Header.Get("Content-Security-Policy"), "script-src 'self'", "CSP not working")
		assert.Contains(t, resp.Header.Get("Content-Security-Policy"), "style-src 'self'", "CSP not working")
	} else {
		log.Println(resp.Header.Get("Content-Security-Policy-Report-Only"))
		assert.Contains(t, resp.Header.Get("Content-Security-Policy-Report-Only"), "report-uri /some-dummy-report-api", "CSP Report Only not working")
	}
}
