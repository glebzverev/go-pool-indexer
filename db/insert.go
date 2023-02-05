package db

import (
	"fmt"

	"github.com/go-pg/pg/v10"
)

func (token *Token) SafetyInsert(db *pg.DB) error {
	inserted, err := db.Model(token).
		Column("network").
		Where("address = ?address").
		OnConflict("DO NOTHING"). // OnConflict is optional
		Returning("network").
		SelectOrInsert()
	if err != nil {
		return err
	}
	fmt.Println(inserted, token)
	return nil
}

func (dex *Dex) SafetyInsert(db *pg.DB) error {
	inserted, err := db.Model(dex).
		Column("network").
		Where("factory_address = ?factory_address").
		OnConflict("DO NOTHING"). // OnConflict is optional
		Returning("network").
		SelectOrInsert()
	if err != nil {
		return err
	}
	fmt.Println(inserted, dex)
	return nil
}

func (pool *Pool) SafetyInsert(db *pg.DB) error {
	inserted, err := db.Model(pool).
		Column("network").
		Where("address = ?address").
		OnConflict("DO NOTHING"). // OnConflict is optional
		Returning("network").
		SelectOrInsert()
	if err != nil {
		return err
	}
	fmt.Println(inserted, pool)
	return nil
}
