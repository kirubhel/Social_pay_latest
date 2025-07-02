package usecase

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/socialpay/socialpay/src/pkg/auth/core/entity"
	"github.com/socialpay/socialpay/src/pkg/erp/usecase"
	"github.com/socialpay/socialpay/src/pkg/jwt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Initiate Authentication
// Responsible for creating a unique sharable pre session token which will be used in auth processes
func (uc Usecase) InitPreSession() (entity.PreSession, error) {

	uc.log.SetPrefix("[AUTH] [USECASE] [InitPreSession] ")

	// Errors
	var ErrFailedToInitiateAuth string = "FAILED_TO_CREATE_PRE_SESSION"

	id := uuid.New()

	token := jwt.Encode(jwt.Payload{
		Exp:    time.Now().Unix() + 1800,
		Public: id,
	}, "pre_session_secret")

	preSession := entity.PreSession{
		Id:        id,
		Token:     token,
		CreatedAt: time.Now(),
	}

	uc.log.Println("created pre session")

	// Store pre session record
	err := uc.repo.StorePreSession(preSession)
	uc.log.Println("storing pre session")
	if err != nil {
		uc.log.Printf("failed to store pre session : %s\n", err)
		return preSession, Error{
			Type:    ErrFailedToInitiateAuth,
			Message: err.Error(),
		}
	}

	uc.log.Println("initiated authentication")
	// Return presession
	return preSession, nil
}

func (uc Usecase) CheckPreSession(token string) error {
	uc.log.SetPrefix("[AUTH] [USECASE] [CheckPreSession] ")
	var ErrInvalidPreSessionToken string = "INVALID_PRE_SESSION_TOKEN"
	uc.log.Println("checking presession")
	_, err := jwt.Decode(token, "pre_session_secret")
	if err != nil {
		uc.log.Printf("failed checking presession : %s\n", err.Error())
		return Error{
			Type:    ErrInvalidPreSessionToken,
			Message: err.Error(),
		}
	}

	uc.log.Println("checked presession")
	return nil
}

// Authenticate Device

// Responsible for authenticating a device
func (uc Usecase) AuthDevice(token string, ip net.IPAddr, name string, agent string) error {
	uc.log.SetPrefix("[AUTH] [USECASE] [AuthDevice] ")

	// Error
	var ErrFailedToAuthenticateDevice string = "FAILED_TO_AUTH_DEVICE"

	// Check token
	err := uc.CheckPreSession(token)
	if err != nil {
		return err
	}

	// Do device validation
	// Create device
	var device entity.Device

	id := uuid.New()
	device = entity.Device{
		Id:        id,
		IP:        ip,
		Name:      name,
		Agent:     agent,
		CreatedAt: time.Now(),
	}

	// Store device info
	err = uc.repo.StoreDevice(device)
	if err != nil {
		return Error{
			Type:    ErrFailedToAuthenticateDevice,
			Message: err.Error(),
		}
	}

	// Create device auth
	var deviceAuth entity.DeviceAuth

	id = uuid.New()
	deviceAuth = entity.DeviceAuth{
		Id:        id,
		Device:    device,
		Token:     token,
		CreatedAt: time.Now(),
	}

	// Store device auth
	err = uc.repo.StoreDeviceAuth(deviceAuth)
	if err != nil {
		return Error{
			Type:    ErrFailedToAuthenticateDevice,
			Message: err.Error(),
		}
	}

	// Update device auth
	err = uc.repo.UpdateDeviceAuthStatus(deviceAuth.Id, true)
	if err != nil {
		return Error{
			Type:    ErrFailedToAuthenticateDevice,
			Message: err.Error(),
		}
	}

	return nil
}

// Check device auth
func (uc Usecase) CheckDeviceAuth(token string) error {
	var ErrFailedToAuthenticateDevice string = "FAILED_TO_AUTH_DEVICE"

	// Check token
	err := uc.CheckPreSession(token)
	if err != nil {
		return Error{
			Type:    ErrFailedToAuthenticateDevice,
			Message: err.Error(),
		}
	}

	// Get device auth from db
	deviceAuth, err := uc.repo.FindDeviceAuth(token)
	if err != nil {
		return Error{
			Type:    ErrFailedToAuthenticateDevice,
			Message: err.Error(),
		}
	}
	// Check the status
	if !deviceAuth.Status {
		return Error{
			Type:    ErrFailedToAuthenticateDevice,
			Message: "device is not authenticated",
		}
	}
	return nil
}

