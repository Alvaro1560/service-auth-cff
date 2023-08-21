package validation_email

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"service-auth-cff/internal/sendmail"

	"math/rand"
	"net/http"
	"service-auth-cff/internal/env"
	"service-auth-cff/internal/logger"
	"service-auth-cff/internal/msgs"
	"service-auth-cff/internal/password"
	"service-auth-cff/internal/response"
	"service-auth-cff/pkg/auth"
	"strconv"
	"time"
)

type handlerValidationEmail struct {
	DB   *sqlx.DB
	TxID string
}

func (h *handlerValidationEmail) sendCode(c *fiber.Ctx) error {
	var parameters = make(map[string]string, 0)
	e := env.NewConfiguration()
	var msg msgs.Model
	res := response.Model{Error: true}
	m := VerificationRequest{}

	err := c.BodyParser(&m)
	if err != nil {
		logger.Error.Printf(h.TxID, "couldn't bind model validate email: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	srvUser := auth.NewServerAuth(h.DB, nil, h.TxID)
	min := 1000
	max := 9999
	rand.Seed(time.Now().UnixNano())
	emailCode := strconv.Itoa(rand.Intn(max-min+1) + min)
	verifiedCode := password.Encrypt(emailCode)

	codVerify, code, err := srvUser.SrvVerificationEmail.CreateVerificationEmail(m.Email, verifiedCode, "", nil)
	if err != nil {
		logger.Error.Printf(h.TxID, "couldn't create verify code: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	parameters["@access-code"] = emailCode
	parameters["@TEMPLATE-PATH"] = e.Template.EmailCode
	tos := []string{m.Email}

	email := sendmail.Model{
		From:        e.Template.EmailSender,
		To:          tos,
		CC:          nil,
		Subject:     e.Template.EmailCodeSubject,
		Body:        fmt.Sprintf("<h1>%s</h1>", emailCode),
		Attach:      "",
		Attachments: nil,
	}
	tpl, err := email.GenerateTemplateMail(parameters)
	if err != nil {
		logger.Error.Println(h.TxID, "error when parse template: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(86)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	email.Body = tpl
	err = email.SendMail()
	if err != nil {
		logger.Error.Println(h.TxID, "error when execute send email: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(86)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = codVerify.ID
	res.Code, res.Type, res.Msg = msg.GetByCode(29)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}

func (h *handlerValidationEmail) verifyCode(c *fiber.Ctx) error {
	res := response.Model{Error: true}
	var msg msgs.Model

	m := VerificationDataRequest{}

	err := c.BodyParser(&m)
	if err != nil {
		logger.Error.Printf(h.TxID, "couldn't bind model verification: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	srvUser := auth.NewServerAuth(h.DB, nil, h.TxID)

	dataVerify, code, err := srvUser.SrvVerificationEmail.GetVerificationEmailByID(m.Id)
	if err != nil {
		logger.Error.Printf(h.TxID, "couldn't get email verification: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if dataVerify.ID == 0 {
		logger.Error.Printf(h.TxID, "couldn't get email verification: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	if !password.Compare(dataVerify.Email, dataVerify.VerificationCode, m.Code) {
		logger.Error.Printf(h.TxID, "the verification code is not correct: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(10)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	if dataVerify.VerificationDate != nil {
		res.Code, res.Type, res.Msg = msg.GetByCode(5)
		return c.Status(http.StatusAccepted).JSON(res)
	}
	dateTime := time.Now()
	_, code, err = srvUser.SrvVerificationEmail.UpdateVerificationEmail(dataVerify.ID, dataVerify.Email, "", dataVerify.Identification, &dateTime)
	if err != nil {
		logger.Error.Printf(h.TxID, "couldn't get email verification: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(code)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = "successful email validation"
	res.Code, res.Type, res.Msg = msg.GetByCode(29)
	res.Error = false
	return c.Status(http.StatusOK).JSON(res)
}
