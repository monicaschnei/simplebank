package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transactions
type Store interface {
	Querier
	TransferTransaction(ctx context.Context, argument TransferTransactionParams) (TransferTransactionResult, error)
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new Store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// executeTransaction executes a function within a database transaction
func (store *SQLStore) executeTransaction(ctx context.Context, fn func(*Queries) error) error {
	startedTransaction, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(startedTransaction)
	err = fn(q)
	if err != nil {
		if rolbackErr := startedTransaction.Rollback(); rolbackErr != nil {
			return fmt.Errorf("startedTransaction error: %v, rolback error: %v", err, rolbackErr)
		}
		return err
	}
	return startedTransaction.Commit()
}

// TransferTransactionParams contains the input parameters of the transfer transaction
type TransferTransactionParams struct {
	FromAccountID int64 `"json:from_account_id"`
	ToAccountID   int64 `"json:to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTransactionResult is the result of the transfer transaction
type TransferTransactionResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var transactionKey = struct{}{}

// TransferTransaction performs a money transfer from one account to the other
// It creates a transfer record, add account entries, and updates accounts' balance within a single database transaction
func (store *SQLStore) TransferTransaction(ctx context.Context, argument TransferTransactionParams) (TransferTransactionResult, error) {
	var resultTransaction TransferTransactionResult

	err := store.executeTransaction(ctx, func(query *Queries) error {
		var err error

		transactioName := ctx.Value(transactionKey)

		fmt.Println(transactioName, "create transfer")
		resultTransaction.Transfer, err = query.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: argument.FromAccountID,
			ToAccountID:   argument.ToAccountID,
			Amount:        argument.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(transactioName, "create entry 1")
		resultTransaction.FromEntry, err = query.CreateEntry(ctx, CreateEntryParams{
			AccountID: argument.FromAccountID,
			Amount:    -argument.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(transactioName, "create entry2")
		resultTransaction.ToEntry, err = query.CreateEntry(ctx, CreateEntryParams{
			AccountID: argument.ToAccountID,
			Amount:    argument.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(transactioName, "update account1")

		if argument.FromAccountID < argument.ToAccountID {
			resultTransaction.FromAccount, resultTransaction.ToAccount, err = addMoney(ctx, query, argument.FromAccountID, -argument.Amount, argument.ToAccountID, argument.Amount)
		} else {

			resultTransaction.ToAccount, resultTransaction.FromAccount, err = addMoney(ctx, query, argument.ToAccountID, argument.Amount, argument.FromAccountID, -argument.Amount)
		}

		return err
	})

	return resultTransaction, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})

	return
}
