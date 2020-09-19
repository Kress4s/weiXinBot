package constant

const (
	EXPEIRE_ACCOUNT_CODE   = 2000 //用户账户登录过期状态码
	EXPEIRE_WXACCOUNT_CODE = 2001 //微信登录状态过期状态码(暂且不用)
)

const (
	TOKEN_KEY = "WX_TOKEN"
	WX_ID     = "WX_ID"
)

const (
	BASE_URL = "49.234.86.244:8080" //微信对接baselink
)

// login
const (
	LOGIN_AUTO_URL   = BASE_URL + "/login/auto"
	LOGIN_AWAKE_URL  = BASE_URL + "/login/awake"
	LOGIN_CHECK_URL  = BASE_URL + "/login/check"
	LOGIN_HEART_URL  = BASE_URL + "/login/heartbeat"
	LOGIN_INIT_URL   = BASE_URL + "/login/init"
	LOGIN_LOGOUT_URL = BASE_URL + "/login/logout"
	LOGIN_QRCODE_URL = BASE_URL + "/login/qr_code"
)

// wxuser
const (
	WXUSER_PROFILE_URL = BASE_URL + "/user/profile"
)

//contact
const (
	CONTACT_ACCEPT_URL     = BASE_URL + "/contact/accept"
	CONTACT_ADD_URL        = BASE_URL + "/contact/add"
	CONTACT_BATCH_URL      = BASE_URL + "/contact/batch"
	CONTACT_LIST_URL       = BASE_URL + "/contact/list/all"
	CONTACT_GROUP_LIST_URL = BASE_URL + "/contact/list/group"
	CONTACT_SEARCH_URL     = BASE_URL + "/contact/search"
)

// group
const (
	GROUP_ACCEPT_URL      = BASE_URL + "/group/accept"
	GROUP_CREATE_URL      = BASE_URL + "/group/create"
	GROUP_DETAIL_URL      = BASE_URL + "/group/get/detail"
	GROUP_INFO_URL        = BASE_URL + "/group/get/info"
	GROUP_MEMBERS_URL     = BASE_URL + "/group/get/members"
	GROUP_ADD_MEMBERS_URL = BASE_URL + "/group/members/add"
	GROUP_DEL_MEMBERS_URL = BASE_URL + "/group/members/delete"
	// GROUP_ADD_MEMBERS_URL  = BASE_URL + "/group/members/invite"
	GROUP_QUIT_URL         = BASE_URL + "/group/quit"
	GROUP_SET_ANNOUNCE_URL = BASE_URL + "/group/set/announcement"
)

// label
const (
	LABEL_ADD_URL         = BASE_URL + "/label/add"
	LABEL_DELETE_URL      = BASE_URL + "/label/delete"
	LABEL_LIST_URL        = BASE_URL + "/label/list"
	LABEL_UPDATE_URL      = BASE_URL + "/label/update"
	LABEL_LIST_UPDATE_URL = BASE_URL + "/label/update/list"
)

// sns
const ()

// message info
const (
	TEXT_TYPE_MESSAGE = iota
	IMAGE_TYPE_MESSAGE
	VIDEO_TYPE_MESSAGE
	CARD_TYPE_MESSAGE
	EMOJI_TYPE_MESSAGE
	SMALL_PROGRAM_TYPE_MESSAGE
)

// message source
const (
	CONTACT_MESSAGE = iota
	GROUP_MESSAGE
	PUBLIC_MESSAGE
	SYSTEM_MESSAGE
)

// resource type
const (
	SOURCE_TEXT = iota
	SOURCE_IMAGE
	SOURCE_VOICE
	SOURCE_VIDEO
	SOURCE_FILE
	SOURCE_LINK
	SOURCE_APP
	SOURCE_EMOJI
)
