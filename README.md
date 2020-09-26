# price-tracker-api-exercise

### config.json
```
{     
   "PORT": <INT>,   
   "DOMAIN": <STRING>,   
   "MONGOHQ_URL": <STRING MongoDb URL>,    
   "SENDER_GMAIL_ADDR": <STRING gmail address>,    
   "SENDER_GMAIL_PASS": <STRING gmail password>,   
   "REFRESH_SECS": <INT seconds>  
}   
```
### Running
go 1.15 is required   
Once the config.json is created, run in project dir:      
```
go install
go run main.go
```


### Overview
there are three endpoints:
* / [GET]     
returns the landing page with a form 
* /subscribe [POST]   
Accepts a JSON object with Url of avito ad, and Email for subscribtion.    
Sends a confirmation link to specified Email.           
The link contains a confirmation key and points to the /subscribe [GET] endpoint   
Returns a JSON object with result or error   

* /subscribe [GET]    
Finishes the confirmation process and returns a simple page with the result or error       



### Severely uninformative Bpmn diagram    


<img src="https://github.com/gagiopapinni/price-tracker-api-exercise/blob/master/Diagram.png" >

