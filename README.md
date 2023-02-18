# Description
This is still a wip, but the end goal will be a cloud file share server that will use grpc to stream file uploads and downloads. 
Right now, basic uploads are possible, and a secure asymmetric encryption base authentication scheme exists. 
Goals to finish a CLI client, and in the more distant future, an android app client as well.
Addtionally, contained inside the repo is a simple up/down migration framework built around postgres.

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
