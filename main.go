// Formal Dinner Sorter
// Jacob Wu
// 02.09.2020

package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
  "sort"
  "strconv"
	"strings"
	"time"
)

type person struct {
	name string
	id int // Unique Identifying Number of the person, used to prevent sitting someone next to the same people
  timesServed int // Number of times as waiter or kitchen staff
	haveMet []int // ID numbers of people seated with
	previousAssignments []string // History of table assignment or staff
  currentAssignment string // Currently task - table with specified number or staff
}

type table struct {
	id int // Unique identifying number of the table
  occupants []person	// All the people seated at the table
  disallow []int // Everyone the table's occupants have met before
}

type ByTimesServed []person // Sort in ascending order by number of times served as a kitchen/waiter staff

func (a ByTimesServed) Len() int           { return len(a) }
func (a ByTimesServed) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTimesServed) Less(i, j int) bool { return a[i].timesServed > a[j].timesServed }

type ByPeopleSeated []table // Sort in ascending order by number of times served as a kitchen/waiter staff

func (a ByPeopleSeated) Len() int           { return len(a) }
func (a ByPeopleSeated) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPeopleSeated) Less(i, j int) bool { return len(a[i].occupants) < len(a[j].occupants) }

var ( // Global variables
	people []person // Everyone from the excel sheet
  tables []table // Keeps track of all tables
  staffList []person // Keeps track of all staff
  numKitchenStaff = 7 // MUTABLE no. ppl needed at kitchen
  numTables = 31 // MUTABLE no. ppl needed for waiting
  printDetailsForDebug = false // Toggle for detailed info
	printEssentials = false
  dataToPrint [][]string // data for CSV export
) // ### Should have used a map for people variable; will make references easier

func resetVariables() { // Clear existing data for next round of assignments
  tables = tables[:0]
  staffList = make([]person, 0)
  dataToPrint = make([][]string, 0)
	for i := 0; i<len(people); i++ {
		people[i].previousAssignments = append(people[i].previousAssignments, people[i].currentAssignment)
	} // Save every student's previous assignment
}

func getStudentNames() { // This function populates student data
  csvFile, _ := os.Open("Dinner Seating - Student List 2018-19.csv") // Opens csv
	reader := csv.NewReader(bufio.NewReader(csvFile)) // Reads csv
  i := 0 // Counter used to assign ID numbers
	for { // For every student name
    i ++ // increment counter
		line, error := reader.Read()
		if error == io.EOF { // Deal with error, if there is one
			break // exit if finished (for loop not while)
		} else if error != nil {
			log.Fatal(error)
		}
		people = append(people, person{line[1] + " " + line[0], i, 0, make([]int, 0), make([]string, 0), ""}) // Add the student
	}
}

func shuffleStudents() { // Randomize order of students
  rand.Seed(time.Now().UnixNano())
	for i := len(people) - 1; i > 0; i-- { // Fisher–Yates shuffle
		j := rand.Intn(i + 1)
		people[i], people[j] = people[j], people[i]
	}
}

func sortTables() {
  sort.Sort(ByTimesServed(people)) // Sort students in order of number of served, big first; this way people who have served the least number of times are more likely to become waiters (since staff are picked last)
  tables = make([]table, numTables) // Make empty table	
	for i := 0; i < numTables; i++ { // Number all the tables
		tables[i].id = i+1
	}
	
	for i := 1; i<=(len(people)-numKitchenStaff-numTables); i++ {
		tryStudent(i)
	} // Assign students to a table until we have left the number of staff
}

func tryStudent(i int) {
	sort.Sort(ByPeopleSeated(tables)) // Sort tables by num. people seated, small first, as we start by assigning students to tables with least num. people  
	tableNo := 0 // Start trying at the first table (smallest num. people)
  for contains(tables[tableNo].disallow, people[0].id) {
    // If the upcoming student has met one of the existing table members
    tableNo ++ // In that's the case, pick another table
		if tableNo >= numTables { // If the student cannot be accomodated at any table
  		people = append(people, person{}) // Try another student by deferring this student
			copy(people[len(people)-i:], people[len(people)-i-1:]) // Move every person back by 1
			people[len(people)-i-1] = people[0] // Copy first person to last unsorted pos
			people = append(people[:0], people[1:]...) // Remove the original copy
			tryStudent(i+1) // Try to seat another student
		}
  }
	// At this point the student is good to be seated at the table
	people[0].currentAssignment = "Table " + strconv.Itoa(tables[tableNo].id)
  tables[tableNo].disallow = append(tables[tableNo].disallow, people[0].haveMet...) // Incorporate this person's haveMet list into the table's haveMet list
  tables[tableNo].occupants = append(tables[tableNo].occupants, people[0]) // Add student to table
	for _, occupant := range tables[tableNo].occupants {
		occupant.haveMet = tables[tableNo].disallow // Update every table occupant's haveMet list
	}
  people = append(people, people[0]) // Moves that person to last, since he's been assigned role
  people = append(people[:0], people[1:]...)
}

func contains(array []int, key int) bool { // Check if there are overlap
	for _, a := range array {
		if a == key {
			return true
		}
	}
	return false
}

