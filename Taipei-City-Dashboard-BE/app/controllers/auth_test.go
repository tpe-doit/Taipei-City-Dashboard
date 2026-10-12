package controllers

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"TaipeiCityDashboardBE/app/models"
	"TaipeiCityDashboardBE/app/util"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const loginTestEmail = "login-test@example.invalid"

func TestLoginPasswordPersistence(t *testing.T) {
	password := "fake login password"
	strong, err := util.HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name string
		hash string
	}{
		{name: "legacy password is upgraded", hash: util.HashString(password)},
		{name: "current password is preserved", hash: strong},
	} {
		t.Run(tc.name, func(t *testing.T) {
			state := newLoginTestDatabase(&tc.hash)
			useLoginTestDatabase(t, state)
			response := requestTestLogin(t, password)
			assertLoginResponse(t, response, http.StatusOK, true)
			if state.updates != 1 || state.commits != 1 {
				t.Fatalf("expected one committed login update, got %d updates and %d commits", state.updates, state.commits)
			}
			if state.password == nil {
				t.Fatal("login removed the stored password")
			}
			if valid, upgrade := util.VerifyPassword(password, *state.password); !valid || upgrade {
				t.Fatal("persisted password must verify in the current format")
			}
			if tc.hash == strong && *state.password != strong {
				t.Fatal("login replaced an already current password hash")
			}
			if !state.loginAt.After(state.previousLoginAt) {
				t.Fatal("successful login did not persist its new timestamp")
			}
			if state.permissionQueries != 2 {
				t.Fatalf("expected both permission queries, got %d", state.permissionQueries)
			}
		})
	}
}

func TestLoginRejectsCredentialsWithoutWrites(t *testing.T) {
	password := "fake login password"
	legacy := util.HashString(password)
	strong, err := util.HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}
	malformed := "$argon2id$v=19$m=19456,t=2,p=1$invalid"
	for _, tc := range []struct {
		name     string
		hash     *string
		password string
		notFound bool
	}{
		{name: "wrong legacy password", hash: &legacy, password: "wrong"},
		{name: "wrong current password", hash: &strong, password: "wrong"},
		{name: "SSO account without password", password: password},
		{name: "malformed stored password", hash: &malformed, password: password},
		{name: "account not found", hash: &legacy, password: password, notFound: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			state := newLoginTestDatabase(tc.hash)
			state.notFound = tc.notFound
			useLoginTestDatabase(t, state)
			assertLoginResponse(t, requestTestLogin(t, tc.password), http.StatusUnauthorized, false)
			if state.updates != 0 || state.commits != 0 || state.permissionQueries != 0 {
				t.Fatal("rejected credentials must not update the account or query permissions")
			}
		})
	}
}

func TestLoginRejectsInactiveAccountWithoutWrites(t *testing.T) {
	password := "fake login password"
	legacy := util.HashString(password)
	for _, tc := range []struct {
		name   string
		active driver.Value
	}{
		{name: "inactive", active: false},
		{name: "NULL activation", active: nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			state := newLoginTestDatabase(&legacy)
			state.active = tc.active
			useLoginTestDatabase(t, state)
			assertLoginResponse(t, requestTestLogin(t, password), http.StatusForbidden, false)
			if state.updates != 0 || state.commits != 0 || state.permissionQueries != 0 {
				t.Fatal("inactive accounts must not upgrade passwords or query permissions")
			}
		})
	}
}

func TestLoginSaveFailureDoesNotReturnToken(t *testing.T) {
	password := "fake login password"
	legacy := util.HashString(password)
	state := newLoginTestDatabase(&legacy)
	state.failUpdate = true
	useLoginTestDatabase(t, state)
	assertLoginResponse(t, requestTestLogin(t, password), http.StatusInternalServerError, false)
	if state.updates != 1 || state.commits != 0 {
		t.Fatal("failed login persistence must not commit")
	}
	if state.password == nil || *state.password != legacy {
		t.Fatal("failed login persistence changed the stored credential")
	}
}

