package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"
)

// This is a simple program I wrote to get to know Go. Here is what it does:
//
// Part 1
// Read a text file to get a list of first names. Read another text file to get a list of last names.
// Define a Person as a struct with the following fields: a first name, a last name, and an age.
// Create a person by randomly selecting a first name from the list of first names, a last name from the list of last
// names, and an age as a random integer between 1 through 100.
// Write PeopleCount number of such person objects in JSON format to a text file.
//
// Part 2
// Read the list of people and analyze the data to print the following stats:
// 1. Number of people in the dataset
// 2. The count of distinct first names
// 3. The oldest three people
// 4. Minimum and maximum ages in the dataset
// 5. The three most popular last names with their frequencies
func main() {
	rand.Seed(time.Now().Unix())

	firstNames, err := getFirstNames()
	checkForError(err)

	lastNames, err := getLastNames()
	checkForError(err)

	err = writePeople(firstNames, lastNames)
	checkForError(err)

	err = analyze()
	checkForError(err)
}

func checkForError(err error) {

	if err != nil {
		log.Println("*** There was an error ***")
		log.Fatal(err) // Any error is fatal. Exit
	}
}

type Person struct {
	FirstName string
	LastName  string
	Age       int
}

func (p Person) String() string {
	return fmt.Sprintf("Person(FirstName=%s, LastName=%s, age=%d)", p.FirstName, p.LastName, p.Age)
}

var workspaceDir = os.Getenv("GOPATH")
var peopleFile = workspaceDir + "/src/github.com/hemantgokhale/learning-go/data/people.txt"

const PeopleCount = 100

func writePeople(firstNames, lastNames []string) (err error) {
	file, err := os.Create(peopleFile)

	if err != nil {
		return
	}

	defer func() {
		err = file.Close()
	}()

	for i := 0; i < PeopleCount; i++ {
		person := &Person{
			FirstName: firstNames[rand.Intn(len(firstNames))],
			LastName:  lastNames[rand.Intn(len(lastNames))],
			Age:       rand.Intn(100) + 1,
		}

		personBytes, err := json.Marshal(person)
		if err != nil {
			break
		}

		_, err = fmt.Fprintln(file, string(personBytes))
		if err != nil {
			break
		}
	}

	fmt.Printf("A dataset with %d people generated.\n", PeopleCount)
	return
}

func getFirstNames() (names []string, err error) {
	lines, err := getLines(workspaceDir + "/src/github.com/hemantgokhale/learning-go/data/firstNames.txt")
	if err != nil {
		return
	}

	for _, line := range lines {
		tokens := strings.Split(line, "\t")
		names = append(names, tokens[1], tokens[3])
	}
	return
}

func getLastNames() (names []string, err error) {
	lines, err := getLines(workspaceDir + "/src/github.com/hemantgokhale/learning-go/data/lastNames.txt")
	if err != nil {
		return
	}

	for _, line := range lines {
		tokens := strings.Split(line, " ")
		names = append(names, strings.Title(strings.ToLower(tokens[0])))
	}
	return
}

func analyze() (err error) {
	people, err := readPeople()
	if err != nil {
		return
	}

	const topCount = 3
	distinctFirstNames(people)
	minAndMaxAge(people)
	oldestPeople(people, topCount)
	popularLastNames(people, topCount)

	return
}

func readPeople() (people []Person, err error) {
	lines, err := getLines(peopleFile)
	if err != nil {
		return
	}

	for _, line := range lines {
		var person Person
		err := json.Unmarshal([]byte(line), &person)

		if err != nil {
			break
		}
		people = append(people, person)
	}
	return
}

func distinctFirstNames(people []Person) {
	set := make(map[string]struct{})
	for _, person := range people {
		set[person.FirstName] = struct{}{}
	}
	fmt.Printf("\nThere are %d distinct first names.\n", len(set))
}

func minAndMaxAge(people []Person) {
	minAge, maxAge := people[0].Age, people[0].Age
	for _, person := range people {
		if minAge > person.Age {
			minAge = person.Age
		}
		if maxAge < person.Age {
			maxAge = person.Age
		}
	}
	fmt.Printf("\nMinimum age is %d.\n", minAge)
	fmt.Printf("Maximum age is %d.\n", maxAge)
}

func oldestPeople(people []Person, topCount int) {
	sort.Slice(people, func(i, j int) bool { return people[i].Age > people[j].Age })
	fmt.Printf("\nTop %d oldest people:\n", topCount)
	for _, p := range people[:topCount] {
		fmt.Println(p)
	}
}

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func popularLastNames(people []Person, topCount int) {
	fmt.Printf("\nTop %d most popular last names with frequency:\n", topCount)
	lastNameCounts := make(map[string]int)
	for _, p := range people {
		count, _ := lastNameCounts[p.LastName]
		lastNameCounts[p.LastName] = count + 1
	}
	pairList := make(PairList, len(lastNameCounts))
	i := 0
	for k, v := range lastNameCounts {
		pairList[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pairList))
	for _, p := range pairList[:topCount] {
		fmt.Printf("%s, %d\n", p.Key, p.Value)
	}
}

func getLines(fileName string) (lines []string, err error) {

	file, err := os.Open(fileName)
	if err != nil {
		return
	}

	defer func() {
		err = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	err = scanner.Err()
	return
}
