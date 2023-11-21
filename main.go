package main

import (
  "net/http"
  "regexp"
)

var (
  CreditCardRe = regexp.MustCompile(`^/credit-cards/*$`)
  CreditCardReWithID = regexp.MustCompile(`^/credit-cards/([a-z0-9]+(?:-[a-z0-9]+)+)$`)
)

func main() {
  mux := http.NewServeMux()

  mux.Handle("/", &homeHandler{})
  mux.Handle("/credit-cards", &CreditCardsHandler{})
  mux.Handle("/credit-cards/", &CreditCardsHandler{})

  http.ListenAndServe(":8000", mux)
}

type homeHandler struct{}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("This is my home page"))
}

type CreditCardsHandler struct{}

func (cc *CreditCardsHandler) CreateCreditCard(w http.ResponseWriter, r *http.Request) {}
func (cc *CreditCardsHandler) ListCreditCards(w http.ResponseWriter, r *http.Request) {}
func (cc *CreditCardsHandler) GetCreditCard(w http.ResponseWriter, r *http.Request) {}
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