func TestLoginCaseVariantEmailCandidates(t *testing.T) {
	password := "later account fake password"
	otherPassword := "earlier account fake password"
	strong, err := util.HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}
	otherStrong, err := util.HashPassword(otherPassword)
	if err != nil {
		t.Fatal(err)
	}
	legacy := util.HashString(password)
	otherLegacy := util.HashString(otherPassword)
	for _, tc := range []struct {
		name         string
		firstHash    *string
		matchingHash string
	}{
		{name: "earlier strong and later legacy", firstHash: &otherStrong, matchingHash: legacy},
		{name: "earlier legacy and later strong", firstHash: &otherLegacy, matchingHash: strong},
		{name: "earlier NULL and later strong", matchingHash: strong},
	} {
		t.Run(tc.name, func(t *testing.T) {
			state := newLoginTestDatabase(nil)
			state.candidates = []loginTestAccount{
				{id: 7, email: loginTestEmail, password: tc.firstHash, active: true, loginAt: state.previousLoginAt},
				{id: 14, email: strings.ToUpper(loginTestEmail), password: &tc.matchingHash, active: true, loginAt: state.previousLoginAt},
			}
			useLoginTestDatabase(t, state)
			response := requestTestLogin(t, password)
			assertLoginResponse(t, response, http.StatusOK, true)
			var result struct {
				User struct {
					ID int `json:"user_id"`
				} `json:"user"`
			}
			if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
				t.Fatal(err)
			}
			if result.User.ID != 14 || state.updatedAccountID != 14 || state.permissionAccountID != 14 {
				t.Fatal("login must return, authorize, and update only the later matching account")
			}
			if state.updates != 1 || state.commits != 1 {
				t.Fatal("matching account must have exactly one committed login update")
			}
			first, matching := state.candidates[0], state.candidates[1]
			if first.password != tc.firstHash || !first.loginAt.Equal(state.previousLoginAt) {
				t.Fatal("login changed the earlier nonmatching account")
			}
			if matching.password == nil || !matching.loginAt.After(state.previousLoginAt) {
				t.Fatal("login did not persist the later matching account")
			}
			if valid, upgrade := util.VerifyPassword(password, *matching.password); !valid || upgrade {
				t.Fatal("later account password must persist in the current format")
			}
			if tc.matchingHash == strong && *matching.password != strong {
				t.Fatal("later account's current password hash was unnecessarily changed")
			}

			permissionQueries := state.permissionQueries
			assertLoginResponse(t, requestTestLogin(t, "matches neither account"), http.StatusUnauthorized, false)
			if state.updates != 1 || state.commits != 1 || state.permissionQueries != permissionQueries {
				t.Fatal("wrong password must reject every candidate without further writes or permission queries")
			}
			if state.candidates[1].password != matching.password || !state.candidates[1].loginAt.Equal(matching.loginAt) {
				t.Fatal("rejected login changed the previously authenticated account")
			}
		})
	}
}

func TestLoginUsesEarliestMatchingCaseVariant(t *testing.T) {
	password := "shared fake password"
	strong, err := util.HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}
	legacy := util.HashString(password)
	for _, tc := range []struct {
		name   string
		active driver.Value
		status int
	}{
		{name: "earliest matching account", active: true, status: http.StatusOK},
		{name: "earliest matching account is inactive", active: false, status: http.StatusForbidden},
		{name: "earliest matching account has NULL activation", active: nil, status: http.StatusForbidden},
	} {
		t.Run(tc.name, func(t *testing.T) {
			state := newLoginTestDatabase(nil)
			// Storage order differs from ID order: the query must preserve the
			// earliest matching account selected by the original login behavior.
			state.candidates = []loginTestAccount{
				{id: 14, email: strings.ToUpper(loginTestEmail), password: &legacy, active: true, loginAt: state.previousLoginAt},
				{id: 7, email: loginTestEmail, password: &strong, active: tc.active, loginAt: state.previousLoginAt},
			}
			useLoginTestDatabase(t, state)
			response := requestTestLogin(t, password)
			assertLoginResponse(t, response, tc.status, tc.status == http.StatusOK)
			if tc.status == http.StatusOK {
				var body struct {
					User struct {
						ID int `json:"user_id"`
					} `json:"user"`
				}
				if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
					t.Fatal(err)
				}
				if body.User.ID != 7 || state.updatedAccountID != 7 || state.permissionAccountID != 7 || state.commits != 1 {
					t.Fatal("login did not select the earliest matching account")
				}
			} else if state.updates != 0 || state.commits != 0 || state.permissionQueries != 0 {
				t.Fatal("login must reject the inactive earliest match without trying a later active account")
			}
			if state.candidates[0].password != &legacy || !state.candidates[0].loginAt.Equal(state.previousLoginAt) {
				t.Fatal("login changed the later account with the same password")
			}
		})
	}
}

