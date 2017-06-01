package goat

import (
	"bytes"
	"html"
	"html/template"
	"net/http"
	"strings"
)

//CSPOptions struct for getting all csp options , more info @ https://content-security-policy.com/
type CSPOptions struct {
	DefaultSrc      []string //The default-src is the default policy for loading content such as JavaScript, Images, CSS, Font's, AJAX requests, Frames, HTML5 Media
	ScriptSrc       []string //Defines valid sources of JavaScript.
	StyleSrc        []string //Defines valid sources of stylesheets.
	ImgSrc          []string //Defines valid sources of images.
	ConnectSrc      []string //Applies to XMLHttpRequest (AJAX), WebSocket or EventSource. If not allowed the browser emulates a 400 HTTP status code.
	FontSrc         []string //	Defines valid sources of fonts.
	ObjectSrc       []string //Defines valid sources of plugins, eg <object>, <embed> or <applet>.
	MediaSrc        []string //Defines valid sources of audio and video, eg HTML5 <audio>, <video> elements.
	ChildSrc        []string //Defines valid sources for web workers and nested browsing contexts loaded using elements such as <frame> and <iframe>
	Sandbox         []string //Enables a sandbox for the requested resource similar to the iframe sandbox attribute. The sandbox applies a same origin policy, prevents popups, plugins and script execution is blocked. You can keep the sandbox value empty to keep all restrictions in place, or add values: allow-forms allow-same-origin allow-scripts allow-popups, allow-modals, allow-orientation-lock, allow-pointer-lock, allow-presentation, allow-popups-to-escape-sandbox, and allow-top-navigation
	ReportURI       string   //Instructs the browser to POST reports of policy failures to this URI. You can also append -Report-Only to the HTTP header name to instruct the browser to only send reports (does not block anything).
	FormAction      []string //Defines valid sources that can be used as a HTML <form> action.
	FrameAncestors  []string //Defines valid sources for embedding the resource using <frame> <iframe> <object> <embed> <applet>. Setting this directive to 'none' should be roughly equivalent to X-Frame-Options: DENY
	PluginTypes     []string //Defines valid MIME types for plugins invoked via <object> and <embed>. To load an <applet> you must specify application/x-java-applet.
	IsHeaderCreated bool
	HeaderString    string
	HeaderStrings   []string
	IsReportOnly    bool //send  Content-Security-Policy-Report-Only header
}

type CSPHandler struct {
	cspOptions CSPOptions
}

// var m = map[string]string{
// 	"defaultSrc": defaultTemplate,
// 	"scriptSrc":  scriptTemplate,
// 	"styleSrc":   styleTemplate,
// }

func getTemplate(name string) string {
	return name + " {{.}}"
}

func createStringFromValues(name string, values []string) string {
	t := template.New(name)
	tem := template.Must(t.Parse(getTemplate(name)))
	s := strings.Join(values, " ")

	buf := &bytes.Buffer{}
	tem.Execute(buf, s)

	str := html.UnescapeString(buf.String())
	return str
}

func createStringFromValue(name string, value string) string {
	t := template.New(name)
	tem := template.Must(t.Parse(getTemplate(name)))

	buf := &bytes.Buffer{}
	tem.Execute(buf, value)

	str := html.UnescapeString(buf.String())
	return str
}

func NewCSP(cspOptions CSPOptions) *CSPHandler {
	csp := &CSPHandler{
		cspOptions: cspOptions,
	}
	return csp
}

