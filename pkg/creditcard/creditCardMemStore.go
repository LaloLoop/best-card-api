package creditcard

import "errors"

var (
  NotFoundErr = errors.New("not found")
)

type MemStore struct {
  list map[string]CreditCard
}

func NewMemStore() *MemStore {
  list := make(map[string]CreditCard)
  return &MemStore{
    list,
  }
}

func (m MemStore) Add(name string, cc CreditCard) error {
  m.list[name] = cc
  return nil
}

func (m MemStore) Get(name string) (CreditCard, error) {
  if val, ok := m.list[name]; ok {
    return val, nil
  }

  return CreditCard{}, NotFoundErr
}

func (m MemStore) List() (map[string]CreditCard, error) {
  return m.list, nil
}

func (m MemStore) Update(name string, cc CreditCard) error {
  if _, ok := m.list[name]; ok {
    m.list[name] = cc
    return nil
  }

  return NotFoundErr
}

func (m MemStore) Remove(name string) error {
  delete(m.list, name)
  return nil
}