func requestTestLogin(t *testing.T, password string) *httptest.ResponseRecorder {
	t.Helper()
	response := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(response)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
	credentials := base64.StdEncoding.EncodeToString([]byte(loginTestEmail + ":" + password))
	ctx.Request.Header.Set("Authorization", "Basic "+credentials)
	Login(ctx)
	return response
}

func assertLoginResponse(t *testing.T, response *httptest.ResponseRecorder, status int, wantToken bool) {
	t.Helper()
	if response.Code != status {
		t.Fatalf("status = %d, want %d; body: %s", response.Code, status, response.Body.String())
	}
	var body map[string]json.RawMessage
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	var token string
	if raw, exists := body["token"]; exists {
		if !wantToken {
			t.Fatal("failed login returned a token field")
		}
		if err := json.Unmarshal(raw, &token); err != nil {
			t.Fatal(err)
		}
	}
	if wantToken && token == "" {
		t.Fatal("successful login did not return a token")
	}
	if wantToken {
		var user map[string]json.RawMessage
		if err := json.Unmarshal(body["user"], &user); err != nil {
			t.Fatal(err)
		}
		if _, exists := user["password"]; exists {
			t.Fatal("login response exposed the stored password")
		}
	}
}

// Exercise real GORM queries and Save through database/sql without a network
// connection. Reject unexpected statements so tests cannot silently skip writes.
type loginTestDatabase struct {
	password            *string
	active              driver.Value
	notFound            bool
	failUpdate          bool
	updates             int
	commits             int
	permissionQueries   int
	previousLoginAt     time.Time
	loginAt             time.Time
	candidates          []loginTestAccount
	updatedAccountID    int64
	permissionAccountID int64
}

type loginTestAccount struct {
	id       int64
	email    string
	password *string
	active   driver.Value
	loginAt  time.Time
}

func newLoginTestDatabase(password *string) *loginTestDatabase {
	previous := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	return &loginTestDatabase{password: password, active: true, previousLoginAt: previous, loginAt: previous}
}

func useLoginTestDatabase(t *testing.T, state *loginTestDatabase) {
	t.Helper()
	sqlDB := sql.OpenDB(loginTestConnector{state})
	sqlDB.SetMaxOpenConns(1)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{
		DisableAutomaticPing: true,
		Logger:               logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		_ = sqlDB.Close()
		t.Fatal(err)
	}
	previous := models.DBManager
	models.DBManager = db
	t.Cleanup(func() {
		models.DBManager = previous
		_ = sqlDB.Close()
	})
}

type loginTestConnector struct{ state *loginTestDatabase }

func (c loginTestConnector) Connect(context.Context) (driver.Conn, error) {
	return &loginTestConnection{state: c.state}, nil
}
func (loginTestConnector) Driver() driver.Driver { return loginTestDriver{} }

type loginTestDriver struct{}

func (loginTestDriver) Open(string) (driver.Conn, error) {
	return nil, errors.New("login tests require their in-memory connector")
}

type loginTestConnection struct {
	state            *loginTestDatabase
	pendingPassword  string
	pendingLoginAt   time.Time
	pendingAccountID int64
}

func (*loginTestConnection) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("unexpected prepared statement in login test")
}
func (*loginTestConnection) Close() error { return nil }
func (c *loginTestConnection) Begin() (driver.Tx, error) {
	return loginTestTransaction{c}, nil
}

