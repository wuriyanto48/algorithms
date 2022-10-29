package main

import (
	"fmt"
	"github.com/wuriyanto48/productrecommendationknn/fileutil"
	"github.com/wuriyanto48/productrecommendationknn/knn"
	"os"
	"strconv"
	"strings"
)

func main() {
	labelledProducts, err := GetLabelledProductData()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	productDatabases, err := GetProductData()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	search := "Oopsy Daisy Birdie Family Growth Chart by Patchi Cancado, 12 by 42-Inch"
	selectedProduct, ok := productDatabases[strings.ToLower(search)]
	if !ok {
		fmt.Printf("product %s not available\n", search)
		os.Exit(1)
	}

	fmt.Println("------- search product ------")
	fmt.Println(selectedProduct)

	fmt.Println("------- recommended product ------")
	kNearest := knn.NewKNN(labelledProducts)

	recommendedProducts := kNearest.Classify(20, selectedProduct)

	for _, recommended := range recommendedProducts {
		fmt.Println("Distance: ", recommended.Distance, "Title: ", recommended.Product.ProductName)
	}

}

func GetLabelledProductData() (knn.LabelledDatasets, error) {
	labelledCSVProductDatas, err := fileutil.ReadCSV("./productdataset/product_with_category.csv")
	if err != nil {
		return nil, err
	}

	var labelledDatasets knn.LabelledDatasets
	for _, record := range labelledCSVProductDatas {
		var (
			product      knn.Product
			labelledData knn.LabelledDataset
		)

		var categories []int
		categoryNums := record[2:]
		for _, cNumStr := range categoryNums {
			catNum, err := strconv.Atoi(cNumStr)
			if err != nil {
				return nil, err
			}

			categories = append(categories, catNum)
		}

		product.ProductName = record[0]
		product.Url = record[1]
		product.CategoryNum = categories

		labelledData.Product = &product
		labelledData.Distance = 0.0
		labelledDatasets = append(labelledDatasets, labelledData)
	}

	return labelledDatasets, nil
}

func GetProductData() (map[string]*knn.Product, error) {
	productDatas, err := fileutil.ReadCSV("./productdataset/product_with_category.csv")
	if err != nil {
		return nil, err
	}

	trash := []byte{194, 160}

	productDataMap := make(map[string]*knn.Product)

	for _, record := range productDatas {
		product := &knn.Product{}

		product.ProductName = strings.Trim(record[0], " ")
		product.Url = strings.Trim(record[1], " ")
		var categories []int
		categoryNums := record[2:]
		for _, cNumStr := range categoryNums {
			catNum, err := strconv.Atoi(cNumStr)
			if err != nil {
				return nil, err
			}

			categories = append(categories, catNum)
		}

		product.CategoryNum = categories

		productDataMap[strings.ToLower(strings.Trim(product.ProductName, string(trash)))] = product
	}

	return productDataMap, nil

}
