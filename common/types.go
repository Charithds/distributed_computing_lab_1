package common

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/subchen/go-trylock/v2"
)

type Vegetable struct {
	Name  string  `csv:"Name"`
	Price float64 `csv:"Price"`
	QTY   float64 `csv:"QTY"`
}

var mu = trylock.New()

type Store struct {
	database map[string]Vegetable // private
}

func (s Vegetable) Details() string {
	price := strconv.FormatFloat(s.Price, 'f', 2, 64)
	qty := strconv.FormatFloat(s.QTY, 'f', 2, 64)
	return s.Name + " Price: " + price + " Available QTY: " + qty
}

// Get methods returns a vegetable with specific id (procedure).
func (c *Store) Get(payload string, reply *Vegetable) error {
	fmt.Println("Getting values")
	all := readCsvFile()
	if all == nil {
		return fmt.Errorf("error while reading the file")
	}

	foundVeg := getVeg(payload, all)
	if foundVeg != nil {
		*reply = *foundVeg
	} else {
		return fmt.Errorf("vegetable with id '%s' does not exist", payload)
	}
	return nil
}

// Get methods returns all vegetables.
func (c *Store) GetAll(payload string, reply *[]*Vegetable) error {
	fmt.Println("Get all values")
	*reply = readCsvFile()
	return nil
}

// Add.
func (c *Store) AddVeg(payload Vegetable, reply *[]*Vegetable) error {
	fmt.Println("Adding values")
	all := readCsvFile()
	foundVeg := getVeg(payload.Name, all)
	if foundVeg != nil {
		return fmt.Errorf("Vegetable already exists")
	}
	*reply = writeCsvFile(all, &payload)

	return nil
}

// Update.
func (c *Store) UpdateVeg(payload Vegetable, reply *[]*Vegetable) error {
	fmt.Println("Updating values")
	all := readCsvFile()
	foundVeg := getVeg(payload.Name, all)
	if foundVeg == nil {
		return fmt.Errorf("Vegetable does not exist")
	}
	foundVeg.Price = payload.Price
	foundVeg.QTY = payload.QTY
	*reply = writeCsvFile(all, nil)

	return nil
}

func writeCsvFile(previousVegs []*Vegetable, newVeg *Vegetable) []*Vegetable {
	if ok := mu.TryLockTimeout(3 * time.Second); !ok {
		return nil
	}

	err := os.Remove("server/data.csv")

	if err != nil {
		fmt.Println(err)
		return nil
	}
	f, err := os.OpenFile("server/data.csv", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Unable to open data file ", err)
	}
	writer := csv.NewWriter(f)
	writer.Write([]string{"Name", "Price", "QTY"})

	for _, veg := range previousVegs {
		row := []string{veg.Name, fmt.Sprintf("%f", veg.Price), fmt.Sprintf("%f", veg.QTY)}
		_ = writer.Write(row)
	}
	if newVeg != nil {
		newRow := []string{newVeg.Name, fmt.Sprintf("%f", newVeg.Price), fmt.Sprintf("%f", newVeg.QTY)}
		_ = writer.Write(newRow)
	}
	writer.Flush()
	f.Close()
	mu.Unlock()
	return readCsvFile()
}

func readCsvFile() []*Vegetable {
	if ok := mu.RTryLockTimeout(1 * time.Second); !ok {
		return nil
	}
	defer mu.RUnlock()

	f, err := os.OpenFile("server/data.csv", os.O_CREATE|os.O_RDONLY, 0644)

	if err != nil {
		log.Fatal("Unable to read input file server/data.csv", err)
	}
	defer f.Close()

	var vegetables []*Vegetable

	if err := gocsv.UnmarshalFile(f, &vegetables); err != nil {
		log.Print(err)
		return nil
	}

	return vegetables
}

func getVeg(name string, vegetables []*Vegetable) *Vegetable {
	for _, veg := range vegetables {
		if strings.EqualFold(name, veg.Name) {
			return veg
		}
	}
	return nil
}

// Store function returns a new instance of College (pointer).
func NewStore() *Store {
	return &Store{
		database: make(map[string]Vegetable),
	}
}
