---
layout: "aws"
page_title: "AWS: ses_email_mail_from"
sidebar_current: "docs-aws-resource-ses-email-mail-from"
description: |-
  Provides an SES email MAIL FROM resource
---

# aws_ses_email_mail_from

Provides an SES email MAIL FROM resource.

~> **NOTE:** For the MAIL FROM email to be fully usable, this resource should be paired with the [aws_ses_email_identity resource](/docs/providers/aws/r/ses_email_identity.html). To validate the MAIL FROM domain, a DNS MX record is required. To pass SPF checks, a DNS TXT record may also be required. See the [Amazon SES MAIL FROM documentation](https://docs.aws.amazon.com/ses/latest/DeveloperGuide/mail-from-set.html) for more information.

## Example Usage

```hcl
resource "aws_ses_email_mail_from" "example" {
  email           = "${aws_ses_email_identity.example.email}"
  mail_from_domain = "bounce.example.com"
}

# Example SES Email Identity
resource "aws_ses_email_identity" "example" {
  email = "email@example.com"
}

# Example Route53 MX record
resource "aws_route53_record" "example_ses_email_mail_from_mx" {
  zone_id = "${aws_route53_zone.example.id}"
  name    = "${aws_ses_email_mail_from.example.mail_from_domain}"
  type    = "MX"
  ttl     = "600"
  records = ["10 feedback-smtp.us-east-1.amazonses.com"] # Change to the region in which `aws_ses_email_identity.example` is created
}

# Example Route53 TXT record for SPF
resource "aws_route53_record" "example_ses_email_mail_from_txt" {
  zone_id = "${aws_route53_zone.example.id}"
  name    = "${aws_ses_email_mail_from.example.mail_from_domain}"
  type    = "TXT"
  ttl     = "600"
  records = ["v=spf1 include:amazonses.com -all"]
}
```

## Argument Reference

The following arguments are required:

* `email` - (Required) Email address to set MAIL FROM domain.
* `mail_from_domain` - (Required) Subdomain which is to be used as MAIL FROM address (Required for DMARC validation)

The following arguments are optional:

* `behavior_on_mx_failure` - (Optional) The action that you want Amazon SES to take if it cannot successfully read the required MX record when you send an email. Defaults to `UseDefaultValue`. See the [SES API documentation](https://docs.aws.amazon.com/ses/latest/APIReference/API_SetIdentityMailFromDomain.html) for more information.

## Attributes Reference

In addition to the arguments, which are exported, the following attributes are exported:

* `id` - The email address.

## Import

MAIL FROM email can be imported using the `email` attribute, e.g.

```
$ terraform import aws_ses_email_mail_from.example email@example.com
```
