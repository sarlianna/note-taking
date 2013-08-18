package main

import (
  "fmt"
  "os"
  "log"
  "flag"
  "strings"
  "bytes"
  _ "github.com/lib/pq"
  "database/sql"
)

//TODO: error checking on input; escaping against sql injection and type checking on arguments to queries.
//      add deleting by id and id retun when adding notes, for automated todo lists.
//CLEANUP: add a function that checks silent flag and prints given arguments if it isn't set.
//CLEANUP: add error checking function? that handles sql.ErrNoRows distinctions. (?)

func main() {

  log.SetFlags(log.Lshortfile)

  db, err := sql.Open("postgres", "user=postgres dbname=notes sslmode=disable")
  if err != nil {
    log.Fatal(err)
  }

  help := flag.Bool("h", false, "help")
  silent := flag.Bool("s", false, "stops any output")

  flag.Parse()

  args := os.Args[1:]

  if *help {
    if ! *silent {
      fmt.Println("Commands:\n  no args                  - print last note stored. \n" +
                    "  <tag>                    - print last note stored for #tag. \n" +
                    "  note  <tag> <message...> - save a note under section tag with body message.\n" +
                    "  list  [tag] [num = 5]    - print last num messages under #tag.  No tag results in last num from all tags. \n" +
                    "  clear [tag]              - delete num notes from tag.  No tag results in num notes from all tags. \n" +
                    "  clear last               - delete last note saved. \n" +
                    "  clear all                - delete all saved notes. \n" +
                    "\nFlags:\n  -h print this help message\n" +
                    "  -s silent mode\n")
    }
  }

  if len(args) == 0 {
    //case: $notes
    var msg string
    var tag string
    err := db.QueryRow("SELECT tag, message FROM notes ORDER BY age DESC LIMIT 1").Scan(&tag, &msg)
    switch {
      case err == sql.ErrNoRows:
        if ! *silent {
          fmt.Println("No recent notes.");
        }
      case err != nil: 
        log.Fatal(err)
      default: 
        if ! *silent {
          fmt.Println("Latest note:\n", tag,":", msg)
        }
    }

    os.Exit(0)

  }

  if strings.EqualFold(args[0], "note") {
    //case: $notes note <tag> <message>

    var message bytes.Buffer
    for i := 2; i < len(args); i++ {
      message.WriteString(" " + args[i])
    }

    _, err := db.Query("INSERT INTO notes ( tag, message ) VALUES ( " + "'" + args[1] + "', '" + message.String() + "');" )
    if err != nil {
      log.Fatal(err)
    }

  } else if strings.EqualFold(args[0], "list") {
    //case: $notes list [tag] [num]

    query := "SELECT tag, message FROM notes"
    if len(args) > 1 {
      query += " WHERE tag='" + args[1] +"'"
    }
    query += " ORDER BY age DESC"
    if len(args) > 2 {
      query += " LIMIT " + args[2]
    } else {
      query += " LIMIT 5"
    }

    rows, err := db.Query(query)

    if err != nil {
      log.Fatal(err)
    }

    var msg string
    var tag string
    for rows.Next() { 
      errc := rows.Scan(&tag, &msg)
      if errc != nil {
        log.Fatal(errc)
      }    

      if ! *silent {
        fmt.Println(tag, ":", msg, "\n")
      }
    }

  } else if strings.EqualFold(args[0], "clear") {
    //case: $notes clear ...

    if len(args) == 1 || strings.EqualFold(args[1], "last") {
      _, err := db.Query("DELETE FROM notes WHERE age= " +
                           " (SELECT max(age) FROM notes)")
      if err != nil {
        log.Fatal(err)
      }

    } else if strings.EqualFold(args[1], "all") {
      //TODO: ask if they're sure.  Add a flag to force it.
      db.Query("TRUNCATE notes")

    } else {
      //case: $notes clear [tag] [num]
      //TODO: FIX THIS, IT DOESN'T WORK
      query :="DELETE FROM notes WHERE age= (SELECT max(age) FROM notes) " +
                           "AND tag='" + args[1] +"'"
      
      _, err := db.Query(query)
      if err != nil {
        log.Fatal(err)
      }
    }

  } else {
    //tag is being passed, print the last note
    var ret string
    err := db.QueryRow("SELECT message FROM notes WHERE tag='" + args[0] + "' ORDER BY age DESC LIMIT 1").Scan(&ret)

    switch {
      case err ==sql.ErrNoRows:
        if ! *silent {
          fmt.Println("No recent notes about " + args[0] + ".\n")
        }
      case err != nil: 
        log.Fatal(err)
      
      default:
        if ! *silent { 
          fmt.Println("Latest ", args[0], " note:\n", ret) 
        }
    }
    os.Exit(0)

  }
 

}
