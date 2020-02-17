package template4app

var (
	ServiceNotificationCodes = `
package notifications

import (
	"crypto/sha1" // #nosec
	"encoding/hex"
	"fmt"
	"github.com/Unknwon/com"

	m "{{.Dir}}/pkg/models"
	"{{.Dir}}/pkg/setting"
	"time"
)

const timeLimitCodeLength = 12 + 6 + 40

// create a time limit code
// code format: 12 length date time string + 6 minutes string + 40 sha1 encoded string
func createTimeLimitCode(data string, minutes int, startInf interface{}) string {
	format := "200601021504"

	var start, end time.Time
	var startStr, endStr string

	if startInf == nil {
		// Use now time create code
		start = time.Now()
		startStr = start.Format(format)
	} else {
		// use start string create code
		startStr = startInf.(string)
		start, _ = time.ParseInLocation(format, startStr, time.Local)
		startStr = start.Format(format)
	}

	end = start.Add(time.Minute * time.Duration(minutes))
	endStr = end.Format(format)

	// create sha1 encode string
	sh := sha1.New()
	sh.Write([]byte(data + setting.SecretKey + startStr + endStr + com.ToStr(minutes)))
	encoded := hex.EncodeToString(sh.Sum(nil))

	code := fmt.Sprintf("%s%06d%s", startStr, minutes, encoded)
	return code
}

// verify time limit code
func validateUserEmailCode(user *m.User, code string) bool {
	if len(code) <= 18 {
		return false
	}

	minutes := setting.EmailCodeValidMinutes
	code = code[:timeLimitCodeLength]

	// split code
	start := code[:12]
	lives := code[12:18]
	if d, err := com.StrTo(lives).Int(); err == nil {
		minutes = d
	}

	// right active code
	data := com.ToStr(user.Id) + user.Email + user.Login + user.Password + user.Rands
	retCode := createTimeLimitCode(data, minutes, start)
	fmt.Printf("code : %s\ncode2: %s", retCode, code)
	if retCode == code && minutes > 0 {
		// check time is expired or not
		before, _ := time.ParseInLocation("200601021504", start, time.Local)
		now := time.Now()
		if before.Add(time.Minute*time.Duration(minutes)).Unix() > now.Unix() {
			return true
		}
	}

	return false
}

func getLoginForEmailCode(code string) string {
	if len(code) <= timeLimitCodeLength {
		return ""
	}

	// use tail hex username query user
	hexStr := code[timeLimitCodeLength:]
	b, _ := hex.DecodeString(hexStr)
	return string(b)
}

func createUserEmailCode(u *m.User, startInf interface{}) string {
	minutes := setting.EmailCodeValidMinutes
	data := com.ToStr(u.Id) + u.Email + u.Login + u.Password + u.Rands
	code := createTimeLimitCode(data, minutes, startInf)

	// add tail hex username
	code += hex.EncodeToString([]byte(u.Login))
	return code
}

`
	ServiceNotificationCodesTest =`
package notifications

import (
	"testing"

	m "{{.Dir}}/pkg/models"
	"{{.Dir}}/pkg/setting"
	. "github.com/smartystreets/goconvey/convey"
)

func TestEmailCodes(t *testing.T) {

	Convey("When generating code", t, func() {
		setting.EmailCodeValidMinutes = 120

		user := &m.User{Id: 10, Email: "t@a.com", Login: "asd", Password: "1", Rands: "2"}
		code := createUserEmailCode(user, nil)

		Convey("getLoginForCode should return login", func() {
			login := getLoginForEmailCode(code)
			So(login, ShouldEqual, "asd")
		})

		Convey("Can verify valid code", func() {
			So(validateUserEmailCode(user, code), ShouldBeTrue)
		})

		Convey("Cannot verify in-valid code", func() {
			code = "ASD"
			So(validateUserEmailCode(user, code), ShouldBeFalse)
		})

	})

}

`
	ServiceNotificationEmail = `
package notifications

import (
	m "{{.Dir}}/pkg/models"
	"{{.Dir}}/pkg/setting"
)

type Message struct {
	To           []string
	From         string
	Subject      string
	Body         string
	Info         string
	EmbededFiles []string
}

func setDefaultTemplateData(data map[string]interface{}, u *m.User) {
	data["AppUrl"] = setting.AppUrl
	data["BuildVersion"] = setting.BuildVersion
	data["BuildStamp"] = setting.BuildStamp
	data["EmailCodeValidHours"] = setting.EmailCodeValidMinutes / 60
	data["Subject"] = map[string]interface{}{}
	if u != nil {
		data["Name"] = u.NameOrFallback()
	}
}


`
	ServiceNotificationMailer = `
// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package notifications

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"net"
	"strconv"

	"{{.Dir}}/pkg/models"
	"{{.Dir}}/pkg/setting"
	"{{.Dir}}/pkg/util/errutil"
	gomail "gopkg.in/mail.v2"
)

func (ns *NotificationService) send(msg *Message) (int, error) {
	dialer, err := ns.createDialer()
	if err != nil {
		return 0, err
	}

	var num int
	for _, address := range msg.To {
		m := gomail.NewMessage()
		m.SetHeader("From", msg.From)
		m.SetHeader("To", address)
		m.SetHeader("Subject", msg.Subject)
		for _, file := range msg.EmbededFiles {
			m.Embed(file)
		}

		m.SetBody("text/html", msg.Body)

		e := dialer.DialAndSend(m)
		if e != nil {
			err = errutil.Wrapf(e, "Failed to send notification to email address: %s", address)
			continue
		}

		num++
	}

	return num, err
}

func (ns *NotificationService) createDialer() (*gomail.Dialer, error) {
	host, port, err := net.SplitHostPort(ns.Cfg.Smtp.Host)

	if err != nil {
		return nil, err
	}
	iPort, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}

	tlsconfig := &tls.Config{
		InsecureSkipVerify: ns.Cfg.Smtp.SkipVerify,
		ServerName:         host,
	}

	if ns.Cfg.Smtp.CertFile != "" {
		cert, err := tls.LoadX509KeyPair(ns.Cfg.Smtp.CertFile, ns.Cfg.Smtp.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("Could not load cert or key file. error: %v", err)
		}
		tlsconfig.Certificates = []tls.Certificate{cert}
	}

	d := gomail.NewDialer(host, iPort, ns.Cfg.Smtp.User, ns.Cfg.Smtp.Password)
	d.TLSConfig = tlsconfig

	if ns.Cfg.Smtp.EhloIdentity != "" {
		d.LocalName = ns.Cfg.Smtp.EhloIdentity
	} else {
		d.LocalName = setting.InstanceName
	}
	return d, nil
}

func (ns *NotificationService) buildEmailMessage(cmd *models.SendEmailCommand) (*Message, error) {
	if !ns.Cfg.Smtp.Enabled {
		return nil, models.ErrSmtpNotEnabled
	}

	var buffer bytes.Buffer
	var err error

	data := cmd.Data
	if data == nil {
		data = make(map[string]interface{}, 10)
	}

	setDefaultTemplateData(data, nil)
	err = mailTemplates.ExecuteTemplate(&buffer, cmd.Template, data)
	if err != nil {
		return nil, err
	}

	subject := cmd.Subject
	if cmd.Subject == "" {
		var subjectText interface{}
		subjectData := data["Subject"].(map[string]interface{})
		subjectText, hasSubject := subjectData["value"]

		if !hasSubject {
			return nil, fmt.Errorf("Missing subject in Template %s", cmd.Template)
		}

		subjectTmpl, err := template.New("subject").Parse(subjectText.(string))
		if err != nil {
			return nil, err
		}

		var subjectBuffer bytes.Buffer
		err = subjectTmpl.ExecuteTemplate(&subjectBuffer, "subject", data)
		if err != nil {
			return nil, err
		}

		subject = subjectBuffer.String()
	}

	return &Message{
		To:           cmd.To,
		From:         fmt.Sprintf("%s <%s>", ns.Cfg.Smtp.FromName, ns.Cfg.Smtp.FromAddress),
		Subject:      subject,
		Body:         buffer.String(),
		EmbededFiles: cmd.EmbededFiles,
	}, nil
}

`
	ServiceNotifications = `
package notifications

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/url"
	"path/filepath"
	"strings"

	"{{.Dir}}/pkg/bus"
	"{{.Dir}}/pkg/events"
	"{{.Dir}}/pkg/infra/log"
	m "{{.Dir}}/pkg/models"
	"{{.Dir}}/pkg/registry"
	"{{.Dir}}/pkg/setting"
	"{{.Dir}}/pkg/util"
)

var mailTemplates *template.Template
var tmplResetPassword = "reset_password.html"
var tmplSignUpStarted = "signup_started.html"
var tmplWelcomeOnSignUp = "welcome_on_signup.html"

func init() {
	registry.RegisterService(&NotificationService{})
	fmt.Println("Initialized notification ....")
}

type NotificationService struct {
	Bus bus.Bus      ` +"`inject:\"\"`"+`
	Cfg *setting.Cfg ` + "`inject:\"\"`"+`

	mailQueue    chan *Message
	webhookQueue chan *Webhook
	log          log.Logger
}

func (ns *NotificationService) Init() error {
	ns.log = log.New("notifications")
	ns.mailQueue = make(chan *Message, 10)
	ns.webhookQueue = make(chan *Webhook, 10)

	ns.Bus.AddHandler(ns.sendResetPasswordEmail)
	ns.Bus.AddHandler(ns.validateResetPasswordCode)
	ns.Bus.AddHandler(ns.sendEmailCommandHandler)

	ns.Bus.AddHandlerCtx(ns.sendEmailCommandHandlerSync)
	ns.Bus.AddHandlerCtx(ns.SendWebhookSync)

	ns.Bus.AddEventListener(ns.signUpStartedHandler)
	ns.Bus.AddEventListener(ns.signUpCompletedHandler)

	mailTemplates = template.New("name")
	mailTemplates.Funcs(template.FuncMap{
		"Subject": subjectTemplateFunc,
	})

	templatePattern := filepath.Join(setting.StaticRootPath, ns.Cfg.Smtp.TemplatesPattern)
	_, err := mailTemplates.ParseGlob(templatePattern)
	if err != nil {
		return err
	}

	if !util.IsEmail(ns.Cfg.Smtp.FromAddress) {
		return errors.New("Invalid email address for SMTP from_address config")
	}

	if setting.EmailCodeValidMinutes == 0 {
		setting.EmailCodeValidMinutes = 120
	}

	return nil
}

func (ns *NotificationService) Run(ctx context.Context) error {
	ns.log.Info("NotificationService does startup ....")
	for {
		select {
		case webhook := <-ns.webhookQueue:
			err := ns.sendWebRequestSync(context.Background(), webhook)

			if err != nil {
				ns.log.Error("Failed to send webrequest ", "error", err)
			}
		case msg := <-ns.mailQueue:
			num, err := ns.send(msg)
			tos := strings.Join(msg.To, "; ")
			info := ""
			if err != nil {
				if len(msg.Info) > 0 {
					info = ", info: " + msg.Info
				}
				ns.log.Error(fmt.Sprintf("Async sent email %d succeed, not send emails: %s%s err: %s", num, tos, info, err))
			} else {
				ns.log.Debug(fmt.Sprintf("Async sent email %d succeed, sent emails: %s%s", num, tos, info))
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (ns *NotificationService) SendWebhookSync(ctx context.Context, cmd *m.SendWebhookSync) error {
	return ns.sendWebRequestSync(ctx, &Webhook{
		Url:         cmd.Url,
		User:        cmd.User,
		Password:    cmd.Password,
		Body:        cmd.Body,
		HttpMethod:  cmd.HttpMethod,
		HttpHeader:  cmd.HttpHeader,
		ContentType: cmd.ContentType,
	})
}

func subjectTemplateFunc(obj map[string]interface{}, value string) string {
	obj["value"] = value
	return ""
}

func (ns *NotificationService) sendEmailCommandHandlerSync(ctx context.Context, cmd *m.SendEmailCommandSync) error {
	message, err := ns.buildEmailMessage(&m.SendEmailCommand{
		Data:         cmd.Data,
		Info:         cmd.Info,
		Template:     cmd.Template,
		To:           cmd.To,
		EmbededFiles: cmd.EmbededFiles,
		Subject:      cmd.Subject,
	})

	if err != nil {
		return err
	}

	_, err = ns.send(message)
	return err
}

func (ns *NotificationService) sendEmailCommandHandler(cmd *m.SendEmailCommand) error {
	message, err := ns.buildEmailMessage(cmd)

	if err != nil {
		return err
	}

	ns.mailQueue <- message
	return nil
}

func (ns *NotificationService) sendResetPasswordEmail(cmd *m.SendResetPasswordEmailCommand) error {
	return ns.sendEmailCommandHandler(&m.SendEmailCommand{
		To:       []string{cmd.User.Email},
		Template: tmplResetPassword,
		Data: map[string]interface{}{
			"Code": createUserEmailCode(cmd.User, nil),
			"Name": cmd.User.NameOrFallback(),
		},
	})
}

func (ns *NotificationService) validateResetPasswordCode(query *m.ValidateResetPasswordCodeQuery) error {
	login := getLoginForEmailCode(query.Code)
	if login == "" {
		return m.ErrInvalidEmailCode
	}

	userQuery := m.GetUserByLoginQuery{LoginOrEmail: login}
	if err := bus.Dispatch(&userQuery); err != nil {
		return err
	}

	if !validateUserEmailCode(userQuery.Result, query.Code) {
		return m.ErrInvalidEmailCode
	}

	query.Result = userQuery.Result
	return nil
}

func (ns *NotificationService) signUpStartedHandler(evt *events.SignUpStarted) error {
	if !setting.VerifyEmailEnabled {
		return nil
	}

	ns.log.Info("User signup started", "email", evt.Email)

	if evt.Email == "" {
		return nil
	}

	err := ns.sendEmailCommandHandler(&m.SendEmailCommand{
		To:       []string{evt.Email},
		Template: tmplSignUpStarted,
		Data: map[string]interface{}{
			"Email":     evt.Email,
			"Code":      evt.Code,
			"SignUpUrl": setting.ToAbsUrl(fmt.Sprintf("signup/?email=%s&code=%s", url.QueryEscape(evt.Email), url.QueryEscape(evt.Code))),
		},
	})

	if err != nil {
		return err
	}

	emailSentCmd := m.UpdateTempUserWithEmailSentCommand{Code: evt.Code}
	return bus.Dispatch(&emailSentCmd)
}

func (ns *NotificationService) signUpCompletedHandler(evt *events.SignUpCompleted) error {
	if evt.Email == "" || !ns.Cfg.Smtp.SendWelcomeEmailOnSignUp {
		return nil
	}

	return ns.sendEmailCommandHandler(&m.SendEmailCommand{
		To:       []string{evt.Email},
		Template: tmplWelcomeOnSignUp,
		Data: map[string]interface{}{
			"Name": evt.Name,
		},
	})
}

`
	ServiceNotificationTest = `
package notifications

import (
	"testing"

	"{{.Dir}}/pkg/bus"
	m "{{.Dir}}/pkg/models"
	"{{.Dir}}/pkg/setting"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNotifications(t *testing.T) {

	Convey("Given the notifications service", t, func() {
		setting.StaticRootPath = "../../../public/"

		ns := &NotificationService{}
		ns.Bus = bus.New()
		ns.Cfg = setting.NewCfg()
		ns.Cfg.Smtp.Enabled = true
		ns.Cfg.Smtp.TemplatesPattern = "emails/*.html"
		ns.Cfg.Smtp.FromAddress = "from@address.com"
		ns.Cfg.Smtp.FromName = "Grafana Admin"

		err := ns.Init()
		So(err, ShouldBeNil)

		Convey("When sending reset email password", func() {
			err := ns.sendResetPasswordEmail(&m.SendResetPasswordEmailCommand{User: &m.User{Email: "asd@asd.com"}})
			So(err, ShouldBeNil)

			sentMsg := <-ns.mailQueue
			So(sentMsg.Body, ShouldContainSubstring, "body")
			So(sentMsg.Subject, ShouldEqual, "Reset your Grafana password - asd@asd.com")
			So(sentMsg.Body, ShouldNotContainSubstring, "Subject")
		})
	})
}

`
	ServiceNotificationSendEmailIntegrationTest = `
package notifications

import (
	"io/ioutil"
	"testing"

	"{{.Dir}}/pkg/bus"
	m "{{.Dir}}/pkg/models"
	"{{.Dir}}/pkg/setting"
	. "github.com/smartystreets/goconvey/convey"
)

func TestEmailIntegrationTest(t *testing.T) {
	SkipConvey("Given the notifications service", t, func() {
		setting.StaticRootPath = "../../../public/"
		setting.BuildVersion = "4.0.0"

		ns := &NotificationService{}
		ns.Bus = bus.New()
		ns.Cfg = setting.NewCfg()
		ns.Cfg.Smtp.Enabled = true
		ns.Cfg.Smtp.TemplatesPattern = "emails/*.html"
		ns.Cfg.Smtp.FromAddress = "from@address.com"
		ns.Cfg.Smtp.FromName = "Grafana Admin"

		err := ns.Init()
		So(err, ShouldBeNil)

		Convey("When sending reset email password", func() {
			cmd := &m.SendEmailCommand{

				Data: map[string]interface{}{
					"Title":         "[CRITICAL] Imaginary timeserie alert",
					"State":         "Firing",
					"Name":          "Imaginary timeserie alert",
					"Severity":      "ok",
					"SeverityColor": "#D63232",
					"Message":       "Alert message that will support markdown in some distant future.",
					"RuleUrl":       "http://localhost:3000/dashboard/db/graphite-dashboard",
					"ImageLink":     "http://localhost:3000/render/dashboard-solo/db/graphite-dashboard?panelId=1&from=1471008499616&to=1471012099617&width=1000&height=500",
					"AlertPageUrl":  "http://localhost:3000/alerting",
					"EmbeddedImage": "test.png",
					"EvalMatches": []map[string]string{
						{
							"Metric": "desktop",
							"Value":  "40",
						},
						{
							"Metric": "mobile",
							"Value":  "20",
						},
					},
				},
				To:       []string{"asdf@asdf.com"},
				Template: "alert_notification.html",
			}

			err := ns.sendEmailCommandHandler(cmd)
			So(err, ShouldBeNil)

			sentMsg := <-ns.mailQueue
			So(sentMsg.From, ShouldEqual, "Grafana Admin <from@address.com>")
			So(sentMsg.To[0], ShouldEqual, "asdf@asdf.com")
			ioutil.WriteFile("../../../tmp/test_email.html", []byte(sentMsg.Body), 0777)
		})
	})
}

`
	ServiceNotificationWebhook = `
package notifications

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/context/ctxhttp"

	"{{.Dir}}/pkg/util"
)

type Webhook struct {
	Url         string
	User        string
	Password    string
	Body        string
	HttpMethod  string
	HttpHeader  map[string]string
	ContentType string
}

var netTransport = &http.Transport{
	TLSClientConfig: &tls.Config{
		Renegotiation: tls.RenegotiateFreelyAsClient,
	},
	Proxy: http.ProxyFromEnvironment,
	Dial: (&net.Dialer{
		Timeout: 30 * time.Second,
	}).Dial,
	TLSHandshakeTimeout: 5 * time.Second,
}
var netClient = &http.Client{
	Timeout:   time.Second * 30,
	Transport: netTransport,
}

func (ns *NotificationService) sendWebRequestSync(ctx context.Context, webhook *Webhook) error {
	ns.log.Debug("Sending webhook", "url", webhook.Url, "http method", webhook.HttpMethod)

	if webhook.HttpMethod == "" {
		webhook.HttpMethod = http.MethodPost
	}

	request, err := http.NewRequest(webhook.HttpMethod, webhook.Url, bytes.NewReader([]byte(webhook.Body)))
	if err != nil {
		return err
	}

	if webhook.ContentType == "" {
		webhook.ContentType = "application/json"
	}

	request.Header.Add("Content-Type", webhook.ContentType)
	request.Header.Add("User-Agent", "Grafana")

	if webhook.User != "" && webhook.Password != "" {
		request.Header.Add("Authorization", util.GetBasicAuthHeader(webhook.User, webhook.Password))
	}

	for k, v := range webhook.HttpHeader {
		request.Header.Set(k, v)
	}

	resp, err := ctxhttp.Do(ctx, netClient, request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode/100 == 2 {
		// flushing the body enables the transport to reuse the same connection
		io.Copy(ioutil.Discard, resp.Body)
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	ns.log.Debug("Webhook failed", "statuscode", resp.Status, "body", string(body))
	return fmt.Errorf("Webhook response status %v", resp.Status)
}

`
)
