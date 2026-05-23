package utils

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Sna-ken/videoweb/config"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/protocol"
)

type Result struct {
	Name    string `json:"name"`
	Snippet string `json:"snippet"`
	URL     string `json:"url"`
}

func Search(ctx context.Context, word string) (string, error) {
	apiKey := config.Search.SearchAPIKey
	searchurl := "https://api.qnaigc.com/v1/search/web"

	// 创建client
	c, err := client.NewClient()
	if err != nil {
		return "", err
	}

	// 创建请求与响应对象
	req := &protocol.Request{}
	resp := &protocol.Response{}

	req.SetMethod("POST")
	req.SetRequestURI(searchurl)

	body := map[string]interface{}{
		"query":       word + "什么意思",
		"max_results": 2,
		"search_type": "web",
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	req.SetBody(bodyBytes)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	err = c.Do(ctx, req, resp)
	if err != nil {
		return "", err
	}

	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("status error:%d", resp.StatusCode())
	}

	var raw struct {
		Data struct {
			Results []Result `json:"results"`
		} `json:"data"`
	}

	if err := json.Unmarshal(resp.Body(), &raw); err != nil {
		return "", err
	}

	if len(raw.Data.Results) == 0 {
		return "", nil
	}
	return raw.Data.Results[0].Snippet, nil
}
