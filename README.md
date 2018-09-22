# lumenaddr

Generate vanity addresses for stellar (addresses with specific suffixes).

## Usage

Just run the command providing the words you're looking for. Just remember that longer the word, the more time it'll take. The following command will lookup for keys that match the words `STELLAR`, `LUMENS`, and `FNANDO`, and output them to the console.

```
lumenaddr STELLAR LUMENS FNANDO
```

If you want to save keys to the database instead (only PostgreSQL is supported), create the table `addresses` using the following SQL:

```sql
CREATE TABLE addresses (
  word text NOT NULL,
  public_key text NOT NULL,
  encrypted_private_key text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);
```

Then you can set the database url and encryption key via env vars.

```bash
DATABASE_URL=postgres:///lumenaddr?sslmode=disable ENCRYPTION_KEY=sekret lumenaddr STELLAR LUMENS
```

To retrieve the private key, use the following query:

```sql
SELECT
  word,
  public_key,
  convert_from(decrypt(encrypted_private_key::bytea, 'sekret', 'aes'), 'SQL_ASCII') private_key
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

