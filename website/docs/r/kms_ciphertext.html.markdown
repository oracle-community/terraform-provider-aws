---
subcategory: "KMS (Key Management)"
layout: "aws"
page_title: "AWS: aws_kms_ciphertext"
description: |-
    Provides ciphertext encrypted using a KMS key
---

# Resource: aws_kms_ciphertext

The KMS ciphertext resource allows you to encrypt plaintext into ciphertext
by using an AWS KMS customer master key. The value returned by this resource
is stable across every apply. For a changing ciphertext value each apply, see
the [`aws_kms_ciphertext` data source](/docs/providers/aws/d/kms_ciphertext.html).

~> **Note:** All arguments including the plaintext be stored in the raw state as plain-text.
[Read more about sensitive data in state](https://www.terraform.io/docs/state/sensitive-data.html).

## Example Usage

```terraform
resource "aws_kms_key" "oauth_config" {
  description = "oauth config"
  is_enabled  = true
}

resource "aws_kms_ciphertext" "oauth" {
  key_id = aws_kms_key.oauth_config.key_id

  plaintext = <<EOF
{
  "client_id": "e587dbae22222f55da22",
  "client_secret": "8289575d00000ace55e1815ec13673955721b8a5"
}
EOF
}
```

## Argument Reference

This resource supports the following arguments:

* `region` - (Optional) Region where this resource will be [managed](https://docs.aws.amazon.com/general/latest/gr/rande.html#regional-endpoints). Defaults to the Region set in the [provider configuration](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#aws-configuration-reference).
* `plaintext` - (Required) Data to be encrypted. Note that this may show up in logs, and it will be stored in the state file.
* `key_id` - (Required) Globally unique key ID for the customer master key.
* `context` - (Optional) An optional mapping that makes up the encryption context.

## Attribute Reference

This resource exports the following attributes in addition to the arguments above:

* `ciphertext_blob` - Base64 encoded ciphertext
