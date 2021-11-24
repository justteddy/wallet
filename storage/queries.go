package storage

import "regexp"

var clearQueryWhitespacesRegex = regexp.MustCompile(`\s+`)

func removeExtraWhitespaces(query string) string {
	return clearQueryWhitespacesRegex.ReplaceAllString(query, " ")
}

var (
	queryInsertWallet = removeExtraWhitespaces(`
		INSERT INTO wallet(id, balance, created_at)
		VALUES ($1, 0, DEFAULT)`,
	)

	queryLockWalletForUpdate = removeExtraWhitespaces(`
		SELECT balance FROM wallet WHERE id = $1 FOR UPDATE`,
	)

	queryInsertOperation = removeExtraWhitespaces(`
		INSERT INTO operations(id, wallet_id, operation_type, amount, created_at)
		VALUES (DEFAULT, $1, $2, $3, DEFAULT)`,
	)

	queryUpdateWallet = removeExtraWhitespaces(`
		UPDATE wallet SET balance = $1 WHERE id = $2`,
	)
)
