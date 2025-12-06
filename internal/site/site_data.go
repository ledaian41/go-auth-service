package site

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

func GetData() []Site {
	file, err := os.Open("./internal/site/siteData.json")
	if err != nil {
		log.Println(fmt.Sprintf("❌ Error opening file: %v", err))
		return nil
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Println(fmt.Sprintf("❌ Error reading file: %v", err))
		return nil
	}

	var sites []Site
	err = json.Unmarshal(bytes, &sites)
	if err != nil {
		log.Println(fmt.Sprintf("❌ Error unmarshalling JSON: %v", err))
		return nil
	}

	return sites
}
