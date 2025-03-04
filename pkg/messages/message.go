package messages

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
)

type Type string

const (
	ChallengeRequestType  Type = "challenge_request"
	ChallengeResponseType Type = "challenge_response"
	ErrorType             Type = "error"
	TextType              Type = "text"
)

type Wrapper struct {
	Type Type
	Body []byte
}

func newWrapper(t Type, body interface{}) *Wrapper {
	b, _ := json.Marshal(body)
	return &Wrapper{
		Type: t,
		Body: b,
	}
}

func GetBody[T any](w *Wrapper) (*T, error) {
	var v T
	if err := json.Unmarshal(w.Body, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

type ChallengeRequest struct {
	Complexity int
	Data       []byte
}

type Text struct {
	Text string
}

func NewError(err error) *Wrapper {
	return newWrapper(ErrorType, &Text{Text: fmt.Sprintf("%v", err)})
}

func NewText(txt string) *Wrapper {
	return newWrapper(TextType, &Text{Text: txt})
}

func NewChallengeResponse(response string) *Wrapper {
	return newWrapper(ChallengeResponseType, &Text{Text: response})
}

func NewChallengeRequest(complexity int, data []byte) (*Wrapper, error) {
	return newWrapper(ChallengeRequestType, &ChallengeRequest{
		Complexity: complexity,
		Data:       data,
	}), nil
}

// Encode implements binary protocol with first 2 bytes message length and then message content in binary format
// json should be replaced with proto for the best performance
func (m *Wrapper) Encode(w io.Writer, b []byte) error {
	buf := bytes.NewBuffer(b[:0])
	en := json.NewEncoder(buf)

	err := en.Encode(m)
	if err != nil {
		return fmt.Errorf("message marshal: %v", err)
	}

	prefixBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(prefixBytes, uint16(buf.Len()))

	_, err = w.Write(prefixBytes)
	if err != nil {
		return fmt.Errorf("message write: %v", err)
	}

	_, err = w.Write(buf.Bytes())

	return err
}

// Decode implements binary protocol with first 2 bytes message length and then message content in binary format
// json should be replaced with proto for the best performance
func Decode(r io.Reader, b []byte) (*Wrapper, error) {
	lenReadBytes := 0

	for lenReadBytes < 2 {
		ln, err := r.Read(b[lenReadBytes:2])
		if err != nil {
			return nil, fmt.Errorf("decode message read: %v", err)
		}

		lenReadBytes += ln
		if lenReadBytes < 2 {
			continue
		}
	}

	msgLen := binary.BigEndian.Uint16(b[:2])
	ln, err := io.ReadAtLeast(r, b, int(msgLen))
	if err != nil {
		return nil, fmt.Errorf("decode message read at least: %v", err)
	}

	m := &Wrapper{}
	err = json.Unmarshal(b[:ln], m)
	if err != nil {
		return nil, fmt.Errorf("json message unmarshal: %v", err)
	}

	return m, nil
}
