package postgres

import (
	"avito-intership/balance"
	"avito-intership/models"
	"database/sql"
	"time"
)

type BalanceRepository struct {
	db *sql.DB
}

func NewBalanceRepository(dbConn *sql.DB) *BalanceRepository {
	return &BalanceRepository{dbConn}
}

type Transaction struct {
	Id       int64
	UserId   int64
	Amount   float32
	TargetId int64
	Type     string
	Time     time.Time
}

func transactionToModel(transaction Transaction) *models.Transaction {
	return &models.Transaction{
		UserId:   transaction.UserId,
		Amount:   transaction.Amount,
		TargetId: transaction.TargetId,
		Type:     transaction.Type,
		Time:     transaction.Time,
	}
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

func (r BalanceRepository) insertTransaction(userId int64, amount float32, target int64, trType string, tx *sql.Tx) error {
	_, err := tx.Exec(`INSERT INTO transactions (user_id, amount, target_id, type) VALUES ($1, $2, $3, $4)`,
		userId, amount, target, trType)
	return err
}

func (r BalanceRepository) ChangeBalance(userId int64, amount float32, productId int64) error {
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
		return err
	}

	var txType string
	if amount < 0 {
		txType = balance.WithdrawType
	} else {
		txType = balance.RefillType
	}

	err = r.insertTransaction(userId, amount, productId, txType, tx)
	if err != nil {
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
		return err
	}

	if currentAmount-amount < 0 {
		err = balance.ErrTooLowBalance
		return err
	}

	_, err = tx.Exec("UPDATE balances SET amount = amount - $1 WHERE id = $2", amount, srcUserId)
	if err != nil {
		return err
	}

	// Если счета у dstUserId нет, то создаем его, иначе обновляем
	_, err = tx.Exec(
		`INSERT INTO balances(id, amount) VALUES ($1, $2) 
		ON CONFLICT(id) DO UPDATE SET amount = balances.amount + EXCLUDED.amount`, dstUserId, amount)
	if err != nil {
		return err
	}

	err = r.insertTransaction(srcUserId, -amount, dstUserId, balance.TransferType, tx)
	if err != nil {
		return err
	}

	err = r.insertTransaction(dstUserId, amount, srcUserId, balance.TransferType, tx)
	if err != nil {
		return err
	}

	return nil
}

func (r BalanceRepository) GetHistory(userId int64, page int64, perPage int64, sort int, desc bool) ([]*models.Transaction, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	orderColumn := "date"
	if sort == balance.SortAmount {
		orderColumn = "amount"
	}

	query := `SELECT user_id, amount, target_id, type, date 
				FROM transactions WHERE user_id = $1 ORDER BY ` + orderColumn
	if desc {
		query += " DESC"
	}
	query += " LIMIT $2 OFFSET $3"

	rows, err := tx.Query(query, userId, perPage, (page-1)*perPage)
	if err != nil {
		return nil, err
	}

	transactions := make([]*models.Transaction, 0)
	for rows.Next() {
		var tx Transaction
		err = rows.Scan(&tx.UserId, &tx.Amount, &tx.TargetId, &tx.Type, &tx.Time)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transactionToModel(tx))
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
