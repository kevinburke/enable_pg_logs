package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var loggingConf = []byte(
	`logging_collector = on
log_rotation_size = 200MB
log_duration = on
log_lock_waits = on
log_statement = 'all'
`)

const include = "include = 'conf.d/logging.conf'\t\t# github.com/kevinburke/enable_pg_logs"

func init() {
	flag.Usage = func() {
		os.Stderr.WriteString(`enable_pg_logs

Turn on logging for your Postgres database. This has only been tested on Mac 
with Homebrew postgresql so far, though I tried to make it portable.

This tool requires access to the "psql" command.
`)
		os.Exit(2)
	}
}

func main() {
	flag.Parse()
	cmd := exec.Command("psql", "--no-psqlrc", "--no-align", "--tuples-only",
		"--command", "SHOW data_directory", "postgres")
	bits, err := cmd.CombinedOutput()
	if err != nil {
		os.Stderr.Write(bits)
		log.Fatal(err)
	}
	dataDir := strings.TrimSpace(string(bits))
	logDir := filepath.Join(dataDir, "pg_log")
	fmt.Fprintf(os.Stderr, "creating log directory %s... ", logDir)
	if err := os.Mkdir(logDir, 0755); err != nil {
		if strings.Contains(err.Error(), "file exists") {
			fmt.Fprintf(os.Stderr, "already created\n")
		} else {
			log.Fatal(err)
		}
	} else {
		os.Stderr.WriteString("\n")
	}

	confDir := filepath.Join(dataDir, "conf.d")
	fmt.Fprintf(os.Stderr, "creating conf.d directory %s... ", confDir)
	if err := os.Mkdir(confDir, 0755); err != nil {
		if strings.Contains(err.Error(), "file exists") {
			fmt.Fprintf(os.Stderr, "already created\n")
		} else {
			log.Fatal(err)
		}
	} else {
		os.Stderr.WriteString("\n")
	}
	pgConf := filepath.Join(dataDir, "postgresql.conf")
	f, err := os.Open(pgConf)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(f)
	foundInclude := false
	for scanner.Scan() {
		text := scanner.Text()
		if text == include {
			foundInclude = true
			break
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	f.Close()
	if foundInclude {
		fmt.Fprintf(os.Stderr, "found include stmt in %s, nothing to do!\n", pgConf)
	} else {
		fmt.Fprintf(os.Stderr, "adding include stmt to %s\n", pgConf)
		f2, err2 := os.OpenFile(pgConf, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err2 != nil {
			log.Fatal(err2)
		}
		if _, err3 := f2.WriteString(include + "\n"); err3 != nil {
			log.Fatal(err3)
		}
		loggingFilename := filepath.Join(confDir, "logging.conf")
		if err := ioutil.WriteFile(loggingFilename, loggingConf, 0600); err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(os.Stderr, "Wrote %s\n", loggingFilename)
		fmt.Fprintf(os.Stderr, "Done. Be sure to restart Postgres\n")
	}
}
