package lumina

type MdKeyFlag uint32
const (
    MDKF_NONE       MdKeyFlag = 0
    MDKF_STR        MdKeyFlag = 1
    MDKF_TYPE       MdKeyFlag = 2
    MDKF_INT64      MdKeyFlag = 3
    MDKF_UINT64     MdKeyFlag = 4
    MDKF_DCSTRLIST  MdKeyFlag = 5
    MDKF_DSVALLIST  MdKeyFlag = 6
    MDKF_FRAME_DESC MdKeyFlag = 7
    MDKF_NLSTRLIST  MdKeyFlag = 8
    MDKF_DOPSLIST   MdKeyFlag = 9
)

type MdKey uint32
const (
    MDK_NONE         MdKey = 0
    MDK_TYPE         MdKey = 1
    MDK_VD_ELAPSED   MdKey = 2
    MDK_FCMT         MdKey = 3
    MDK_FRPTCMT      MdKey = 4
    MDK_CMTS         MdKey = 5
    MDK_RPTCMTS      MdKey = 6
    MDK_EXTRACMTS    MdKey = 7
    MDK_USER_STKPNTS MdKey = 8
    MDK_FRAME_DESC   MdKey = 9
    MDK_OPS          MdKey = 10
    MDK_LAST         MdKey = 11
)
