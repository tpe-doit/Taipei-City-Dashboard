package models

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// This connector records the SQL and bound values actually passed through GORM
// to database/sql. It never opens a socket or executes a database query.
type componentDataConnector struct {
	columns []string
	rows    [][]driver.Value
	queries []componentDataRecordedQuery
}

type componentDataRecordedQuery struct {
	sql  string
	args []driver.NamedValue
}

func (c *componentDataConnector) Connect(context.Context) (driver.Conn, error) {
	return &componentDataConnection{connector: c}, nil
}
func (c *componentDataConnector) Driver() driver.Driver { return componentDataDriver{} }

type componentDataDriver struct{}

func (componentDataDriver) Open(string) (driver.Conn, error) {
	return nil, errors.New("unexpected driver open")
}

type componentDataConnection struct{ connector *componentDataConnector }

func (*componentDataConnection) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("unexpected prepare")
}
func (*componentDataConnection) Close() error { return nil }
func (*componentDataConnection) Begin() (driver.Tx, error) {
	return nil, errors.New("unexpected transaction")
}
func (c *componentDataConnection) QueryContext(_ context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	c.connector.queries = append(c.connector.queries, componentDataRecordedQuery{
		sql: query, args: append([]driver.NamedValue(nil), args...),
	})
	return &componentDataRows{columns: c.connector.columns, rows: c.connector.rows}, nil
}

type componentDataRows struct {
	columns []string
	rows    [][]driver.Value
	next    int
}

func (r *componentDataRows) Columns() []string { return r.columns }
func (*componentDataRows) Close() error        { return nil }
func (r *componentDataRows) Next(dest []driver.Value) error {
	if r.next >= len(r.rows) {
		return io.EOF
	}
	copy(dest, r.rows[r.next])
	r.next++
	return nil
}

func componentDataTestDB(t *testing.T, columns []string, rows [][]driver.Value) (*gorm.DB, *componentDataConnector) {
	t.Helper()
	connector := &componentDataConnector{columns: columns, rows: rows}
	connection := sql.OpenDB(connector)
	t.Cleanup(func() { _ = connection.Close() })
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: connection}), &gorm.Config{
		DisableAutomaticPing: true, Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatal(err)
	}
	return db, connector
}

func TestComponentDataBindsDatesForEveryChartType(t *testing.T) {
	sampleTime := time.Date(2026, 10, 12, 8, 0, 0, 0, time.FixedZone("UTC+8", 8*60*60))
	types := []struct {
		name    string
		columns []string
		rows    [][]driver.Value
		query   func(*string, string, string) (any, error)
		want    string
	}{
		{"two_d", []string{"x_axis", "data"}, [][]driver.Value{{"A", float64(17)}}, func(q *string, from, to string) (any, error) { return GetTwoDimensionalData(q, from, to) }, `[{"data":[{"x":"A","y":17}]}]`},
		{"three_d_and_percent", []string{"x_axis", "y_axis", "icon", "data"}, [][]driver.Value{{"A", "Series", "%", int64(17)}}, func(q *string, from, to string) (any, error) {
			data, _, err := GetThreeDimensionalData(q, from, to)
			return data, err
		}, `[{"name":"Series","icon":"%","data":[17]}]`},
		{"time", []string{"x_axis", "y_axis", "data"}, [][]driver.Value{{sampleTime, "Series", float64(17)}}, func(q *string, from, to string) (any, error) { return GetTimeSeriesData(q, from, to) }, `[{"name":"Series","data":[{"x":"2026-10-12T08:00:00+08:00","y":17}]}]`},
		{"map_legend", []string{"name", "type", "icon", "value"}, [][]driver.Value{{"A", "fill", "circle", float64(17)}}, func(q *string, from, to string) (any, error) { return GetMapLegendData(q, from, to) }, `[{"name":"A","type":"fill","icon":"circle","value":17}]`},
		{"bubble", []string{"y_axis", "x", "y", "z"}, [][]driver.Value{{"Series", "5", float64(6), float64(7)}}, func(q *string, from, to string) (any, error) {
			data, _, err := GetBubbleData(q, from, to)
			return data, err
		}, `[{"name":"Series","data":[{"x":5,"y":6,"z":7}]}]`},
		{"layered_flow", []string{"source", "source_layer", "target", "target_layer", "value"}, [][]driver.Value{{"A", int64(1), "B", int64(2), float64(17)}}, func(q *string, from, to string) (any, error) {
			data, _, err := GetLayeredFlowData(q, from, to)
			return data, err
		}, `[{"source":"A","source_layer":1,"target":"B","target_layer":2,"value":17}]`},
	}
	for _, chart := range types {
		t.Run(chart.name, func(t *testing.T) {
			for _, test := range []struct{ name, from, to, template string }{
				{"normal_quoted", "2026-10-01T00:00:00+08:00", "2026-10-12T08:00:00+08:00", "SELECT * FROM chart WHERE sampled_at BETWEEN '%s' AND '%s'"},
				{"normal_bare", "2026-10-01T00:00:00+08:00", "2026-10-12T08:00:00+08:00", "SELECT * FROM chart WHERE sampled_at BETWEEN %s::timestamptz AND %s::timestamptz"},
				{"statement_injection", "2026-10-01T00:00:00+08:00'; SELECT pg_sleep(10); --", "2026-10-12T08:00:00+08:00", "SELECT * FROM chart WHERE sampled_at BETWEEN '%s' AND '%s'"},
				{"expression_injection", "2026-10-01T00:00:00+08:00", "x' OR TRUE --", "SELECT * FROM chart WHERE sampled_at BETWEEN '%s' AND '%s'"},
				{"sql_operators_and_literals", "2026-10-01T00:00:00+08:00", "2026-10-12T08:00:00+08:00", "SELECT * FROM chart WHERE sampled_at BETWEEN '%s' AND '%s' AND payload ? 'key' AND note = '?' AND label LIKE '50%'"},
			} {
				t.Run(test.name, func(t *testing.T) {
					db, recorder := componentDataTestDB(t, chart.columns, chart.rows)
					previousDB := DBDashboard
					DBDashboard = db
					t.Cleanup(func() { DBDashboard = previousDB })
					data, err := chart.query(&test.template, test.from, test.to)
					if err != nil {
						t.Fatal(err)
					}
					output, err := json.Marshal(data)
					if err != nil || string(output) != chart.want {
						t.Fatalf("output = %s, err = %v, want %s", output, err, chart.want)
					}
					if len(recorder.queries) != 1 {
						t.Fatalf("queries = %d, want 1", len(recorder.queries))
					}
					query := recorder.queries[0]
					wantSQL := strings.Replace(test.template, "'%s'", "$1", 1)
					wantSQL = strings.Replace(wantSQL, "'%s'", "$2", 1)
					wantSQL = strings.Replace(wantSQL, "%s", "$1", 1)
					wantSQL = strings.Replace(wantSQL, "%s", "$2", 1)
					if query.sql != wantSQL {
						t.Errorf("SQL = %q, want %q", query.sql, wantSQL)
					}
					var args []any
					for _, arg := range query.args {
						args = append(args, arg.Value)
					}
					if !reflect.DeepEqual(args, []any{test.from, test.to}) {
						t.Errorf("bound args = %#v, want %#v", args, []any{test.from, test.to})
					}
				})
			}
		})
	}
}