// Authenticate Phone
func (uc Usecase) InitPhoneAuth(token, prefix, number string) (*entity.PhoneAuth, error) {
	// Error
	var ErrFailedToInitPhoneAuth string = "FAILED_TO_INITIATE_PHONE_AUTH"

	// Validate token
	err := uc.CheckPreSession(token)
	if err != nil {
		return nil, err
	}

	// [TODO] Validate Phone

	// Find / create phone
	phone, err := uc.repo.FindPhone(prefix, number)
	if err != nil {
		// Create phone
		return nil, Error{
			Type:    ErrFailedToInitPhoneAuth,
			Message: err.Error(),
		}
	}

	if phone == nil {
		id := uuid.New()

		phone = &entity.Phone{
			Id:        id,
			Prefix:    prefix,
			Number:    number,
			CreatedAt: time.Now(),
		}

		err := uc.repo.StorePhone(*phone)
		if err != nil {
			return nil, Error{
				Type:    ErrFailedToInitPhoneAuth,
				Message: err.Error(),
			}
		}
	}

	// Generate and send OTP
	phoneAuth := entity.PhoneAuth{
		Id:      uuid.New(),
		Token:   token,
		Phone:   *phone,
		Method:  "SMS",
		Length:  6,
		Timeout: 120,
	}

	otp := rand.Intn(999999-100000) + 100000
	otpStr := fmt.Sprint(otp)
	phoneAuth.Code = jwt.Encode(jwt.Payload{
		Exp:    time.Now().Unix() + 30*60,
		Public: otpStr,
	}, "otp_verification_secret")

	err = uc.repo.StorePhoneAuth(phoneAuth)
	if err != nil {
		return nil, Error{
			Type:    ErrFailedToInitPhoneAuth,
			Message: err.Error(),
		}
	}

	// Send code
	go func() {
		err := uc.sms.SendSMS(phone.String(), fmt.Sprintf("Your SocialPay verification code is %s. Do not share this code with anyone.", strconv.Itoa(otp)))
		if err != nil {
			uc.log.Println("[InitPhoneAuth] Failed to send SMS:", err)
		}
	}()

	return &phoneAuth, nil
}

func (uc Usecase) AuthPhone(token, prefix, number, otp string) error {
	// Error
	var ErrFailedToAuthPhone string = "FAILED_TO_AUTH_PHONE"

	// Validate token
	err := uc.CheckPreSession(token)
	if err != nil {
		uc.log.Println("Invalid pre-session token:", err)
		return err
	}

	// Get Phone Auth
	phoneAuth, err := uc.repo.FindPhoneAuth(token)
	if err != nil {
		uc.log.Println("Failed to find phone auth:", err)
		return Error{
			Type:    ErrFailedToAuthPhone,
			Message: "OTP verification failed",
		}
	}

	// Validate phone
	if phoneAuth.Phone.Prefix != prefix || phoneAuth.Phone.Number != number {
		uc.log.Println("Phone mismatch:", phoneAuth.Phone.Prefix, phoneAuth.Phone.Number, "!=", prefix, number)
		return Error{
			Type:    ErrFailedToAuthPhone,
			Message: "Phone number doesn't match verification request",
		}
	}

	// Check if already verified
	if phoneAuth.Status {
		return nil
	}

	// Decode and verify OTP
	decodedPayload, err := jwt.Decode(phoneAuth.Code, "otp_verification_secret")
	if err != nil {
		uc.log.Println("Failed to decode OTP:", err)
		return Error{
			Type:    ErrFailedToAuthPhone,
			Message: "Invalid OTP format",
		}
	}

	decodedCode, ok := decodedPayload.Public.(string)
	if !ok {
		uc.log.Println("Failed to extract OTP code from payload")
		return Error{
			Type:    ErrFailedToAuthPhone,
			Message: "Invalid OTP format",
		}
	}

	// Verify OTP matches
	if decodedCode != otp {
		uc.log.Println("OTP mismatch:", decodedCode, "!=", otp)
		return Error{
			Type:    ErrFailedToAuthPhone,
			Message: "Incorrect OTP code",
		}
	}

	// Check if OTP is expired
	if time.Now().Unix() > decodedPayload.Exp {
		uc.log.Println("OTP expired")
		return Error{
			Type:    ErrFailedToAuthPhone,
			Message: "OTP has expired",
		}
	}

	// Update Status
	err = uc.repo.UpdatePhoneAuthStatus(phoneAuth.Id, true)
	if err != nil {
		uc.log.Println("Failed to update phone auth status:", err)
		return Error{
			Type:    ErrFailedToAuthPhone,
			Message: "Failed to complete verification",
		}
	}

	return nil
}

