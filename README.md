```sql
CREATE TABLE addresses (
  word text NOT NULL,
  public_key text NOT NULL,
  encrypted_private_key text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);
```

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
