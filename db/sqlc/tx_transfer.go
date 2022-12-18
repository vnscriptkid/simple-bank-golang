package db

import "context"

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to the other.
// It creates the transfer, add account entries, and update accounts' balance within a database transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// Update user's balance
		// For ref: GET+SET pattern
		/*
			fromAcc, err := q.GetAccountForUpdate(ctx, int64(arg.FromAccountID))

			if err != nil {
				return err
			}

			result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:      arg.FromAccountID,
				Balance: fromAcc.Balance - arg.Amount,
			})

			if err != nil {
				return err
			}

			toAcc, err := q.GetAccount(ctx, int64(arg.ToAccountID))

			if err != nil {
				return err
			}

			result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
				ID:      arg.ToAccountID,
				Balance: toAcc.Balance + arg.Amount,
			})

			if err != nil {
				return err
			}
		*/

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, err = q.UpdateAccountByAmount(ctx, UpdateAccountByAmountParams{
				ID:     arg.FromAccountID,
				Amount: -arg.Amount,
			})

			if err != nil {
				return err
			}

			result.ToAccount, err = q.UpdateAccountByAmount(ctx, UpdateAccountByAmountParams{
				ID:     arg.ToAccountID,
				Amount: arg.Amount,
			})

			if err != nil {
				return err
			}
		} else {
			result.ToAccount, err = q.UpdateAccountByAmount(ctx, UpdateAccountByAmountParams{
				ID:     arg.ToAccountID,
				Amount: arg.Amount,
			})

			if err != nil {
				return err
			}

			result.FromAccount, err = q.UpdateAccountByAmount(ctx, UpdateAccountByAmountParams{
				ID:     arg.FromAccountID,
				Amount: -arg.Amount,
			})

			if err != nil {
				return err
			}
		}

		// TODO: May throw if ID does not exist and returns empty

		return nil
	})

	return result, err
}
