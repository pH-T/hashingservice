# HashService

A simple bcrypt hashingservice to show the features of Go:

- testing
- http(s)
- crypto
- interfaces
- flag
- json

## Endpoints

- `GET /health` returns `returns`
- `POST /hash` with body `{"pw":"foobar"}` returns `{"hash":"$2a$14$DewuCqBaOSjOVwQ3bhBsnORYdZUeXQ5i00D5b9l1NYgd1ND6zisq2"}`
- `POST /verify` with body `{"pw":"foobar", "hash":"$2a$14$DewuCqBaOSjOVwQ3bhBsnORYdZUeXQ5i00D5b9l1NYgd1ND6zisq2"}` returns `{"verified":true,"selfmade":true}` (where `verified` indicates that bcrypt(pw)==hash and where `selfmade` indicates that this services has seen this hash before (just to get a usecase for "storage"))




