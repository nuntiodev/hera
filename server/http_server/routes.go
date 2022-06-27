package http_server

const (
	userBasePath               = "/hera/user"
	userCreate                 = "/create"
	userGet                    = "/user"
	usersGet                   = "/users"
	userList                   = "/list"
	userDelete                 = "/delete"
	userLogin                  = "/login"
	userResetPassword          = "/reset-password"
	userSearch                 = "/search"
	userSendResetPasswordEmail = "/send-reset-password-email"
	userSendResetPasswordText  = "/send-reset-password-text"
	userSendVerificationEmail  = "/send-verification-email"
	userSendVerificationText   = "/send-verification-text"
	userUpdateContact          = "/update-contact"
	userUpdateProfile          = "/update-profile"
	userValidateCredentials    = "/validate-credentials"
	userVerifyEmail            = "/verify-email"
	userVerifyPhone            = "/verify-phone"
)
const (
	configBasePath        = "/hera/config"
	configCreateNamespace = "/create-namespace"
	configDelete          = "/delete"
	configNamespace       = "/delete-namespace"
	configGet             = "/config"
	configRemovePublicKey = "/remove-public-key"
	configRegisterRsaKey  = "/register-rsa-key"
	configUpdate          = "/update"
)
const (
	tokenBasePath = "/hera/token"
)
