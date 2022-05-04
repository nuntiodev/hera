package email

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestSendEmail(t *testing.T) {
	os.Setenv("EMAIL_FROM", "oscar@softcorp.io")
	os.Setenv("POSTMARK_SERVER_TOKEN", "a474a6fb-94d5-4a45-ae91-eb0ad8db8bdb")
	os.Setenv("POSTMARK_ACCOUNT_TOKEN", "ab25188e-c661-4587-a376-8a199b817c75")
	os.Setenv("EMAIL_TEMPLATE_PATHS", "./../email_templates/verify_email.html")
	email, err := New(nil)
	assert.NoError(t, err)
	assert.NotNil(t, email)
	assert.NoError(t, email.SendVerificationEmail("oscar@softcorp.io", "Test nuntio mail", "./../email_templates/verify_email.html", &VerificationData{
		TemplateData: TemplateData{
			LogoUrl:        "https://raw.githubusercontent.com/softcorp-io/website/main/nuntio/nuntio_text_white.png",
			Title:          "Verify your email",
			WelcomeMessage: "Hello",
			NameOfUser:     "oscar@nuntio.io",
			BodyMessage:    "We are happy that you have signed up for Nuntio. In order to get started, we ask of you to confirm your email by entering the following numbers in your Nuntio app.",
			FooterMessage:  "Best, the Nuntio team.",
		},
		Code: "6 7 8 9 2",
	}))
}