func TestComponentHistoryBindsEveryRepeatedGroup(t *testing.T) {
	from := "2026-01-01T00:00:00+08:00"
	for _, period := range []struct{ name, to, unit string }{
		{"hour_at_24_hours", "2026-01-02T00:00:00+08:00", "hour"},
		{"day_below_32_days", "2026-02-01T00:00:00+08:00", "day"},
		{"week_at_32_days", "2026-02-02T00:00:00+08:00", "week"},
		{"month_at_93_days", "2026-04-04T00:00:00+08:00", "month"},
		{"year_at_740_days", "2028-01-11T00:00:00+08:00", "year"},
	} {
		for _, groups := range []int{1, 2, 5} {
			t.Run(period.name+strings.Repeat("_group", groups), func(t *testing.T) {
				group := "SELECT date_trunc('%s', sampled_at::timestamp) AS x_axis, 'A' AS y_axis, 1 AS data FROM chart WHERE sampled_at BETWEEN '%s' AND '%s'"
				template := group + strings.Repeat(" UNION ALL "+group, groups-1)
				managerDB, manager := componentDataTestDB(t, []string{"query_history"}, [][]driver.Value{{template}})
				dashboardDB, dashboard := componentDataTestDB(t, []string{"x_axis", "y_axis", "data"}, nil)
				oldManager, oldDashboard := DBManager, DBDashboard
				DBManager, DBDashboard = managerDB, dashboardDB
				t.Cleanup(func() { DBManager, DBDashboard = oldManager, oldDashboard })
				query, err := GetComponentHistoryDataQuery(42, "taipei", from, period.to)
				if err != nil || query != template {
					t.Fatalf("history template = %q, err = %v, want unchanged template", query, err)
				}
				if len(manager.queries) != 1 || len(manager.queries[0].args) != 2 || manager.queries[0].args[0].Value != int64(42) || manager.queries[0].args[1].Value != "taipei" {
					t.Fatalf("history lookup did not bind component ID and city: %#v", manager.queries)
				}
				if _, err = GetTimeSeriesData(&query, from, period.to); err != nil {
					t.Fatal(err)
				}
				if len(dashboard.queries) != 1 {
					t.Fatalf("queries = %d, want 1", len(dashboard.queries))
				}
				got := dashboard.queries[0]
				if strings.Contains(got.sql, "%s") || strings.Contains(got.sql, from) || strings.Contains(got.sql, period.to) || strings.Contains(got.sql, "'"+period.unit+"'") {
					t.Fatalf("request values remain in history SQL: %q", got.sql)
				}
				var args, want []any
				for _, arg := range got.args {
					args = append(args, arg.Value)
				}
				for i := 0; i < groups; i++ {
					want = append(want, period.unit, from, period.to)
				}
				if !reflect.DeepEqual(args, want) {
					t.Fatalf("history args = %#v, want %#v", args, want)
				}
			})
		}
	}
}

