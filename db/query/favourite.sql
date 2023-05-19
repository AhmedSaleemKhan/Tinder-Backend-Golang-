-- name: AddFavourite :one
INSERT INTO favourites (
  fav_id,
  user_id,
  target_id
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetFavorite :one
SELECT true From favourites
WHERE target_id = $1 AND user_id = $2;

-- name: GetAllFavourites :many
SELECT 
	AC.id, 
	AC.first_name, 
	AC.birth_date,
	FAV.fav_id
FROM 
	accounts AS AC
	INNER JOIN favourites AS FAV ON AC.id = $1
ORDER BY
    Fav.fav_at DESC
LIMIT $2
OFFSET $3;

-- name: FavoritesMetadata :one
SELECT 
	count(*)
FROM 
	accounts
	INNER JOIN favourites ON accounts.id = $1;
