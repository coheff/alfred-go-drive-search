package main

import (
	"encoding/json"

	"golang.org/x/oauth2"
)

const key = "token"

// cacheToken adds a token to Keychain. If a token already exists, it is replaced.
func cacheToken(tok *oauth2.Token) error {
	jToken, err := json.Marshal(tok)
	if err != nil {
		return err
	}

	err = kc.Set(key, string(jToken))
	if err != nil {
		return err
	}
	return nil
}

// cachedToken retrieves a token from Keychain.
func cachedToken() *oauth2.Token {
	jToken, err := kc.Get(key)
	if err != nil {
		return nil
	}

	var tok oauth2.Token
	err = json.Unmarshal([]byte(jToken), &tok)
	if err != nil {
		return nil
	}
	return &tok
}
