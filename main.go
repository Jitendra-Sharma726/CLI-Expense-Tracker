package main 

import (
  "encoding/csv"
  "encoding/json"
  "fmt"
  "os"
  "strconv"
  "time"
  bolt "go.etcd.io/bbolt"
)

const dbFile = "expenses.db"
const budgetLimit = 1000.00
var bucketName = []byte("expenses")

//Expense model
type Expense struct {
  ID   uint64   `json:"id"`
  Amount float64  `json:"amount"`
  Category string  `json:"category"`
  Description string `json:"description"`
  Date        string `json:"date"`
}

//Intialize DB and bucket
func initDB() *bolt.DB {
  // Open DB with default options
  db,err := bolt.Open(dbFile, 0600, nil)
  if err != nil {
    panic(fmt.Sprintf("Failed to open database: %v", err))
  }

  //Ensure bucket exists safely 
  err = db.Update(func(tx *bolt.Tx) error {
    _, err := tx.CreateBucketIfNotExists(bucketName)
    return err
    })

  if err != nil {
    panic(fmt.Sprintf("Failed to create bucket: %v", err))
  }

  return db
}

//Convert ID to zero-padded byte slice for sorted keys
func idToBytes(id uint64) []byte {
  return []byte(fmt.Sprintf("%05d", id))
}

//check if spending exceeds budget limit
func checkBudgetAlert(db *bolt.DB) {
  var total float64

  //Read-only transaction
  db.View(func(tx *bolt.Tx) error {
    b := tx.Bucket(bucketName)

  //Iterate through bucket
  return b.ForEach(func(k,v []byte) error {
    var e Expense
    if err := json.Unmarshal(v, &e); err == nil {
      total += e.Amount
    }
    return nil
  })
})

if total > budgetLimit {
  fmt.Printf("BUDGET ALERT: You have spent $%.2f. Limit is $%.2f.\n", total, budgetLimit)
  }
}


//Add a new expense
func addExpense(db *bolt.DB amount float64, cateogory, description string) {
  err := db.Update(func(tx *bolt.Tx) error {
    b := tx.Bucket(bucketName)

    //Auto-increment ID
    id, _ := b.NextSequence()

    expense := Expense{
      ID: id,
      Amount:  amount,
      Category: category,
      Description:description,
      Date       :time.Now().Format("2006-01-02"),
      }

    //Marshal to JSON
    data, err := json.Marshal(expense)
    if err != nil {
      return err
    }

    //save expense
    return b.Put(idToBytes(id), data)
    })

  if err != nil {
    fmt.Printf("Failed to add expense: %v\n", err)
    return
  }

  fmt.Printf("Added $%.2f for %s (%s)\n", amount, category, description)
  checkBudgetAlert(db)
}

//Fetch and display all expenses
func listExpenses(db *bolt.DB) {
  fmt.Println("===   Your Expenses ===")

  var total float64
  count := 0 

  db.View(func(tx *bolt.Tx) error {
    b := tx.Bucket(bucketName)


    return b.ForEach(func(k, v []byte) error {
      var e Expense
      if err := json.Unmarshal(v, &e); err == nil {
        total += e.Amount
      }
      return nil
   })   
})

if count == 0 {
  fmt.Println("No expenses found. You are rich.")
  return
 }
fmt.Println("------------------------------------")
fmt.Printf("Total Spent: $%.2f\n", total)  
}  
  
//Delete an expense by ID
func deleteExpense(db *bolt.DB, id uint64) {
  err := db.Update(func(tx  *bolt.Tx) error {
    b := tx.Bucket(bucketName)
    key := idToBytes(id)

// Verify existence before deletion
if b.Get(key) == nil {
  return fmt.Errorf("no expense found with ID %d", id)
  }
    
  return b.Delete(key)
})

if err != nil {
  fmt.Printf("Failed to delete : %v\n", err)
  return
 }

fmt.Printf("Deleted expense ID %d\n", id)
}


//Generate a CSV report
func exportCSV(db *bolt.DB) {
  file, err := os.Create("expenses_report.csv")
  if err != nil {
    fmt.Printf("Failed to create file: %v\n", err)
    return
  }
  defer file.Close()

  writer := csv.NewWriter(file)
  defer writer.Flush()

  //Write Headers
  writer.Write([]string{"ID", "Date", "Category", "Amount", "Descrption"})

  count :=0
  db.View(func(tx *bolt.Tx) error {
    b := tx.Bucket(bucketName)

    return b.Foreach(func(k, v []byte) error {
      var e Expense 
      if err := json.Unmarshal(v, &e); err != nil {
        return nil
      }

    writer.Write([]string{
      fmt.Sprintf("%d", e.ID),
      e.Date,
      e.Category,
      fmt.Sprintf("%.2f", e.Amount),
      e.Desciption,
      })
      count++
      return nil
    })
  })  
  fmt.Printf("Sucessfully exported %d expenses to expenses_report.csv\n", count)
}


func main() {
  //Validate arguments
  if len(os.Args) < 2 {
    fmt.Println("Usage:")
    fmt.Println("   go run main.go add <amount> <category>  \"<description>\"")
    fmt.Println("   go run main.go list")
    fmt.Println("   go run main.go delete <id>")
    fmt.Println("   go run main.go export")
    return
  }

  db := initDB()
  defer db.Close()
  command := os.Args[1]

  //CLI Router
  switch command {
  case "add":
    if len(os.Args) < 5 {
      fmt.Println("Error: Missing arguments. Usage: add <amount> <cateogory> \"<description>\"")
      return
    }
  amount, err := strconv.ParseFloat(os.Args[2], 64)
  if err != nil {
    fmt.Println("Error: Amount must be a valid number.")
    return 
  }

 fmt.Println("=== Expense Tracker ===")
 addExpense(db, amount, os.Args[3], os.Args[4]

Case "list":
    listExpenses(db)
            
Case "delete":
   if len(os.Args) < 3 {
     fmt.Println("Error: Missing ID. Usage: delete <id>")
     return
    }

id, err := strconv.ParseUint(os.Args[2], 10, 64)
if err != nil {
  fmt.Println("Error: ID must be a whole positive number.")
  return
}


fmt.Println("=== Expense Tracker ===")
deleteExpense(db, id)

default:
  fmt.Printf("Unknown command: %s\n", command)
}

}    


            











    
            
            
            

            
















      


















  



  
      
  


  




  

  

  





















  







  



  




  
  










    





























    
