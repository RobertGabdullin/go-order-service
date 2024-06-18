package storage

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

type Storage interface {
	AddOrder(Order) error
	ChangeStatus(int, string) error
	FindOrders([]int) ([]Order, error)
	ListOrders(int) ([]Order, error)
	GetReturns() ([]Order, error)
	UpdateHash(hash string) error
}

type JSONStorage struct {
	fileName string
}

func New(fileName string) (JSONStorage, error) {
	s := JSONStorage{fileName}
	if _, err := os.Stat(s.fileName); errors.Is(err, os.ErrNotExist) {

		if errCreateFile := s.createFile(); errCreateFile != nil {
			return JSONStorage{}, errCreateFile
		}
	}

	return s, nil
}

func (s JSONStorage) AddOrder(ord Order) error {

	records, err := s.readFile()
	if err != nil {
		return err
	}

	for _, elem := range records.Orders {
		if elem.Id == ord.Id {
			return errors.New("order with such id already exists")
		}
	}

	records.Orders = append(records.Orders, ord)

	bWrite, errMarshal := json.MarshalIndent(records, "  ", "  ")
	if errMarshal != nil {
		return errMarshal
	}

	return os.WriteFile(s.fileName, bWrite, 0666)

}

func (s JSONStorage) ChangeStatus(id int, status string) error {

	records, err := s.readFile()
	if err != nil {
		return err
	}

	ok := false

	for i, elem := range records.Orders {
		if elem.Id == id {
			records.Orders[i].Status = status
			if status == "delivered" {
				records.Orders[i].DeliviredAt = time.Now()
			}
			if status == "returned" {
				records.Orders[i].ReturnedAt = time.Now()
			}
			ok = true
		}
	}

	if !ok {
		return errors.New("order with such id does not exist")
	}

	bWrite, errMarshal := json.MarshalIndent(records, "  ", "  ")
	if errMarshal != nil {
		return errMarshal
	}

	return os.WriteFile(s.fileName, bWrite, 0666)

}

func exists(orderID int, orderIDs []int) bool {
	for _, elem := range orderIDs {
		if elem == orderID {
			return true
		}
	}
	return false
}

func (s JSONStorage) FindOrders(ids []int) ([]Order, error) {

	records, err := s.readFile()
	if err != nil {
		return nil, err
	}

	ans := make([]Order, 0)

	for _, elem := range records.Orders {
		if exists(elem.Id, ids) {
			ans = append(ans, elem)
		}
	}

	return ans, nil

}

func (s JSONStorage) ListOrders(recipient int) ([]Order, error) {

	records, err := s.readFile()
	if err != nil {
		return nil, err
	}

	ans := make([]Order, 0)

	for _, elem := range records.Orders {
		if elem.Recipient == recipient {
			ans = append(ans, elem)
		}
	}

	return ans, nil
}

func (s JSONStorage) GetReturns() ([]Order, error) {
	if _, err := os.Stat(s.fileName); errors.Is(err, os.ErrNotExist) {

		if errCreateFile := s.createFile(); errCreateFile != nil {
			return nil, errCreateFile
		}
	}

	b, err := os.ReadFile(s.fileName)
	if err != nil {
		return nil, err
	}

	var records OrdersDTO
	if errUnmarshal := json.Unmarshal(b, &records); errUnmarshal != nil {
		return nil, errUnmarshal
	}

	ans := make([]Order, 0)

	for _, elem := range records.Orders {
		if elem.Status == "returned" {
			ans = append(ans, elem)
		}
	}

	return ans, nil
}

func (s JSONStorage) UpdateHash(hash string) error {

	b, err := os.ReadFile(s.fileName)
	if err != nil {
		return err
	}

	var records OrdersDTO
	if errUnmarshal := json.Unmarshal(b, &records); errUnmarshal != nil {
		return errUnmarshal
	}

	records.CurHash = hash

	bWrite, errMarshal := json.MarshalIndent(records, "  ", "  ")
	if errMarshal != nil {
		return errMarshal
	}

	return os.WriteFile(s.fileName, bWrite, 0666)
}

func (s JSONStorage) createFile() error {
	f, err := os.Create(s.fileName)
	os.WriteFile(s.fileName, []byte("{}"), 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

func (s JSONStorage) readFile() (OrdersDTO, error) {
	b, err := os.ReadFile(s.fileName)
	if err != nil {
		return OrdersDTO{nil, ""}, err
	}

	var records OrdersDTO
	if errUnmarshal := json.Unmarshal(b, &records); errUnmarshal != nil {
		return OrdersDTO{nil, ""}, errUnmarshal
	}

	return records, nil
}
