package app

import (
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"

	config "github.com/BIQ-Cat/easyserver/config/base"
)

const toLower = 'a' - 'A'

type CORS struct {
	// AllowedOrigins is a list of origins a cross-domain request can be executed from.
	// If the special "*" value is present in the list, all origins will be allowed.
	// An origin may contain wildcards: * to replace 0 or more characters
	// (i.e. http://*.domain.* is http://www.domain.com or http://us.domain.co.uk)
	// or ? to match any character
	// (i.e. http://u?.sourc?.com is http://us.sourcy.com or http://uk.source.com).
	// Usage of wildcards implies a small performance penalty.
	// Default value is ["*"]
	AllowedOrigins []string

	// AllowOriginFunc is a custom function to validate the origin. It takes the origin
	// as argument and returns true if allowed or false otherwise. If this option is
	// set, the content of AllowedOrigins is ignored.
	AllowOriginFunc func(r *http.Request, origin string) bool

	// AllowedMethods is a list of methods the client is allowed to use with
	// cross-domain requests. Default value is simple methods (HEAD, GET and POST).
	AllowedMethods []string

	// AllowedHeaders is list of non simple headers the client is allowed to use with
	// cross-domain requests.
	// If the special "*" value is present in the list, all headers will be allowed.
	// Default value is [] but "Origin" is always appended to the list.
	AllowedHeaders []string

	// ExposedHeaders indicates which headers are safe to expose to the API of a CORS
	// API specification
	ExposedHeaders []string

	// AllowCredentials indicates whether the request can include user credentials like
	// cookies, HTTP authentication or client side SSL certificates.
	AllowCredentails bool

	// MaxAge indicates how long (in seconds) the results of a preflight request
	// can be cached
	MaxAge int

	// OptionsPassthrough instructs preflight to let other potential next handlers to
	// process the OPTIONS method. Turn this on if your application handles OPTIONS.
	OptionsPassthrough bool
}

func AllowAll() *CORS {
	return &CORS{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodPatch,
			http.MethodHead,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentails: false,
	}
}

func (c *CORS) handlePreflight(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	origin := r.Header.Get("Origin")

	if r.Method != http.MethodOptions {
		if config.Config.Debug {
			log.Printf("Preflight aborted: %s!=OPTIONS", r.Method)
		}

		return
	}
	// Always set Vary headers
	// see https://github.com/rs/cors/issues/10,
	//     https://github.com/rs/cors/commit/dbdca4d95feaa7511a46e6f1efb3b3aa505bc43f#commitcomment-12352001
	headers.Add("Vary", "Origin")
	headers.Add("Vary", "Access-Control-Request-Method")
	headers.Add("Vary", "Access-Control-Request-Headers")

	if origin == "" {
		if config.Config.Debug {
			log.Printf("Preflight aborted: empty origin")
		}

		return
	}

	if !c.isOriginAllowed(r, origin) {
		if config.Config.Debug {
			log.Printf("Preflight aborted: origin '%s' not allowed", origin)
		}

		return
	}

	reqMethod := r.Header.Get("Access-Control-Request-Method")
	if !c.isMethodAllowed(reqMethod) {
		if config.Config.Debug {
			log.Printf("Preflight aborted: method '%s' not allowed", reqMethod)
		}

		return
	}

	reqHeaders := parseAndFormatHeaderList(r.Header.Get("Access-Control-Request-Headers"))
	if !c.areHeadersAllowed(reqHeaders) {
		if config.Config.Debug {
			log.Printf("Preflight aborted: method '%s' not allowed", reqMethod)
		}

		return
	}

	if slices.Contains(c.AllowedOrigins, "*") {
		headers.Set("Access-Control-Allow-Origin", "*")
	} else {
		headers.Set("Access-Control-Allow-Origin", origin)
	}

	// Spec says: Since the list of methods can be unbounded, simply returning the method indicated
	// by Access-Control-Request-Method (if supported) can be enough
	headers.Set("Access-Control-Allow-Methods", strings.ToUpper(reqMethod))

	if c.AllowCredentails {
		headers.Set("Access-Control-Allow-Credentials", "true")
	}

	if len(reqHeaders) > 0 {

		// Spec says: Since the list of headers can be unbounded, simply returning supported headers
		// from Access-Control-Request-Headers can be enough
		headers.Set("Access-Control-Allow-Headers", strings.Join(reqHeaders, ", "))
	}

	if c.MaxAge > 0 {
		headers.Set("Access-Control-Max-Age", strconv.Itoa(c.MaxAge))
	}
}

