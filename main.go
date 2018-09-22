package main

import (
  "fmt"
  "log"
  "flag"
  "strings"
  "os"

  "database/sql"
  _ "github.com/lib/pq"
  "github.com/stellar/go/keypair"
)

var databaseUrl string = os.Getenv("DATABASE_URL")
var encryptionKey string = os.Getenv("ENCRYPTION_KEY")
var saveToDatabase bool = databaseUrl != "" && encryptionKey != ""
var maxConcurrency int
var db *sql.DB
var throttle chan bool

func main() {
  var err error
  var words []string
  var verbose bool

  flag.IntVar(&maxConcurrency, "concurrency", 10, "Specify the concurrency")
  flag.BoolVar(&verbose, "verbose", false, "Output iterations count")

  flag.Parse()
  words = flag.Args()

  if len(words) == 0 {
    fmt.Fprintf(os.Stderr, "ERROR: You need to provide at least one word.\n")
    os.Exit(1)
  }

  if len(words) == 1 {
    words = strings.Split(words[0], " ")
  }

  throttle = make(chan bool, maxConcurrency)

  if saveToDatabase {
    if verbose {
      fmt.Printf("\x1b[32m=> NOTICE: saving keys to the database. To view them, use the following SQL:\x1b[0m\n")
      fmt.Printf("\x1b[1;30m   SELECT word, public_key, convert_from(decrypt(encrypted_private_key::bytea, '<encryption key>', 'aes'), 'SQL_ASCII') private_key FROM addresses ORDER BY length(word) DESC LIMIT 10;\x1b[0m\n\n")
    }

    db, err = sql.Open("postgres", databaseUrl)
    panicWithError(err)
  } else {
    if verbose {
      fmt.Fprintf(os.Stderr, "\x1b[31mNOTICE: DATABASE_URL and ENCRYPTION_KEY config vars not set; outputting keys instead.\x1b[0m\n\n")
    }
  }

  index := 0

  for {
    index += 1

    if verbose {
      fmt.Printf("\r=> Lookups %d", index)
    }

    throttle <- true
    go generatePair(words)
  }
}

func panicWithError(err error) {
  if err != nil {
    panic(err)
  }
}

func logError(err error) {
  if err != nil {
    log.Fatal(err)
  }
}

func matchingWord(address string, words []string)(string) {
  var match string

  for _, suffix := range words {
    suffix := strings.ToUpper(suffix)

    if strings.HasSuffix(address, suffix) {
      match = suffix
      break
    }
  }

  return match
}

func generatePair(words []string) {
  var err error
  pair, err := keypair.Random()

  logError(err)

  address := pair.Address()
  seed := pair.Seed()
  match := matchingWord(address, words)

  if match == "" {
    <-throttle
    return
  }

  if saveToDatabase {
    result := db.QueryRow("INSERT INTO addresses (word, public_key, encrypted_private_key) VALUES ($1, $2, encrypt($3, $4, 'aes'))", match, address, seed, encryptionKey)
    _ = result
    panicWithError(err)
  } else {
    lastIndex := strings.LastIndex(address, match)
    fmt.Printf("\n\n%s\x1b[32m%s\x1b[0m\n%s\n\n", address[0:lastIndex], match, seed)
  }

  <-throttle
}
