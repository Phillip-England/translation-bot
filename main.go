package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Message struct {
	Text string `json:"text"`
	GroupID string `json:"group_id"`
}

type TranslationResponse struct {
	Translations []struct {
		DetectedSourceLanguage string `json:"detected_source_language"`
		Text                   string `json:"text"`
	} `json:"translations"`
}

type GroupmeRequestBody struct {
	Text  string `json:"text"`
	BotID string `json:"bot_id"`
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func main() {

	godotenv.Load()

	// setting up http server with cors
	r := gin.New()
	r.Use(CORSMiddleware())

	// building route to handle all incoming groupme messages
	r.POST("/", func(c *gin.Context) {

		// grabbing message from incoming groupme request
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		// unwrapping message and getting raw text
		var message Message
		err = json.Unmarshal(body, &message)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		text := message.Text
		groupID := message.GroupID

		// setting up our bot ID based off the groupID
		var groupmeBotID string

		// SOUTHROADS LEADERSHIP
		if groupID == os.Getenv("SOUTHROADS_LEADERSHIP_GROUP_ID") {
			groupmeBotID = os.Getenv("SOUTHROADS_LEADERSHIP_BOT_ID")
		}

		// SOUTHRODS AMAZON/SUPPLY CENTRAL
		if groupID == os.Getenv("SOUTHROADS_SUPPLY_ORDER_GROUP_ID") {
			groupmeBotID = os.Getenv("SOUTHROADS_SUPPLY_ORDER_BOT_ID")
		}

		// UTICA KITCHEN LEADERSHIP
		if groupID == os.Getenv("UTICA_KITCHEN_LEADERSHIP_GROUP_ID") {
			groupmeBotID = os.Getenv("UTICA_KITCHEN_LEADERSHIP_BOT_ID")
		}

		// SOUTHROADS KITCHEN LEADERSHIP
		if groupID == os.Getenv("SOUTHROADS_KITCHEN_LEADERSHIP_GROUP_ID") {
			groupmeBotID = os.Getenv("SOUTHROADS_KITCHEN_LEADERSHIP_BOT_ID")
		}

		// UTICA FOH LEADERSHIP
		if groupID == os.Getenv("UTICA_FOH_LEADERSHIP_GROUP_ID") {
			groupmeBotID = os.Getenv("UTICA_FOH_LEADERSHIP_BOT_ID")
		}

		// UTICA SUPPLY ORDER
		if groupID == os.Getenv("UTICA_SUPPLY_ORDER_GROUP_ID") {
			groupmeBotID = os.Getenv("UTICA_SUPPLY_ORDER_BOT_ID")
		}

		// UTICA REPAIRS
		if groupID == os.Getenv("UTICA_REPAIRS_GROUP_ID") {
			groupmeBotID = os.Getenv("UTICA_REPAIRS_BOT_ID")
		}

		// SOUTHROADS TRAINING
		if groupID == os.Getenv("SOUTHROADS_TRAINING_GROUP_ID") {
			groupmeBotID = os.Getenv("SOUTHROADS_TRAINING_BOT_ID")
		}
	
		// checking if we are translating to spanish or english
		keyword := string(text[:8])
		toSpanish := false
		toEnglish := false
		if keyword == "$spanish" || keyword == "$Spanish" {
			toSpanish = true
		}
		if keyword == "$english" || keyword == "$English" {
			toEnglish = true
		}

		// if we dont get a keyword, exit
		if !toSpanish && !toEnglish {
			c.JSON(http.StatusBadRequest, gin.H{"message": "no keyword provided"})
			return
		}

		// grabbing substring (excluding the #spanish or #english from message translation)
		var subString string
		if toEnglish || toSpanish {
			subString = string(text[9:])
		}

		// setting the target language
		var targetLanguage string
		if toEnglish {
			targetLanguage = "EN"
		} else if toSpanish {
			targetLanguage = "ES"
		}

		// setting up variables for translation api request
		apiurl := "https://api-free.deepl.com/"
		resource := "/v2/translate"
		data := url.Values{
			"text":        {subString},
			"target_lang": {targetLanguage},
		}
		u, _ := url.ParseRequestURI(apiurl)
		u.Path = resource
		urlStr := u.String()

		// building request object
		req, err := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode()))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// setting headers
		key := os.Getenv("DEEPL_API_KEY")
		authHeader := fmt.Sprintf("DeepL-Auth-Key %s", key)
		req.Header.Set("Authorization", authHeader)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// building client and performing request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		// Decode the response into a TranslationResponse struct
		var translationResponse TranslationResponse
		err = json.NewDecoder(resp.Body).Decode(&translationResponse)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Access the translation text
		translatedText := translationResponse.Translations[0].Text

		// setting up variables for groupme request
		groupmeRequestURL := "https://api.groupme.com/v3/bots/post"

		// creating json body for groupme request
		groupmeRequestBody := GroupmeRequestBody{
			Text:  translatedText,
			BotID: groupmeBotID,
		}
		requestBody, err := json.Marshal(groupmeRequestBody)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// creating http request to post groupme message
		req, err = http.NewRequest("POST", groupmeRequestURL, bytes.NewBuffer(requestBody))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		req.Header.Set("Content-Type", "application/json")

		// Make the request with an HTTP client
		client = &http.Client{}
		resp, err = client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

	})

	r.Run()

}
