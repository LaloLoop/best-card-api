package main

import (
  "bytes"
  "io"
  "net/http"
  "net/http/httptest"
  "os"
  "testing"

  "github.com/stretchr/testify/assert"
  "github.com/laloloop/best-card-api/pkg/creditcard"
)

func readTestData(t *testing.T, name string) []byte {
  t.Helper()
  content, err := os.ReadFile("./testdata/" + name)
  if err != nil {
    t.Errorf("Could not read %v", name)
  }

  return content
}

func TestCreditCardHandlersCRUD_Integration(t *testing.T) {
  // Create a MemStore and Recipe Handler
  store := creditcard.NewMemStore()
  creditCardsHandler := NewCreditCardsHandler(store)
  
  // Testdata
  masterCard := readTestData(t, "master_card.json")
  masterCardReader := bytes.NewReader(masterCard)

  // CREATE - add new card
  req := httptest.NewRequest(http.MethodPost, "/credit-cards", masterCardReader)
  w := httptest.NewRecorder()
  creditCardsHandler.ServeHTTP(w, req)

  res := w.Result()
  defer res.Body.Close()
  assert.Equal(t, 200, res.StatusCode)

  saved, _ := store.List()
  assert.Len(t, saved, 1)

  // GET - find the record we just added

  req = httptest.NewRequest(http.MethodGet, "/credit-cards/master-card", nil)
  w = httptest.NewRecorder()
  creditCardsHandler.ServeHTTP(w, req)

  res = w.Result()
  defer res.Body.Close()
  assert.Equal(t, 200, res.StatusCode)

  data, err := io.ReadAll(res.Body)
  if err != nil {
    t.Errorf("unexpected error: %v", err)
  }

  assert.JSONEq(t, string(masterCard), string(data))
}
