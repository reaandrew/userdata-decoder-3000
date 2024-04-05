terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.66.1"
    }
  }

  backend "s3" {
    bucket  = "state.andrewrea.co.uk"
    key     = "state/userdata-decoder-3000/terraform.tfstate"
    region  = "eu-west-2"
    encrypt = true
  }
}

provider "aws" {
  # us-east-1 instance
  region = "us-east-1"
  alias  = "use1"
}

# The attribute `${data.aws_caller_identity.current.account_id}` will be current account number.
data "aws_caller_identity" "current" {}

# The attribute `${data.aws_region.current.name}` will be current region
data "aws_region" "current" {}

locals {
  content_types = {
    ".html" : "text/html",
    ".css" : "text/css",
    ".js" : "text/javascript",
    ".png" : "image/x-png",
    ".jpg" : "image/jpeg",
    ".woff2": "application/font-woff2",
    ".woff": "application/font-woff",
    ".ttf": "font/ttf",
    ".svg": "image/svg+xml",
    ".map": "application/json"
  }
}

resource "aws_s3_bucket" "bucket" {
  bucket = var.bucket_name
}

data "aws_iam_policy_document" "cloud_front_access_policy" {
  statement {
    sid       = "AllowCloudFrontServicePrincipal"
    actions   = ["s3:GetObject"]
    resources = ["${aws_s3_bucket.bucket.arn}/*"]

    condition {
      test     = "StringEquals"
      variable = "AWS:SourceArn"
      values   = [aws_cloudfront_distribution.distribution.arn]
    }

    principals {
      type        = "Service"
      identifiers = ["cloudfront.amazonaws.com"]
    }
  }
}

resource "aws_s3_bucket_policy" "cdn-oac-bucket-policy" {
  bucket = aws_s3_bucket.bucket.id
  policy = data.aws_iam_policy_document.cloud_front_access_policy.json
}

resource "aws_s3_object" "file" {
  for_each     = fileset(path.module, "dist/**/*.{html,css,js,png,jpg,woff,woff2,ttf,svg,map}")
  bucket       = aws_s3_bucket.bucket.id
  key          = replace(each.value, "/^dist//", "")
  source       = each.value
  content_type = lookup(local.content_types, regex("\\.[^.]+$", each.value), null)
  etag         = filemd5(each.value)
}

resource "aws_cloudfront_origin_access_control" "hosting" {
  name                              = var.sub_domain_name
  description                       = "UserData Decoder 3000 Website Policy"
  origin_access_control_origin_type = "s3"
  signing_behavior                  = "always"
  signing_protocol                  = "sigv4"

}

resource "aws_cloudfront_origin_access_identity" "hosting" {
}

resource "aws_cloudfront_response_headers_policy" "security_headers_policy" {
  name = "userdata-decoder-3000-security-headers-policy"

  custom_headers_config {
    items {
      header = "permissions-policy"
      override = true
      value = "accelerometer=(), camera=(), geolocation=(), gyroscope=(), magnetometer=(), microphone=(), payment=(), usb=()"
    }
  }

  security_headers_config {
    content_type_options {
      override = true
    }
    frame_options {
      frame_option = "DENY"
      override = true
    }
    referrer_policy {
      referrer_policy = "same-origin"
      override = true
    }
    xss_protection {
      mode_block = true
      protection = true
      override = true
    }
    strict_transport_security {
      access_control_max_age_sec = "63072000"
      include_subdomains = true
      preload = true
      override = true
    }
    content_security_policy {
      content_security_policy = "default-src 'none'; frame-src https://td.doubleclick.net; img-src 'self' https://www.google.com https://www.google.co.uk https://fonts.gstatic.com https://www.googletagmanager.com https://region1.google-analytics.com; script-src 'self' https://googleads.g.doubleclick.net https://www.googletagmanager.com https://cmp.osano.com https://9cd02f31fa91.eu-west-2.captcha-sdk.awswaf.com https://9cd02f31fa91.a7cfe63f.eu-west-2.captcha.awswaf.com https://9cd02f31fa91.a7cfe63f.eu-west-2.token.awswaf.com https://cdn.jsdelivr.net; style-src 'self' https://www.googletagmanager.com 'unsafe-inline' https://fonts.googleapis.com https://cdn.jsdelivr.net; worker-src 'self' blob:; connect-src 'self' https://www.google.com https://googleads.g.doubleclick.net https://pagead2.googlesyndication.com https://www.google-analytics.com https://disclosure.api.osano.com https://consent.api.osano.com/record https://region1.google-analytics.com https://tattle.api.osano.com https://9cd02f31fa91.a7cfe63f.eu-west-2.token.awswaf.com https://api.secronyx.com; font-src 'self' data: fonts.gstatic.com; style-src 'self' https://www.googletagmanager.com 'unsafe-inline' fonts.googleapis.com https://cdn.jsdelivr.net; object-src 'none';"
      override = true
    }
  }
}


