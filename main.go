package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gosimple/slug"
	"github.com/laloloop/best-card-api/pkg/creditcard"
	"io"
	"log"
	"net/http"
	"regexp"
)

var (
	CreditCardRe       = regexp.MustCompile(`^/credit-cards/*$`)
	CreditCardReWithID = regexp.MustCompile(`^/credit-cards/([a-z0-9]+(?:-[a-z0-9]+)+)$`)
)

func main() {

	store := creditcard.NewMemStore()
	creditCardsHandler := NewCreditCardsHandler(store)

	mux := http.NewServeMux()

	mux.Handle("/", &homeHandler{})
	mux.Handle("/credit-cards", creditCardsHandler)
	mux.Handle("/credit-cards/", creditCardsHandler)

	http.ListenAndServe(":8000", mux)
}

// HTTP Handlers
type homeHandler struct{}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my home page"))
}

type CreditCardsHandler struct {
	store creditCardStore
}

func NewCreditCardsHandler(s creditCardStore) *CreditCardsHandler {
	return &CreditCardsHandler{
		store: s,
	}
}

func (h *CreditCardsHandler) CreateCreditCard(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Body: %s\n", b)

	r.Body = io.NopCloser(bytes.NewReader(b))

	var dto creditcard.CreditCardDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	resourceID := slug.Make(dto.Name)

	fmt.Printf("CC_DTO: %v\n", dto)

  daysToPay, err := dto.DaysToPay.Int64()
  if err != nil {
    log.Printf("Failed to get DaysToPay from DTO: %v\n", err)
    InternalServerErrorHandler(w, r)
    return
  }
  
  cutOffDay, err := dto.CutOffDay.Int64()
  if err != nil {
    log.Printf("Failed to get CutOffDay from DTO: %v\n", err)
    InternalServerErrorHandler(w, r)
    return
  }

  cc := creditcard.CreditCard {
    Name: dto.Name,
    DaysToPay: uint8(daysToPay),
    CutOffDay: uint8(cutOffDay),
  }

	if err := h.store.Add(resourceID, cc); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func (h *CreditCardsHandler) ListCreditCards(w http.ResponseWriter, r *http.Request) {
	resources, err := h.store.List()

	jsonBytes, err := json.Marshal(resources)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
func (cc *CreditCardsHandler) GetCreditCard(w http.ResponseWriter, r *http.Request)    {}
func (cc *CreditCardsHandler) UpdateCreditCard(w http.ResponseWriter, r *http.Request) {}
func (cc *CreditCardsHandler) DeleteCreditCard(w http.ResponseWriter, r *http.Request) {}

func (h *CreditCardsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && CreditCardRe.MatchString(r.URL.Path):
		h.CreateCreditCard(w, r)
		return
	case r.Method == http.MethodGet && CreditCardRe.MatchString(r.URL.Path):
		h.ListCreditCards(w, r)
		return
	case r.Method == http.MethodGet && CreditCardReWithID.MatchString(r.URL.Path):
		h.GetCreditCard(w, r)
		return
	case r.Method == http.MethodDelete && CreditCardReWithID.MatchString(r.URL.Path):
		h.DeleteCreditCard(w, r)
		return
	default:
		return
	}
}

// Storage

type creditCardStore interface {
	Add(name string, cc creditcard.CreditCard) error
	Get(name string) (creditcard.CreditCard, error)
	Update(name string, cc creditcard.CreditCard) error
	List() (map[string]creditcard.CreditCard, error)
	Remove(name string) error
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}
