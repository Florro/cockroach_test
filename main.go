package main

import (
	"fmt"
	"database/sql"
	"log"
	"net/url"
	"math/big"
	"time"
	// "github.com/cockroachdb/cockroach-go"
	_ "github.com/lib/pq"

)

func BitCount(n *big.Int) int {
	count := 0
	for _, v := range n.Bits() {
		count += popcount(uint64(v))
	}
	return count
}

// Straight and simple C to Go translation from https://en.wikipedia.org/wiki/Hamming_weight
func popcount(x uint64) int {
	const (
		m1  = 0x5555555555555555 //binary: 0101...
		m2  = 0x3333333333333333 //binary: 00110011..
		m4  = 0x0f0f0f0f0f0f0f0f //binary:  4 zeros,  4 ones ...
		h01 = 0x0101010101010101 //the sum of 256 to the power of 0,1,2,3...
	)
	x -= (x >> 1) & m1             //put count of each 2 bits into those 2 bits
	x = (x & m2) + ((x >> 2) & m2) //put count of each 4 bits into those 4 bits
	x = (x + (x >> 4)) & m4        //put count of each 8 bits into those 8 bits
	return int((x * h01) >> 56)    //returns left 8 bits of x + (x<<8) + (x<<16) + (x<<24) + ...
}

func main() {

	dbURL := "postgres://root@localhost:26257?sslmode=disable"
	parsedURL, err := url.Parse(dbURL)
	if err != nil {
		log.Fatal(err)
	}
	parsedURL.Path = "tester"
	fmt.Println(parsedURL.String())

	db, err := sql.Open("postgres", parsedURL.String())
	fmt.Println("open")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("create db")
	if _, err := db.Exec("CREATE DATABASE IF NOT EXISTS tester"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("create done")

	db.SetMaxOpenConns(5 + 1)
	fmt.Println("create table")
	if _, err = db.Exec("CREATE TABLE IF NOT EXISTS test (id BIGINT PRIMARY KEY, features BIGINT NOT NULL)"); err != nil {
		log.Fatal(err)
	}

	// stmt, err := db.Prepare("INSERT INTO test VALUES($1, $2)")
	// fmt.Println("prep")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer stmt.Close()

	// fmt.Println("exec")
	// for i := 0; i < 500000; i++ {
	// 	_, err = stmt.Exec(i, i)
	// 	if i % 1000 == 0 {
	// 		fmt.Println(i)
	// 	}
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }


	start := time.Now()
	rows, err := db.Query("SELECT features, 1 ^ features AS tmp FROM test")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	fmt.Println("exec time query: ", time.Since(start))

	var (
		id int
		// features int
		tmp int
	)

	start = time.Now()
	for rows.Next() {
		err := rows.Scan(&id, &tmp)
		if err != nil {
			log.Fatal(err)
		}
		// log.Println(id, features, tmp)
		// fmt.Println(BitCount(&tmp))
		// fmt.Println(popcount(uint64(tmp)))
	}
	fmt.Println("exec time: ", time.Since(start))
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
