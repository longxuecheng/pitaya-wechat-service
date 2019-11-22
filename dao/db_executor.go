package dao

import "database/sql"

type TxFunc func(input interface{}, tx *sql.Tx) (interface{}, error)

type TxExecutor struct {
	tx    *sql.Tx
	chain []TxFunc
}

func (t *TxExecutor) AppendFunc(f TxFunc) {
	t.chain = append(t.chain, f)
}

func (t *TxExecutor) Execute() error {
	var err error
	var in interface{}
	var out interface{}
	defer func() {
		if err != nil {
			t.tx.Rollback()
		} else {
			t.tx.Commit()
		}
	}()
	for _, f := range t.chain {
		out, err = f(in, t.tx)
		if err != nil {
			break
		}
		in = out
	}
	return err
}
