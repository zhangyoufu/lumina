package lumina

type Serializable interface {
	readFrom(Reader) error
	writeTo(Writer) error
}