func (uc Usecase) CheckPhoneAuth(token string) error {
	// Error
	var ErrFailedToAuthPhone = "FAILED_TO_AUTH_PHONE"

	// Validate token
	err := uc.CheckPreSession(token)
	if err != nil {
		return err
	}

	// Get Phone Auth (lightweight version without phone details)
	phoneAuth, err := uc.repo.FindPhoneAuthWithoutPhone(token)
	if err != nil {
		uc.log.Println("Failed to find phone auth:", err)
		return Error{
			Type:    ErrFailedToAuthPhone,
			Message: "Verification session not found",
		}
	}

	// Check the status
	if !phoneAuth.Status {
		return Error{
			Type:    ErrFailedToAuthPhone,
			Message: "Phone number not verified",
		}
	}

	return nil
}

func (uc Usecase) LoginFindPhone(prefix, number string) (*entity.Phone, error) {
	const ErrFailedToInitPhoneAuth = "FAILED_TO_INITIATE_PHONE_AUTH"

	phone, err := uc.repo.LoginFindPhone(prefix, number)
	if err != nil {
		return nil, usecase.Error{
			Type:    ErrFailedToInitPhoneAuth,
			Message: err.Error(),
		}
	}
	return phone, nil
}

// Password

func (uc Usecase) InitPasswordAuth(token string, password string, hint string) (*entity.PasswordAuth, error) {
	// Validate token
	err := uc.CheckPreSession(token)
	if err != nil {
		return nil, Error{
			Type:    "",
			Message: err.Error(),
		}
	}

	// Find user
	phoneAuth, err := uc.repo.FindPhoneAuth(token)
	if err != nil {
		return nil, Error{
			Type:    "ERRAUTHPASS",
			Message: err.Error(),
		}
	}

	user, err := uc.repo.FindUserUsingPhoneIdentity(phoneAuth.Phone.Id)
	if err != nil {
		return nil, Error{
			Type:    "ERRAUTHPASS",
			Message: err.Error(),
		}
	}

	uc.log.Println("user.Id")
	uc.log.Println(user.Id)

	pass, err := uc.CreatePasswordIdentity(user.Id, password, hint)
	if err != nil {
		return nil, Error{
			Type:    "ERRAUTHPASS",
			Message: err.Error(),
		}
	}

	passAuth := entity.PasswordAuth{
		Id:        uuid.New(),
		Token:     token,
		Password:  *pass,
		Status:    true,
		CreatedAt: time.Now(),
	}

	// Store pass auth
	err = uc.repo.StorePasswordAuth(passAuth)

	return &passAuth, err
}

