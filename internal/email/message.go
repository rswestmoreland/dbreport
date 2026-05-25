package email

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime"
	"net/mail"
	"path/filepath"
	"strings"
	"time"

	"github.com/rswestmoreland/dbreport/internal/config"
)

const crlf = "\r\n"

func BuildMessage(cfg config.EmailConfig, html []byte, reportPath string, generatedAt time.Time) ([]byte, error) {
	if len(html) == 0 {
		return nil, fmt.Errorf("email html body is empty")
	}
	if err := validateMessageConfig(cfg); err != nil {
		return nil, err
	}

	plain := plainTextBody(reportPath, generatedAt)
	attachmentName := filepath.Base(reportPath)
	if attachmentName == "." || attachmentName == string(filepath.Separator) || strings.TrimSpace(attachmentName) == "" {
		attachmentName = "report.html"
	}

	var body bytes.Buffer
	writeHeaders(&body, cfg, generatedAt)

	switch {
	case cfg.SendHTMLBody && cfg.AttachHTML:
		writeMixedMessage(&body, plain, html, attachmentName)
	case cfg.SendHTMLBody:
		writeAlternativeMessage(&body, plain, html)
	case cfg.AttachHTML:
		writeAttachmentOnlyMessage(&body, plain, html, attachmentName)
	default:
		return nil, fmt.Errorf("email.send_html_body or email.attach_html must be true")
	}

	return body.Bytes(), nil
}

func validateMessageConfig(cfg config.EmailConfig) error {
	if strings.TrimSpace(cfg.From) == "" {
		return fmt.Errorf("email.from is required")
	}
	if _, err := mail.ParseAddress(cfg.From); err != nil {
		return fmt.Errorf("email.from is invalid: %w", err)
	}
	if len(cfg.To) == 0 {
		return fmt.Errorf("email.to must contain at least one recipient")
	}
	for _, recipient := range cfg.To {
		if strings.TrimSpace(recipient) == "" {
			return fmt.Errorf("email.to contains an empty recipient")
		}
		if _, err := mail.ParseAddress(recipient); err != nil {
			return fmt.Errorf("email.to recipient %q is invalid: %w", recipient, err)
		}
	}
	if strings.TrimSpace(cfg.Subject) == "" {
		return fmt.Errorf("email.subject is required")
	}
	if containsHeaderBreak(cfg.Subject) {
		return fmt.Errorf("email.subject must not contain newlines")
	}
	if containsHeaderBreak(cfg.From) {
		return fmt.Errorf("email.from must not contain newlines")
	}
	for _, recipient := range cfg.To {
		if containsHeaderBreak(recipient) {
			return fmt.Errorf("email.to recipient must not contain newlines")
		}
	}
	return nil
}

func writeHeaders(body *bytes.Buffer, cfg config.EmailConfig, generatedAt time.Time) {
	body.WriteString("From: ")
	body.WriteString(cfg.From)
	body.WriteString(crlf)
	body.WriteString("To: ")
	body.WriteString(strings.Join(cfg.To, ", "))
	body.WriteString(crlf)
	body.WriteString("Subject: ")
	body.WriteString(mime.QEncoding.Encode("utf-8", cfg.Subject))
	body.WriteString(crlf)
	body.WriteString("Date: ")
	body.WriteString(generatedAt.Format(time.RFC1123Z))
	body.WriteString(crlf)
	body.WriteString("MIME-Version: 1.0")
	body.WriteString(crlf)
}

func writeAlternativeMessage(body *bytes.Buffer, plain string, html []byte) {
	boundary := "dbreport-alt-boundary"
	body.WriteString("Content-Type: multipart/alternative; boundary=\"")
	body.WriteString(boundary)
	body.WriteString("\"")
	body.WriteString(crlf + crlf)
	writePlainPart(body, boundary, plain)
	writeHTMLPart(body, boundary, html)
	closeBoundary(body, boundary)
}

