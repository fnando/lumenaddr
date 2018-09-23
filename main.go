package main

import (
  "flag"
  "fmt"
  "os"
  "strings"
  "time"

  "database/sql"
  _ "github.com/lib/pq"
  "github.com/stellar/go/keypair"
  "github.com/dustin/go-humanize"
  "github.com/dustin/go-humanize/english"
  "github.com/hako/durafmt"
)

var databaseUrl string = os.Getenv("DATABASE_URL")
var encryptionKey string = os.Getenv("ENCRYPTION_KEY")
var saveToDatabase bool = databaseUrl != "" && encryptionKey != ""
var printKeys bool
var maxConcurrency int
var db *sql.DB
var throttle chan bool
var totalKeys int64
var matchedKeys int64
var starts time.Time
var words []string

const findSql = `
  SELECT
    word,
    public_key,
    convert_from(decrypt(encrypted_private_key::bytea, $1, 'aes'), 'SQL_ASCII') private_key,
    created_at
  FROM
    addresses
  WHERE
    created_at > $2
  ORDER BY
    created_at ASC
  LIMIT 1000
`

func main() {
  flag.IntVar(&maxConcurrency, "concurrency", 10, "Specify the concurrency")
  flag.BoolVar(&printKeys, "print", false, "Output saved keys")

  flag.Parse()
  words = flag.Args()
  throttle = make(chan bool, maxConcurrency)

  if saveToDatabase {
    connectToDatabase()
  }

  if printKeys {
    runPrintKeys()
  } else {
    runFindKeys()
  }
}

func connectToDatabase() {
  var err error
  db, err = sql.Open("postgres", databaseUrl)
  panicWithError(err)
}

func runPrintKeys() {
  var timeCondition time.Time
  var err error
  timeCondition, _ = time.Parse(time.RFC3339, "1900-01-1")

  panicWithError(err)
  total := 0

  for {
    rows, err := db.Query(findSql, encryptionKey, timeCondition)
    panicWithError(err)

    count := 0

    for rows.Next() {
      count += 1
      total += 1

      var suffix string
      var address string
      var seed string
      var createdAt time.Time

      err := rows.Scan(&suffix, &address, &seed, &createdAt)
      panicWithError(err)

      formattedAddress := formatAddress(address, suffix)
      template := "\r🔑 %s\n🔐 %s\n⏱  %s\n\n"
      fmt.Printf(template, formattedAddress, seed, humanize.Time(createdAt))

      timeCondition = createdAt
    }

    if count == 0 {
      break
    }
  }

  if total == 0 {
    fmt.Printf("\n😞 No keys found so far.\n")
  } else {
    fmt.Printf("🙂 %s %s found so far.\n", humanize.Comma(int64(total)), english.PluralWord(total, "key", ""))
  }
}

func runFindKeys() {
  if len(words) == 0 {
    fmt.Fprintf(os.Stderr, "💣 \x1b[31mERROR: You need to provide at least one word.\x1b[0m\n")
    os.Exit(1)
  }

  if len(words) == 1 {
    words = strings.Split(words[0], " ")
  }

  if saveToDatabase {
    fmt.Printf("\x1b[32m✅ Matching keys will be saved to the database. Use `lumenaddr --print` to view them.\x1b[0m\n\n")
  } else {
    fmt.Fprintf(os.Stderr, "\x1b[31m⚠️  DATABASE_URL and ENCRYPTION_KEY config vars not set; outputting keys instead.\x1b[0m\n\n")
  }

  index := 0
  starts = time.Now()

  print("🔎 Searching. Press CTRL-C to stop...")

  for {
    index += 1
    throttle <- true
    go generatePair(words)
  }
}

func panicWithError(err error) {
  if err != nil {
    fmt.Printf("\r💣 \x1b[31mERROR: %s\x1b[0m\n", err)
    os.Exit(1)
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

  panicWithError(err)

  address := pair.Address()
  seed := pair.Seed()
  elapsed := time.Since(starts)
  totalKeys += 1
  match := matchingWord(address, words)

  if match == "" {
    if totalKeys % 1000 == 0 {
      printStatsMessage()
    }

    <-throttle
    return
  }

  matchedKeys += 1

  if saveToDatabase {
    _, err = db.Exec("INSERT INTO addresses (word, public_key, encrypted_private_key) VALUES ($1, $2, encrypt($3, $4, 'aes'))", match, address, seed, encryptionKey)
    panicWithError(err)
    printStatsMessage()
  } else {
    formattedAddress := formatAddress(address, match)
    duration, _ := durafmt.ParseString(elapsed.String())
    _ = duration
    template := "\r🔑 %s\n🔐 %s\n\n"
    fmt.Printf(template, formattedAddress, seed)
  }

  <-throttle
}

func printStatsMessage() {
  template := "\r🔎 Found %s out of %s %s. Press CTRL-C to stop..."
  fmt.Printf(template, humanize.Comma(matchedKeys), humanize.Comma(totalKeys), english.PluralWord(int(totalKeys), "key", ""))
}

func formatAddress(address string, suffix string)(string) {
  lastIndex := strings.LastIndex(address, suffix)
  return fmt.Sprintf("\x1b[34m%s\x1b[44m\x1b[37m%s\x1b[0m", address[0:lastIndex], suffix)
}