func TestComponentDataRejectsAmbiguousTemplatesBeforeQuery(t *testing.T) {
	for _, template := range []string{
		"SELECT '%s suffix', '%s'", "SELECT 'prefix %s', '%s'", "SELECT \"%s\", '%s'",
		"SELECT $tag$%s$tag$, '%s'", "SELECT %%s, '%s'", "SELECT prefix%s, '%s'",
		"SELECT %s_suffix, '%s'", "SELECT '%s'", "SELECT '%s', '%s', '%s', '%s'",
		"SELECT '%s', '%s' /* unclosed", "SELECT '%s', '%s",
		"SELECT TIMESTAMP '%s', '%s'", "SELECT INTERVAL '%s', '%s'", "SELECT DATE '%s', '%s'",
		"SELECT TIME WITH TIME ZONE '%s', '%s'",
	} {
		t.Run(template, func(t *testing.T) {
			db, recorder := componentDataTestDB(t, nil, nil)
			previousDB := DBDashboard
			DBDashboard = db
			t.Cleanup(func() { DBDashboard = previousDB })
			if _, err := GetTimeSeriesData(&template, "2026-10-01T00:00:00+08:00", "2026-10-12T00:00:00+08:00"); err == nil {
				t.Fatal("invalid template accepted")
			}
			if len(recorder.queries) != 0 {
				t.Fatalf("invalid template reached SQL driver: %#v", recorder.queries)
			}
		})
	}
}

func TestComponentDataPreservesStaticSQLAndComments(t *testing.T) {
	for _, template := range []string{
		"SELECT '50%' AS x_axis, 1 AS data WHERE payload ? 'key' AND note = '?'",
		"SELECT $$a ? and 'quote'$$ AS x_axis, 1 AS data -- %s ignored\n /* %s /* nested */ %s */",
		"SELECT 'it''s fine' AS x_axis, 1 AS data",
		"SELECT foo$tag$ AS x_axis, 1 AS data",
	} {
		t.Run(template, func(t *testing.T) {
			db, recorder := componentDataTestDB(t, []string{"x_axis", "data"}, nil)
			previousDB := DBDashboard
			DBDashboard = db
			t.Cleanup(func() { DBDashboard = previousDB })
			if _, err := GetTwoDimensionalData(&template, "unused", "unused"); err != nil {
				t.Fatal(err)
			}
			if len(recorder.queries) != 1 || recorder.queries[0].sql != template || len(recorder.queries[0].args) != 0 {
				t.Fatalf("static template changed: %#v", recorder.queries)
			}
		})
	}
}

func TestComponentDataBindsEscapedStringLiterals(t *testing.T) {
	template := "SELECT * FROM chart WHERE sampled_at BETWEEN E'%s'::timestamp AND e'%s'::timestamp AND note = E'it\\'s ? fine'"
	db, recorder := componentDataTestDB(t, []string{"x_axis", "data"}, nil)
	previousDB := DBDashboard
	DBDashboard = db
	t.Cleanup(func() { DBDashboard = previousDB })
	if _, err := GetTwoDimensionalData(&template, "2026-10-01", "2026-10-12"); err != nil {
		t.Fatal(err)
	}
	wantSQL := "SELECT * FROM chart WHERE sampled_at BETWEEN $1::timestamp AND $2::timestamp AND note = E'it\\'s ? fine'"
	if len(recorder.queries) != 1 || recorder.queries[0].sql != wantSQL || len(recorder.queries[0].args) != 2 {
		t.Fatalf("escaped literal template changed: %#v", recorder.queries)
	}
}

func TestComponentHistoryRejectsInvalidDatesAndTemplates(t *testing.T) {
	for _, test := range []struct{ template, from, to string }{
		{"SELECT '%s', '%s', '%s'", "x'; SELECT 1--", "2026-10-12T00:00:00+08:00"},
		{"SELECT '%s', '%s', '%s'", "2026-10-01T00:00:00+08:00", "x' OR TRUE--"},
		{"SELECT '%s', '%s'", "2026-10-01T00:00:00+08:00", "2026-10-12T00:00:00+08:00"},
	} {
		db, _ := componentDataTestDB(t, []string{"query_history"}, [][]driver.Value{{test.template}})
		previousDB := DBManager
		DBManager = db
		if _, err := GetComponentHistoryDataQuery(42, "taipei", test.from, test.to); err == nil {
			t.Errorf("invalid history input accepted: %#v", test)
		}
		DBManager = previousDB
	}
}