func writeMixedMessage(body *bytes.Buffer, plain string, html []byte, attachmentName string) {
	mixedBoundary := "dbreport-mixed-boundary"
	altBoundary := "dbreport-alt-boundary"

	body.WriteString("Content-Type: multipart/mixed; boundary=\"")
	body.WriteString(mixedBoundary)
	body.WriteString("\"")
	body.WriteString(crlf + crlf)

	body.WriteString("--")
	body.WriteString(mixedBoundary)
	body.WriteString(crlf)
	body.WriteString("Content-Type: multipart/alternative; boundary=\"")
	body.WriteString(altBoundary)
	body.WriteString("\"")
	body.WriteString(crlf + crlf)
	writePlainPart(body, altBoundary, plain)
	writeHTMLPart(body, altBoundary, html)
	closeBoundary(body, altBoundary)
	body.WriteString(crlf)

	writeAttachmentPart(body, mixedBoundary, attachmentName, html)
	closeBoundary(body, mixedBoundary)
}

func writeAttachmentOnlyMessage(body *bytes.Buffer, plain string, html []byte, attachmentName string) {
	boundary := "dbreport-mixed-boundary"
	body.WriteString("Content-Type: multipart/mixed; boundary=\"")
	body.WriteString(boundary)
	body.WriteString("\"")
	body.WriteString(crlf + crlf)
	writePlainPart(body, boundary, plain)
	writeAttachmentPart(body, boundary, attachmentName, html)
	closeBoundary(body, boundary)
}

func writePlainPart(body *bytes.Buffer, boundary string, plain string) {
	body.WriteString("--")
	body.WriteString(boundary)
	body.WriteString(crlf)
	body.WriteString("Content-Type: text/plain; charset=utf-8")
	body.WriteString(crlf)
	body.WriteString("Content-Transfer-Encoding: 7bit")
	body.WriteString(crlf + crlf)
	body.WriteString(normalizePlainBody(plain))
	body.WriteString(crlf)
}

func writeHTMLPart(body *bytes.Buffer, boundary string, html []byte) {
	body.WriteString("--")
	body.WriteString(boundary)
	body.WriteString(crlf)
	body.WriteString("Content-Type: text/html; charset=utf-8")
	body.WriteString(crlf)
	body.WriteString("Content-Transfer-Encoding: base64")
	body.WriteString(crlf + crlf)
	writeBase64Lines(body, html)
	body.WriteString(crlf)
}

func writeAttachmentPart(body *bytes.Buffer, boundary string, attachmentName string, data []byte) {
	encodedName := mime.QEncoding.Encode("utf-8", attachmentName)
	body.WriteString("--")
	body.WriteString(boundary)
	body.WriteString(crlf)
	body.WriteString("Content-Type: text/html; charset=utf-8; name=\"")
	body.WriteString(encodedName)
	body.WriteString("\"")
	body.WriteString(crlf)
	body.WriteString("Content-Disposition: attachment; filename=\"")
	body.WriteString(encodedName)
	body.WriteString("\"")
	body.WriteString(crlf)
	body.WriteString("Content-Transfer-Encoding: base64")
	body.WriteString(crlf + crlf)
	writeBase64Lines(body, data)
	body.WriteString(crlf)
}

func closeBoundary(body *bytes.Buffer, boundary string) {
	body.WriteString("--")
	body.WriteString(boundary)
	body.WriteString("--")
	body.WriteString(crlf)
}

func writeBase64Lines(body *bytes.Buffer, data []byte) {
	encoded := base64.StdEncoding.EncodeToString(data)
	for len(encoded) > 76 {
		body.WriteString(encoded[:76])
		body.WriteString(crlf)
		encoded = encoded[76:]
	}
	body.WriteString(encoded)
	body.WriteString(crlf)
}

func plainTextBody(reportPath string, generatedAt time.Time) string {
	if strings.TrimSpace(reportPath) == "" {
		reportPath = "report.html"
	}
	return fmt.Sprintf("dbreport generated an HTML report.\n\nReport file: %s\nGenerated: %s\n", reportPath, generatedAt.Format("2006-01-02 15:04:05"))
}

func normalizePlainBody(value string) string {
	replacer := strings.NewReplacer("\r\n", "\n", "\r", "\n")
	return replacer.Replace(value)
}

func envelopeAddress(value string) (string, error) {
	addr, err := mail.ParseAddress(value)
	if err != nil {
		return "", err
	}
	return addr.Address, nil
}

func containsHeaderBreak(value string) bool {
	return strings.ContainsAny(value, "\r\n")
}
