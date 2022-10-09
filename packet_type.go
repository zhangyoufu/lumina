package lumina

import "fmt"

type PacketType uint8

const (
	PKT_RPC_OK                   PacketType = 0x0A
	PKT_RPC_FAIL                 PacketType = 0x0B
	PKT_RPC_NOTIFY               PacketType = 0x0C
	PKT_HELO                     PacketType = 0x0D
	PKT_PULL_MD                  PacketType = 0x0E
	PKT_PULL_MD_RESULT           PacketType = 0x0F
	PKT_PUSH_MD                  PacketType = 0x10
	PKT_PUSH_MD_RESULT           PacketType = 0x11
	PKT_GET_POP                  PacketType = 0x12
	PKT_GET_POP_RESULT           PacketType = 0x13
	PKT_LIST_PEERS               PacketType = 0x14
	PKT_LIST_PEERS_RESULT        PacketType = 0x15
	PKT_KILL_SESSIONS            PacketType = 0x16
	PKT_KILL_SESSIONS_RESULT     PacketType = 0x17
	PKT_DEL_HISTORY              PacketType = 0x18
	PKT_DEL_HISTORY_RESULT       PacketType = 0x19
	PKT_SHOW_HISTORY             PacketType = 0x1A
	PKT_SHOW_HISTORY_RESULT      PacketType = 0x1B
	PKT_DUMP_MD                  PacketType = 0x1C
	PKT_DUMP_MD_RESULT           PacketType = 0x1D
	PKT_CLEAN_DB                 PacketType = 0x1E
	PKT_DEBUGCTL                 PacketType = 0x1F
	PKT_DECOMPILE                PacketType = 0x20
	PKT_DECOMPILE_RESULT         PacketType = 0x21
	PKT_PUSH_TLM                 PacketType = 0x22
	PKT_SHOW_TLM_SESSIONS        PacketType = 0x23
	PKT_SHOW_TLM_SESSIONS_RESULT PacketType = 0x24
	PKT_SHOW_PUSHES              PacketType = 0x25
	PKT_SHOW_PUSHES_RESULT       PacketType = 0x26
	PKT_USER_OPERATION           PacketType = 0x27
	PKT_SHOW_USERS               PacketType = 0x28
	PKT_SHOW_USERS_RESULT        PacketType = 0x29
	PKT_SET_PASSWORD             PacketType = 0x2A
	PKT_GET_LUMINA_INFO          PacketType = 0x2B
	PKT_GET_LUMINA_INFO_RESULT   PacketType = 0x2C
	PKT_GET_LUMINA_STATS         PacketType = 0x2D
	PKT_GET_LUMINA_STATS_RESULT  PacketType = 0x2E
)

func (t PacketType) String() string {
	switch t {
	case PKT_RPC_OK:
		return "PKT_RPC_OK"
	case PKT_RPC_FAIL:
		return "PKT_RPC_FAIL"
	case PKT_RPC_NOTIFY:
		return "PKT_RPC_NOTIFY"
	case PKT_HELO:
		return "PKT_HELO"
	case PKT_PULL_MD:
		return "PKT_PULL_MD"
	case PKT_PULL_MD_RESULT:
		return "PKT_PULL_MD_RESULT"
	case PKT_PUSH_MD:
		return "PKT_PUSH_MD"
	case PKT_PUSH_MD_RESULT:
		return "PKT_PUSH_MD_RESULT"
	case PKT_GET_POP:
		return "PKT_GET_POP"
	case PKT_GET_POP_RESULT:
		return "PKT_GET_POP_RESULT"
	case PKT_LIST_PEERS:
		return "PKT_LIST_PEERS"
	case PKT_LIST_PEERS_RESULT:
		return "PKT_LIST_PEERS_RESULT"
	case PKT_KILL_SESSIONS:
		return "PKT_KILL_SESSIONS"
	case PKT_KILL_SESSIONS_RESULT:
		return "PKT_KILL_SESSIONS_RESULT"
	case PKT_DEL_HISTORY:
		return "PKT_DEL_HISTORY"
	case PKT_DEL_HISTORY_RESULT:
		return "PKT_DEL_HISTORY_RESULT"
	case PKT_SHOW_HISTORY:
		return "PKT_SHOW_HISTORY"
	case PKT_SHOW_HISTORY_RESULT:
		return "PKT_SHOW_HISTORY_RESULT"
	case PKT_DUMP_MD:
		return "PKT_DUMP_MD"
	case PKT_DUMP_MD_RESULT:
		return "PKT_DUMP_MD_RESULT"
	case PKT_CLEAN_DB:
		return "PKT_CLEAN_DB"
	case PKT_DEBUGCTL:
		return "PKT_DEBUGCTL"
	case PKT_DECOMPILE:
		return "PKT_DECOMPILE"
	case PKT_DECOMPILE_RESULT:
		return "PKT_DECOMPILE_RESULT"
	case PKT_PUSH_TLM:
		return "PKT_PUSH_TLM"
	case PKT_SHOW_TLM_SESSIONS:
		return "PKT_SHOW_TLM_SESSIONS"
	case PKT_SHOW_TLM_SESSIONS_RESULT:
		return "PKT_SHOW_TLM_SESSIONS_RESULT"
	case PKT_SHOW_PUSHES:
		return "PKT_SHOW_PUSHES"
	case PKT_SHOW_PUSHES_RESULT:
		return "PKT_SHOW_PUSHES_RESULT"
	case PKT_USER_OPERATION:
		return "PKT_USER_OPERATION"
	case PKT_SHOW_USERS:
		return "PKT_SHOW_USERS"
	case PKT_SHOW_USERS_RESULT:
		return "PKT_SHOW_USERS_RESULT"
	case PKT_SET_PASSWORD:
		return "PKT_SET_PASSWORD"
	case PKT_GET_LUMINA_INFO:
		return "PKT_GET_LUMINA_INFO"
	case PKT_GET_LUMINA_INFO_RESULT:
		return "PKT_GET_LUMINA_INFO_RESULT"
	case PKT_GET_LUMINA_STATS:
		return "PKT_GET_LUMINA_STATS"
	case PKT_GET_LUMINA_STATS_RESULT:
		return "PKT_GET_LUMINA_STATS_RESULT"
	default:
		return fmt.Sprintf("0x%02X", uint8(t))
	}
}
