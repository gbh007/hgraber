package agent

import "time"

type BookToHandle struct {
	ID       int       `json:"id"`
	URL      string    `json:"url"`
	CreateAt time.Time `json:"create_at"`
}

type PageToHandle struct {
	BookID     int       `json:"book_id"`
	PageNumber int       `json:"page_number"`
	CreateAt   time.Time `json:"create_at"`
	BookURL    string    `json:"book_url"`
	PageURL    string    `json:"page_url"`
	Ext        string    `json:"ext"`
}

type BookToUpdate struct {
	ID         int            `json:"id"`
	Name       string         `json:"name"`
	Attributes []Attribute    `json:"attributes,omitempty"`
	Pages      []PageToUpdate `json:"pages,omitempty"`
}

type PageToUpdate struct {
	PageNumber int    `json:"page_number"`
	URL        string `json:"url"`
	Ext        string `json:"ext"`
}

// PageInfoToUpload - данные для обновления страницы при заливке,
// важно: расширение и адрес страницы могут измениться в результате обработки.
type PageInfoToUpload struct {
	BookID     int
	PageNumber int
	URL        string
	Ext        string
}

type Attribute struct {
	// По значениям 1 в 1 domain.Attribute
	Code   string   `json:"code"`
	Parsed bool     `json:"parsed"`
	Values []string `json:"values,omitempty"`
}

type CreateBookResult struct {
	ID          int    `json:"id"`
	URL         string `json:"url"`
	IsDuplicate bool   `json:"is_duplicate"`
	IsHandled   bool   `json:"is_handled"`
	ErrorReason string `json:"error_reason"`
}

type CreateBooksResult struct {
	TotalCount     int64              `json:"total_count"`
	LoadedCount    int64              `json:"loaded_count"`
	DuplicateCount int64              `json:"duplicate_count"`
	ErrorCount     int64              `json:"error_count"`
	NotHandled     []string           `json:"not_handled,omitempty"`
	Details        []CreateBookResult `json:"details"`
}