func (c *CORS) handleActualRequest(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	origin := r.Header.Get("Origin")

	// Always set Vary, see https://github.com/rs/cors/issues/10
	headers.Add("Vary", "Origin")
	if origin == "" {
		return
	}
	if !c.isOriginAllowed(r, origin) {
		return
	}

	// Note that spec does define a way to specifically disallow a simple method like GET or
	// POST. Access-Control-Allow-Methods is only used for pre-flight requests and the
	// spec doesn't instruct to check the allowed methods for simple cross-origin requests.
	// We think it's a nice feature to be able to have control on those methods though.
	if !c.isMethodAllowed(r.Method) {
		return
	}
	if slices.Contains(c.AllowedOrigins, "*") {
		headers.Set("Access-Control-Allow-Origin", "*")
	} else {
		headers.Set("Access-Control-Allow-Origin", origin)
	}
	if len(c.ExposedHeaders) > 0 {
		headers.Set("Access-Control-Expose-Headers", strings.Join(c.ExposedHeaders, ", "))
	}
	if c.AllowCredentails {
		headers.Set("Access-Control-Allow-Credentials", "true")
	}
}

func (c *CORS) isOriginAllowed(r *http.Request, origin string) bool {
	if c.AllowOriginFunc != nil {
		return c.AllowOriginFunc(r, origin)
	}

	if c.AllowedOrigins == nil {
		c.AllowedOrigins = []string{"*"}
	}

	for _, orig := range c.AllowedOrigins {
		if isMatch(origin, orig) {
			return true
		}
	}

	return false
}

func (c *CORS) isMethodAllowed(method string) bool {
	if c.AllowedMethods == nil {
		c.AllowedMethods = []string{http.MethodGet, http.MethodPost, http.MethodHead}
	}

	if len(c.AllowedMethods) == 0 {
		// Disabled (even for preflight!)
		return false
	}

	method = strings.ToUpper(method)
	if method == http.MethodOptions {
		// Preflight
		return true
	}

	c.AllowedMethods = convert(c.AllowedMethods, strings.ToUpper)
	if slices.Contains(c.AllowedMethods, "*") || slices.Contains(c.AllowedMethods, method) {
		return true
	}

	return false
}

func (c *CORS) areHeadersAllowed(requestHeaders []string) bool {
	if c.AllowedHeaders == nil {
		c.AllowedHeaders = []string{"Origin", "Accept", "Content-Type"}
	}
	if slices.Contains(c.AllowedHeaders, "*") || len(requestHeaders) == 0 {
		return true
	}
	for _, header := range requestHeaders {
		header = http.CanonicalHeaderKey(header)

		if !slices.Contains(convert(c.AllowedHeaders, http.CanonicalHeaderKey), header) {
			return false
		}
	}
	return true
}

func isMatch(orig string, patt string) bool {
	runeInput := []rune(orig)
	runePattern := []rune(patt)

	lenInput := len(runeInput)
	lenPattern := len(runePattern)

	matchMatrix := make([][]bool, lenInput+1)

	for i := range matchMatrix {
		matchMatrix[i] = make([]bool, lenPattern+1)
	}

	matchMatrix[0][0] = true

	for i := 1; i <= lenInput; i++ {
		matchMatrix[i][0] = false
	}

	if lenPattern > 0 {
		if runePattern[0] == '*' {
			matchMatrix[0][1] = true
		}
	}

	for j := 2; j <= lenPattern; j++ {
		if runePattern[j-1] == '*' {
			matchMatrix[0][j] = matchMatrix[0][j-1]
		}
	}

	for i := 1; i <= lenInput; i++ {
		for j := 1; j <= lenPattern; j++ {
			if runePattern[j-1] == '*' {
				matchMatrix[i][j] = matchMatrix[i-1][j] || matchMatrix[i][j-1]
			}

			if runePattern[j-1] == '?' || runeInput[i-1] == runePattern[j-1] {
				matchMatrix[i][j] = matchMatrix[i-1][j-1]
			}
		}
	}

	return matchMatrix[lenInput][lenPattern]
}

func convert(s []string, c func(string) string) []string {
	out := []string{}
	for _, i := range s {
		out = append(out, c(i))
	}
	return out
}

func parseAndFormatHeaderList(headerList string) []string {
	headersLength := len(headerList)
	h := make([]byte, 0, headersLength)
	upper := true

	headers := make([]string, strings.Count(headerList, ","))

	for i := 0; i < headersLength; i++ {
		b := headerList[i]
		if b >= 'a' && b <= 'z' {
			if upper {
				h = append(h, b-toLower)
			} else {
				h = append(h, b)
			}
		} else if b >= 'A' && b <= 'Z' {
			if !upper {
				h = append(h, b+toLower)
			} else {
				h = append(h, b)
			}
		} else if b == '-' || b == '_' || b == '.' || (b >= '0' && b <= '9') {
			h = append(h, b)
		}

		if b == ' ' || b == ',' || i == headersLength-1 {
			if len(h) > 0 {
				// Flush the found header
				headers = append(headers, string(h))
				h = h[:0]
				upper = true
			}
		} else {
			upper = b == '-'
		}
	}

	return headers
}

func (c *CORS) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
			c.handlePreflight(w, r)
			// Preflight requests are standalone and should stop the chain as some other
			// middleware may not handle OPTIONS requests correctly. One typical example
			// is authentication middleware ; OPTIONS requests won't carry authentication
			// headers (see #1)
			if c.OptionsPassthrough {
				next.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusOK)
			}
		} else {
			c.handleActualRequest(w, r)
			next.ServeHTTP(w, r)
		}
	})
}
