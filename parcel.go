package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	addSQL := `
              INSERT INTO parcel (client, status, address, created_at)
              VALUES (:client, :status, :address, :created_at)
              `
	res, err := s.db.Exec(addSQL,
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))

	if err != nil {
		return 0, fmt.Errorf("add query error: %w", err)
	}

	// верните идентификатор последней добавленной записи
	id, err := res.LastInsertId()

	if err != nil {
		fmt.Printf("last insertion id error: %v", err)
		return 0, err
	}

	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	getByNumberSQL := `
                      SELECT number,
	                         client,
	                         status,
	                         address,
	                         created_at
                      FROM parcel p
                      WHERE p.number = :number
                      `
	row := s.db.QueryRow(getByNumberSQL, sql.Named("number", number))

	// заполните объект Parcel данными из таблицы
	p := Parcel{}
	err := row.Scan(
		&p.Number,
		&p.Client,
		&p.Status,
		&p.Address,
		&p.CreatedAt)

	if err != nil {
		return p, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	getByClientSQL := `
                      SELECT number,
	                         client,
	                         status,
	                         address,
	                         created_at
                      FROM parcel p
                      WHERE p.client = :client
                      `
	rows, err := s.db.Query(getByClientSQL, sql.Named("client", client))
	if err != nil {
		fmt.Printf("last insertion id error: %v", err)
		return nil, err
	}
	defer rows.Close()

	// заполните срез Parcel данными из таблицы
	var res []Parcel
	for rows.Next() {
		p := Parcel{}
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	setStatusSQL := `
                    UPDATE parcel
                    SET status = :status
                    WHERE number = :number
                    `
	_, err := s.db.Exec(setStatusSQL,
		sql.Named("status", status),
		sql.Named("number", number))
	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	setAddressSQL := `
                     UPDATE parcel 
                     SET address = :address 
                     WHERE number = :number AND status = :status
                     `
	_, err := s.db.Exec(setAddressSQL,
		sql.Named("address", address),
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))
	return err
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	deleteSQL := `
	             DELETE FROM parcel
	             WHERE number = :number AND status = :status
	             `
	_, err := s.db.Exec(deleteSQL,
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))
	return err
}
