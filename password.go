package pk

type Encoder interface {
	Encode(password string) ([]byte, error)
	Decode(encoded []byte) (string, error)
}
