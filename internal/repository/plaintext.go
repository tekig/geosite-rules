package repository

import (
	"fmt"
	"strings"

	router "github.com/v2fly/v2ray-core/v5/app/router/routercommon"
)

func PlainText(rules []*router.Domain) (string, error) {
	var buf strings.Builder
	for _, r := range rules {
		line, err := PlainLine(r)
		if err != nil {
			return "", fmt.Errorf("line: %w", err)
		}

		if _, err := buf.WriteString(line); err != nil {
			return "", fmt.Errorf("write line: %w", err)
		}

		if _, err := buf.WriteRune('\n'); err != nil {
			return "", fmt.Errorf("write new line: %w", err)
		}
	}

	return buf.String(), nil
}

func PlainLine(rule *router.Domain) (string, error) {
	var line string
	switch rule.Type {
	case router.Domain_Plain:
		line += "DOMAIN-KEYWORD,"
	case router.Domain_Regex:
		line += "URL-REGEX,"
	case router.Domain_RootDomain:
		line += "DOMAIN-SUFFIX,"
	case router.Domain_Full:
		line += "DOMAIN,"
	default:
		return "", ErrUnknownDomainType
	}

	line += rule.Value

	return line, nil
}
