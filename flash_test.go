// Copyright 2015 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flash

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/lunny/tango"
	"github.com/tango-contrib/session"
)

type FlashAction struct {
	Flash
}

func (a *FlashAction) Get() string {
	a.Flash.Set("test", "test")
	return "test"
}

func (a *FlashAction) Post() string {
	s := a.Flash.Get("test")
	if s != nil {
		return s.(string)
	}
	return ""
}

func (a *FlashAction) Put() string {
	s := a.Flash.Get("test")
	if s != nil {
		return s.(string)
	}
	return ""
}

func TestFlash(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	tg := tango.Classic()
	sessions := session.New()
	tg.Use(Flashes(sessions))
	tg.Any("/", new(FlashAction))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	tg.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "test")

	buff.Reset()

	req, err = http.NewRequest("POST", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	cks := readSetCookies(recorder.Header())
	for _, ck := range cks {
		req.AddCookie(ck)
	}
	recorder = httptest.NewRecorder()
	recorder.Body = buff

	tg.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "test")

	buff.Reset()

	req, err = http.NewRequest("PUT", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	cks = readSetCookies(recorder.Header())
	for _, ck := range cks {
		req.AddCookie(ck)
	}
	recorder = httptest.NewRecorder()
	recorder.Body = buff

	tg.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	expect(t, len(buff.String()), 0)
	expect(t, buff.String(), "")
}

/*func TestFlashRedis(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()
	recorder.Body = buff

	tg := tango.Classic()
	sessions := session.New(session.Options{
		Store: redistore.New(redistore.Options{
			Host:    "127.0.0.1",
			DbIndex: 0,
			MaxAge:  30 * time.Minute,
		}),
	})
	tg.Use(Flashes(sessions))
	tg.Any("/", new(FlashAction))

	req, err := http.NewRequest("GET", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	tg.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "test")

	buff.Reset()

	req, err = http.NewRequest("POST", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	cks := readSetCookies(recorder.Header())
	for _, ck := range cks {
		req.AddCookie(ck)
	}
	recorder.HeaderMap = make(map[string][]string)

	tg.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	refute(t, len(buff.String()), 0)
	expect(t, buff.String(), "test")

	buff.Reset()

	req, err = http.NewRequest("PUT", "http://localhost:8000/", nil)
	if err != nil {
		t.Error(err)
	}

	cks = readSetCookies(recorder.Header())
	for _, ck := range cks {
		req.AddCookie(ck)
	}
	recorder.HeaderMap = make(map[string][]string)

	tg.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusOK)
	expect(t, len(buff.String()), 0)
	expect(t, buff.String(), "")
}*/

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func parseCookieValue(raw string, allowDoubleQuote bool) (string, bool) {
	// Strip the quotes, if present.
	if allowDoubleQuote && len(raw) > 1 && raw[0] == '"' && raw[len(raw)-1] == '"' {
		raw = raw[1 : len(raw)-1]
	}
	for i := 0; i < len(raw); i++ {
		if !validCookieValueByte(raw[i]) {
			return "", false
		}
	}
	return raw, true
}

var isTokenTable = [127]bool{
	'!':  true,
	'#':  true,
	'$':  true,
	'%':  true,
	'&':  true,
	'\'': true,
	'*':  true,
	'+':  true,
	'-':  true,
	'.':  true,
	'0':  true,
	'1':  true,
	'2':  true,
	'3':  true,
	'4':  true,
	'5':  true,
	'6':  true,
	'7':  true,
	'8':  true,
	'9':  true,
	'A':  true,
	'B':  true,
	'C':  true,
	'D':  true,
	'E':  true,
	'F':  true,
	'G':  true,
	'H':  true,
	'I':  true,
	'J':  true,
	'K':  true,
	'L':  true,
	'M':  true,
	'N':  true,
	'O':  true,
	'P':  true,
	'Q':  true,
	'R':  true,
	'S':  true,
	'T':  true,
	'U':  true,
	'W':  true,
	'V':  true,
	'X':  true,
	'Y':  true,
	'Z':  true,
	'^':  true,
	'_':  true,
	'`':  true,
	'a':  true,
	'b':  true,
	'c':  true,
	'd':  true,
	'e':  true,
	'f':  true,
	'g':  true,
	'h':  true,
	'i':  true,
	'j':  true,
	'k':  true,
	'l':  true,
	'm':  true,
	'n':  true,
	'o':  true,
	'p':  true,
	'q':  true,
	'r':  true,
	's':  true,
	't':  true,
	'u':  true,
	'v':  true,
	'w':  true,
	'x':  true,
	'y':  true,
	'z':  true,
	'|':  true,
	'~':  true,
}

func isToken(r rune) bool {
	i := int(r)
	return i < len(isTokenTable) && isTokenTable[i]
}

func isNotToken(r rune) bool {
	return !isToken(r)
}

func validCookieValueByte(b byte) bool {
	return 0x20 <= b && b < 0x7f && b != '"' && b != ';' && b != '\\'
}

func isCookieNameValid(raw string) bool {
	return strings.IndexFunc(raw, isNotToken) < 0
}

// readSetCookies parses all "Set-Cookie" values from
// the header h and returns the successfully parsed Cookies.
func readSetCookies(h http.Header) []*http.Cookie {
	cookies := []*http.Cookie{}
	for _, line := range h["Set-Cookie"] {
		parts := strings.Split(strings.TrimSpace(line), ";")
		if len(parts) == 1 && parts[0] == "" {
			continue
		}
		parts[0] = strings.TrimSpace(parts[0])
		j := strings.Index(parts[0], "=")
		if j < 0 {
			continue
		}
		name, value := parts[0][:j], parts[0][j+1:]
		if !isCookieNameValid(name) {
			continue
		}
		value, success := parseCookieValue(value, true)
		if !success {
			continue
		}
		c := &http.Cookie{
			Name:  name,
			Value: value,
			Raw:   line,
		}
		for i := 1; i < len(parts); i++ {
			parts[i] = strings.TrimSpace(parts[i])
			if len(parts[i]) == 0 {
				continue
			}

			attr, val := parts[i], ""
			if j := strings.Index(attr, "="); j >= 0 {
				attr, val = attr[:j], attr[j+1:]
			}
			lowerAttr := strings.ToLower(attr)
			val, success = parseCookieValue(val, false)
			if !success {
				c.Unparsed = append(c.Unparsed, parts[i])
				continue
			}
			switch lowerAttr {
			case "secure":
				c.Secure = true
				continue
			case "httponly":
				c.HttpOnly = true
				continue
			case "domain":
				c.Domain = val
				continue
			case "max-age":
				secs, err := strconv.Atoi(val)
				if err != nil || secs != 0 && val[0] == '0' {
					break
				}
				if secs <= 0 {
					c.MaxAge = -1
				} else {
					c.MaxAge = secs
				}
				continue
			case "expires":
				c.RawExpires = val
				exptime, err := time.Parse(time.RFC1123, val)
				if err != nil {
					exptime, err = time.Parse("Mon, 02-Jan-2006 15:04:05 MST", val)
					if err != nil {
						c.Expires = time.Time{}
						break
					}
				}
				c.Expires = exptime.UTC()
				continue
			case "path":
				c.Path = val
				continue
			}
			c.Unparsed = append(c.Unparsed, parts[i])
		}
		cookies = append(cookies, c)
	}
	return cookies
}
