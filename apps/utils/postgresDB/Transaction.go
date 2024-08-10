package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

type Transaction struct {
	ID       string   `json:"id" pg:"id,pk"`
	Type     *string  `json:"type" pg:"type,notnull"`
	Amount   *float64 `json:"amount" pg:"amount,notnull"`
	ParentId *uint64  `json:"parent_id" pg:"parent_id"`
}

func (pg *Postgres) InsertSingleTransaction(t *Transaction) error {
	if t == nil {
		return fmt.Errorf("Transaction cannot be nil")
	}
	err := pg.Db.QueryRow(context.Background(), "INSERT INTO transactions (type, amount, parent_id) VALUES ($1, $2, $3)  RETURNING *", t.Type, t.Amount, t.ParentId).Scan(&t.ID, &t.Type, &t.Amount, &t.ParentId)
	fmt.Println(t)
	return err
}

func (pg *Postgres) SelectSingleTransaction(id string) (*Transaction, error) {
	txList, err := pg.SelectMultipleTransaction([]string{id})
	if err != nil {
		return nil, err
	}
	if len(txList) == 0 {
		return nil, fmt.Errorf("Transaction with id %s not found", id)
	}
	return txList[0], nil
}

func (pg *Postgres) SelectChildTransactions(id []string) ([]*Transaction, error) {
	rows, err := pg.Db.Query(context.Background(), "SELECT id, type, amount, parent_id FROM transactions where parent_id = ANY($1)", id)
	if err != nil {
		fmt.Println("Error querying data:", err)
		return nil, err
	}
	defer rows.Close()
	return processRows(rows)
}

func (pg *Postgres) SelectMultipleTransaction(id []string) ([]*Transaction, error) {
	var err error
	var rows pgx.Rows
	if len(id) == 0 {
		rows, err = pg.Db.Query(context.Background(), "SELECT id, type, amount, parent_id FROM transactions")
	} else {
		rows, err = pg.Db.Query(context.Background(), "SELECT id, type, amount, parent_id FROM transactions where id = ANY($1)", id)
	}
	if err != nil {
		fmt.Println("Error querying data:", err)
		return nil, err
	}
	defer rows.Close()
	return processRows(rows)
}

func processRows(rows pgx.Rows) ([]*Transaction, error) {
	txList := []*Transaction{}
	for rows.Next() {
		var tx Transaction
		err := rows.Scan(&tx.ID, &tx.Type, &tx.Amount, &tx.ParentId)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		if err := rows.Err(); err != nil {
			fmt.Println("Error iterating over rows:", err)
			continue
		}
		txList = append(txList, &tx)
	}
	return txList, nil
}

func (pg *Postgres) SelectAllTransaction() ([]*Transaction, error) {
	return pg.SelectMultipleTransaction(nil)
}

func (pg *Postgres) UpdateSingleTransaction(patch *Transaction) error {

	// Build the SQL command dynamically
	var setClauses []string
	var args []interface{}
	argIndex := 1
	hasChange := false

	if patch.Type != nil {
		setClauses = append(setClauses, fmt.Sprintf("type = $%d", argIndex))
		args = append(args, *patch.Type)
		argIndex++
		hasChange = true
	}
	if patch.Amount != nil {
		setClauses = append(setClauses, fmt.Sprintf("amount = $%d", argIndex))
		args = append(args, *patch.Amount)
		argIndex++
		hasChange = true
	}
	if patch.ParentId != nil {
		setClauses = append(setClauses, fmt.Sprintf("parent_id = $%d", argIndex))
		args = append(args, *patch.ParentId)
		argIndex++
		hasChange = true
	}

	if !hasChange {
		fmt.Println("No fields to update")
		return nil
	}

	setClause := strings.Join(setClauses, ", ")
	updateSQL := fmt.Sprintf("UPDATE transactions SET %s WHERE id = $%d returning *", setClause, argIndex)
	args = append(args, patch.ID)

	// Execute the SQL command
	err := pg.Db.QueryRow(context.Background(), updateSQL, args...).Scan(&patch.ID, &patch.Type, &patch.Amount, &patch.ParentId)
	if err != nil {
		return err
	}

	fmt.Println("Data updated successfully")
	return nil
}

func (pg *Postgres) SelectTransactionFromType(transactionType string) ([]*Transaction, error) {
	rows, err := pg.Db.Query(context.Background(), "SELECT id, type, amount, parent_id FROM transactions where type = $1", transactionType)
	if err != nil {
		fmt.Println("Error querying data:", err)
		return nil, err
	}
	defer rows.Close()
	return processRows(rows)
}
