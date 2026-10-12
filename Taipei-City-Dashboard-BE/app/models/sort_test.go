package models

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Any accidental database access fails the test instead of reaching a server.
type sortTestConnection struct {
	calls int
}

func (c *sortTestConnection) PrepareContext(context.Context, string) (*sql.Stmt, error) {
	c.calls++
	return nil, errors.New("unexpected database access")
}

func (c *sortTestConnection) ExecContext(context.Context, string, ...interface{}) (sql.Result, error) {
	c.calls++
	return nil, errors.New("unexpected database access")
}

func (c *sortTestConnection) QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error) {
	c.calls++
	return nil, errors.New("unexpected database access")
}

func (c *sortTestConnection) QueryRowContext(context.Context, string, ...interface{}) *sql.Row {
	c.calls++
	return &sql.Row{}
}

type sortTestQuery struct {
	sql  string
	vars []interface{}
}

func captureSortQueries(t *testing.T) *[]sortTestQuery {
	t.Helper()
	connection := &sortTestConnection{}
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: connection}), &gorm.Config{
		DryRun:               true,
		DisableAutomaticPing: true,
		Logger:               logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatal(err)
	}

	var queries []sortTestQuery
	err = db.Callback().Query().After("gorm:query").Register("test:capture_sort_queries", func(tx *gorm.DB) {
		if tx.Statement.Dest == nil {
			return // Keep nested component subqueries available to their parent.
		}
		queries = append(queries, sortTestQuery{
			sql:  tx.Statement.SQL.String(),
			vars: append([]interface{}(nil), tx.Statement.Vars...),
		})
		// DryRun retains SQL, unlike normal execution. Reset it so a subsequent
		// Find on the same query builder generates its own SQL after Count.
		tx.Statement.SQL.Reset()
		tx.Statement.Vars = nil
	})
	if err != nil {
		t.Fatal(err)
	}
	previousDB := DBManager
	DBManager = db
	t.Cleanup(func() {
		DBManager = previousDB
		if connection.calls != 0 {
			t.Errorf("query attempted %d database calls", connection.calls)
		}
	})
	return &queries
}

