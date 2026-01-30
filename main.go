package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"slices"
	"time"
)

const (
	PATH         = "./data.json"
	TODAY_EXISTS = true
)

type Entry struct {
	Date     time.Time
	Count    int
	Duration int
	Calls    int
}

// Function will check if the first entry is today or yesterday
// Will return error if first is in the future since should not be possible
func DateIsToday(now, first time.Time) (bool, error) {
	d1, m1, y1 := first.Date()
	d2, m2, y2 := now.Date()

	if y1 == y2 && m1 == m2 && d1 == d2 {
		// Date is the same
		return true, nil
	} else if y1 > y2 || m1 > m2 || d1 > d2 {
		// Date is in the future -> error
		return false, fmt.Errorf("Entries[0] cannot be in the future")
	}

	// Date is in the past
	return false, nil
}

func CreateNewEntry(d int) Entry {
	return Entry{
		Date:     time.Now().UTC().Truncate(24 * time.Hour),
		Duration: d,
		Count:    1,
		Calls:    0,
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("duration argument mandatory (type: int)")
		return
	}
	duration, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("Error converting '%s' to integer", os.Args[1])
		return
	}

	file, err := os.ReadFile(PATH)
	if err != nil {
		fmt.Printf("Error reading file: %+v", err)
		return
	}
	entries := make([]Entry, 0)
	err = json.Unmarshal(file, &entries)
	if err != nil {
		fmt.Printf("Error parsing to json: %+v", err)
		return
	}

	// Read the first one (it should be the newest one)
	first := &entries[0]
	today := time.Now()
	todayExists, err := DateIsToday(today, first.Date)
	if err != nil {
		fmt.Printf("Error with dates: %+v", err)
		return
	}

	f, err := os.OpenFile(PATH, os.O_WRONLY, 0600)
	if err != nil {
		fmt.Printf("Error opening file at %s \n", PATH)
	}
	defer f.Close()

	var n int
	if todayExists {
		first.Count += 1
		first.Duration += duration
	} else {
		e := CreateNewEntry(duration)
		entries = append([]Entry{e}, entries...)

		dataToWrite, err := json.Marshal(e)
		if err != nil {
			fmt.Printf("%+v\n", err)
			return
		}

		bytesToWrite := slices.Concat(dataToWrite, []byte{44, 10})
		n, err = f.WriteAt(bytesToWrite, 2)
		if err != nil {
			fmt.Printf("Error writing to a file: %+v\n", err)
			return
		}
	}

	fmt.Printf("Wrote %d bytes to file", n)
}
