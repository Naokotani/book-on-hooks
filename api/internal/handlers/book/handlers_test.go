package book

import (
	"booksonhooks.ca/internal/logger"
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"booksonhooks.ca/internal/domain"
	"github.com/julienschmidt/httprouter"
)

// --- Mocks ---

type mockBookRepo struct {
	getBooksFn         func() ([]domain.Book, error)
	getBookByIDFn      func(id int64) (*domain.Book, error)
	getBookLocationsFn func(bookID int64) (*domain.BookLocation, error)
	insertBookFn       func(book *domain.Book, file multipart.File, header *multipart.FileHeader) (int64, error)
	updateBookFn       func(book *domain.Book) error
	updateBookImageFn  func(id int64, file multipart.File, header *multipart.FileHeader) (string, error)
	deleteBookFn       func(id int64) (string, error)
	insertBookMetricFn func(bookID, machineID int64, qr bool, source string) (int64, error)
}

func (m *mockBookRepo) GetBooks(ctx context.Context) ([]domain.Book, error) {
	return m.getBooksFn()
}

func (m *mockBookRepo) GetBookByID(ctx context.Context, id int64) (*domain.Book, error) {
	return m.getBookByIDFn(id)
}

func (m *mockBookRepo) GetBookLocations(ctx context.Context, bookID int64) (*domain.BookLocation, error) {
	if m.getBookLocationsFn == nil {
		return &domain.BookLocation{}, nil
	}
	return m.getBookLocationsFn(bookID)
}

func (m *mockBookRepo) InsertBook(ctx context.Context, book *domain.Book, file multipart.File, header *multipart.FileHeader) (int64, error) {
	if m.insertBookFn == nil {
		return 0, nil
	}
	return m.insertBookFn(book, file, header)
}

func (m *mockBookRepo) UpdateBook(ctx context.Context, book *domain.Book) error {
	if m.updateBookFn == nil {
		return nil
	}
	return m.updateBookFn(book)
}

func (m *mockBookRepo) UpdateBookImage(ctx context.Context, id int64, file multipart.File, header *multipart.FileHeader) (string, error) {
	if m.updateBookImageFn == nil {
		return "", nil
	}
	return m.updateBookImageFn(id, file, header)
}

func (m *mockBookRepo) DeleteBook(ctx context.Context, id int64) (string, error) {
	if m.deleteBookFn == nil {
		return "", nil
	}
	return m.deleteBookFn(id)
}

func (m *mockBookRepo) InsertBookMetric(ctx context.Context, bookID, machineID int64, qr bool, source string) (int64, error) {
	if m.insertBookMetricFn == nil {
		return 0, nil
	}
	return m.insertBookMetricFn(bookID, machineID, qr, source)
}

// --- Tests ---

