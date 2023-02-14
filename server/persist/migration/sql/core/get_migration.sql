SELECT (
  id,
  name,
  description,
  migration_state
)
FROM migrations
WHERE id = $1;
