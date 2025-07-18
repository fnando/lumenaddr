package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/dustin/go-humanize/english"
	"github.com/hako/durafmt"
	"github.com/stellar/go/keypair"
	"github.com/tyler-smith/go-bip39"
)

var generateMnemonic bool
var printKeys bool
var throttle chan bool
var totalKeys int64
var matchedKeys int64
var starts time.Time
var words []string

func main() {
	flag.BoolVar(&generateMnemonic, "mnemonic", false, "Generate keys out of mnemonic (slower)")
	flag.BoolVar(&printKeys, "print", false, "Output saved keys")

	flag.Parse()
	words = flag.Args()
	throttle = make(chan bool, 100)

	runFindKeys()
}

func runFindKeys() {
	if len(words) == 0 {
		fmt.Fprintf(os.Stderr, "💣 \x1b[31mERROR: You need to provide at least one word.\x1b[0m\n")
		os.Exit(1)
	}

	if len(words) == 1 {
		words = strings.Split(words[0], " ")
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

func matchingWord(address string, words []string) string {
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

// SLIP-0010 implementation for ed25519 (matches stellar-hd-wallet)
func deriveEd25519Key(mnemonic string) ([]byte, error) {
	seed := bip39.NewSeed(mnemonic, "")

	// Create master key using SLIP-0010 for ed25519
	h := hmac.New(sha512.New, []byte("ed25519 seed"))
	h.Write(seed)
	masterKey := h.Sum(nil)

	// Stellar derivation path: m/44'/148'/0'
	indices := []uint32{
		0x80000000 + 44,  // 44' (hardened)
		0x80000000 + 148, // 148' (hardened, Stellar coin type)
		0x80000000 + 0,   // 0' (hardened, account 0)
	}

	key := masterKey[:32]
	chainCode := masterKey[32:]

	// Derive each level in the path
	for _, index := range indices {
		key, chainCode = deriveChildKeyEd25519(key, chainCode, index)
	}

	return key, nil
}

func deriveChildKeyEd25519(parentKey, parentChainCode []byte, index uint32) ([]byte, []byte) {
	h := hmac.New(sha512.New, parentChainCode)

	// For hardened derivation (all Stellar keys are hardened)
	h.Write([]byte{0x00})
	h.Write(parentKey)

	// Write the index as big-endian 32-bit integer
	indexBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(indexBytes, index)
	h.Write(indexBytes)

	digest := h.Sum(nil)
	return digest[:32], digest[32:]
}

func generatePair(words []string) {
	var err error
	var mnemonic string
	var pair *keypair.Full

	if generateMnemonic {
		entropy, err := bip39.NewEntropy(256)
		panicWithError(err)

		mnemonic, err = bip39.NewMnemonic(entropy)
		panicWithError(err)

		// Use SLIP-0010 ed25519 derivation instead of BIP32
		privateKeyBytes, err := deriveEd25519Key(mnemonic)
		panicWithError(err)

		var privateKey [32]byte
		copy(privateKey[:], privateKeyBytes)

		pair, err = keypair.FromRawSeed(privateKey)
	} else {
		pair, err = keypair.Random()
	}

	panicWithError(err)

	address := pair.Address()
	seed := pair.Seed()
	elapsed := time.Since(starts)
	totalKeys += 1
	match := matchingWord(address, words)

	if match == "" {
		if totalKeys%1000 == 0 {
			printStatsMessage()
		}

		<-throttle
		return
	}

	matchedKeys += 1

	formattedAddress := formatAddress(address, match)
	duration, _ := durafmt.ParseString(elapsed.String())
	_ = duration

	if generateMnemonic {
		template := "\r🔑 %s\n🔐 %s\n📄 %s\n\n"
		fmt.Printf(template, formattedAddress, seed, mnemonic)
	} else {
		template := "\r🔑 %s\n🔐 %s\n\n"
		fmt.Printf(template, formattedAddress, seed)
	}

	<-throttle
}

func printStatsMessage() {
	template := "\r🔎 Found %s out of %s %s. Press CTRL-C to stop..."
	fmt.Printf(template, humanize.Comma(matchedKeys), humanize.Comma(totalKeys), english.PluralWord(int(totalKeys), "key", ""))
}

func formatAddress(address string, suffix string) string {
	lastIndex := strings.LastIndex(address, suffix)
	return fmt.Sprintf("\x1b[34m%s\x1b[44m\x1b[37m%s\x1b[0m", address[0:lastIndex], suffix)
}
