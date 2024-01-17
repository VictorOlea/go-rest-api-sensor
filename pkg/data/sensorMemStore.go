package data

import "errors"

var NotFoundErr = errors.New("Not Found")

type MemStore struct {
	list map[string]Sensor
}

func NewMemStore() *MemStore {
	list := make(map[string]Sensor)
	return &MemStore{
		list,
	}
}

func (m MemStore) Add(name string, sensor Sensor) error {
	m.list[name] = sensor
	return nil
}

func (m MemStore) Get(name string) (Sensor, error) {
	if val, ok := m.list[name]; ok {
		return val, nil
	}

	return Sensor{}, NotFoundErr
}

func (m MemStore) List() (map[string]Sensor, error) {
	return m.list, nil
}