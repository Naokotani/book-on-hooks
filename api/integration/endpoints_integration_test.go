//go:build integration

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"booksonhooks.ca/internal/app"
	"booksonhooks.ca/internal/domain"
	"booksonhooks.ca/internal/logger"
	"booksonhooks.ca/internal/repository"
	"booksonhooks.ca/internal/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

func testDBURL() string {
	if v := os.Getenv("TEST_DATABASE_URL"); v != "" {
		return v
	}
	return "postgres://postgres:secret@localhost:5433/books_test?sslmode=disable"
}

func connectTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	ctx := context.Background()
	url := testDBURL()
	var conn *pgxpool.Pool
	var err error
	for i := 0; i < 30; i++ {
		conn, err = pgxpool.New(ctx, url)
		if err == nil {
			if pingErr := conn.Ping(ctx); pingErr == nil {
				return conn
			} else {
				err = pingErr
				conn.Close()
			}
		}
		time.Sleep(500 * time.Millisecond)
	}
	t.Fatalf("failed to connect to test database %q: %v", url, err)
	return nil
}

func ensureSchema(t *testing.T, conn *pgxpool.Pool) {
	t.Helper()

	ctx := context.Background()
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS book (
			id BIGSERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			author TEXT NOT NULL,
			summary TEXT NOT NULL,
			image TEXT NOT NULL,
			price TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS machine (
			id BIGSERIAL PRIMARY KEY,
			location TEXT NOT NULL,
			rows INTEGER NOT NULL,
			columns INTEGER NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS book_machine (
			machine_id BIGINT NOT NULL REFERENCES machine(id) ON DELETE CASCADE,
			book_id BIGINT NOT NULL REFERENCES book(id) ON DELETE CASCADE,
			row INTEGER NOT NULL,
			col INTEGER NOT NULL,
			PRIMARY KEY(machine_id, book_id)
		);`,
		`DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint WHERE conname = 'uq_book_machine_machine_id_row_col'
			) THEN
				ALTER TABLE book_machine
				ADD CONSTRAINT uq_book_machine_machine_id_row_col UNIQUE (machine_id, row, col);
			END IF;
		END $$;`,
	}

	for _, stmt := range stmts {
		if _, err := conn.Exec(ctx, stmt); err != nil {
			t.Fatalf("failed to ensure schema: %v", err)
		}
	}
}

func truncateAll(t *testing.T, conn *pgxpool.Pool) {
	t.Helper()
	if _, err := conn.Exec(context.Background(), `TRUNCATE TABLE book_machine, book, machine RESTART IDENTITY CASCADE;`); err != nil {
		t.Fatalf("failed to truncate tables: %v", err)
	}
}

func newTestServer(t *testing.T) (*httptest.Server, *pgxpool.Pool) {
	t.Helper()

	conn := connectTestDB(t)
	ensureSchema(t, conn)
	truncateAll(t, conn)

	log := logger.NewLogger("warn")
	db := &repository.Database{Db: conn, Q: sqlc.New(conn)}
	srv := app.CreateApp(strPtr(":0"), &log, db)
	ts := httptest.NewServer(srv.Handler)

	t.Cleanup(func() {
		ts.Close()
		conn.Close()
	})

	return ts, conn
}

func strPtr(v string) *string { return &v }

func seedBook(t *testing.T, conn *pgxpool.Pool, title string) int64 {
	t.Helper()
	var id int64
	err := conn.QueryRow(context.Background(),
		`INSERT INTO book (title, author, summary, image, price)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		title, "Author", "Summary", "seed.jpg", "$9.99",
	).Scan(&id)
	if err != nil {
		t.Fatalf("failed to seed book: %v", err)
	}
	return id
}

func reqJSON(t *testing.T, method, url string, body any) *http.Response {
	t.Helper()
	var reader *bytes.Reader
	if body == nil {
		reader = bytes.NewReader(nil)
	} else {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal request body: %v", err)
		}
		reader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("http request failed: %v", err)
	}
	return resp
}

func decodeJSON[T any](t *testing.T, resp *http.Response, out *T) {
	t.Helper()
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		t.Fatalf("decode response json: %v", err)
	}
}

func TestBookCRUD_Integration(t *testing.T) {
	ts, conn := newTestServer(t)
	bookID := seedBook(t, conn, "Seed Book")

	resp := reqJSON(t, http.MethodGet, fmt.Sprintf("%s/api/books/%d", ts.URL, bookID), nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET book expected 200, got %d", resp.StatusCode)
	}
	var got domain.Book
	decodeJSON(t, resp, &got)
	if got.ID != bookID {
		t.Fatalf("expected book id %d, got %d", bookID, got.ID)
	}

	resp = reqJSON(t, http.MethodPatch, fmt.Sprintf("%s/api/books/%d", ts.URL, bookID), map[string]any{
		"title":   "Updated Title",
		"author":  "Updated Author",
		"summary": "Updated Summary",
		"price":   "$19.99",
	})
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("PATCH book expected 200, got %d", resp.StatusCode)
	}
	var updated domain.Book
	decodeJSON(t, resp, &updated)
	if updated.Title != "Updated Title" {
		t.Fatalf("expected updated title, got %q", updated.Title)
	}

	resp = reqJSON(t, http.MethodDelete, fmt.Sprintf("%s/api/books/%d", ts.URL, bookID), nil)
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("DELETE book expected 204, got %d", resp.StatusCode)
	}
	_ = resp.Body.Close()

	resp = reqJSON(t, http.MethodGet, fmt.Sprintf("%s/api/books/%d", ts.URL, bookID), nil)
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("GET deleted book expected 404, got %d", resp.StatusCode)
	}
	_ = resp.Body.Close()
}

