-- name: GetAccountForId :one
SELECT account.* FROM account WHERE id = $1 LIMIT 1;
