package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"encoding/csv"
	"fmt"
	"hash"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type config struct {
	in  string
	out string
	alg string
}

func main() {

	var cArgs []string

	cArgs = os.Args

	if len(cArgs) > 1 {

		switch len(cArgs) {
		case 2:

			if cArgs[1] == "help" {
				outputUsage()
				break
			}

			if strings.Contains(cArgs[1], ".csv") {
				cArgs = append(cArgs, timeToString(time.Now()), "-sha1")
				convert(cArgs)
				break
			}

			cArgs = append(cArgs, timeToString(time.Now()), "-sha1")

			cArgs[1] += ".csv"

			convert(cArgs)

			break
		case 3:
			if strings.Contains(cArgs[2], ".csv") {
				cArgs = append(cArgs, "-sha1")
				if !strings.Contains(cArgs[1], ".csv") {
					cArgs[1] += ".csv"
				}
				convert(cArgs)
				break
			}
			cArgs = append(cArgs, "-sha1")
			cArgs[2] += ".csv"
			if !strings.Contains(cArgs[1], ".csv") {
				cArgs[1] += ".csv"
			}
			convert(cArgs)
			break
		case 4:
			if cArgs[3] != "-md5" && cArgs[3] != "-sha1" {
				fmt.Println("\n", "Err: incorrect hashing method. Only '-md5' and '-sha1' are acceptable parameters.")
				outputUsage()
			}
			if !strings.Contains(cArgs[1], ".csv") {
				cArgs[1] += ".csv"
			}
			if !strings.Contains(cArgs[2], ".csv") {
				cArgs[2] += ".csv"
			}
			convert(cArgs)
			break
		default:
			break
		}
	} else {
		fmt.Println("Use './hashbear help' to learn how to use the tool.")
		return
	}
}

func timeToString(t time.Time) string {
	r := fmt.Sprintf("%x", truncateString(t.String(), 10))
	r += ".csv"
	return r
}

func truncateString(s string, c int) string {
	r := s
	if len(r) > c {
		r = s[len(s)-c:]
	}
	return r
}

func checkCsv(s string) bool {
	if strings.Contains(s, ".csv") {
		return true
	}
	return false
}

func checkFormat(s string) string {
	if strings.Contains(s, "sha1") {
		return "sha1"
	} else if strings.Contains(s, "md5") {
		return "md5"
	} else {
		return "invalid"
	}
}

func convert(opt []string) {
	st := time.Now()
	source, err := os.Open(opt[1])
	if err != nil {
		log.Fatalf("Couldn't open source file: %x", err)
	}
	r := csv.NewReader(bufio.NewReader(source))
	defer source.Close()
	dest, err := os.Create(opt[2])
	if err != nil {
		log.Fatalf("Couldn't create destination file: %x", err)
	}
	defer dest.Close()
	w := csv.NewWriter(dest)
	defer w.Flush()
	totalRecords := 0
	for {
		line, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		var h hash.Hash
		if opt[3] == "-sha1" {
			h = sha1.New()
		} else if opt[3] == "-md5" {
			h = md5.New()
		}
		st := strings.Join(line, "")
		io.WriteString(h, st)
		val := []string{fmt.Sprintf("%x", h.Sum(nil))}
		w.Write(val)
		totalRecords++
	}
	fmt.Println("Total of", totalRecords, "records to", opt[2], "in format", opt[3], "elapsed time", time.Since(st))
}

func outputUsage() {
	fmt.Println(
		"\n",
		"| Hashbear v1.0 \n",
		"\n",
		"Hashbear takes a source CSV file, formatted as a single column of values to hash, and writes hashes to an output CSV file. Hashbear can accept either md5 or sha1 algorithsm.\n",
		"\n",
		"Usage:\n",
		"./hashbear {source.csv}*\n",
		"./hashbear {source.csv}* {destination.csv}\n",
		"./hashbear {source.csv}* {destination.csv} {algorithm}\n",
		"\n",
		"* is required\n",
		"the type flag can be either -md5 or -sha1")
}
