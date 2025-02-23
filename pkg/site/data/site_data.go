package site_data

import (
	"encoding/json"
	"fmt"
	"go-auth-service/pkg/site/model"
	"io"
	"log"
	"os"
)

func SiteData() []site_model.Site {
	file, err := os.Open("./pkg/site/data/siteData.json")
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

	var sites []site_model.Site
	err = json.Unmarshal(bytes, &sites)
	if err != nil {
		log.Println(fmt.Sprintf("❌ Error unmarshalling JSON: %v", err))
		return nil
	}

	return sites
}
