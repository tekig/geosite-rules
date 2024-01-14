package repository

import (
	"context"
	"fmt"
	"os"
	"strings"

	router "github.com/v2fly/v2ray-core/v5/app/router/routercommon"
	"google.golang.org/protobuf/proto"
)

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

func (r *DLC) Rules(ctx context.Context, name string) ([]*router.Domain, error) {
	rules, ok := r.db[strings.ToUpper(name)]
	if !ok {
		return nil, ErrRuleNotFound
	}

	return rules, nil
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
