package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/wuriyanto48/productrecommendationknn/fileutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	// writeCategory()

	categoryList := readCategoryToList()

	records, err := fileutil.ReadCSV("amazon_products.csv")
	if err != nil {
		log.Fatal(err)
	}

	cleanedData, err := os.Create("product_with_category.csv")
	defer cleanedData.Close()

	csvOut := csv.NewWriter(cleanedData)
	defer csvOut.Flush()

	var headers []string
	headers = append(headers, []string{"Product Name", "Url"}...)
	headers = append(headers, categoryList...)

	csvOut.Write(headers)

	for i := 0; i < len(records); i++ {
		record := records[i]

		var newRecord []string

		productName := strings.Trim(record[1], " ")
		productUrl := strings.Trim(record[18], " ")

		// append to new record
		newRecord = append(newRecord, productName)
		newRecord = append(newRecord, productUrl)

		productCategory := strings.Trim(record[4], " ")
		categorySplit := strings.Split(productCategory, "|")

		fmt.Printf("name: %s | url: %s | category: %s\n", productName, productUrl, productCategory)

		for _, availableCategory := range categoryList {
			var available bool = false
			for _, c := range categorySplit {
				c = strings.Trim(c, " ")
				if availableCategory == c {
					available = true
				}
			}

			if available {
				newRecord = append(newRecord, "1")
			} else {
				newRecord = append(newRecord, "0")
			}
		}

		csvOut.Write(newRecord)
	}
}

func writeCategory() {
	records, err := fileutil.ReadCSV("amazon_products.csv")
	if err != nil {
		log.Fatal(err)
	}

	var categories = make(map[string]bool)
	var categoryNums = make(map[string]int)

	for i := 0; i < len(records); i++ {
		record := records[i]
		// 9 10 11

		productName := record[1]
		productUrl := record[18]
		productCategory := record[4]

		fmt.Printf("name: %s | url: %s | category: %s\n", productName, productUrl, productCategory)

		categorySplits := strings.Split(productCategory, "|")
		for _, category := range categorySplits {
			category = strings.Trim(category, " ")
			categories[category] = true
		}
	}

	categoryNum := 0
	for k, _ := range categories {
		categoryNums[k] = categoryNum
		categoryNum = categoryNum + 1
	}

	categoryFile, err := os.Create("category.txt")
	defer categoryFile.Close()

	fmt.Println(len(categoryNums))

	for k, v := range categoryNums {
		categoryFile.WriteString(fmt.Sprintf("%s|%d \n", k, v))
	}
}

func readCategoryToMap() map[string]int {

	f, err := os.Open("category.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	fileScanner := bufio.NewScanner(f)

	fileScanner.Split(bufio.ScanLines)

	categoryMap := make(map[string]int)
	for fileScanner.Scan() {
		lineSplits := strings.Split(fileScanner.Text(), " ")

		categoryNum, err := strconv.Atoi(lineSplits[1])
		if err != nil {
			log.Fatal(err)
		}

		categoryMap[lineSplits[0]] = categoryNum
	}

	return categoryMap
}

func readCategoryToList() []string {

	f, err := os.Open("category.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	fileScanner := bufio.NewScanner(f)

	fileScanner.Split(bufio.ScanLines)

	var categoryList []string
	for fileScanner.Scan() {
		lineSplits := strings.Split(fileScanner.Text(), "|")

		categoryList = append(categoryList, strings.Trim(lineSplits[0], " "))
	}

	return categoryList
}
