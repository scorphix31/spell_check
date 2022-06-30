package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"spell_check/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const serviceURL = "https://speller.yandex.net/services/spellservice.json/checkText"

type SpellController struct {
	logger *zap.Logger
}

func NewSpellController(logger *zap.Logger) *SpellController {
	return &SpellController{logger: logger}
}

func spellCheck(prop string, ch chan<- string, logger *zap.Logger) {
	resp, err := http.PostForm(serviceURL, url.Values{
		"text": {prop},
	})
	if err != nil {
		logger.Sugar().Error(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Sugar().Error(err)
	}

	err = resp.Body.Close()
	if err != nil {
		logger.Sugar().Error(err)
	}

	var respText []models.Text
	if err = json.Unmarshal(body, &respText); err != nil {
		logger.Sugar().Error(err)
	}

	builder := strings.Builder{}
	runeProp := []rune(prop)
	builder.Grow(len(prop) + len(prop)/2)
	var addText string
	var lastIndexProp int
	for _, str := range respText {
		addText = string(runeProp[lastIndexProp:str.Pos])
		builder.WriteString(addText)
		builder.WriteString(str.Suggestions[0])
		lastIndexProp = str.Pos + str.Len
	}
	endPath := string(runeProp[lastIndexProp:])
	builder.WriteString(endPath)
	//logger.Info("get correct proposal", zap.String("string", builder.String()))
	ch <- builder.String()
	close(ch)
}

func (sc *SpellController) CheckText(c *gin.Context) {
	var request models.RequestTask
	if err := c.ShouldBindJSON(&request); err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}
	props := request.Props
	channels := make([]chan string, len(props))
	for i := range channels {
		channels[i] = make(chan string)
	}
	for i := range props {
		go spellCheck(props[i], channels[i], sc.logger)
	}
	response := make([]string, 0, len(props))
	for _, c := range channels {
		response = append(response, <-c)
	}
	c.JSON(200, response)
}
