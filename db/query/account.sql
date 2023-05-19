-- name: CreateAccount :one
INSERT INTO accounts (
  id,
  first_name,
  email,
  phone,
  birth_date,
  gender,
  show_me,
  university,
  nsfw,
  ethnicity,
  interests
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: CheckAccountExists :one
SELECT * FROM accounts
WHERE phone =$1 or email = $2 LIMIT 1;

-- name: UpdateAccount :one
UPDATE accounts
SET
  verify_yourself = COALESCE(sqlc.narg(verify_yourself),verify_yourself),
  about_me = COALESCE(sqlc.narg(about_me),about_me),
  interests = COALESCE(sqlc.narg(interests),interests),
  gender = COALESCE(sqlc.narg(gender),gender),
  time_zone = COALESCE(sqlc.narg(time_zone),time_zone),
  ethnicity = COALESCE(sqlc.narg(ethnicity),ethnicity),
  nsfw = COALESCE(sqlc.narg(nsfw),nsfw),
  picture = COALESCE(sqlc.narg(picture),picture)
WHERE
  id = sqlc.arg(id)
RETURNING *;


-- name: DiscoverAccountsWithFilter :many
SELECT
  id,
  first_name,
  birth_date,
  picture
FROM accounts 
WHERE gender = $1 
    AND ethnicity = $2 
    AND birth_date <= $3
    AND birth_date >= $4
LIMIT $5
OFFSET $6;

-- name: DiscoverdAccountsMetadata :one
SELECT 
	count(*)
FROM 
	accounts
WHERE gender = $1 
    AND ethnicity = $2 
    AND birth_date <= $3
    AND birth_date >= $4;

-- name: DeleteAccount :one
DELETE FROM accounts WHERE id = $1 RETURNING id;


-- name: AddAccountPictures :one
UPDATE accounts
SET 
picture = sqlc.arg(picture)
WHERE id = $1
RETURNING picture;

-- name: GetAccountPictures :one
SELECT picture FROM accounts WHERE id=$1;



