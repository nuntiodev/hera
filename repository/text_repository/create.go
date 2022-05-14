package text_repository

import (
	"context"
	"errors"
	"github.com/nuntiodev/block-proto/go_block"
)

func (t *defaultTextRepository) Create(ctx context.Context, text *go_block.Text) (*go_block.Text, error) {
	if text == nil {
		return nil, errors.New("missing required text")
	} else if text.Id == go_block.LanguageCode_INVALID_LANGUAGE_CODE {
		return nil, errors.New("missing required language code")
	}
	switch text.Id {
	case go_block.LanguageCode_EN:
		defaultEngText := getDefaultEngText()
		text.GeneralText = defaultEngText.GeneralText
		text.WelcomeText = defaultEngText.WelcomeText
		text.RegisterText = defaultEngText.RegisterText
		text.LoginText = defaultEngText.LoginText
		text.ProfileText = defaultEngText.ProfileText
	case go_block.LanguageCode_DK:
		defaultDkText := getDefaultDkText()
		text.GeneralText = defaultDkText.GeneralText
		text.WelcomeText = defaultDkText.WelcomeText
		text.RegisterText = defaultDkText.RegisterText
		text.LoginText = defaultDkText.LoginText
		text.ProfileText = defaultDkText.ProfileText
	}
	create := ProtoTextToText(text)
	if len(t.internalEncryptionKeys) > 0 {
		if err := t.EncryptText(actionCreate, create); err != nil {
			return nil, err
		}
		create.InternalEncryptionLevel = int32(len(t.internalEncryptionKeys))
	}
	if _, err := t.collection.InsertOne(ctx, create); err != nil {
		return nil, err
	}
	// set created fields
	text.InternalEncryptionLevel = create.InternalEncryptionLevel
	return text, nil
}

func getDefaultDkText() *go_block.Text {
	text := &go_block.Text{}
	// set default text
	text.GeneralText = &go_block.GeneralText{
		MissingPasswordTitle:   "Mangler adgangskode",
		MissingPasswordDetails: "Du skal skrive en adgangskode for at oprette dig/logge ind.",
		MissingEmailTitle:      "Mangler email",
		MissingEmailDetails:    "Du skal skrive en email for at oprette dig/logge ind.",
		CreatedBy:              "Skabt med Nuntio.",
		PasswordHint:           "Skriv din adgangskode",
		EmailHint:              "Skriv din email",
		ErrorTitle:             "Der skete en fejl",
		ErrorDescription:       "Noget gik galt. Prøv igen.",
		NoWifiTitle:            "Ingen forbindelse",
		NoWifiDescription:      "Ingen stabil mobil- eller wifi-forbindelse er tilgængelig.",
	}
	text.WelcomeText = &go_block.WelcomeText{
		WelcomeTitle:       "Velkommen",
		WelcomeDetails:     "Velkommen til denne platform",
		ContinueWithNuntio: "Forsæt med",
	}
	text.RegisterText = &go_block.RegisterText{
		RegisterButton:            "Opret",
		RegisterTitle:             "Opret",
		RegisterDetails:           "Udfyld felterne for at oprette en konto",
		PasswordDoNotMatchTitle:   "Kodeordene er ikke ens",
		PasswordDoNotMatchDetails: "De to adgangskoder er ikke de samme.",
		RepeatPasswordHint:        "Skriv dit kodeord igen",
		ContainsSpecialChar:       "Adgangskoden skal indeholde et særligt tegn",
		ContainsNumberChar:        "Adgangskoden skal indeholde et nummer",
		PasswordMustMatch:         "De to adgangskoder skal matche",
		ContainsEightChars:        "Adgangskoden skal være på mindst 8 tegn",
	}
	text.LoginText = &go_block.LoginText{
		LoginButton:    "Log ind",
		LoginTitle:     "Log ind",
		LoginDetails:   "Udfyld detaljerne nedenfor for at logge ind på din konto",
		ForgotPassword: "Glemt din adgangskode?",
	}
	text.ProfileText = &go_block.ProfileText{
		ProfileTitle:              "Profil",
		Logout:                    "Log ud",
		ChangeEmailTitle:          "Skift din nuværende e-mail",
		ChangeEmailDescription:    "Indtast en ny e-mail nedenfor for at ændre din nuværende e-mail.",
		ChangePasswordTitle:       "Skift din nuværende adgangskode",
		ChangePasswordDescription: "Indtast en ny adgangskode nedenfor for at ændre din nuværende adgangskode.",
	}
	return text
}

func getDefaultEngText() *go_block.Text {
	text := &go_block.Text{}
	// set default text
	text.GeneralText = &go_block.GeneralText{
		MissingPasswordTitle:   "Missing required password",
		MissingPasswordDetails: "You need to provide a password to create/login to an account.",
		MissingEmailTitle:      "Missing required email",
		MissingEmailDetails:    "You need to provide an email to create/login to an account.",
		CreatedBy:              "Powered by Nuntio.",
		PasswordHint:           "Enter your password",
		EmailHint:              "Enter your email",
		ErrorTitle:             "An error occurred",
		ErrorDescription:       "Something went wrong. Please try again.",
		NoWifiTitle:            "No connection",
		NoWifiDescription:      "No stable cellular or wifi connection is available.",
	}
	text.WelcomeText = &go_block.WelcomeText{
		WelcomeTitle:       "Welcome",
		WelcomeDetails:     "Welcome to this awesome platform.",
		ContinueWithNuntio: "Continue with",
	}
	text.RegisterText = &go_block.RegisterText{
		RegisterButton:            "Register",
		RegisterTitle:             "Register",
		RegisterDetails:           "Fill in the fields in order to register for an account.",
		PasswordDoNotMatchTitle:   "Passwords do not match",
		PasswordDoNotMatchDetails: "The provided and repeat passwords are not the same.",
		RepeatPasswordHint:        "Enter your password again",
		ContainsSpecialChar:       "Password must contain a special char",
		ContainsNumberChar:        "Password must contain a number",
		PasswordMustMatch:         "The two passwords must match",
		ContainsEightChars:        "Password must be at least 8 chars long",
	}
	text.LoginText = &go_block.LoginText{
		LoginButton:    "Login",
		LoginTitle:     "Login",
		LoginDetails:   "Fill in the details below to login to your account.",
		ForgotPassword: "Forgot your password?",
	}
	text.ProfileText = &go_block.ProfileText{
		ProfileTitle:              "Profile",
		Logout:                    "Log out",
		ChangeEmailTitle:          "Change your current email",
		ChangeEmailDescription:    "Enter a new email below to change your current email.",
		ChangePasswordTitle:       "Change your current password",
		ChangePasswordDescription: "Enter a new password below to change your current password.",
	}
	return text
}
