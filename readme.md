# string encryptor

## using aes-gcm and pbkdf2-sha256

## commands

- create
- set
- generate
- get
- keys
- delete
- change \[password|iterations\]

## flags

| name | default  |            hint             |
| :--: | :------: | :-------------------------: |
|  -f  | .cryptex |        cryptex file         |
|  -k  |          |             key             |
|  -s  |          |           string            |
|  -p  |          |          password           |
|  -i  | 600_000  |         iterations          |
|  -l  |    8     |      generation length      |
| -cm  |   ludo   | generation charset modifier |
