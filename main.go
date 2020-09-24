package main

import (
	"fmt"
        "github.com/gagiopapinni/price-tracker-api-exercise/email"
        "github.com/gagiopapinni/price-tracker-api-exercise/helper"
        "github.com/gagiopapinni/price-tracker-api-exercise/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
) 


var collection *mongo.Collection
var config helper.Config
var cache models.AdCache
var confirmPool email.ConfirmationPool

func RefreshAndNotify (){
        fmt.Print(".")
	for url, ad := range cache {
		changed, err := ad.FetchPrice()
	 	if err != nil{
			fmt.Printf("\nUnable to fetch price:%s\n", err.Error())
		} else if changed {
			message := fmt.Sprintf(`<p>Price for</p> 		
                                                <p>%s</p> 
                                                <p>has changed, it's now %d RUB </p>`, url, ad.Price)
			err := email.Send(message, ad.Subscribers)
			if err!=nil { 
				fmt.Printf("\nUnable to notify:%s\n", err.Error()) 
			}
			ad.PushToDb(collection)
		}

	}

}


func startNotificationDaemon() {
	for {
		RefreshAndNotify()
		time.Sleep(time.Duration(config.REFRESH_SECS) * time.Second)
	}
}


func main() {
        config.Load()

	email.Configure(config.SENDER_GMAIL_ADDR, config.SENDER_GMAIL_PASS)
	collection = helper.ConnectDb(config.MONGOHQ_URL)
	cache.Init(collection)

	go startNotificationDaemon()

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.POST("/subscribe", func(c *gin.Context) {
		req := struct {
			Url string 
			Email string
		}{}

		err_bind := c.ShouldBind(&req)
		if err_bind != nil {
			c.HTML(400,"Default.html", gin.H{"msg":"Error..."})
			return 
		}

		if !helper.IsValidUrl(req.Url) {
			c.HTML(400,"Default.html", gin.H{"msg":"invalid url"})
			return
		}	

		key := confirmPool.GenerateKey(req.Email)
		url := fmt.Sprintf("http://%s:%d/subscribe?Url=%s&Email=%s&Key=%s",config.DOMAIN,
										   config.PORT,
										   req.Url,
										   req.Email,key)
		message := fmt.Sprintf(`<p> Email Confirmation for %s</p>
					<p> Click the following link to subscribe: %s</p>`,req.Url, url)

		err := email.Send(message, []string{req.Email})
		if err!=nil { 
			c.HTML(400,"Default.html", gin.H{"msg":"could not send confirmation message"})
			return
		}
		
		c.HTML(200,"Default.html", gin.H{"msg":"confirmation message sent"})
	})

	r.GET("/subscribe", func(c *gin.Context) {
		req := struct {
			Url string 
			Email string
			Key string
		}{}

		err_bind := c.ShouldBindQuery(&req)
		if err_bind != nil {
			c.HTML(404,"Default.html", gin.H{"msg":"Not found"})
			return 
		}

		confirmed := confirmPool.Confirm(req.Email,req.Key)
		if !confirmed {
			c.HTML(400,"Default.html", gin.H{"msg":"Expired link"})
			return 
		}

		if ad, ok := cache[req.Url]; ok {
			done := ad.AddSubscriber(req.Email)	
			if !done {
				c.HTML(200, "Default.html", gin.H{"msg": "already subscribed"})
				return	
			}

			err_push := ad.PushToDb(collection)
			if err_push != nil {
				c.HTML(400, "Default.html", gin.H{"msg": err_push.Error()})
				return			
			}

			c.HTML(200,"Default.html", gin.H{"msg": "Subscribed!"})
			return		
		} else {
			ad := models.Ad{ 
				Url: req.Url,
				Subscribers: []string{req.Email},
			}
			
			_, err_fetch := ad.FetchPrice()
			if err_fetch != nil {
				c.HTML(400, "Default.html", gin.H{"msg": err_fetch.Error()})
				return
			}

			err_push := ad.PushToDb(collection)
			if err_push != nil {
				c.HTML(400, "Default.html", gin.H{"msg": err_push.Error()})
				return
			}

			cache[ad.Url] = &ad

			c.HTML(200,"Default.html", gin.H{"msg": "Subscribed!"})
			return		

		}

				
	})

	r.GET("/", func(c *gin.Context) {
		c.HTML(200,"landing.html", gin.H{ "action":"/subscribe" })
	})

	r.Run(fmt.Sprintf(":%d",config.PORT)) 
} 