resource "aws_cloudfront_distribution" "distribution" {
  enabled             = true
  is_ipv6_enabled     = true
  default_root_object = "index.html"

  origin {
    domain_name              = aws_s3_bucket.bucket.bucket_regional_domain_name
    origin_id                = var.sub_domain_name
    origin_access_control_id = aws_cloudfront_origin_access_control.hosting.id
  }

  viewer_certificate {
    acm_certificate_arn      = aws_acm_certificate_validation.website_cert.certificate_arn
    ssl_support_method       = "sni-only"
    minimum_protocol_version = "TLSv1.1_2016"
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
      locations        = []
    }
  }

  aliases = [var.sub_domain_name]

  default_cache_behavior {
    cache_policy_id        = "658327ea-f89d-4fab-a63d-7e88639e58f6"
    viewer_protocol_policy = "redirect-to-https"
    compress               = true
    allowed_methods        = ["DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT"]
    cached_methods         = ["GET", "HEAD"]
    target_origin_id       = var.sub_domain_name
    response_headers_policy_id = aws_cloudfront_response_headers_policy.security_headers_policy.id
  }

  depends_on = [aws_s3_bucket.bucket]
}

data "aws_route53_zone" "selected" {
  name         = "${var.domain_name}."
  private_zone = false
}

resource "aws_acm_certificate" "website_domain_name_cert" {
  provider          = aws.use1
  domain_name       = var.sub_domain_name
  validation_method = "DNS"
  lifecycle {
    create_before_destroy = true
  }
  depends_on = [data.aws_route53_zone.selected]
}

resource "aws_acm_certificate_validation" "website_cert" {
  provider                = aws.use1
  certificate_arn         = aws_acm_certificate.website_domain_name_cert.arn
  validation_record_fqdns = [
    aws_route53_record.website_cert_validation.fqdn,
  ]
  depends_on = [aws_route53_record.website_cert_validation]
}

resource "aws_route53_record" "website_cert_validation" {
  name    = aws_acm_certificate.website_domain_name_cert.domain_validation_options.*.resource_record_name[0]
  type    = aws_acm_certificate.website_domain_name_cert.domain_validation_options.*.resource_record_type[0]
  zone_id = data.aws_route53_zone.selected.id
  records = [
    aws_acm_certificate.website_domain_name_cert.domain_validation_options.*.resource_record_value[0]
  ]
  ttl             = 60
  allow_overwrite = true
  depends_on      = [aws_acm_certificate.website_domain_name_cert]
}

resource "aws_route53_record" "www" {
  zone_id = data.aws_route53_zone.selected.zone_id
  name    = var.sub_domain_name
  type    = "A"

  alias {
    name                   = aws_cloudfront_distribution.distribution.domain_name
    zone_id                = aws_cloudfront_distribution.distribution.hosted_zone_id
    evaluate_target_health = false
  }

  depends_on = [aws_cloudfront_distribution.distribution, aws_route53_record.website_cert_validation]
}