func (csp *CSPHandler) CSP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isStringAvailable := false
		if csp.cspOptions.IsHeaderCreated {
			if csp.cspOptions.IsReportOnly && csp.cspOptions.ReportURI != "" {
				w.Header().Set("Content-Security-Policy-Report-Only", csp.cspOptions.HeaderString)
			} else {
				w.Header().Set("Content-Security-Policy", csp.cspOptions.HeaderString)
			}
			next.ServeHTTP(w, r)
			return
		}

		if len(csp.cspOptions.DefaultSrc) != 0 {
			isStringAvailable = true
			csp.cspOptions.HeaderStrings = append(csp.cspOptions.HeaderStrings, createStringFromValues("default-src", csp.cspOptions.DefaultSrc))
		}

		if len(csp.cspOptions.ScriptSrc) != 0 {
			isStringAvailable = true
			csp.cspOptions.HeaderStrings = append(csp.cspOptions.HeaderStrings, createStringFromValues("script-src", csp.cspOptions.ScriptSrc))
		}

		if len(csp.cspOptions.StyleSrc) != 0 {
			isStringAvailable = true
			csp.cspOptions.HeaderStrings = append(csp.cspOptions.HeaderStrings, createStringFromValues("style-src", csp.cspOptions.StyleSrc))
		}

		if len(csp.cspOptions.ImgSrc) != 0 {
			isStringAvailable = true
			csp.cspOptions.HeaderStrings = append(csp.cspOptions.HeaderStrings, createStringFromValues("img-src", csp.cspOptions.ImgSrc))
		}

		if len(csp.cspOptions.ConnectSrc) != 0 {
			isStringAvailable = true
			csp.cspOptions.HeaderStrings = append(csp.cspOptions.HeaderStrings, createStringFromValues("connect-src", csp.cspOptions.ConnectSrc))
		}
		if len(csp.cspOptions.FontSrc) != 0 {
			isStringAvailable = true
			csp.cspOptions.HeaderStrings = append(csp.cspOptions.HeaderStrings, createStringFromValues("font-src", csp.cspOptions.FontSrc))
		}

		if len(csp.cspOptions.ObjectSrc) != 0 {
			isStringAvailable = true
			csp.cspOptions.HeaderStrings = append(csp.cspOptions.HeaderStrings, createStringFromValues("object-src", csp.cspOptions.ObjectSrc))
		}

		if len(csp.cspOptions.MediaSrc) != 0 {
			isStringAvailable = true
			csp.cspOptions.HeaderStrings = append(csp.cspOptions.HeaderStrings, createStringFromValues("media-src", csp.cspOptions.MediaSrc))
		}

		if len(csp.cspOptions.ChildSrc) != 0 {
			isStringAvailable = true
			csp.cspOptions.HeaderStrings = append(csp.cspOptions.HeaderStrings, createStringFromValues("child-src", csp.cspOptions.ChildSrc))
		}

		if len(csp.cspOptions.FormAction) != 0 {
			isStringAvailable = true
			csp.cspOptions.HeaderStrings = append(csp.cspOptions.HeaderStrings, createStringFromValues("form-action", csp.cspOptions.FormAction))
		}

		if len(csp.cspOptions.FrameAncestors) != 0 {
			isStringAvailable = true
			csp.cspOptions.HeaderStrings = append(csp.cspOptions.HeaderStrings, createStringFromValues("frame-ancestors", csp.cspOptions.FrameAncestors))
		}

		if len(csp.cspOptions.PluginTypes) != 0 {
			isStringAvailable = true
			csp.cspOptions.HeaderStrings = append(csp.cspOptions.HeaderStrings, createStringFromValues("plugin-types", csp.cspOptions.PluginTypes))
		}

		if len(csp.cspOptions.Sandbox) != 0 {
			isStringAvailable = true
			csp.cspOptions.HeaderStrings = append(csp.cspOptions.HeaderStrings, createStringFromValues("sandbox", csp.cspOptions.Sandbox))
		}
		if csp.cspOptions.ReportURI != "" {
			isStringAvailable = true
			csp.cspOptions.HeaderStrings = append(csp.cspOptions.HeaderStrings, createStringFromValue("report-uri", csp.cspOptions.ReportURI))
		}

		if csp.cspOptions.IsReportOnly && csp.cspOptions.ReportURI != "" {
			csp.cspOptions.HeaderString = strings.Join(csp.cspOptions.HeaderStrings, "; ")
			csp.cspOptions.IsHeaderCreated = true

			w.Header().Set("Content-Security-Policy-Report-Only", csp.cspOptions.HeaderString)
			next.ServeHTTP(w, r)
			return
		}

		if isStringAvailable {
			csp.cspOptions.HeaderString = strings.Join(csp.cspOptions.HeaderStrings, "; ")
			csp.cspOptions.IsHeaderCreated = true

			w.Header().Set("Content-Security-Policy", csp.cspOptions.HeaderString)
		}
		next.ServeHTTP(w, r)
	})

}
