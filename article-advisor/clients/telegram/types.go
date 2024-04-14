package telegram

// UpdatesResponse is the structure with response from tg bot API with results.
// The response to each request is formed in such structures.
type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

// Update is the structure with the necessary data, received from the tg bot API.
type Update struct {
	ID      int              `json:"update_id"`
	Message *IncomingMessage `json:"message"`
}

// IncomingMessage is the structure for Update structure. It consists of three fields - Text, From, Chat.
// Where Text is string data received from the user, From is nickname of user, Chat - chat ID.
type IncomingMessage struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

// From is nickname of user.
type From struct {
	Username string `json:"username"`
}

// Chat is chat ID.
type Chat struct {
	ID int `json:"id"`
}
