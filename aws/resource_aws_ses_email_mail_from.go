package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAwsSesEmailMailFrom() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsSesEmailMailFromSet,
		Read:   resourceAwsSesEmailMailFromRead,
		Update: resourceAwsSesEmailMailFromSet,
		Delete: resourceAwsSesEmailMailFromDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"mail_from_domain": {
				Type:     schema.TypeString,
				Required: true,
			},
			"behavior_on_mx_failure": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  ses.BehaviorOnMXFailureUseDefaultValue,
			},
		},
	}
}

func resourceAwsSesEmailMailFromSet(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).sesConn

	behaviorOnMxFailure := d.Get("behavior_on_mx_failure").(string)
	emailAddress := d.Get("email").(string)
	mailFromDomain := d.Get("mail_from_domain").(string)

	input := &ses.SetIdentityMailFromDomainInput{
		BehaviorOnMXFailure: aws.String(behaviorOnMxFailure),
		Identity:            aws.String(emailAddress),
		MailFromDomain:      aws.String(mailFromDomain),
	}

	_, err := conn.SetIdentityMailFromDomain(input)
	if err != nil {
		return fmt.Errorf("Error setting MAIL FROM domain: %s", err)
	}

	d.SetId(emailAddress)

	return resourceAwsSesEmailMailFromRead(d, meta)
}

func resourceAwsSesEmailMailFromRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).sesConn

	emailAddress := d.Id()

	readOpts := &ses.GetIdentityMailFromDomainAttributesInput{
		Identities: []*string{
			aws.String(emailAddress),
		},
	}

	out, err := conn.GetIdentityMailFromDomainAttributes(readOpts)
	if err != nil {
		log.Printf("error fetching MAIL FROM domain attributes for %s: %s", emailAddress, err)
		return err
	}

	d.Set("email", emailAddress)

	if v, ok := out.MailFromDomainAttributes[emailAddress]; ok {
		d.Set("behavior_on_mx_failure", v.BehaviorOnMXFailure)
		d.Set("mail_from_domain", v.MailFromDomain)
	} else {
		d.Set("behavior_on_mx_failure", v.BehaviorOnMXFailure)
		d.Set("mail_from_domain", "")
	}

	return nil
}

func resourceAwsSesEmailMailFromDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).sesConn

	emailAddress := d.Id()

	deleteOpts := &ses.SetIdentityMailFromDomainInput{
		Identity:       aws.String(emailAddress),
		MailFromDomain: nil,
	}

	_, err := conn.SetIdentityMailFromDomain(deleteOpts)
	if err != nil {
		return fmt.Errorf("Error deleting SES Email Mail From: %s", err)
	}

	return nil
}