func TestMachineCRUDAndLoad_Integration(t *testing.T) {
	ts, conn := newTestServer(t)
	book1 := seedBook(t, conn, "Book One")
	book2 := seedBook(t, conn, "Book Two")

	resp := reqJSON(t, http.MethodPost, ts.URL+"/api/machines", map[string]any{
		"location": "HQ",
		"rows":     4,
		"cols":     5,
	})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create machine expected 201, got %d", resp.StatusCode)
	}
	var created struct {
		ID int64 `json:"id"`
	}
	decodeJSON(t, resp, &created)
	if created.ID <= 0 {
		t.Fatalf("expected machine id > 0, got %d", created.ID)
	}

	resp = reqJSON(t, http.MethodGet, fmt.Sprintf("%s/api/machines/%d", ts.URL, created.ID), nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("get machine expected 200, got %d", resp.StatusCode)
	}
	var machine struct {
		ID       int64  `json:"id"`
		Location string `json:"location"`
		Rows     int    `json:"rows"`
		Cols     int    `json:"cols"`
	}
	decodeJSON(t, resp, &machine)
	if machine.Location != "HQ" {
		t.Fatalf("unexpected machine location %q", machine.Location)
	}

	resp = reqJSON(t, http.MethodPatch, fmt.Sprintf("%s/api/machines/%d", ts.URL, created.ID), map[string]any{
		"location": "Warehouse",
		"rows":     6,
		"cols":     7,
	})
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("update machine expected 200, got %d", resp.StatusCode)
	}
	var updatedMachine struct {
		Location string `json:"location"`
	}
	decodeJSON(t, resp, &updatedMachine)
	if updatedMachine.Location != "Warehouse" {
		t.Fatalf("expected updated location, got %q", updatedMachine.Location)
	}

	resp = reqJSON(t, http.MethodPut, fmt.Sprintf("%s/api/machines/%d/books", ts.URL, created.ID), map[string]any{
		"books": []map[string]any{
			{"book_id": book1, "row": 0, "col": 0},
			{"book_id": book2, "row": 0, "col": 1},
		},
	})
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("load machine expected 200, got %d", resp.StatusCode)
	}
	_ = resp.Body.Close()

	resp = reqJSON(t, http.MethodGet, fmt.Sprintf("%s/api/machines/%d/books", ts.URL, created.ID), nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("get machine with books expected 200, got %d", resp.StatusCode)
	}
	var withBooks struct {
		Machine domain.Machine `json:"machine"`
		Books   []domain.Book  `json:"books"`
	}
	decodeJSON(t, resp, &withBooks)
	if withBooks.Machine.ID != created.ID {
		t.Fatalf("expected machine id %d, got %d", created.ID, withBooks.Machine.ID)
	}
	if len(withBooks.Books) != 2 {
		t.Fatalf("expected 2 books in machine, got %d", len(withBooks.Books))
	}

	resp = reqJSON(t, http.MethodDelete, fmt.Sprintf("%s/api/machines/%d/books", ts.URL, created.ID), nil)
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("clear machine books expected 204, got %d", resp.StatusCode)
	}
	_ = resp.Body.Close()

	resp = reqJSON(t, http.MethodGet, fmt.Sprintf("%s/api/machines/%d/books", ts.URL, created.ID), nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("get machine books after clear expected 200, got %d", resp.StatusCode)
	}
	var cleared struct {
		Books []domain.Book `json:"books"`
	}
	decodeJSON(t, resp, &cleared)
	if len(cleared.Books) != 0 {
		t.Fatalf("expected 0 books after clear, got %d", len(cleared.Books))
	}

	resp = reqJSON(t, http.MethodDelete, fmt.Sprintf("%s/api/machines/%d", ts.URL, created.ID), nil)
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("delete machine expected 204, got %d", resp.StatusCode)
	}
	_ = resp.Body.Close()

	resp = reqJSON(t, http.MethodGet, fmt.Sprintf("%s/api/machines/%d", ts.URL, created.ID), nil)
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("get deleted machine expected 404, got %d", resp.StatusCode)
	}
	_ = resp.Body.Close()
}

func TestMain(m *testing.M) {
	// Run from module root (api/) so template/static relative paths resolve.
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		os.Exit(1)
	}
	root := filepath.Dir(filepath.Dir(filename))
	err := os.Chdir(root)
	if err != nil {
		os.Exit(1)
	}
	os.Exit(m.Run())
}
