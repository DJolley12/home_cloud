  UPDATE migrations 
  SET
    migration_state = $2
  WHERE id = $1;