func drawStaff(n int, role string) { // Take last 38 students
	for i := 1;  i<=n; i++ { // Until we get number of staff we need
    if printEssentials { fmt.Println(role + " <- " + people[0].name) }
    people[0].currentAssignment = role
    people[0].timesServed ++
    staffList = append(staffList, people[0])
    people = append(people, people[0]) // Moves first person to last, since he's been assigned role
    people = append(people[:0], people[1:]...)
  }
  if printDetailsForDebug { fmt.Print("\n", people, "\n\n") }
}

func writeToCSV() { // Export kitchen crew at row 33 and waiters row 34
	for tableNo, table := range tables { // For every table
		tableNames := []string{"Table " + strconv.Itoa(tableNo+1)}
		for _, occupant := range table.occupants { // For every person at the table
			tableNames = append(tableNames, occupant.name) // Include their name
		}
		dataToPrint = append(dataToPrint, tableNames) // And add all table names to queue
	}
	kitchenNames := []string{"Kitchen Staff"} // Prepare a variable to put all names in
	waiterNames := []string{"Waiter Staff"} // Prepare a variable to put all names in
	for _, occupant := range staffList { // For every person at the table
		if occupant.currentAssignment == "Kitchen" {
			kitchenNames = append(kitchenNames, occupant.name) // Include their name
		} else {
			waiterNames = append(waiterNames, occupant.name)
		}
	}
	dataToPrint = append(dataToPrint, make([]string, 0)) // Include empty line to demarcate diners and staffs
	dataToPrint = append(dataToPrint, kitchenNames) // Add names of kitchen staff
	dataToPrint = append(dataToPrint, waiterNames) // Add names of waiter staff
}

func exportCSV(n string) {
  fileName := "result" + n + ".csv" // Name the file: Seating 1, 2, 3; or Master
  file, err := os.Create(fileName)
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close() // Wait until we finished, then close file

  writer := csv.NewWriter(file)
  defer writer.Flush()

  for _, value := range dataToPrint { // For line in data to print
    err := writer.Write(value) // Write the line
    if err != nil {
      log.Fatal(err)
    }
  }
}

func writeMasterCSV() {
	sortPeopleByName() // Sort students by their name
	dataToPrint = make([][]string, 0)
	heading := []string{"Name", "ID", "No. Times as Staff", "Seatings 1-3"}
	dataToPrint = append(dataToPrint, heading)
	for _, person := range people { // For every student
		personInfo := []string{person.name, strconv.Itoa(person.id), strconv.Itoa(person.timesServed), strings.Join(person.previousAssignments, " ")}
		dataToPrint = append(dataToPrint, personInfo) // Include these information
	}
}

// Below include substantive reference to golang - sort package website //

type By func(p1, p2 *person) bool

type personSorter struct {
	people []person
	by      func(p1, p2 *person) bool // Closure used in the Less method.
}

func (by By) Sort(people []person) {
	ps := &personSorter{
		people: people,
		by:      by,
	}
	sort.Sort(ps)
}

func (s *personSorter) Len() int {
	return len(s.people)
}

func (s *personSorter) Swap(i, j int) {
	s.people[i], s.people[j] = s.people[j], s.people[i]
}

func (s *personSorter) Less(i, j int) bool {
	return s.by(&s.people[i], &s.people[j])
}

func sortPeopleByName() {
	name := func(p1, p2 *person) bool {
		return p1.name < p2.name
	}
	By(name).Sort(people)
}

// Above include substantive reference to golang - sort package website //

func main() {
	getStudentNames()
  if printDetailsForDebug { fmt.Print(people, "\n\n") }

	for i := 1; i<=3; i++ {
		shuffleStudents()
  	sortTables()

  	drawStaff(numKitchenStaff, "Kitchen")
  	drawStaff(numTables, "Waiter")

  	writeToCSV()
  	exportCSV(strconv.Itoa(i))

  	if printEssentials { fmt.Print("Finished run without error\n\n") }

		if printDetailsForDebug {
			for _, table := range tables {
				fmt.Println(table)
			}
		}

	  resetVariables()
	}
	writeMasterCSV()
	exportCSV("Master")
}

// Below are auxiliary functions helpful for debugging
func printStudentID() {
	for _, element := range people {
		fmt.Print(element.id, " ")
	}
	fmt.Print("\n\n")
}

// References
// {1} https://www.thepolyglotdeveloper.com/2017/03/parse-csv-data-go-programming-language/
// {2} https://yourbasic.org/golang/shuffle-slice-array/
// {3} https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-a-slice-in-golang
// {4} https://www.cyberciti.biz/faq/golang-for-loop-examples/
// {5} https://www.digitalocean.com/community/tutorials/how-to-do-math-in-go-with-operators
// {6} https://golangcode.com/write-data-to-a-csv-file/
// {7} https://golang.org/pkg/sort/
// {8} https://ispycode.com/GO/Collections/Arrays/Check-if-item-is-in-array
// {9} https://stackoverflow.com/questions/10105935/how-to-convert-an-int-value-to-string-in-go
// {10} https://stackoverflow.com/questions/46128016/insert-value-in-a-slice-at-given-index
// {11} https://www.dotnetperls.com/convert-slice-string-go

// Note: There is significant reference from the sort package on the Go website (7)