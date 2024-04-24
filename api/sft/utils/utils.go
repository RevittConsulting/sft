package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

type ErrResponse struct {
	Error string `json:"error"`
}

func TxBegin(ctx context.Context, pool *pgxpool.Pool) (pgx.Tx, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}
	return tx, nil
}

func TxDefer(tx pgx.Tx, ctx context.Context) {
	err := tx.Rollback(ctx)
	if err != nil {
		if err != pgx.ErrTxClosed {
			fmt.Println("TxDefer error...")
		}
	}

}

func WriteJSON(w http.ResponseWriter, v interface{}) {
	if v == nil {
		WriteErr(w, fmt.Errorf("not found"), http.StatusNotFound)
		return
	}

	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetIndent("", "\t")

	if err := encoder.Encode(v); err != nil {
		WriteErr(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, err := w.Write(buffer.Bytes())
	if err != nil {
		fmt.Errorf("error writing JSON to response: %v", err)
	}
}

func WriteErr(w http.ResponseWriter, err error, code int) {
	if err == nil {
		WriteErr(w, fmt.Errorf(http.StatusText(code)), code)
		return
	}

	switch err.Error() {
	case "unauthorised":
		w.WriteHeader(http.StatusUnauthorized)
		WriteJSON(w, ErrResponse{Error: err.Error()})
	// case *errors.UnauthorizedError:
	// 	w.WriteHeader(e.ResponseCode)
	// 	WriteJSON(w, ErrResponse{Error: e.ClientError})
	default:
		w.WriteHeader(code)
		WriteJSON(w, ErrResponse{Error: err.Error()})
	}
}
