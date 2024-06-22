package models

import "database/sql"

type Wrapper struct {
	Id        int
	Type      string
	MaxWeight sql.NullInt64
	Markup    int
}

func NewWrapper(Id int, Type string, MaxWeight sql.NullInt64, Markup int) Wrapper {
	return Wrapper{
		Id:        Id,
		Type:      Type,
		MaxWeight: MaxWeight,
		Markup:    Markup,
	}
}