func TestGetBooksHandler(t *testing.T) {
	mockRepo := &mockBookRepo{
		getBooksFn: func() ([]domain.Book, error) {
			return []domain.Book{{Title: "Book 1"}, {Title: "Book 2"}}, nil
		},
	}

	logger := logger.NewLogger("info")
	h := New(&logger, mockRepo)

	req := httptest.NewRequest("GET", "/api/books", nil)
	rr := httptest.NewRecorder()

	h.GetBooks(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestGetBookHandler(t *testing.T) {
	mockRepo := &mockBookRepo{
		getBooksFn: func() ([]domain.Book, error) { return []domain.Book{}, nil },
		getBookByIDFn: func(id int64) (*domain.Book, error) {
			return &domain.Book{ID: id, Title: "Book 1"}, nil
		},
	}

	logger := logger.NewLogger("info")
	h := New(&logger, mockRepo)

	req := httptest.NewRequest("GET", "/api/books/book/12", nil)
	params := httprouter.Params{{Key: "id", Value: "12"}}
	req = req.WithContext(context.WithValue(req.Context(), httprouter.ParamsKey, params))
	rr := httptest.NewRecorder()

	h.GetBook(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var resp domain.Book
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.ID != 12 {
		t.Fatalf("expected id=12, got %d", resp.ID)
	}
}

func TestGetBookLocationsHandler_ValidQueryString(t *testing.T) {
	mockRepo := &mockBookRepo{
		insertBookMetricFn: func(bookID, machineID int64, qr bool, source string) (int64, error) {
			if bookID != 7 || machineID != 9 || !qr || source != "location-grid" {
				t.Fatalf("unexpected metric payload bookID=%d machineID=%d qr=%t source=%q", bookID, machineID, qr, source)
			}
			return 1, nil
		},
		getBookLocationsFn: func(bookID int64) (*domain.BookLocation, error) {
			return &domain.BookLocation{BookID: bookID}, nil
		},
	}

	logger := logger.NewLogger("info")
	h := New(&logger, mockRepo)

	req := httptest.NewRequest("GET", "/api/books/summary/7?is_qr=true&machine=9&source=location-grid", nil)
	params := httprouter.Params{{Key: "id", Value: "7"}}
	req = req.WithContext(context.WithValue(req.Context(), httprouter.ParamsKey, params))
	rr := httptest.NewRecorder()

	h.GetBookSummary(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestGetBookLocationsHandler_InvalidQueryString(t *testing.T) {
	mockRepo := &mockBookRepo{
		insertBookMetricFn: func(bookID, machineID int64, qr bool, source string) (int64, error) {
			t.Fatal("InsertBookMetric should not be called for invalid query params")
			return 0, nil
		},
		getBookLocationsFn: func(bookID int64) (*domain.BookLocation, error) {
			return &domain.BookLocation{BookID: bookID}, nil
		},
	}

	logger := logger.NewLogger("info")
	h := New(&logger, mockRepo)

	req := httptest.NewRequest("GET", "/api/books/summary/7?is_qr=nope&machine=9", nil)
	params := httprouter.Params{{Key: "id", Value: "7"}}
	req = req.WithContext(context.WithValue(req.Context(), httprouter.ParamsKey, params))
	rr := httptest.NewRecorder()

	h.GetBookSummary(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestGetBookLocationsHandler_InvalidSource(t *testing.T) {
	mockRepo := &mockBookRepo{
		insertBookMetricFn: func(bookID, machineID int64, qr bool, source string) (int64, error) {
			t.Fatal("InsertBookMetric should not be called for invalid source")
			return 0, nil
		},
		getBookLocationsFn: func(bookID int64) (*domain.BookLocation, error) {
			return &domain.BookLocation{BookID: bookID}, nil
		},
	}

	logger := logger.NewLogger("info")
	h := New(&logger, mockRepo)

	req := httptest.NewRequest("GET", "/api/books/summary/7?is_qr=true&machine=9&source=bad/source", nil)
	params := httprouter.Params{{Key: "id", Value: "7"}}
	req = req.WithContext(context.WithValue(req.Context(), httprouter.ParamsKey, params))
	rr := httptest.NewRecorder()

	h.GetBookSummary(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestGetBookLocationsHandler_MetricInsertErrorDoesNotFailRequest(t *testing.T) {
	mockRepo := &mockBookRepo{
		insertBookMetricFn: func(bookID, machineID int64, qr bool, source string) (int64, error) {
			return 0, errors.New("insert failed")
		},
		getBookLocationsFn: func(bookID int64) (*domain.BookLocation, error) {
			return &domain.BookLocation{BookID: bookID}, nil
		},
	}

	logger := logger.NewLogger("info")
	h := New(&logger, mockRepo)

	req := httptest.NewRequest("GET", "/api/books/summary/7?is_qr=true&machine=9&source=location-grid", nil)
	params := httprouter.Params{{Key: "id", Value: "7"}}
	req = req.WithContext(context.WithValue(req.Context(), httprouter.ParamsKey, params))
	rr := httptest.NewRecorder()

	h.GetBookSummary(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}
