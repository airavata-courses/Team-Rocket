package controllers

import (
     "encoding/json"	 
	 "net/http"
	 "io/ioutil"
	 "fmt"
     "RecommendationEngine/utils"
     "os"
)

/* OLD APPROACH - gets the details from ENV through Kubernetes now
type Configuration struct {
    API_GATEWAY_DOMAIN string
    API_GATEWAY_PORT   string
}
*/

// Struct for converting the JSON response to obj
type UserHistory struct {
	Message string `json:"message"`
	Data    []struct {
		ID      string        `json:"_id"`
		UserID  string        `json:"userId"`
		History []interface{} `json:"history"`
		Mood    []interface{} `json:"mood"`
	} `json:"data"`
}

func PreflightGetRecommendedValence(w http.ResponseWriter, r *http.Request){

	 w.Header().Set("Access-Control-Allow-Origin", "*")

    w.Header().Set("Access-Control-Allow-Headers", "*")
    // return "OK"
    json.NewEncoder(w).Encode("OK")
	
}

// Hits UserMgmt micro-service to fetch recent mood/valence counts.
// Finds avg valence and returns that value to caller
// Returns a JSON with avg value as data
// NB: JWT verification happens at the second micro-service for this func
func GetRecommendedValence(w http.ResponseWriter, r *http.Request)  {
	
	 w.Header().Set("Access-Control-Allow-Origin", "*")

	 /* OLD APPROACH - replaced by Kubernetes ENV 
	 // GETS THE IP ADDRESS OF API GATEWAY
		file, _ := os.Open("config.json")
		defer file.Close()
		decoder := json.NewDecoder(file)
		config := Configuration{}
		err := decoder.Decode(&config)
		if err != nil {
		  fmt.Println("error:", err)
		}
	// END OF READING FROM CONFIG.JSON
	*/

	incomingAuthHeader := r.Header.Get("Authorization")
	
	// Extracting token (removes "Bearer ")
	token := incomingAuthHeader[7:]
	fmt.Printf("%s TOKEN: ", token)
	
	// Forming the URL to hit for API GATEWAY - gets from the ENV created by Kubernetes
	USER_MGMT_URL := "http://"+ os.Getenv("APIGATEWAY_SERVICE_PORT_3003_TCP_ADDR")+":" + os.Getenv("APIGATEWAY_SERVICE_PORT_3003_TCP_PORT")
	fmt.Printf("%s API GATEWAY URL: ", USER_MGMT_URL)
	
	// OLD APPROACH: Forming the URL to hit for API GATEWAY
	// USER_MGMT_URL := "http://"+ config.API_GATEWAY_DOMAIN+":" + config.API_GATEWAY_PORT

	// Passing the JWT token along to the next micro-service...
	request, _ := http.NewRequest("GET", USER_MGMT_URL+ "/profile_service/get_history_and_mood", nil)
	request.Header.Set("Authorization", "Bearer "+ token)
	client := &http.Client{}
	response, err := client.Do(request)
	
	
	// Checking if token is valid happens at the other micro-service!!
	// If status code is 200, it means that JWT is correctly verified
	if err != nil || response.StatusCode != 200 {
		 data := map[string]string{"message":"error", "data": "Request failed! Check connection/header! "}
		 utils.Respond(w, data, "FAIL")
		 return
	} else{
		data, _ := ioutil.ReadAll(response.Body)
		var userHist UserHistory	// Struct obj
		
		// Declaring vars - for finding avg
		var sum float64
		var avg float64
		
		json.Unmarshal([]byte(data), &userHist) //now we can use userHist like a normal JSON object
		
		// If no history, returns 200 Success
		if(len(userHist.Data) == 0){
			no_history_data := map[string]string{"message":"No history", "data":"User has no history"}
			utils.Respond(w, no_history_data, "OK")
			return
		}
		
		// There is some history.. Find avg now!
		// Loops through float interface and calculates sum
		for i:=0; i < len(userHist.Data[0].Mood); i++{
			var temp float64
			temp,_ = utils.GetFloat(userHist.Data[0].Mood[i]) // have to convert to float for the add operation below
			sum += temp
		}
		length_of_moods,_ := utils.GetFloat(len(userHist.Data[0].Mood)) //Counts how many elements are present in array
		
		avg = sum/length_of_moods
		
		success_data := map[string]string{"message":"success", "data":fmt.Sprintf("%f", avg) }
		
		utils.Respond(w, success_data, "OK")		
	} 
}