func (uc Usecase) AuthPassword(token string, password string) error {
	// Validate token
	err := uc.CheckPreSession(token)
	if err != nil {
		return Error{
			Type:    "",
			Message: err.Error(),
		}
	}

	// Find user
	phoneAuth, err := uc.repo.FindPhoneAuth(token)
	if err != nil {
		uc.log.Println(err)
		return Error{
			Type:    "ERRAUTHPASS",
			Message: err.Error(),
		}
	}

	user, err := uc.repo.FindUserUsingPhoneIdentity(phoneAuth.Phone.Id)
	if err != nil {
		uc.log.Println(err)
		return Error{
			Type:    "ERRAUTHPASS",
			Message: err.Error(),
		}
	}

	// Get user password
	pass, err := uc.repo.FindPasswordIdentityByUser(user.Id)
	if err != nil {
		uc.log.Println(err)
		return Error{
			Type:    "ERRAUTHPASS",
			Message: err.Error(),
		}
	}

	if pass == nil {
		uc.log.Println("err")
		return Error{
			Type:    "ERRAUTHPASS",
			Message: "No password found",
		}
	}

	// Compare password using bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(pass.Password), []byte(password))
	if err != nil {
		return Error{
			Type:    "INCORRECT_PASSWORD",
			Message: "Password is incorrect",
		}
	}

	passAuth := entity.PasswordAuth{
		Id:        uuid.New(),
		Token:     token,
		Password:  *pass,
		Status:    true,
		CreatedAt: time.Now(),
	}

	// Store pass auth
	err = uc.repo.StorePasswordAuth(passAuth)

	return err
}

func (uc Usecase) CheckPasswordAuth(userId uuid.UUID, token string) error {
	// Error
	var ErrFailedToAuthPassword = "FAILED_TO_AUTH_PASSWORD"

	// Validate token
	err := uc.CheckPreSession(token)
	if err != nil {
		return err
	}

	// Get Password Auth
	passAuth, err := uc.repo.FindPasswordAuth(token)
	if err != nil {
		uc.log.Println(err)
		return Error{
			Type:    ErrFailedToAuthPassword,
			Message: err.Error(),
		}
	}

	uc.log.Println(passAuth)

	if passAuth == nil {

		uc.log.Println("pass auth nil")

		// Get user password
		pass, err := uc.repo.FindPasswordIdentityByUser(userId)
		if err != nil {
			return Error{
				Type:    "ERRAUTHPASS",
				Message: err.Error(),
			}
		}

		if pass == nil {
			return Error{
				Type:    "SET_PASSWORD",
				Message: "No password found for the requested user",
			}
		} else {
			return Error{
				Type:    "CHECK_PASSWORD",
				Message: "Password is set and must be verified before authenticating",
			}
		}
	}

	// Check the status
	if !passAuth.Status {
		return Error{
			Type:    ErrFailedToAuthPassword,
			Message: "Password unauthenticated",
		}
	}

	return nil
}

func (uc Usecase) CreateSession(token string) (*entity.Session, string, error) {

	// Error
	var (
		ErrCreatingSession string = "FAILED_TO_CREATE_SESSION"
		ErrSignUp          string = "SIGN_UP"
	)

	var session entity.Session
	var activeToken string

	uc.log.SetPrefix("[AUTH] [USECASE] [CreateSession] ")
	uc.log.Printf("Starting CreateSession with token: %s", token)

	// Validate token
	uc.log.Println("Checking pre-session token validity")
	err := uc.CheckPreSession(token)
	if err != nil {
		uc.log.Printf("Pre-session token check failed: %v", err)
		return &session, activeToken, err
	}

	// Phone auth check
	uc.log.Println("Checking phone authentication for token")
	err = uc.CheckPhoneAuth(token)
	if err != nil {
		uc.log.Printf("Phone authentication check failed: %v", err)
		return &session, activeToken, err
	}

	// Check user existence
	uc.log.Println("Finding phone authentication record")
	phoneAuth, err := uc.repo.FindPhoneAuth(token)
	if err != nil {
		uc.log.Printf("Failed to find phone auth: %v", err)
		return &session, token, Error{
			Type:    ErrCreatingSession,
			Message: err.Error(),
		}
	}

	uc.log.Printf("Finding user using phone identity: %s", phoneAuth.Phone.Id)
	user, err := uc.repo.FindUserUsingPhoneIdentity(phoneAuth.Phone.Id)
	if err != nil {
		uc.log.Printf("Failed to find user using phone identity: %v", err)
		return &session, token, Error{
			Type:    ErrSignUp,
			Message: err.Error(),
		}
	}

	if user == nil {
		uc.log.Printf("No user found for phone: %s", phoneAuth.Phone.Id)
		return &session, token, Error{
			Type:    "SIGN_UP",
			Message: "there is no associated user with the provided phone",
		}
	}

	id := uuid.New()
	uc.log.Printf("Creating new session with ID: %s", id)

	// Generate tokens
	uc.log.Println("Generating active and refresh tokens")
	active := jwt.Encode(jwt.Payload{
		Exp:    time.Now().Unix() + (3 * 24 * 60 * 60),
		Public: id,
	}, "active")

	refresh := jwt.Encode(jwt.Payload{
		Exp:    time.Now().Unix() + (30 * 24 * 60 * 60),
		Public: id,
	}, "active")

	session = entity.Session{
		Id:        id,
		User:      *user,
		Token:     refresh,
		CreatedAt: time.Now(),
	}

	uc.log.Printf("Storing session for user: %s", user.Id)
	err = uc.repo.StoreSession(session)
	if err != nil {
		uc.log.Printf("Failed to store session: %v", err)
		return &session, activeToken, Error{
			Type:    ErrCreatingSession,
			Message: err.Error(),
		}
	}

	uc.log.Printf("Session created successfully for user: %s, session ID: %s", user.Id, id)
	return &session, active, nil
}

