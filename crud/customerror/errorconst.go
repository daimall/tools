package customerror

var (
	// HTTP错误常量
	Unauthorized         = New(401, "Unauthorized")
	Forbidden            = New(403, "Forbidden")
	NotFound             = New(404, "Not Found") // NotFound 表示资源未找到
	MethodNotAllowed     = New(405, "Method Not Allowed")
	RequestTimeout       = New(408, "Request Timeout")
	Conflict             = New(409, "Conflict")
	Gone                 = New(410, "Gone")
	UnsupportedMediaType = New(415, "Unsupported Media Type")
	TooManyRequests      = New(429, "Too Many Requests")
	InternalServerError  = New(500, "Internal Server Error") // InternalServerError 表示内部服务器错误

	// 认证鉴权相关错误常量
	InvalidToken              = New(600, "Invalid Token")
	ExpiredToken              = New(601, "Expired Token")
	TokenNotFound             = New(602, "Token Not Found")                // 表示TOKEN不存在
	RefreshTokenExpire        = New(603, "Refresh Token Error or Expired") // 表示帐号App的刷新令牌过期
	OauthCodeExpire           = New(604, "OAuth Code Expired")             // 表示OAuth授权码过期
	AuthCodeInvalid           = New(605, "Invalid Auth Code")              // 表示Auth授权码不正确
	AuthStateInvalid          = New(606, "Invalid Auth State")             // 表示Auth状态不正确
	InvalidParameters         = New(607, "Invalid Parameters")
	MissingParameters         = New(608, "Missing Parameters")
	InvalidCredentials        = New(609, "Invalid Credentials")
	AccountDisabled           = New(610, "Account Disabled")
	InvalidUsernameFormat     = New(611, "Invalid Username Format")
	UsernameAlreadyExists     = New(612, "Username Already Exists")
	UsernameNotFound          = New(613, "Username Not Found")
	InvalidPasswordFormat     = New(614, "Invalid Password Format")
	PasswordMismatch          = New(615, "Password Mismatch")
	PasswordTooWeak           = New(616, "Password Too Weak")
	VerificationCodeError     = New(617, "Verification Code Error")             // 表示验证码错误
	DisabledUser              = New(618, "Disabled User")                       // 表示禁用的用户
	AnotherClientLogin        = New(619, "Another Client Already Logged In")    // 表示其他客户端已登录
	GetCodeFrequently         = New(620, "Too Many Verification Code Requests") // 表示请求验证码次数过多
	UserNameOrPasswordInvalid = New(621, "Invalid Username or Password")        // 表示用户名或密码错误
	NicknameError             = New(622, "Nickname Contains Forbidden Words")   // 表示昵称包含禁用词

	// CRUD 业务错误常量
	ActionNotFound       = New(1000, "Action Not Found")             // 表示Action不存在
	UpdateActionNotFound = New(1001, "Update Action Not Found")      // 表示Update Action不存在
	StepTypeNotFound     = New(1002, "Flow Step Not Found")          // 表示Flow Step不存在
	ServiceIdNotInt      = New(1003, "Incorrect Service ID Format")  // 表示Service ID格式不正确
	SignError            = New(1004, "Sign Error")                   // 表示签名错误
	AppIDError           = New(1005, "AppID Error")                  // 表示AppID错误
	TimestampError       = New(1006, "Timestamp Error")              // 表示时间戳错误
	LdapErr              = New(1007, "LDAP Account Login Exception") // 表示域账号登录异常
	ImageError           = New(1008, "Unsupported Image Format")     // 表示图片格式不支持
	QueryCondErr         = New(1009, "Incorrect Query Conditions")   // 表示查询条件不正确
	UploadErr            = New(1010, "Failed to Upload File")        // 表示上传文件失败
	ServiceNotFound      = New(1011, "Service Not Found")            // 表示服务没有注册
	CRUDContextNotFound  = New(1012, "CRUD Context Not Found")       // 表示gin上下文中缺少	CRUDContext
	MethodNotImplement   = New(1013, "method is not implemented")    // 表示接口未实现
	ServiceLoadFailed    = New(1014, "load service instance failed") // 获取server新实例失败
)
