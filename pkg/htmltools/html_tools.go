package htmltools

import (
	"fmt"
	"net/http"
	"strings"
)

const replacement = "<br>\n"

var replacer = strings.NewReplacer(
	"\r\n", replacement,
	"\r", replacement,
	"\n", replacement,
	"\v", replacement,
	"\f", replacement,
	"\u0085", replacement,
	"\u2028", replacement,
	"\u2029", replacement,
)

func GenerateHtmlHeader(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/html, charset=UTF-8")

	fmt.Fprintln(w, "<!DOCTYPE html>")
	fmt.Fprintln(w, "<html lang=\"en\">")
	fmt.Fprintln(w, "<head>")
	fmt.Fprintln(w, "</head>")
	fmt.Fprintln(w, "<body>")
}

func GenerateHtmlFooter(w http.ResponseWriter) {
	fmt.Fprintln(w, "</body>")
	fmt.Fprintln(w, "</html>")

}

func ConvertLineBreaksToHtml(s string) string {
	return replacer.Replace(s)
}
