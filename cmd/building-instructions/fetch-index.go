//+build ignore

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func fetchIndexPage(offset int, out interface{}) error {
	const indexURL = "https://www.lego.com//service/biservice/searchbytheme?fromIndex=%d&onlyAlternatives=false&theme=10000-20229"
	var url = fmt.Sprintf(indexURL, offset)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("%s: unexpected status %d", url, resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func cleanURL(url string) string {
	if i := strings.Index(url, "?l.r2="); i > 0 {
		url = url[:i]
	}
	return url
}

func cleanURLs(a []map[string]interface{}) {
	for _, o := range a {
		for k, v := range o {
			if v, isString := v.(string); isString {
				o[k] = cleanURL(v)
			}
		}
	}
}

func main() {

	type Product struct {
		ProductId            string                   `json:"productId"`
		ProductName          string                   `json:"productName"`
		ProductImage         string                   `json:"productImage,omitempty"`
		BuildingInstructions []map[string]interface{} `json:"buildingInstructions"`
	}

	var productsIndex []*Product
	offset := int(0)
	var pageLen int
	for {
		var r struct {
			Products []*Product `json:"products"`
		}
		err := fetchIndexPage(offset, &r)
		if err != nil {
			log.Println(err)
			break
		}
		if len(r.Products) == 0 {
			break
		}
		for i, p := range r.Products {
			// https://www.lego.com/r/www/r/service/-/media/service/service%202015/help%20images/help-placeholder-product.jpg?l.r2=451408231
			if strings.Contains(p.ProductImage, "help-placeholder") {
				p.ProductImage = ""
			}
			cleanURLs(p.BuildingInstructions)
			fmt.Printf("%3d [%s] %s\n", offset+i, p.ProductId, p.ProductName)
		}

		productsIndex = append(productsIndex, r.Products...)

		// Dynamically learn the size of a page.
		// In practice this is just 10.
		if len(r.Products) > 0 && pageLen == 0 {
			pageLen = len(r.Products)
		}

		if len(r.Products) < pageLen {
			break
		}
		offset += len(r.Products)
		time.Sleep(200 * time.Millisecond)
	}
	if len(productsIndex) == 0 {
		log.Fatal("Failed to fetch index.")
	}

	if _, err := os.Stat("data"); err != nil {
		err = os.Mkdir("data", 0766)
		if err != nil {
			log.Fatalln(err)
		}
	}
	indexFile, err := os.Create("data/index.json")
	if err != nil {
		log.Fatalln(err)
	}
	defer indexFile.Close()
	enc := json.NewEncoder(indexFile)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	enc.Encode(productsIndex)
}
