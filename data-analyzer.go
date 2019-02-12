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

func main() {
	rand.Seed(time.Now().Unix())
	writePeople(getFirstNames(), getLastNames())
	analyze()
}

type Person struct {
	FirstName string
	LastName  string
	Age       int
}

func (p Person) String() string {
	return fmt.Sprintf("Person(FirstName=%s, LastName=%s, age=%d)", p.FirstName, p.LastName, p.Age)
}

const PeopleFile = "data/people.txt"
const PeopleCount = 1000

func writePeople(firstNames, lastNames []string) {
	file, err := os.Create(PeopleFile)

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		closeErr := file.Close()

		if closeErr != nil {
			log.Fatal(closeErr)
		}
	}()

	for i := 0; i < PeopleCount; i++ {
		person := &Person{
			FirstName: firstNames[rand.Intn(len(firstNames))],
			LastName:  lastNames[rand.Intn(len(lastNames))],
			Age:       rand.Intn(100) + 1,
		}

		personBytes, err := json.Marshal(person)
		if err != nil {
			fmt.Println(err)
			continue
		}

		_, err = fmt.Fprintln(file, string(personBytes))
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Printf("A dataset with %d people generated.\n", PeopleCount)
}

func getFirstNames() (names []string) {
	lines := GetLines("data/firstNames.txt")

	for _, line := range lines {
		tokens := strings.Split(line, "\t")
		names = append(names, tokens[1], tokens[3])
	}
	return
}

func getLastNames() (names []string) {
	lines := GetLines("data/lastNames.txt")

	for _, line := range lines {
		tokens := strings.Split(line, " ")
		names = append(names, strings.Title(strings.ToLower(tokens[0])))
	}
	return
}

func analyze() {
	people := readPeople()

	const topCount = 3
	distinctFirstNames(people)
	minAndMaxAge(people)
	oldestPeople(people, topCount)
	popularLastNames(people, topCount)
}

func readPeople() (people []Person) {
	lines := GetLines(PeopleFile)

	for _, line := range lines {
		var person Person
		err := json.Unmarshal([]byte(line), &person)

		if err == nil {
			people = append(people, person)
		}
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

func GetLines(fileName string) (lines []string) {

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		closeErr := file.Close()

		if closeErr != nil {
			log.Fatal(closeErr)
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return
}
