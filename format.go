package mapps

import (
	"bytes"
	"encoding/xml"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	Br = "<br/>"

	Ussd     = "ussd"
	Telegram = "telegram"
	Wap      = "wap"

	XmlPush = "xmlpush"
)

func formatElement(element string, s string) string {
	return "<" + element + ">" + s + "</" + element + ">"
}

func Bold(s ...string) string {
	return formatElement("b", strings.Join(s, " "))
}

func Italic(s ...string) string {
	return formatElement("i", strings.Join(s, " "))
}

func Code(s ...string) string {
	return formatElement("code", strings.Join(s, " "))
}

func Href(href string, s ...string) string {
	return "<a href=\"" + href + "\">" + strings.Join(s, " ") + "</a>"
}

func Data(s string) string {
	return "<![CDATA[" + s + "]]>"
}

var (
	esc_quot = []byte("&#34;") // shorter than "&quot;"
	esc_apos = []byte("&#39;") // shorter than "&apos;"
	esc_amp  = []byte("&amp;")
	esc_lt   = []byte("&lt;")
	esc_gt   = []byte("&gt;")
	esc_tab  = []byte("&#x9;")
	esc_nl   = []byte("&#xA;")
	esc_cr   = []byte("&#xD;")
	esc_fffd = []byte("\uFFFD") // Unicode replacement character
)

func isInCharacterRange(r rune) (inrange bool) {
	return r == 0x09 ||
		r == 0x0A ||
		r == 0x0D ||
		r >= 0x20 && r <= 0xDF77 ||
		r >= 0xE000 && r <= 0xFFFD ||
		r >= 0x10000 && r <= 0x10FFFF
}

func EscapeString(s string) (result string) {
	var buffer bytes.Buffer
	var esc []byte
	last := 0
	for i := 0; i < len(s); {
		r, width := utf8.DecodeRuneInString(s[i:])
		i += width

		switch r {
		case '"':
			esc = esc_quot
		case '\'':
			esc = esc_apos
		case '&':
			esc = esc_amp
		case '<':
			esc = esc_lt
		case '>':
			esc = esc_gt
		case '\t':
			esc = esc_tab
		case '\n':
			esc = esc_nl
		case '\r':
			esc = esc_cr
		default:
			if !isInCharacterRange(r) || (r == 0xFFFD && width == 1) {
				esc = esc_fffd
				break
			}
			continue
		}
		buffer.WriteString(s[last : i-width])
		buffer.Write(esc)
		last = i
	}

	buffer.WriteString(s[last:])
	return buffer.String()
}

func EscapeStringNotBr(s string) (result string) {
	var buffer bytes.Buffer
	var esc []byte
	last := 0
	for i := 0; i < len(s); {
		r, width := utf8.DecodeRuneInString(s[i:])
		i += width

		switch r {
		case '"':
			esc = esc_quot
		case '\'':
			esc = esc_apos
		case '&':
			esc = esc_amp
		case '<':
			esc = esc_lt
		case '>':
			esc = esc_gt
		case '\t':
			esc = esc_tab
		case '\r':
			esc = esc_cr
		default:
			if !isInCharacterRange(r) || (r == 0xFFFD && width == 1) {
				esc = esc_fffd
				break
			}
			continue
		}
		buffer.WriteString(s[last : i-width])
		buffer.Write(esc)
		last = i
	}

	buffer.WriteString(s[last:])
	return buffer.String()
}

func EscapeStringYesBr(s string) (result string) {
	var buffer bytes.Buffer
	var esc []byte
	last := 0
	for i := 0; i < len(s); {
		r, width := utf8.DecodeRuneInString(s[i:])
		i += width

		switch r {
		case '"':
			esc = esc_quot
		case '\'':
			esc = esc_apos
		case '&':
			esc = esc_amp
		case '<':
			esc = esc_lt
		case '>':
			esc = esc_gt
		case '\t':
			esc = esc_tab
		case '\n':
			esc = []byte("<br/>")
		case '\r':
			esc = esc_cr
		default:
			if !isInCharacterRange(r) || (r == 0xFFFD && width == 1) {
				esc = esc_fffd
				break
			}
			continue
		}
		buffer.WriteString(s[last : i-width])
		buffer.Write(esc)
		last = i
	}

	buffer.WriteString(s[last:])
	return buffer.String()
}

func formatConstruct(cons string, args string, text string) string {
	return "<" + join(cons, args) + ">" + text + "</" + cons + ">"
}

func OnlyTextDefault(text string) string {
	return Page("",
		Div("",
			EscapeString(text),
		),
	)
}

func OnlyTextTelegram(text string) string {
	return Page("",
		Div("",
			EscapeString(EscapeString(text)),
		),
	)
}

func Page(args string, text ...string) string {
	return xml.Header + formatConstruct("page", join("version=\"2.0\"", args), strings.Join(text, " "))
}

func Div(args string, text ...string) string {
	return formatConstruct("div", args, strings.Join(text, " "))
}

func Title(args string, text ...string) string {
	return formatConstruct("title", args, strings.Join(text, " "))
}

func Input(navigationId string, fieldName string, title string) string {
	return "<input " +
		FormatAttr("navigationId", navigationId) + " " +
		FormatAttr("name", fieldName) + " " +
		FormatAttr("title", title) + "/>"
}

func FormatAttr(key string, val string) string {
	return key + "=\"" + val + "\""
}

func AttrProtocol(s string) string {
	return FormatAttr("protocol", s)
}

func Attributes(s ...string) string {
	return FormatAttr("attributes", strings.Join(s, "; "))
}

//telegram.links.realignment.threshold
func TelegramLinksRealignmentThreshold(count int) string {
	return AttrAgr("telegram.links.realignment.threshold", strconv.Itoa(count))
}

//telegram.links.realignment.enabled
func TelegramLinksRealignmentEnabled(flag bool) (s string) {
	s = "telegram.links.realignment.enabled"
	if flag {
		return AttrAgr(s, "true")
	}
	return AttrAgr(s, "false")
}

func AttrAgr(key string, val string) string {
	return key + ": " + val
}

func Navigation(args string, text ...string) string {
	return formatConstruct("navigation", args, strings.Join(text, " "))
}

func Link(args string, pageId string, accessKey string) string {
	return formatConstruct("link", join(FormatAttr("pageId", pageId), args), accessKey)
}

func join(s1 string, s2 string) string {
	if s2 == "" {
		return s1
	}
	return s1 + " " + s2
}
