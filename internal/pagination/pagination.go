package pagination

import (
	"encoding/base64"
	"strconv"
)

func Encode(unixTs int64) string {
	return base64.StdEncoding.EncodeToString([]byte(strconv.FormatInt(unixTs, 10)))
}

func Decode(token string) (int64, error) {
	if token == "" {
		return 0, nil
	}
	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(string(decoded), 10, 64)
}
