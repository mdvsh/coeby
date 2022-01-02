package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func SaveDB(fname string, data interface{}) {
	jsonFile, err := os.Create("data/" + fname)
	if err != nil {
		log.Printf("Error creating json file\n %s", err)
		panic(err)
	}
	defer jsonFile.Close()
	jsonWriter := bufio.NewWriter(jsonFile)
	defer jsonWriter.Flush()
	enc := json.NewEncoder(jsonWriter)
	enc.SetIndent("", "  ")
	enc.Encode(data)
	fmt.Printf("Saved %s\n", fname)
}

func LoadDeptDB(depts interface{}) {
	f, err := os.Open("data/depts.json")
	if err != nil {
		log.Printf("Error opening json file\n %s", err)
		panic(err)
	}
	defer f.Close()
	jsonParser := json.NewDecoder(f)
	jsonParser.Decode(&depts)
}
