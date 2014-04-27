package main

import (
    "log"
    _ "code.google.com/p/odbc" 
    "database/sql"
    "encoding/csv"
    "os"
    "fmt"
    "strings"
    "path/filepath"
)

func main() {
    dbPath := os.Args[1] // first argument
    
    if !filepath.IsAbs(dbPath) {
        dbPath, _ = filepath.Abs(dbPath)
    }

    tables := os.Args[2:] // the rest of the arguments
    Transform(dbPath, tables)
}


func Transform (dbPath string, tables []string) {
    if len(tables) < 1 {
        panic("No tables provided!")
    }

    driver := "{Microsoft Access Driver (*.mdb, *.accdb)}"
    connStr := fmt.Sprintf("Driver=%s;Dbq=%s;", driver, dbPath)

    db, err := sql.Open("odbc", connStr)
    if err != nil { panic(err) }
    defer db.Close()

    done := make(chan bool, len(tables))

    for _, table := range tables {
        rowChannel := make(chan []string)
        go readTable(db, table, rowChannel)

        file, err := os.Create(table + ".csv")
        if err != nil { panic(err) }
        defer file.Close()
        go writeToFile(file, rowChannel, done)
    }

    for i := 0; i < len(tables); i++ {
        <-done
    }
}

func readTable(db *sql.DB, table string, rowChannel chan<- []string) {
    defer close(rowChannel)

    rows, err := db.Query("SELECT * FROM " + table)
    if err != nil { panic(err) }

    cols, err := rows.Columns()
    if err != nil { panic(err) }
    
    rowChannel <- cols

    // Mostly from http://stackoverflow.com/questions/14477941
    rawResult := make([][]byte, len(cols))
    
    dest := make([]interface{}, len(cols)) // A temporary interface{} slice
    for i, _ := range rawResult {
        dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
    }

    for rows.Next() {
        err := rows.Scan(dest...);
        if err != nil { panic(err) }

        result := make([]string, len(cols))
        for i, raw := range rawResult {
            if raw == nil {
                result[i] = ""
            } else {
                result[i] = strings.TrimSpace(string(raw))
            }
        }
        
        rowChannel <- result
    }

    if err := rows.Err(); err != nil {
        log.Fatal(err)
    }
}

func writeToFile(file *os.File, rowChannel <-chan []string, done chan<- bool) {
    writer := csv.NewWriter(file)
    defer writer.Flush()

    for row := range rowChannel {
        writer.Write(row)
    }

    done <- true
}
