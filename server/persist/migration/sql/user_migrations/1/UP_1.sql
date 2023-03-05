CREATE TABLE hc_user (
  id BIGSERIAL PRIMARY KEY
);

CREATE TABLE ip_address (
  id BIGSERIAL PRIMARY KEY,
  ip_addr INET NOT NULL,
  user_id BIGINT NOT NULL
);

CREATE TABLE permission (
  id BIGSERIAL PRIMARY KEY,
  dir_path VARCHAR(500) NOT NULL,
  user_id BIGINT NOT NULL
);

CREATE TABLE payload (
  id BIGSERIAL PRIMARY KEY,
  path VARCHAR(500) NOT NULL,
  upload_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  upload_user_id BIGINT NOT NULL
);

CREATE TABLE payload_download_history (
  id BIGSERIAL PRIMARY KEY,
  download_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  ip_id BIGINT NOT NULL,
  user_id BIGINT NOT NULL
);

CREATE TABLE key_set (
  id BIGSERIAL PRIMARY KEY,
  public_key BYTEA NOT NULL,
  private_key BYTEA NOT NULL,
  public_sign_key BYTEA NOT NULL,
  private_sign_key BYTEA NOT NULL,
  user_public_key BYTEA NOT NULL,
  user_sign_key BYTEA NOT NULL,
  user_id BIGINT NOT NULL
);

CREATE TABLE refresh_token (
  token BYTEA NOT NULL,
  expiry TIMESTAMP NOT NULL,
  user_id BIGINT NOT NULL
);
