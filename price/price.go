package price

import(
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"errors"
)

func Extract(url string) (uint64, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0,errors.New("Could not get specified page")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)


	priceLine := regexp.MustCompile("avito.item.price *=.*")
	line := priceLine.Find(body)

        priceValue := regexp.MustCompile("[0-9]+")
        value := priceValue.Find(line)
 
	price, err := strconv.ParseUint(string(value),10,0)
	if err != nil{
		return 0, errors.New("Unknown error")
	}

	return price, nil
}
