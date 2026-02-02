package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	FILENAME     = "data.json"
	TODAY_EXISTS = true

	PURPLE = "\033[38;2;125;86;244m"
	PINK = "\033[38;2;247;128;226m"
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

func GetFilePath() string {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	dir := filepath.Dir(exePath)
	dataPath := filepath.Join(dir, "data.json")

	return dataPath
}

// fish pomodoro timer will call this with $work variable which is
// either '25m' or '50m' in the current setup.
//
// This program will also work with plain integer inputs
func main() {
	if len(os.Args) < 2 {
		fmt.Println("duration argument mandatory (type: int)")
		return
	}
	numStr := strings.TrimSuffix(os.Args[1], "m")
	duration, err := strconv.Atoi(numStr)
	if err != nil {
		fmt.Printf("Error converting '%s' to integer", os.Args[1])
		return
	}

	path := GetFilePath()
	file, err := os.ReadFile(path)
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

	if todayExists {
		first.Count += 1
		first.Duration += duration
	} else {
		e := CreateNewEntry(duration)
		entries = append([]Entry{e}, entries...)
	}

	dataToWrite, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling: %+v\n", err)
		return
	}

	err = os.WriteFile(path, dataToWrite, 0600)
	if err != nil {
		fmt.Printf("Error writing file: %+v\n", err)
		return
	}

	fmt.Printf("👷🏼 %s %d %s worked minutes added. Great job! 👷🏼\n", PURPLE, duration, PINK)
	fmt.Printf("🚧 Your total worked minutes is now %s %d 🚧\n", PURPLE, first.Duration)
	fmt.Printf("\033[0m")
}
