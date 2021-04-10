package main

import (
	"database/sql"
	"log"
	"net"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/cenkalti/backoff/v4"
	"github.com/spf13/pflag"
	"github.com/xo/dburl"
)

func main() {
	var (
		dsn     = pflag.StringP("dsn", "d", "", "postgres URL")
		timeout = pflag.DurationP("timeout", "t", 10*time.Second, "timeout")
		verbose = pflag.BoolP("verbose", "v", false, "verbose errors")
	)

	pflag.Parse()

	go func() {
		<-time.After(*timeout)
		log.Println("exit by timeout")
		os.Exit(1)
	}()

	uri, err := dburl.Parse(*dsn)
	if err != nil {
		log.Fatalf("cannot parse postgres dsn")
	}

	b := backoff.NewExponentialBackOff()

	err = backoff.Retry(func() (err error) {
		_, err = net.Dial("tcp", uri.Host)
		if *verbose && err != nil {
			log.Printf("try tcp connection: %s", err)
		}

		return err
	}, b)
	if err != nil {
		log.Fatalf("backoff error: %s", err)
	}

	db, err := sql.Open("postgres", uri.String())
	if err != nil {
		log.Fatalf("cannot open connection: %s", err)
	}

	err = backoff.Retry(func() error {
		err := db.Ping()
		if *verbose && err != nil {
			log.Printf("try ping: %s", err)
		}

		if err == nil {
			_, err = db.Exec("select 1")
			if *verbose && err != nil {
				log.Printf("try execute 'select 1': %s", err)
			}
		}

		return err
	}, b)
	if err != nil {
		log.Fatalf("backoff error: %s", err)
	}

	if *verbose {
		log.Println("success!")
	}
}
