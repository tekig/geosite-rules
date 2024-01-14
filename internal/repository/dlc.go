package repository

import (
	"context"
	"fmt"
	"os"
	"strings"

	router "github.com/v2fly/v2ray-core/v5/app/router/routercommon"
	"google.golang.org/protobuf/proto"
)

var _ Repository = (*DLC)(nil)

type DLC struct {
	db map[string][]*router.Domain
}

func NewDLC(path string) (*DLC, error) {
	r := &DLC{}

	db, err := r.unmarshal(path)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return &DLC{
		db: db,
	}, nil
}

func (r *DLC) Rules(ctx context.Context, name string) (string, error) {
	rules, ok := r.db[strings.ToUpper(name)]
	if !ok {
		return "", ErrRuleNotFound
	}

	list, err := r.plainText(rules)
	if err != nil {
		return "", fmt.Errorf("plain text: %w", err)
	}

	return list, nil
}

func (r *DLC) plainText(rules []*router.Domain) (string, error) {
	var buf strings.Builder
	for _, rule := range rules {
		line, err := r.plainLine(rule)
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

func (r *DLC) plainLine(rule *router.Domain) (string, error) {
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

func (r *DLC) unmarshal(path string) (map[string][]*router.Domain, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}

	var list = &router.GeoSiteList{}
	if err := proto.Unmarshal(data, list); err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}

	var db = make(map[string][]*router.Domain, len(list.Entry))
	for _, entry := range list.Entry {
		db[entry.CountryCode] = entry.Domain
	}

	return db, nil
}
