# dbreport Email Support

`dbreport` can optionally send the generated HTML report over SMTP.

Email is not required for normal report generation.

## Modes

### Disabled

```yaml
email:
  enabled: false
```

No email is sent unless `--email` is passed and email config is valid.

### HTML body

```yaml
email:
  enabled: true
  send_html_body: true
  attach_html: false
```

The generated HTML report is sent as the email body.

### Attachment

```yaml
email:
  enabled: true
  send_html_body: false
  attach_html: true
```

A small plain-text body is sent and the generated HTML report is attached.

### Body and attachment

```yaml
email:
  enabled: true
  send_html_body: true
  attach_html: true
```

The report is included as the body and also attached.

## SMTP config

```yaml
email:
  enabled: false
  smtp_host: "smtp.example.com"
  smtp_port: 587
  starttls: true
  username_env: "DBREPORT_SMTP_USER"
  password_env: "DBREPORT_SMTP_PASSWORD"
  from: "reports@example.com"
  to:
    - "ops@example.com"
  subject: "Daily MariaDB Summary"
  send_html_body: true
  attach_html: false
```

## Secrets

Credentials are loaded from environment variables:

```sh
export DBREPORT_SMTP_USER='smtp-user'
export DBREPORT_SMTP_PASSWORD='smtp-password'
```

Do not place SMTP passwords in `report.yml`.

## STARTTLS

Set `starttls: true` for SMTP servers that support STARTTLS.

## Failure behavior

The report file is generated before email is sent. If email sending fails, the
generated HTML report is preserved and the command exits with the email error
exit code.

## Operational notes

- Test with a non-production recipient first.
- Avoid emailing reports with sensitive data unless recipients are approved.
- Prefer STARTTLS when available.
- Consider using a dedicated SMTP account for reports.
