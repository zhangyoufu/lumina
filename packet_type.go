package lumina

import "fmt"

type PacketType uint8

// Most values are documented at https://www.hex-rays.com/products/ida/support/idapython_docs/ida_lumina-module.html.
// PKT_TELEMETRY_[123] are undocumented logging facility, for debugging or
// telemetry purpose.
const (
	PKT_RPC_OK               PacketType = 0x0A
	PKT_RPC_FAIL             PacketType = 0x0B
	PKT_RPC_NOTIFY           PacketType = 0x0C // unused
	PKT_HELO                 PacketType = 0x0D
	PKT_PULL_MD              PacketType = 0x0E
	PKT_PULL_MD_RESULT       PacketType = 0x0F
	PKT_PUSH_MD              PacketType = 0x10
	PKT_PUSH_MD_RESULT       PacketType = 0x11
	PKT_GET_POP              PacketType = 0x12 // unused
	PKT_GET_POP_RESULT       PacketType = 0x13 // unused
	PKT_LIST_PEERS           PacketType = 0x14 // unused
	PKT_LIST_PEERS_RESULT    PacketType = 0x15 // unused
	PKT_KILL_SESSIONS        PacketType = 0x16 // unused
	PKT_KILL_SESSIONS_RESULT PacketType = 0x17 // unused
	PKT_DUMP_MD              PacketType = 0x18 // unused
	PKT_DUMP_MD_RESULT       PacketType = 0x19 // unused
	PKT_CLEAN_DB             PacketType = 0x1A // unused
	PKT_DEBUG_SLEEP          PacketType = 0x1B
	PKT_TELEMETRY_1          PacketType = 0x34 // level < 0, exit(1)
	PKT_TELEMETRY_2          PacketType = 0x35 // level = 0
	PKT_TELEMETRY_3          PacketType = 0x36 // level > 0
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
	case PKT_DUMP_MD:
		return "PKT_DUMP_MD"
	case PKT_DUMP_MD_RESULT:
		return "PKT_DUMP_MD_RESULT"
	case PKT_CLEAN_DB:
		return "PKT_CLEAN_DB"
	case PKT_DEBUG_SLEEP:
		return "PKT_DEBUG_SLEEP"
	case PKT_TELEMETRY_1:
		return "PKT_TELEMETRY_1"
	case PKT_TELEMETRY_2:
		return "PKT_TELEMETRY_2"
	case PKT_TELEMETRY_3:
		return "PKT_TELEMETRY_3"
	default:
		return fmt.Sprintf("0x%02X", t)
	}
}
