package har

import (
	"encoding/json"
	"net/url"

	"github.com/linn221/RequesterBackend/models"
	"github.com/linn221/RequesterBackend/utils"
)

// AI generated struct
type HAR struct {
	Log struct {
		Entries []struct {
			StartedDateTime string  `json:"startedDateTime"`
			Time            float64 `json:"time"`
			Request         struct {
				Method  string `json:"method"`
				URL     string `json:"url"`
				Headers []struct {
					Name  string `json:"name"`
					Value string `json:"value"`
				} `json:"headers"`
				PostData struct {
					Text string `json:"text"`
				} `json:"postData"`
			} `json:"request"`
			Response struct {
				Status  int `json:"status"`
				Headers []struct {
					Name  string `json:"name"`
					Value string `json:"value"`
				} `json:"headers"`
				Content struct {
					Text     string `json:"text"`
					Encoding string `json:"encoding,omitempty"`
				} `json:"content"`
			} `json:"response"`
		} `json:"entries"`
	} `json:"log"`
}

func ParseHAR(bs []byte, resHashFunc func(*models.MyRequest) (string, string)) ([]models.MyRequest, error) {
	var har HAR
	if err := json.Unmarshal(bs, &har); err != nil {
		return nil, err
	}

	var result []models.MyRequest

	for i, entry := range har.Log.Entries {
		reqHeaders := make([]models.Header, 0, len(entry.Request.Headers))
		for _, h := range entry.Request.Headers {
			reqHeaders = append(reqHeaders, models.Header{Name: h.Name, Value: h.Value})
		}

		resHeaders := make([]models.Header, 0, len(entry.Response.Headers))
		for _, h := range entry.Response.Headers {
			resHeaders = append(resHeaders, models.Header{Name: h.Name, Value: h.Value})
		}

		u, err := url.Parse(entry.Request.URL)
		domain := ""
		if err == nil {
			domain = u.Hostname()
		}

		resBody := entry.Response.Content.Text
		// Decode base64 if needed
		// if strings.ToLower(entry.Response.Content.Encoding) == "base64" {
		// 	decoded, err := decodeBase64(resBody)
		// 	if err == nil {
		// 		resBody = decoded
		// 	}
		// }

		// Convert HeaderSlice to JSON strings
		reqHeadersJSON, err := models.HeaderSlice(reqHeaders).ToJSON()
		if err != nil {
			return nil, err
		}

		resHeadersJSON, err := models.HeaderSlice(resHeaders).ToJSON()
		if err != nil {
			return nil, err
		}

		my := models.MyRequest{
			Sequence:    i + 1,
			URL:         entry.Request.URL,
			Domain:      domain,
			ReqHeaders:  reqHeadersJSON,
			ReqBody:     entry.Request.PostData.Text,
			ResHeaders:  resHeadersJSON,
			ResStatus:   entry.Response.Status,
			ResBody:     resBody,
			RespSize:    len(resBody),
			LatencyMs:   int64(entry.Time),
			RequestTime: entry.StartedDateTime,
			Method:      entry.Request.Method,
		}

		requestText, responseText := resHashFunc(&my)

		my.ReqHash = utils.HashString(requestText)
		my.ResHash = utils.HashString(responseText)
		my.ResBodyHash = utils.HashString(my.ResBody)
		reqHeadersFromJSON, _ := models.HeaderSliceFromJSON(my.ReqHeaders)
		my.ReqHash1 = utils.HashString(my.Method + " " + my.URL + " " + my.ReqBody + " " + reqHeadersFromJSON.EchoAll())

		result = append(result, my)
	}

	return result, nil
}
