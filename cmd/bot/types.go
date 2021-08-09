package main

type ivrAPIChatDelayResponse struct {
	Status   int    `json:"status"`
	Error    string `json:"error"`
	Username string `json:"username"`
	Delay    int    `json:"delay"`
}
