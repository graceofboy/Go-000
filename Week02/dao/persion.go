package dao

import (
	"database/sql"
	"errors"
)

type DBService struct {
	db *Db
}

func NewDBService() *DBService {
	return &DBService{db: &Db{}}
}

func (this *DBService) FindPersion() ([]model.Persion, error) {
	data, err := this.db.query("select * from persion")
	if err != nil {
		if errors.Is(sql.ErrNoRows) {
			// 当sql.ErrNoRows时，这里把error吞掉可不可以，直接返回一个nil
			return nil, nil
		}
		return nil, errors.Wrap(err, "query persion error: ")
	}
	return data, err
}
