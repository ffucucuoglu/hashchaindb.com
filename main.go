package main

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"log"
	"strconv"
	"time"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type chain struct {
	value     []byte
	prvHash   []byte
	timestamp int64
	no        int64
}

var newchain, prevchain chain

func addDb(inp string, db *sql.DB) {

	var val, phsh string

	if prevchain.no == 0 {
		log.Printf("prevchain.no = 0\n")
		row, _ := db.Query("SELECT value,prvHash,timestamp,no FROM hashchain WHERE no = (SELECT MAX(no) FROM hashchain)")
		row.Next()
		row.Scan(&val, &phsh, &prevchain.timestamp, &prevchain.no)
		prevchain.value = []byte(val)
		prevchain.prvHash = []byte(phsh)
		row.Close()
	}

	newchain.value = []byte(inp)
	newchain.timestamp = time.Now().Unix()
	newchain.no = prevchain.no + 1
	newchain.prvHash = hashofprev(prevchain)

	stmt, err := db.Prepare("INSERT INTO hashchain(value,prvHash, timestamp, no) values(?,?,?,?)")
	checkErr(err)//input types must be controlled in backend

	res, err := stmt.Exec(string(newchain.value[:]), string(newchain.prvHash[:]), newchain.timestamp, newchain.no)
	checkErr(err)

	_, err = res.LastInsertId()
	checkErr(err)

	prevchain = newchain
}

func hashofprev(pc chain) []byte {

	var prvchnall bytes.Buffer

	prvchnall.WriteString(string(pc.value))
	prvchnall.WriteString(string(pc.prvHash))
	prvchnall.WriteString(strconv.FormatInt(pc.timestamp, 10))
	prvchnall.WriteString(strconv.FormatInt(pc.no, 10))

	sha256 := sha256.Sum256([]byte(prvchnall.String()))

	return []byte(hex.EncodeToString(sha256[:]))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
