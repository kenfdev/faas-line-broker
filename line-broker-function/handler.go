package function

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

type LineMessage struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type LineWebhookRequest struct {
	Events []*struct {
		ReplyToken string       `json:"replyToken"`
		Message    *LineMessage `json:"message"`
	} `json:"events"`
}

type LineReplyRequest struct {
	ReplyToken string         `json:"replyToken"`
	Messages   []*LineMessage `json:"messages"`
}

type DialogFlowResponse struct {
	Speech string `json:"speech"`
}

// Handle a serverless request
func Handle(req []byte, wg *sync.WaitGroup) string {

	var webhookReq LineWebhookRequest
	err := json.Unmarshal(req, &webhookReq)
	if err != nil {
		panic(err)
	}

	wg.Add(1)
	go func() {
		handleLineRequest(webhookReq)
		wg.Done()
	}()

	// immediate response
	return "OK"
}

func postDialogFlow(text string) DialogFlowResponse {
	dfBrokerURL := "http://gateway:8080/function/dummy-df-broker-function"

	values := map[string]string{"text": text}
	jsonValue, _ := json.Marshal(values)
	resp, err := http.Post(dfBrokerURL, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var dfResp DialogFlowResponse
	err = json.Unmarshal(body, &dfResp)
	if err != nil {
		panic(err)
	}

	return dfResp
}

func postLineReply(token string, speech string) {
	replyRequest := &LineReplyRequest{
		ReplyToken: token,
		Messages: []*LineMessage{
			&LineMessage{
				Type: "text",
				Text: speech,
			},
		},
	}

	postData, _ := json.Marshal(replyRequest)
	client := &http.Client{}

	lineReplyURL := "https://api.line.me/v2/bot/message/reply"
	req, err := http.NewRequest("POST", lineReplyURL, bytes.NewReader(postData))
	if err != nil {
		panic(err)
	}

	lineToken := os.Getenv("LINE_ACCESS_TOKEN")
	req.Header.Add("Authorization", "Bearer "+lineToken)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

}

func handleLineRequest(r LineWebhookRequest) {
	event := r.Events[0]

	dfResp := postDialogFlow(event.Message.Text)

	postLineReply(event.ReplyToken, dfResp.Speech)

	return
}
