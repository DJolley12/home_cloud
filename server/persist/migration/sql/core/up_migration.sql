  UPDATE migrations 
  SET
    name = $3, 
    description = $2,
    migration_state = $4
  WHERE id = $1;
