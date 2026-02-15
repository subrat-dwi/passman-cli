package agent

type Request struct {
	Cmd        string `json:"cmd"`
	Key        string `json:"key,omitempty"`
	TTL        int    `json:"ttl,omitempty"`
	Ciphertext string `json:"ciphertext,omitempty"`
	Nonce      string `json:"nonce,omitempty"`
	Plaintext  string `json:"plaintext,omitempty"`
}

type Response struct {
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
	Data  any    `json:"data,omitempty"`
	Nonce string `json:"nonce,omitempty"`
}
