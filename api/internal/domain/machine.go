package domain

import "time"

// Domain Models
type Machine struct {
	ID       int64  `json:"id"`
	Location string `json:"location"`
	Rows     int    `json:"rows"`
	Cols     int    `json:"cols"`
}

type BookMachine struct {
	MachineID int64 `json:"machine_id"`
	BookID    int64 `json:"book_id"`
	Row       int   `json:"row"`
	Col       int   `json:"col"`
}

// Request DTOs
type MachineRequest struct {
	Location string `json:"location"`
	Rows     int    `json:"rows"`
	Cols     int    `json:"cols"`
}

type MachineUpsert struct {
	Location string `json:"location"`
	Rows     int    `json:"rows"`
	Cols     int    `json:"cols"`
}

type MachineRowColUpdateRequest struct {
	Rows int `json:"rows"`
	Cols int `json:"cols"`
}

type MachineLoadRequest struct {
	Books []BookMachine `json:"books"`
}

// Response DTOs
type MachineLocation struct {
	MachineID int64  `json:"machine_id"`
	Location  string `json:"location"`
}

type MachineWithBooks struct {
	Machine Machine      `json:"machine"`
	Books   []LoadedBook `json:"books"`
}

type MachineLoadResponse struct {
	MachineID int64 `json:"machine_id"`
	Count     int   `json:"count"`
}

type MachineMetric struct {
	MachineID int64     `json:"machine_id"`
	Date      time.Time `json:"date"`
	Qr        bool      `json:"qr"`
	Source    string    `json:"source"`
	Admin     bool      `json:"admin"`
}