func (c *loginTestConnection) QueryContext(_ context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	switch {
	case strings.Contains(query, `FROM "auth_users"`):
		if len(args) != 1 || args[0].Value != loginTestEmail {
			return nil, errors.New("login must look up the account by email before verifying its password")
		}
		rows := &loginTestRows{columns: []string{"id", "name", "email", "password", "is_admin", "is_active", "login_at"}}
		if !c.state.notFound {
			accounts := c.state.candidates
			if accounts == nil {
				accounts = []loginTestAccount{{id: 7, email: loginTestEmail, password: c.state.password, active: c.state.active, loginAt: c.state.loginAt}}
			}
			if regexp.MustCompile(`(?i)ORDER BY (?:"auth_users"\.)?"?id"?(?: ASC)?(?:\s|$|,)`).MatchString(query) {
				accounts = append([]loginTestAccount(nil), accounts...)
				sort.SliceStable(accounts, func(i, j int) bool { return accounts[i].id < accounts[j].id })
			}
			for _, account := range accounts {
				if !strings.EqualFold(account.email, loginTestEmail) {
					continue
				}
				var password driver.Value
				if account.password != nil {
					password = *account.password
				}
				rows.values = append(rows.values, []driver.Value{account.id, "Login test", account.email, password, false, account.active, account.loginAt})
				if strings.Contains(query, "LIMIT 1") {
					break
				}
			}
		}
		return rows, nil
	case strings.Contains(query, `FROM "auth_user_group_roles"`):
		if len(args) != 1 {
			return nil, errors.New("permission query must specify the authenticated account")
		}
		id, ok := args[0].Value.(int64)
		if !ok {
			return nil, errors.New("permission query did not bind an account ID")
		}
		c.state.permissionAccountID = id
		c.state.permissionQueries++
		return &loginTestRows{columns: []string{"group_id", "role_id"}}, nil
	case strings.Contains(query, `FROM "groups"`):
		c.state.permissionQueries++
		return &loginTestRows{columns: []string{"id"}}, nil
	default:
		return nil, errors.New("unexpected query in login test")
	}
}

func (c *loginTestConnection) ExecContext(_ context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	if !strings.HasPrefix(query, `UPDATE "auth_users" SET `) {
		return nil, errors.New("unexpected write in login test")
	}
	c.state.updates++
	if c.state.failUpdate {
		return nil, errors.New("simulated login persistence failure")
	}
	password, err := loginTestUpdateValue(query, "password", args)
	if err != nil {
		return nil, err
	}
	loginAt, err := loginTestUpdateValue(query, "login_at", args)
	if err != nil {
		return nil, err
	}
	idMatch := regexp.MustCompile(`WHERE "id"\s*=\s*\$(\d+)`).FindStringSubmatch(query)
	if len(idMatch) != 2 {
		return nil, errors.New("login update must select exactly one account by ID")
	}
	idPosition, err := strconv.Atoi(idMatch[1])
	if err != nil || idPosition < 1 || idPosition > len(args) {
		return nil, errors.New("invalid account ID update parameter")
	}
	var ok bool
	c.pendingAccountID, ok = args[idPosition-1].Value.(int64)
	if !ok {
		return nil, errors.New("account ID update did not bind an integer")
	}
	c.pendingPassword, ok = password.(string)
	if !ok {
		return nil, errors.New("password update did not bind a string")
	}
	c.pendingLoginAt, ok = loginAt.(time.Time)
	if !ok {
		return nil, errors.New("login timestamp update did not bind a time")
	}
	return driver.RowsAffected(1), nil
}

func loginTestUpdateValue(query, column string, args []driver.NamedValue) (driver.Value, error) {
	match := regexp.MustCompile(`"` + regexp.QuoteMeta(column) + `"=\$(\d+)`).FindStringSubmatch(query)
	if len(match) != 2 {
		return nil, errors.New("required login field missing from update")
	}
	position, err := strconv.Atoi(match[1])
	if err != nil || position < 1 || position > len(args) {
		return nil, errors.New("invalid login update parameter")
	}
	return args[position-1].Value, nil
}

type loginTestTransaction struct{ connection *loginTestConnection }

func (tx loginTestTransaction) Commit() error {
	password := tx.connection.pendingPassword
	state := tx.connection.state
	if state.candidates == nil {
		if tx.connection.pendingAccountID != 7 {
			return errors.New("login updated an unknown account")
		}
		state.password = &password
		state.loginAt = tx.connection.pendingLoginAt
	} else {
		found := false
		for i := range state.candidates {
			if state.candidates[i].id == tx.connection.pendingAccountID {
				state.candidates[i].password = &password
				state.candidates[i].loginAt = tx.connection.pendingLoginAt
				found = true
				break
			}
		}
		if !found {
			return errors.New("login updated an unknown account")
		}
	}
	state.updatedAccountID = tx.connection.pendingAccountID
	state.commits++
	return nil
}
func (tx loginTestTransaction) Rollback() error { return nil }

type loginTestRows struct {
	columns []string
	values  [][]driver.Value
}

func (r *loginTestRows) Columns() []string { return r.columns }
func (*loginTestRows) Close() error        { return nil }
func (r *loginTestRows) Next(values []driver.Value) error {
	if len(r.values) == 0 {
		return io.EOF
	}
	copy(values, r.values[0])
	r.values = r.values[1:]
	return nil
}
