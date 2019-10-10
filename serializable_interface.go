package lumina

type Serializable interface {
	ReadFrom(Reader) error
	WriteTo(Writer) error
}
