# Roadmap

## DB

- setup permission structure
- setup user/password system

## Server

- Download method
  - should have queries to check permission (ip address? passphrase? both?)
- Encryption
- Https

## Auth Process

- Server generates user with one time passphrase
- Client gen keys, sends pub key to server, with passphrase
- Server verifies one time passphrase, returns server pub key, refresh token encrypted with client pub key
- Client decrypts with priv key, saves token, encrypts with server pub key -> server
- Server decrypt with priv key,
