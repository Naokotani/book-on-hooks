package machine

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"booksonhooks.ca/internal/domain"
	"booksonhooks.ca/internal/logger"
	"github.com/julienschmidt/httprouter"
)

type mockMachineRepo struct {
	getMachinesFn           func() ([]domain.Machine, error)
	getMachineByIDFn        func(id int64) (*domain.Machine, error)
	getMachineWithBooksFn   func(id int64) (*domain.MachineWithBooks, error)
	hasOutsideBoundsFn      func(machineID int64, rows int, cols int) (bool, error)
	insertMachineFn         func(machine *domain.Machine) (int64, error)
	insertMachineMetricFn   func(machineID int64, qr bool, source string, admin bool, sessionID string) (int64, error)
	loadMachineFn           func(machineID int64, books []domain.BookMachine) error
	updateMachineFn         func(machine *domain.Machine) error
	updateMachineRowsColsFn func(machineID int64, rows int, cols int) error
	deleteMachineFn         func(id int64) error
	clearMachineBooksFn     func(machineID int64) error
}

func (m *mockMachineRepo) GetMachines(ctx context.Context) ([]domain.Machine, error) {
	return m.getMachinesFn()
}

func (m *mockMachineRepo) GetMachineById(ctx context.Context, id int64) (*domain.Machine, error) {
	return m.getMachineByIDFn(id)
}

func (m *mockMachineRepo) GetMachineWithBooks(ctx context.Context, id int64) (*domain.MachineWithBooks, error) {
	if m.getMachineWithBooksFn == nil {
		return &domain.MachineWithBooks{}, nil
	}
	return m.getMachineWithBooksFn(id)
}

func (m *mockMachineRepo) MachineHasLoadedBooksOutsideBounds(ctx context.Context, machineID int64, rows int, cols int) (bool, error) {
	if m.hasOutsideBoundsFn == nil {
		return false, nil
	}
	return m.hasOutsideBoundsFn(machineID, rows, cols)
}

func (m *mockMachineRepo) InsertMachine(ctx context.Context, machine *domain.Machine) (int64, error) {
	return m.insertMachineFn(machine)
}

func (m *mockMachineRepo) InsertMachineMetric(ctx context.Context, machineID int64, qr bool, source string, admin bool, sessionID string) (int64, error) {
	if m.insertMachineMetricFn == nil {
		return 0, nil
	}
	return m.insertMachineMetricFn(machineID, qr, source, admin, sessionID)
}

func (m *mockMachineRepo) LoadMachine(ctx context.Context, machineID int64, books []domain.BookMachine) error {
	return m.loadMachineFn(machineID, books)
}

func (m *mockMachineRepo) UpdateMachine(ctx context.Context, machine *domain.Machine) error {
	if m.updateMachineFn == nil {
		return nil
	}
	return m.updateMachineFn(machine)
}

func (m *mockMachineRepo) UpdateMachineRowsCols(ctx context.Context, machineID int64, rows int, cols int) error {
	if m.updateMachineRowsColsFn == nil {
		return nil
	}
	return m.updateMachineRowsColsFn(machineID, rows, cols)
}

func (m *mockMachineRepo) DeleteMachine(ctx context.Context, id int64) error {
	if m.deleteMachineFn == nil {
		return nil
	}
	return m.deleteMachineFn(id)
}

func (m *mockMachineRepo) ClearMachineBooks(ctx context.Context, machineID int64) error {
	if m.clearMachineBooksFn == nil {
		return nil
	}
	return m.clearMachineBooksFn(machineID)
}

func TestMachineCreateHandler(t *testing.T) {
	mockRepo := &mockMachineRepo{
		getMachinesFn: func() ([]domain.Machine, error) {
			return []domain.Machine{}, nil
		},
		getMachineByIDFn: func(id int64) (*domain.Machine, error) {
			return &domain.Machine{ID: id, Rows: 6, Cols: 5}, nil
		},
		insertMachineFn: func(machine *domain.Machine) (int64, error) {
			if machine.Location != "HQ" || machine.Rows != 3 || machine.Cols != 4 {
				t.Fatalf("unexpected machine payload: %+v", machine)
			}
			return 42, nil
		},
		loadMachineFn: func(machineID int64, books []domain.BookMachine) error {
			return nil
		},
	}

	logger := logger.NewLogger("info")
	h := New(&logger, mockRepo)

	body := []byte(`{"location":"HQ","rows":3,"cols":4}`)
	req := httptest.NewRequest("POST", "/api/machines", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	h.CreateMachine(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rr.Code)
	}
}

