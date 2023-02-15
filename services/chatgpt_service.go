package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/ptonlix/officialaccount-chatgpt/init/chatgpt"
)

type ChatGptService struct {
}

const BASEURL = "https://api.openai.com/v1/"

// ChatGPTResponseBody 请求体
type ChatGPTResponseBody struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int                    `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChoiceItem           `json:"choices"`
	Usage   map[string]interface{} `json:"usage"`
}

type ChoiceItem struct {
	Text         string `json:"text"`
	Index        int    `json:"index"`
	Logprobs     int    `json:"logprobs"`
	FinishReason string `json:"finish_reason"`
}

// ChatGPTRequestBody 响应体
type ChatGPTRequestBody struct {
	Model            string  `json:"model"`
	Prompt           string  `json:"prompt"`
	MaxTokens        uint    `json:"max_tokens"`
	Temperature      float64 `json:"temperature"`
	TopP             int     `json:"top_p"`
	FrequencyPenalty int     `json:"frequency_penalty"`
	PresencePenalty  int     `json:"presence_penalty"`
	Echo             bool    `json:"echo"`
}

//Completions gtp文本模型回复
//curl https://api.openai.com/v1/completions
//-H "Content-Type: application/json"
//-H "Authorization: Bearer your chatGPT key"
//-d '{"model": "text-davinci-003", "prompt": "give me good song", "temperature": 0, "max_tokens": 7}'
func (c *ChatGptService) Completions(msg string) (string, error) {
	requestBody := ChatGPTRequestBody{
		Model:            chatgpt.ChatGptConf.Model,
		Prompt:           msg,
		MaxTokens:        uint(chatgpt.ChatGptConf.MaxTokens),
		Temperature:      chatgpt.ChatGptConf.Temperature,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		Echo:             true,
	}
	requestData, err := json.Marshal(requestBody)

	if err != nil {
		return "", err
	}
	logs.Debug("request gpt json string : %v", string(requestData))
	req, err := http.NewRequest("POST", BASEURL+"completions", bytes.NewBuffer(requestData))
	if err != nil {
		return "", err
	}

	apiKey := chatgpt.ChatGptConf.Key
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	client := &http.Client{Timeout: 60 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		logs.Warn("请求GPT出错了, details: %v", err)
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		body, _ := ioutil.ReadAll(response.Body)
		logs.Warn("请求GPT出错了, http code: %d, details: %v", response.StatusCode, string(body))
		return "", fmt.Errorf("gpt api status code not equals 200,code is %d", response.StatusCode)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logs.Warn("GTP读取数据失败, details: %v")
		return "", err
	}

	gptResponseBody := &ChatGPTResponseBody{}

	err = json.Unmarshal(body, gptResponseBody)
	if err != nil {
		logs.Warn("GTP数据解析错误")
		return "", err
	}

	var reply string
	if len(gptResponseBody.Choices) > 0 {
		reply = gptResponseBody.Choices[0].Text
	}
	reply = c.handleText(reply)
	logs.Debug("gpt response text: %s ", reply)
	return reply, nil
}

// 处理回复的数据，除去Prompt数据，使回复更合理
func (c *ChatGptService) handleText(s string) string {
	index := strings.Index(s, "\n\n") //定位两个换行符
	return s[index+2:]
}
