package telegram

import (
	"article-advisor/lib/er"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

// Client collects event for Fetcher using telegram bot API.
type Client struct {
	//host is a host telegram API.
	host string
	//prefix at the beginning of each request.
	basePath string
	//client is a http.Client. Clients should be reused instead of created as needed.
	client http.Client
}

const (
	// getUpdatesMethod is name of the method for generating the URL request for Update.
	getUpdatesMethod = "getUpdates"
	// sendMessageMethod is name of the method for generating the URL request for send message to user.
	sendMessageMethod = "sendMessage"
)

// Constructor telegram.client.
func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

// newBasePath generates path.
func newBasePath(token string) string {
	return "bot" + token
}

//Clients methods. methods that the client will execute when the application is running:

// Updates receives new messages. Returns the slice of structure Update from field UpdatesResponse.Result
// structure UpdatesResponse and error.
func (c *Client) Updates(offset int, limit int) (updates []Update, err error) {
	defer func() { err = er.WrapIfErr("can't get updates", err) }()

	//generation of request parameters.
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	//limit is a count of updates received per one request.
	q.Add("limit", strconv.Itoa(limit))

	//Getting data from the response.
	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}

	var res UpdatesResponse

	//parsing the response from json.
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

// SendMessage sends messages to user.
func (c *Client) SendMessage(chatID int, text string) error {

	//generation of request parameters.
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return er.Wrap("can't send message", err)
	}

	return nil
}

// doRequest sends a formed request. Method returns the response data and error.
func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	defer func() { err = er.WrapIfErr("can't do request", err) }()

	//generating the URL to which the request is sent.
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	//preparing the Request Object.
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	//passing request parameters(query) to request req.
	req.URL.RawQuery = query.Encode()

	//send request.
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	//close response body of request.
	defer func() { _ = resp.Body.Close() }()

	//getting body of response.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
