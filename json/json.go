package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

func Decode(dest interface{}, r io.Reader) error {
	var buff bytes.Buffer
	if _, err := io.Copy(&buff, r); err != nil {
		return fmt.Errorf("decode json failed: %w", err)
	}
	if err := json.Unmarshal(buff.Bytes(), dest); err != nil {
		return fmt.Errorf("decode json failed: %w", err)
	}
	return nil
}
