package models

import (
	"errors"
	"slices"
	"strings"

	"gorm.io/gorm/clause"
)

// buildSortOrder keeps request parameters out of raw SQL ordering expressions.
func buildSortOrder(sort, order, table string, columns []string) (*clause.OrderByColumn, error) {
	var desc bool
	switch strings.ToLower(order) {
	case "", "asc":
	case "desc":
		desc = true
	default:
		return nil, errors.New("invalid sort order: expected asc or desc")
	}

	if sort == "" {
		return nil, nil
	}
	if !slices.Contains(columns, sort) {
		return nil, errors.New("invalid sort column")
	}

	return &clause.OrderByColumn{
		Column: clause.Column{Table: table, Name: sort},
		Desc:   desc,
	}, nil
}