func TestModelSortQueries(t *testing.T) {
	models := []struct {
		name              string
		columns           []string
		table             func(string) string
		defaultOrder      string
		ignoreUnknownSort bool
		query             func(string, string) error
	}{
		{
			name: "contributors",
			columns: []string{
				"id", "user_id", "user_name", "image", "link", "identity", "description", "include", "created_at", "updated_at",
			},
			table:        func(string) string { return "contributors" },
			defaultOrder: `ORDER BY "contributors"."id"`,
			query: func(sort, order string) error {
				_, _, err := GetAllContributors(10, 2, sort, order)
				return err
			},
		},
		{
			name: "issues",
			columns: []string{
				"id", "title", "user_name", "user_id", "context", "description", "decision_desc", "status", "updated_by", "created_at", "updated_at",
			},
			table: func(string) string { return "issues" },
			query: func(sort, order string) error {
				_, _, _, err := GetAllIssues(10, 2, "", sort, order)
				return err
			},
		},
		{
			name:    "incidents",
			columns: []string{"id", "type", "description", "distance", "latitude", "longitude", "place", "time", "status"},
			table:   func(string) string { return "incidents" },
			query: func(sort, order string) error {
				_, _, _, err := GetAllIncident(10, 2, "", sort, order)
				return err
			},
		},
		{
			name: "users",
			columns: []string{
				"id", "name", "email", "password", "idno", "uuid", "tp_account", "member_type", "verify_level",
				"is_admin", "is_active", "is_whitelist", "is_blacked", "expired_at", "created_at", "login_at",
			},
			table: func(string) string { return "auth_users" },
			query: func(sort, order string) error {
				_, _, _, err := GetAllUsers(10, 2, sort, order, "", "")
				return err
			},
		},
		{
			name:              "components",
			ignoreUnknownSort: true,
			columns: []string{
				"id", "index", "name", "update_freq", "source", "short_desc", "long_desc", "use_case", "links", "contributors", "query_type",
			},
			table: func(column string) string {
				switch column {
				case "id", "index", "name":
					return "components"
				default:
					return "query_charts"
				}
			},
			query: func(sort, order string) error {
				_, _, _, err := GetAllComponents("taipei", 10, 2, sort, order, "", "", "", "", "")
				return err
			},
		},
	}

	for _, model := range models {
		t.Run(model.name, func(t *testing.T) {
			tests := []struct {
				name, sort, order, wantOrder string
				wantError                    bool
			}{
				{name: "default", wantOrder: model.defaultOrder},
				{name: "no_sort_with_direction", order: "desc", wantOrder: model.defaultOrder},
				{name: "empty_direction", sort: "id", wantOrder: `ORDER BY "` + model.table("id") + `"."id"`},
				{name: "uppercase_ascending", sort: "id", order: "ASC", wantOrder: `ORDER BY "` + model.table("id") + `"."id"`},
				{name: "uppercase_descending", sort: "id", order: "DESC", wantOrder: `ORDER BY "` + model.table("id") + `"."id" DESC`},
				{name: "unknown_column", sort: "unknown_column", order: "asc", wantError: true},
				{name: "sort_statement", sort: "id; SELECT 1--", order: "asc", wantError: true},
				{name: "sort_expression", sort: "id, (SELECT 1)", order: "asc", wantError: true},
				{name: "sort_direction", sort: "id desc", wantError: true},
				{name: "order_statement", sort: "id", order: "desc; SELECT 1--", wantError: true},
				{name: "order_expression", sort: "id", order: "desc, (SELECT 1)", wantError: true},
				{name: "invalid_direction", sort: "id", order: "sideways", wantError: true},
				{name: "invalid_direction_without_sort", order: "desc; SELECT 1--", wantError: true},
			}
			for _, column := range model.columns {
				tests = append(tests, struct {
					name, sort, order, wantOrder string
					wantError                    bool
				}{
					name:      "allowed_" + column,
					sort:      column,
					order:     "desc",
					wantOrder: `ORDER BY "` + model.table(column) + `"."` + column + `" DESC`,
				})
			}
			if model.ignoreUnknownSort {
				for index := range tests {
					switch tests[index].name {
					case "unknown_column", "sort_statement", "sort_expression", "sort_direction":
						tests[index].wantError = false
					}
				}
			}
			for _, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					queries := captureSortQueries(t)
					err := model.query(test.sort, test.order)
					if test.wantError {
						if err == nil {
							t.Fatal("expected invalid ordering to be rejected")
						}
						if len(*queries) != 0 {
							t.Fatalf("invalid ordering generated %d queries", len(*queries))
						}
						return
					}
					if err != nil {
						t.Fatal(err)
					}
					if len(*queries) == 0 {
						t.Fatal("model did not build a query")
					}
					query := (*queries)[len(*queries)-1].sql
					actualOrder := ""
					if index := strings.Index(query, "ORDER BY "); index >= 0 {
						actualOrder = strings.Split(query[index:], " LIMIT ")[0]
					}
					if actualOrder != test.wantOrder {
						t.Errorf("ordering = %q, want %q; SQL: %s", actualOrder, test.wantOrder, query)
					}
					if !strings.HasSuffix(query, " LIMIT 10 OFFSET 10") {
						t.Errorf("pagination missing from query: %s", query)
					}
					if model.name == "components" && (!strings.Contains(query, "SELECT components.id") || !strings.Contains(query, "SELECT query_charts.index")) {
						t.Errorf("component subqueries missing from query: %s", query)
					}
				})
			}
		})
	}
}

func TestComponentSortKeepsExistingAllowlist(t *testing.T) {
	queries := captureSortQueries(t)
	_, _, _, err := GetAllComponents("taipei", 0, 0, "updated_at", "desc", "", "", "", "", "")
	if err != nil || len(*queries) == 0 {
		t.Fatalf("unsupported component sort must retain its ignored behavior: %v", err)
	}
	if query := (*queries)[len(*queries)-1].sql; strings.Contains(query, "ORDER BY") {
		t.Fatalf("column outside the component allowlist reached ordering: %s", query)
	}
}
