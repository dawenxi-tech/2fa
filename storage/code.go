package storage

import (
	"encoding/json"
	"log/slog"
	"net/url"
	"os"
	"strings"
)

type Code struct {
	Type      string       `json:"type"`
	Name      string       `json:"name"`
	Algorithm string       `json:"algorithm"`
	Digits    string       `json:"digits"`
	Issuer    string       `json:"issuer"`
	Period    string       `json:"period"`
	Secret    *EncryptData `json:"secret"`
}

func ParseCode(rawUri string) (*Code, error) {
	uri, err := url.Parse(rawUri)
	if err != nil {
		return nil, err
	}
	value := uri.Query()
	code := &Code{
		Type:      uri.Hostname(),
		Name:      strings.TrimPrefix(uri.Path, "/"),
		Algorithm: value.Get("algorithm"),
		Digits:    value.Get("digits"),
		Issuer:    value.Get("issuer"),
		Period:    value.Get("period"),
		Secret:    NewEncryptData(value.Get("secret")),
	}
	return code, nil
}

func (c *Code) Encode() string {
	values := url.Values{
		"algorithm": []string{c.Algorithm},
		"digits":    []string{c.Digits},
		"issuer":    []string{c.Issuer},
		"period":    []string{c.Period},
		"secret":    []string{c.Secret.Val()},
	}
	uri := url.URL{
		Scheme:   "otpauth",
		Host:     c.Type,
		Path:     c.Name,
		RawQuery: values.Encode(),
	}
	return uri.String()
}

type Codes []Code

func LoadCodes() Codes {
	var code Codes
	data, err := os.ReadFile(codePath())
	if err != nil {
		slog.Error("error to read code file", slog.Any("err", err))
		return nil
	}
	err = json.Unmarshal(data, &code)
	if err != nil {
		slog.Error("error to unmarshal codes", slog.Any("err", err))
		return nil
	}
	return code
}

func SaveCode(codes Codes) {
	data, err := json.Marshal(codes)
	if err != nil {
		slog.Error("error to marsha codes", slog.Any("err", err))
		return
	}
	fp, err := os.Create(codePath())
	if err != nil {
		slog.Error("error to create file", slog.Any("err", err))
		return
	}
	_, err = fp.Write(data)
	if err != nil {
		slog.Error("error to write codes", slog.Any("err", err))
		return
	}
}
