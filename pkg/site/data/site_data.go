package site_data

import (
	"encoding/json"
	"fmt"
	site_model "go-auth-service/pkg/site/model"
	"io"
	"os"
)

func SiteData() []site_model.Site {
	file, err := os.Open("./pkg/site/data/siteData.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}

	var sites []site_model.Site
	err = json.Unmarshal(bytes, &sites)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil
	}

	return sites
}
