package lumina

type OpResult int32

const (
	PDRES_BADPTN    OpResult = -3
	PDRES_NOT_FOUND OpResult = -2
	PDRES_ERROR     OpResult = -1
	PDRES_OK        OpResult = 0
	PDRES_ADDED     OpResult = 1
)
