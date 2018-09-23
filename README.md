# lumenaddr

Generate vanity addresses for stellar (addresses with specific suffixes).

## Downloads

I have already built [executables available for download](https://github.com/fnando/lumenaddr/releases).

## Usage

Just run the command providing the words you're looking for. Just remember that longer the word, the more time it'll take. The following command will lookup for keys that match the words `STELLAR` and `LUMENS`, and output them to the console.

```
lumenaddr STELLAR LUMENS
```

If you want to save keys to the database instead (only PostgreSQL is supported), create the table `addresses` using the following SQL:

```sql
CREATE TABLE addresses (
  word text NOT NULL,
  public_key text NOT NULL,
  encrypted_private_key text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE EXTENSION pgcrypto;
```

Then you can set the database url and encryption key via env vars.

```bash
DATABASE_URL=postgres:///lumenaddr?sslmode=disable ENCRYPTION_KEY=sekret lumenaddr STELLAR LUMENS
```

To retrieve the private key, you can use `DATABASE_URL=... ENCRYPTION_KEY=... lumenaddr --print`, which uses the following query:

```sql
SELECT
  word,
  public_key,
  convert_from(decrypt(encrypted_private_key::bytea, 'sekret', 'aes'), 'SQL_ASCII') private_key,
  created_at
FROM
  addresses
ORDER BY
  length(word) DESC,
  created_at DESC
LIMIT 10;
```

## Screenshots

![](https://github.com/fnando/lumenaddr/raw/master/screenshots/database.png)

![](https://github.com/fnando/lumenaddr/raw/master/screenshots/terminal.png)

![](https://github.com/fnando/lumenaddr/raw/master/screenshots/lumenaddr.gif)


## License

(The MIT License)

Copyright (c) 2018 Nando Vieira

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
'Software'), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
