package models

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm/clause"
)

// Dashboard queries are trusted SQL templates stored in query_charts. Request
// dates are always values, never SQL fragments. Building the expression in
// pieces also preserves PostgreSQL '?' operators and '?' inside SQL literals.
type dashboardQuery struct {
	parts  []string
	values []any
}

func (query dashboardQuery) Build(builder clause.Builder) {
	for i, part := range query.parts {
		builder.WriteString(part)
		if i < len(query.values) {
			builder.AddVar(builder, query.values[i])
		}
	}
}

func bindDashboardQuery(template, timeFrom, timeTo string) (dashboardQuery, error) {
	parts, err := splitDashboardQuery(template)
	if err != nil {
		return dashboardQuery{}, err
	}
	query := dashboardQuery{parts: parts}
	count := len(parts) - 1
	switch {
	case count == 0:
		return query, nil
	case count == 2:
		query.values = []any{timeFrom, timeTo}
	case count%3 == 0:
		unit, err := dashboardHistoryTimeStep(timeFrom, timeTo)
		if err != nil {
			return dashboardQuery{}, err
		}
		for i := 0; i < count/3; i++ {
			query.values = append(query.values, unit, timeFrom, timeTo)
		}
	default:
		return dashboardQuery{}, fmt.Errorf("invalid query string")
	}
	return query, nil
}

// Legacy templates use either bare %s or the complete string literal '%s'.
// Reject placeholders embedded in identifiers/strings instead of interpolating
// untrusted values or silently changing the template's meaning.
func splitDashboardQuery(template string) ([]string, error) {
	var parts []string
	start := 0
	invalid := func() ([]string, error) { return nil, fmt.Errorf("invalid query string") }
	for i := 0; i < len(template); {
		switch {
		case strings.HasPrefix(template[i:], "--"):
			end := strings.IndexByte(template[i:], '\n')
			if end < 0 {
				i = len(template)
			} else {
				i += end + 1
			}
		case strings.HasPrefix(template[i:], "/*"):
			depth := 1
			i += 2
			for i < len(template) && depth > 0 {
				if strings.HasPrefix(template[i:], "/*") {
					depth++
					i += 2
				} else if strings.HasPrefix(template[i:], "*/") {
					depth--
					i += 2
				} else {
					i++
				}
			}
			if depth != 0 {
				return invalid()
			}
		case template[i] == '\'' || template[i] == '"':
			quote, begin := template[i], i
			escaped := quote == '\'' && i > 0 && (template[i-1] == 'E' || template[i-1] == 'e') && (i == 1 || !dashboardSQLWord(template[i-2]))
			i++
			closed := false
			for i < len(template) {
				if escaped && template[i] == '\\' {
					i += 2
				} else if template[i] == quote {
					i++
					if i < len(template) && template[i] == quote {
						i++
						continue
					}
					closed = true
					break
				} else {
					i++
				}
			}
			if !closed {
				return invalid()
			}
			literal := template[begin:i]
			if quote == '\'' && literal == "'%s'" {
				if dashboardTypedLiteralPrefix(template[:begin]) {
					return invalid() // TIMESTAMP $1 is not valid PostgreSQL syntax.
				}
				if escaped {
					begin-- // E'%s' is also a complete value literal.
				}
				parts = append(parts, template[start:begin])
				start = i
			} else if strings.Contains(literal, "%s") {
				return invalid()
			}
		case template[i] == '$' && (i == 0 || !dashboardSQLWord(template[i-1])) && dashboardDollarDelimiter(template[i:]) != "":
			delimiter := dashboardDollarDelimiter(template[i:])
			end := strings.Index(template[i+len(delimiter):], delimiter)
			if end < 0 {
				return invalid()
			}
			end += i + 2*len(delimiter)
			if strings.Contains(template[i:end], "%s") {
				return invalid()
			}
			i = end
		case strings.HasPrefix(template[i:], "%s"):
			if (i > 0 && (dashboardSQLWord(template[i-1]) || template[i-1] == '%')) || (i+2 < len(template) && dashboardSQLWord(template[i+2])) {
				return invalid()
			}
			parts = append(parts, template[start:i])
			i += 2
			start = i
		default:
			i++
		}
	}
	return append(parts, template[start:]), nil
}

func dashboardTypedLiteralPrefix(prefix string) bool {
	prefix = strings.ToLower(strings.TrimSpace(prefix))
	for _, name := range []string{"date", "time", "timetz", "timestamp", "timestamptz", "interval", "time with time zone", "time without time zone", "timestamp with time zone", "timestamp without time zone"} {
		if strings.HasSuffix(prefix, name) && (len(prefix) == len(name) || !dashboardSQLWord(prefix[len(prefix)-len(name)-1])) {
			return true
		}
	}
	return false
}

func dashboardSQLWord(char byte) bool {
	return char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' || char >= '0' && char <= '9' || char == '_' || char >= 0x80
}

func dashboardDollarDelimiter(sql string) string {
	for i := 1; i < len(sql); i++ {
		if sql[i] == '$' {
			return sql[:i+1]
		}
		if !dashboardSQLWord(sql[i]) || (i == 1 && sql[i] >= '0' && sql[i] <= '9') {
			break
		}
	}
	return ""
}

func dashboardHistoryTimeStep(timeFrom, timeTo string) (string, error) {
	const layout = "2006-01-02T15:04:05+08:00"
	from, err := time.Parse(layout, timeFrom)
	if err != nil {
		return "", err
	}
	to, err := time.Parse(layout, timeTo)
	if err != nil {
		return "", err
	}
	hours := to.Sub(from).Hours()
	switch {
	case hours <= 24:
		return "hour", nil
	case hours < 24*32:
		return "day", nil
	case hours < 24*93:
		return "week", nil
	case hours < 24*740:
		return "month", nil
	default:
		return "year", nil
	}
}

func executeDashboardRawQuery(query *string, timeFrom, timeTo string, dest any) error {
	if query == nil {
		return fmt.Errorf("query is nil")
	}
	bound, err := bindDashboardQuery(*query, timeFrom, timeTo)
	if err != nil {
		return err
	}
	return DBDashboard.Raw("?", bound).Scan(dest).Error
}
