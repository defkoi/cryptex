# string encryptor

## using aes-cbc and pbkdf2-sha512

## commands

- encrypt
- decrypt

## flags

| name | default  |
| :--: | :------: |
|  -f  | .cryptex |

## .cryptex file

```ts
type Cryptex = {
  string: string;
  initial_vector: string;
  salt: string;
  iterations: number = 1_048_576;
};
```
