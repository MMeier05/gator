-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetUserFeeds :many
SELECT 
  feeds.name AS feed_name,
  users.name AS user_name,
  feeds.id AS feed_id,
  users.id AS user_id,
  feeds.url
FROM
  feeds
LEFT JOIN 
  users ON feeds.user_id = users.id;

-- name: GetFeed :one
SELECT
  *
FROM
  feeds
WHERE
  url = $1;

-- name: MarkFeedFetched :exec
UPDATE 
  feeds
SET
  last_fetched_at = $2,
  updated_at = $2
WHERE
  id = $1;

-- name: GetNextFeedToFetch :one
SELECT
  feeds.id, feeds.url
FROM 
  feeds
LEFT JOIN 
  feed_follows ON feeds.user_id = $1
ORDER BY
  last_fetched_at ASC NULLS FIRST
LIMIT 1;
