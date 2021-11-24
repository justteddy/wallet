package storage

import "regexp"

var clearQueryWhitespacesRegex = regexp.MustCompile(`\s+`)

func removeExtraWhitespaces(query string) string {
	return clearQueryWhitespacesRegex.ReplaceAllString(query, " ")
}

var (
	queryCreateWallet = removeExtraWhitespaces(`
		INSERT INTO public.wallet(id, balance, created_at)
		VALUES ($1, 0, DEFAULT)`,
	)
)