func TestLoadMachineHandler(t *testing.T) {
	mockRepo := &mockMachineRepo{
		getMachinesFn: func() ([]domain.Machine, error) {
			return []domain.Machine{}, nil
		},
		getMachineByIDFn: func(id int64) (*domain.Machine, error) {
			return &domain.Machine{ID: id, Rows: 6, Cols: 5}, nil
		},
		insertMachineFn: func(machine *domain.Machine) (int64, error) {
			return 1, nil
		},
		loadMachineFn: func(machineID int64, books []domain.BookMachine) error {
			if machineID != 7 {
				t.Fatalf("expected machineID=7, got %d", machineID)
			}
			if len(books) != 2 {
				t.Fatalf("expected 2 books, got %d", len(books))
			}
			if books[0].BookID != 11 || books[0].Row != 1 || books[0].Col != 2 {
				t.Fatalf("unexpected first book payload: %+v", books[0])
			}
			if books[1].BookID != 12 || books[1].Row != 3 || books[1].Col != 4 {
				t.Fatalf("unexpected second book payload: %+v", books[1])
			}
			return nil
		},
	}

	logger := logger.NewLogger("info")
	h := New(&logger, mockRepo)

	body := []byte(`{"books":[{"book_id":11,"row":1,"col":2},{"book_id":12,"row":3,"col":4}]}`)
	req := httptest.NewRequest("PUT", "/api/machines/7/books", bytes.NewReader(body))
	params := httprouter.Params{{Key: "id", Value: "7"}}
	req = req.WithContext(context.WithValue(req.Context(), httprouter.ParamsKey, params))
	rr := httptest.NewRecorder()

	h.LoadMachine(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var resp struct {
		MachineID int64 `json:"machine_id"`
		Count     int   `json:"count"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response json: %v", err)
	}

	if resp.MachineID != 7 || resp.Count != 2 {
		t.Fatalf("unexpected response payload: %+v", resp)
	}
}

func TestLoadMachineHandler_RejectsDuplicateBookID(t *testing.T) {
	mockRepo := &mockMachineRepo{
		getMachinesFn: func() ([]domain.Machine, error) {
			return []domain.Machine{}, nil
		},
		getMachineByIDFn: func(id int64) (*domain.Machine, error) {
			return &domain.Machine{ID: id, Rows: 6, Cols: 5}, nil
		},
		insertMachineFn: func(machine *domain.Machine) (int64, error) {
			return 1, nil
		},
		loadMachineFn: func(machineID int64, books []domain.BookMachine) error {
			t.Fatal("repo.LoadMachine should not be called for invalid payload")
			return nil
		},
	}

	logger := logger.NewLogger("info")
	h := New(&logger, mockRepo)

	body := []byte(`{"books":[{"book_id":11,"row":1,"col":1},{"book_id":11,"row":1,"col":2}]}`)
	req := httptest.NewRequest("PUT", "/api/machines/7/books", bytes.NewReader(body))
	params := httprouter.Params{{Key: "id", Value: "7"}}
	req = req.WithContext(context.WithValue(req.Context(), httprouter.ParamsKey, params))
	rr := httptest.NewRecorder()

	h.LoadMachine(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", rr.Code)
	}
}

func TestMachineCreateHandler_RejectsInvalidDimensions(t *testing.T) {
	mockRepo := &mockMachineRepo{
		getMachinesFn: func() ([]domain.Machine, error) { return []domain.Machine{}, nil },
		getMachineByIDFn: func(id int64) (*domain.Machine, error) {
			return &domain.Machine{ID: id}, nil
		},
		insertMachineFn: func(machine *domain.Machine) (int64, error) {
			t.Fatal("InsertMachine should not be called for invalid payload")
			return 0, nil
		},
		loadMachineFn: func(machineID int64, books []domain.BookMachine) error { return nil },
	}

	logger := logger.NewLogger("info")
	h := New(&logger, mockRepo)

	body := []byte(`{"location":"HQ","rows":0,"cols":-1}`)
	req := httptest.NewRequest("POST", "/api/machines", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	h.CreateMachine(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", rr.Code)
	}
}

func TestUpdateMachineRowsColsHandler(t *testing.T) {
	mockRepo := &mockMachineRepo{
		getMachineByIDFn: func(id int64) (*domain.Machine, error) {
			return &domain.Machine{ID: id, Location: "HQ", Rows: 3, Cols: 4}, nil
		},
		hasOutsideBoundsFn: func(machineID int64, rows int, cols int) (bool, error) {
			if machineID != 7 || rows != 5 || cols != 6 {
				t.Fatalf("unexpected bounds check payload machineID=%d rows=%d cols=%d", machineID, rows, cols)
			}
			return false, nil
		},
		updateMachineRowsColsFn: func(machineID int64, rows int, cols int) error {
			if machineID != 7 || rows != 5 || cols != 6 {
				t.Fatalf("unexpected row/col update payload machineID=%d rows=%d cols=%d", machineID, rows, cols)
			}
			return nil
		},
	}

	logger := logger.NewLogger("info")
	h := New(&logger, mockRepo)

	body := []byte(`{"rows":5,"cols":6}`)
	req := httptest.NewRequest("PATCH", "/api/machines/7/grid", bytes.NewReader(body))
	params := httprouter.Params{{Key: "id", Value: "7"}}
	req = req.WithContext(context.WithValue(req.Context(), httprouter.ParamsKey, params))
	rr := httptest.NewRecorder()

	h.UpdateMachineRowsCols(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var resp domain.Machine
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response json: %v", err)
	}

	if resp.ID != 7 || resp.Location != "HQ" || resp.Rows != 5 || resp.Cols != 6 {
		t.Fatalf("unexpected response payload: %+v", resp)
	}
}

func TestUpdateMachineRowsColsHandler_RejectsShrinkWhenRemovedSlotsContainBooks(t *testing.T) {
	mockRepo := &mockMachineRepo{
		getMachineByIDFn: func(id int64) (*domain.Machine, error) {
			return &domain.Machine{ID: id, Location: "HQ", Rows: 4, Cols: 4}, nil
		},
		hasOutsideBoundsFn: func(machineID int64, rows int, cols int) (bool, error) {
			if machineID != 7 || rows != 2 || cols != 2 {
				t.Fatalf("unexpected bounds check payload machineID=%d rows=%d cols=%d", machineID, rows, cols)
			}
			return true, nil
		},
		updateMachineRowsColsFn: func(machineID int64, rows int, cols int) error {
			t.Fatal("UpdateMachineRowsCols should not be called when shrink would remove occupied slots")
			return nil
		},
	}

	logger := logger.NewLogger("info")
	h := New(&logger, mockRepo)

	body := []byte(`{"rows":2,"cols":2}`)
	req := httptest.NewRequest("PATCH", "/api/machines/7/grid", bytes.NewReader(body))
	params := httprouter.Params{{Key: "id", Value: "7"}}
	req = req.WithContext(context.WithValue(req.Context(), httprouter.ParamsKey, params))
	rr := httptest.NewRecorder()

	h.UpdateMachineRowsCols(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", rr.Code)
	}

	var resp struct {
		Error       string            `json:"error"`
		FieldErrors map[string]string `json:"field_errors"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response json: %v", err)
	}

	expected := "Cannot shrink machine while removed rows or columns still contain books"
	if resp.FieldErrors["rows"] != expected || resp.FieldErrors["cols"] != expected {
		t.Fatalf("unexpected field errors: %+v", resp.FieldErrors)
	}
}

func TestLoadMachineHandler_RejectsOutOfBoundsSlot(t *testing.T) {
	mockRepo := &mockMachineRepo{
		getMachinesFn: func() ([]domain.Machine, error) { return []domain.Machine{}, nil },
		getMachineByIDFn: func(id int64) (*domain.Machine, error) {
			return &domain.Machine{ID: id, Rows: 2, Cols: 2}, nil
		},
		insertMachineFn: func(machine *domain.Machine) (int64, error) { return 1, nil },
		loadMachineFn: func(machineID int64, books []domain.BookMachine) error {
			t.Fatal("LoadMachine should not be called for out-of-bounds slot")
			return nil
		},
	}

	logger := logger.NewLogger("info")
	h := New(&logger, mockRepo)

	body := []byte(`{"books":[{"book_id":11,"row":2,"col":0}]}`)
	req := httptest.NewRequest("PUT", "/api/machines/7/books", bytes.NewReader(body))
	params := httprouter.Params{{Key: "id", Value: "7"}}
	req = req.WithContext(context.WithValue(req.Context(), httprouter.ParamsKey, params))
	rr := httptest.NewRecorder()

	h.LoadMachine(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", rr.Code)
	}
}

func TestGetMachineWithBooks_RecordsMetric(t *testing.T) {
	mockRepo := &mockMachineRepo{
		getMachineWithBooksFn: func(id int64) (*domain.MachineWithBooks, error) {
			return &domain.MachineWithBooks{
				Machine: domain.Machine{ID: id, Location: "A1", Rows: 2, Cols: 2},
				Books:   []domain.LoadedBook{},
			}, nil
		},
		insertMachineMetricFn: func(machineID int64, qr bool, source string, admin bool, sessionID string) (int64, error) {
			if machineID != 7 || !qr || source != "admin-load" || !admin || sessionID != "" {
				t.Fatalf("unexpected metric payload machineID=%d qr=%t source=%q admin=%t sessionID=%q", machineID, qr, source, admin, sessionID)
			}
			return 1, nil
		},
	}

	logger := logger.NewLogger("info")
	h := New(&logger, mockRepo)

	req := httptest.NewRequest("GET", "/api/machines/7/books?is_qr=true&source=admin-load", nil)
	params := httprouter.Params{{Key: "id", Value: "7"}}
	req = req.WithContext(context.WithValue(req.Context(), httprouter.ParamsKey, params))
	rr := httptest.NewRecorder()

	h.GetMachineWithBooks(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestGetMachineWithBooks_InvalidSource(t *testing.T) {
	mockRepo := &mockMachineRepo{
		getMachineWithBooksFn: func(id int64) (*domain.MachineWithBooks, error) {
			return &domain.MachineWithBooks{
				Machine: domain.Machine{ID: id, Location: "A1", Rows: 2, Cols: 2},
				Books:   []domain.LoadedBook{},
			}, nil
		},
		insertMachineMetricFn: func(machineID int64, qr bool, source string, admin bool, sessionID string) (int64, error) {
			t.Fatal("InsertMachineMetric should not be called for invalid source")
			return 0, nil
		},
	}

	logger := logger.NewLogger("info")
	h := New(&logger, mockRepo)

	req := httptest.NewRequest("GET", "/api/machines/7/books?source=bad/source", nil)
	params := httprouter.Params{{Key: "id", Value: "7"}}
	req = req.WithContext(context.WithValue(req.Context(), httprouter.ParamsKey, params))
	rr := httptest.NewRecorder()

	h.GetMachineWithBooks(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestGetMachineWithBooks_MetricInsertErrorDoesNotFailRequest(t *testing.T) {
	mockRepo := &mockMachineRepo{
		getMachineWithBooksFn: func(id int64) (*domain.MachineWithBooks, error) {
			return &domain.MachineWithBooks{
				Machine: domain.Machine{ID: id, Location: "A1", Rows: 2, Cols: 2},
				Books:   []domain.LoadedBook{},
			}, nil
		},
		insertMachineMetricFn: func(machineID int64, qr bool, source string, admin bool, sessionID string) (int64, error) {
			return 0, errors.New("insert failed")
		},
	}

	logger := logger.NewLogger("info")
	h := New(&logger, mockRepo)

	req := httptest.NewRequest("GET", "/api/machines/7/books?is_qr=true&source=location-grid", nil)
	params := httprouter.Params{{Key: "id", Value: "7"}}
	req = req.WithContext(context.WithValue(req.Context(), httprouter.ParamsKey, params))
	rr := httptest.NewRecorder()

	h.GetMachineWithBooks(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}
