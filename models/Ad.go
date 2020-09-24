package models

import(
        "github.com/gagiopapinni/price-tracker-api-exercise/helper"
        "github.com/gagiopapinni/price-tracker-api-exercise/price"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"context"
	"errors"
	
)


type Ad struct {
	Url string `json:"url"`
	Subscribers []string `json:"subscribers"`
        Price uint64 `json:"price"`
}

type AdCache map[string]*Ad


func (c *AdCache) Init(collection *mongo.Collection) error {
	if *c==nil{
		*c = make( AdCache )
	}

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
   		return errors.New("Unknown Error")
	}

	var results []Ad
	if err = cursor.All(context.TODO(), &results); err != nil {
   		return errors.New("Unknown Error")
	}

	for _, ad := range results {
		copyOfAd := ad
		(*c)[ad.Url] = &copyOfAd
	}

	return nil
}


func (ad Ad) IsValid() bool {
	if err := ad.validate(); err != nil {
		return false
	}
	return true
}


func (ad Ad) doesSubscriberExist(email string) (int, bool) {
    for i, item := range ad.Subscribers {
        if item == email {
            return i, true
        }
    }
    return -1, false
}

func (ad *Ad) AddSubscriber(email string) bool {
	_, exists := ad.doesSubscriberExist(email)
	if !exists {
		ad.Subscribers = append(ad.Subscribers, email)
		return true
	}
	return false
}

func (ad Ad) validate() error {
	if !helper.IsValidUrl(ad.Url) {
		return errors.New("invalid url")
	}
	if ad.Price < 0 {
		return errors.New("negative price")
	}
	return nil
}

func (ad *Ad) FetchPrice() (bool, error) {
	currentPrice, err := price.Extract(ad.Url)
	if err != nil {
		return false, err
	}

	oldPrice := ad.Price
        ad.Price = currentPrice

	if oldPrice==currentPrice {
		return false, nil
	}

	return true, nil
}

func (ad Ad) PushToDb(collection *mongo.Collection) error {
	err_valid := ad.validate()
	if err_valid != nil {
		return err_valid
	}
	   
	opts := options.Replace().SetUpsert(true)
	filter := bson.D{{"url", ad.Url}}
	
	_, err := collection.ReplaceOne(context.TODO(), filter, ad, opts)
	if err != nil {
		return errors.New("Internal Error")
	}
       
	return nil
}
