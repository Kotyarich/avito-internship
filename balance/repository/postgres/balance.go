package postgres

import (
	"avito-intership/balance"
	"database/sql"
	"log"
)

type BalanceRepository struct {
	db *sql.DB
}

func NewBalanceRepository(dbConn *sql.DB) *BalanceRepository {
	return &BalanceRepository{dbConn}
}

func (r BalanceRepository) GetBalance(userId int64) (float32, error) {
	var currentAmount float32

	tx, err := r.db.Begin()
	if err != nil {
		return currentAmount, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	row := tx.QueryRow("SELECT amount FROM balances WHERE id = $1", userId)
	err = row.Scan(&currentAmount)
	// Предполагается, что отсутствие записи в таблице означает нулевой баланс, а не ошибку
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	return currentAmount, nil
}

func (r BalanceRepository) ChangeBalance(userId int64, amount float32) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	var currentAmount float32
	row := tx.QueryRow("SELECT amount FROM balances WHERE id = $1 FOR UPDATE", userId)
	err = row.Scan(&currentAmount)
	if err != nil && err != sql.ErrNoRows {
		log.Println("1: ", err.Error())
		return err
	}

	if currentAmount+amount < 0 {
		err = balance.ErrTooLowBalance
		return err
	}
	// Если счета у пользователя нет, то создаем его, иначе обновляем
	_, err = tx.Exec(
		`INSERT INTO balances(id, amount) VALUES ($1, $2) 
		ON CONFLICT(id) DO UPDATE SET amount = balances.amount + EXCLUDED.amount`, userId, amount)
	if err != nil {
		log.Println("2: ", err.Error())
		return err
	}

	return nil
}

/* Перевод денег от пользователя srcUserId пользователю dstUserId
   amount - положительное количество переводимых денег */
func (r BalanceRepository) TransferMoney(srcUserId int64, dstUserId int64, amount float32) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	// Проверяем, что у пользователя srcUserId достаточно денег для перевода
	var currentAmount float32
	row := tx.QueryRow("SELECT amount FROM balances WHERE id = $1 FOR UPDATE", srcUserId)
	err = row.Scan(&currentAmount)
	if err != nil && err != sql.ErrNoRows {
		log.Println("3: ", err.Error())
		return err
	}

	if currentAmount-amount < 0 {
		err = balance.ErrTooLowBalance
		return err
	}

	_, err = tx.Exec("UPDATE balances SET amount = amount - $1 WHERE id = $2", amount, srcUserId)
	if err != nil {
		log.Println("4: ", err.Error())
		return err
	}

	// Если счета у dstUserId нет, то создаем его, иначе обновляем
	_, err = tx.Exec(
		`INSERT INTO balances(id, amount) VALUES ($1, $2) 
		ON CONFLICT(id) DO UPDATE SET amount = balances.amount + EXCLUDED.amount`, dstUserId, amount)
	if err != nil {
		log.Println("5: ", err.Error())
		return err
	}

	return nil
}
