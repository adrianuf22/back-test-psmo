package request

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	Get    = Method{"GET"}
	Post   = Method{"POST"}
	Put    = Method{"PUT"}
	Delete = Method{"DELETE"}
	Head   = Method{"HEAD"}
)

type Method struct {
	name string
}

func (m Method) WithPath(s ...string) string {
	re := regexp.MustCompile(`\/{2,}`)
	repl := re.ReplaceAll([]byte(strings.Join(s, "/")), []byte("/"))

	return fmt.Sprintf("%s %s", m.name, string(repl))
}
