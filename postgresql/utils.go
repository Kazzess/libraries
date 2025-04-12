package psql

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgerrcode"
)

func IsErrUniqueViolation(err error) (*PgError, bool) {
	var pgErr *PgError
	if ok := errors.As(err, &pgErr); ok {
		return pgErr, pgErr.Code == pgerrcode.UniqueViolation
	}

	return nil, false
}

func ParsePgError(err error) error {
	var pgErr *PgError
	if ok := errors.As(err, &pgErr); ok {
		return fmt.Errorf(
			"database error. message:%s, detail:%s, where:%s, sqlstate:%s",
			pgErr.Message,
			pgErr.Detail,
			pgErr.Where,
			pgErr.SQLState(),
		)
	}

	return err
}

func PrettySQL(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}
