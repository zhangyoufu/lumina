// This file is automatically generated. DO NOT MODIFY!

package lumina

func (*GetFuncHistoriesPacket) getType() PacketType {
	return PKT_GET_FUNC_HISTORIES
}

func (*GetFuncHistoriesPacket) getResponseType() PacketType {
	return PKT_GET_FUNC_HISTORIES_RESULT
}

func (this *GetFuncHistoriesPacket) readFrom(r Reader) (err error) {
	// Field this.PatternIds
	// Slice []PatternId
	var v1 uint32
	if v1, err = readUint32(r); err != nil {
		return
	}
	this.PatternIds = make([]PatternId, v1)
	for v2 := uint32(0); v2 < v1; v2++ {
		// Field this.PatternIds[v2]
		// Struct PatternId
		if err = this.PatternIds[v2].readFrom(r); err != nil {
			return
		}
	}
	// Field this.Unknown
	// Basic uint32
	if this.Unknown, err = readUint32(r); err != nil {
		return
	}
	return
}

func (this *GetFuncHistoriesPacket) writeTo(w Writer) (err error) {
	// Field this.PatternIds
	// Slice []PatternId
	if len(this.PatternIds) > 0x7FFFFFFF {
		err = errTooLong
		return
	}
	var v1 = uint32(len(this.PatternIds))
	if err = writeUint32(w, v1); err != nil {
		return
	}
	for v2 := uint32(0); v2 < v1; v2++ {
		// Field this.PatternIds[v2]
		// Struct PatternId
		if err = this.PatternIds[v2].writeTo(w); err != nil {
			return
		}
	}
	// Field this.Unknown
	// Basic uint32
	if err = writeUint32(w, this.Unknown); err != nil {
		return
	}
	return
}
