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

## examples

```shell
cryptex create -f my_cryptex -i 2100000
cryptex generate -f my_cryptex -k name.secret -p superPassword123 -l 12 -cm ldo
cryptex keys -f my_cryptex
cryptex get -f my_cryptex -k name.secret -p superPassword123 >> secret.txt
cryptex change iterations -f my_cryptex -i 4200000
```