// Check Session
func (uc Usecase) CheckSession(token string) (*entity.Session, error) {

	// Check session token
	fmt.Println("||||||||||||||||||||||||| check session")
	pld, err := jwt.Decode(token, "active")
	if err != nil {
		return nil, Error{
			Type:    "UNAUTHORIZED",
			Message: err.Error(),
		}
	}

	fmt.Println("////////// one ", pld.Public)
	fmt.Println("////////// two ", pld)

	// Find Session by Id
	userId, err := uuid.Parse(pld.Public.(string))
	if err != nil {
		return nil, Error{
			Type:    "UNAUTHORIZED",
			Message: err.Error(),
		}
	}
	session, err := uc.repo.FindSessionById(userId)
	if err != nil {
		return nil, Error{
			Type:    "UNAUTHORIZED",
			Message: err.Error(),
		}
	}

	return session, nil
}

// Get User
func (uc Usecase) GetUserById(id uuid.UUID) (*entity.User, error) {
	var user *entity.User
	fmt.Println("################################################### , usecase")

	user, err := uc.repo.FindUserById(id)

	return user, err
}

func (uc Usecase) CheckPermission(userID uuid.UUID, requiredPermission entity.Permission) (bool, error) {
	hasPermission, err := uc.repo.CheckPermission(userID, requiredPermission)
	if err != nil {
		return false, err
	}
	return hasPermission, nil
}

func (uc Usecase) CreateUser(
	Title string,
	FirstName string,
	LastName string,
	PhonePrefix string,
	PhoneNumber string,
	Password string,
	PasswordHint string,
	UserType string,
) (*entity.User, error) {

	uc.log.Printf("Attempting to create user: %s %s (%s%s)",
		FirstName, LastName, PhonePrefix, PhoneNumber)

	user, err := uc.repo.CreateUser(
		Title,
		FirstName,
		LastName,
		PhonePrefix,
		PhoneNumber,
		Password,
		PasswordHint,
		UserType,
	)

	if err != nil {
		uc.log.Printf("User creation failed: %v", err)

		if repoErr, ok := err.(*entity.Error); ok {
			return nil, repoErr
		}

		if strings.Contains(err.Error(), "value too long for type character varying(3)") {
			return nil, &entity.Error{
				Type:    entity.ErrInvalidPhoneFormat,
				Message: entity.MsgInvalidPhoneFormat,
			}
		}

		if strings.Contains(err.Error(), "null value in column") {
			return nil, &entity.Error{
				Type:    entity.ErrMissingRequiredData,
				Message: entity.MsgMissingData,
			}
		}
		if strings.Contains(err.Error(), "phone number must be between 9 digits long like +251911234567") {
			return nil, &entity.Error{
				Type:    entity.ErrInvalidPhoneNumberFormat,
				Message: entity.ErrInvalidPhoneNumber,
			}
		}
		return nil, &entity.Error{
			Type:    entity.ErrAccountCreation,
			Message: entity.MsgAccountCreation,
		}
	}

	uc.log.Printf("Successfully created user ID: %s", user.Id)
	return user, nil
}
