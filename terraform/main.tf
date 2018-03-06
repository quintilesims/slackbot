resource "aws_dynamodb_table" "mod" {
  name           = "${var.dynamodb_table_name}"
  read_capacity  = "${var.dynamodb_read_capacity}"
  write_capacity = "${var.dynamodb_write_capacity}"
  hash_key       = "Key"

  attribute {
    name = "Key"
    type = "S"
  }
}

resource "aws_iam_user" "mod" {
  name = "${var.iam_user_name}"
}

resource "aws_iam_access_key" "mod" {
  user = "${aws_iam_user.mod.name}"
}

data "template_file" "user_policy" {
  template = "${file("${path.module}/user_policy.json")}"

  vars {
    dynamodb_table_arn = "${aws_dynamodb_table.mod.arn}"
  }
}

resource "aws_iam_user_policy" "mod" {
  name   = "${aws_iam_user.mod.name}"
  user   = "${aws_iam_user.mod.name}"
  policy = "${data.template_file.user_policy.rendered}"
}

data "aws_region" "current" {}

data "template_file" "deploy" {
  template = "${file("${path.module}/Dockerrun.aws.json")}"

  vars {
    docker_image   = "${var.docker_image}"
    slack_token    = "${var.slack_token}"
    aws_access_key = "${aws_iam_access_key.mod.id}"
    aws_secret_key = "${aws_iam_access_key.mod.secret}"
    aws_region     = "${data.aws_region.current.name}"
    dynamodb_table = "${aws_dynamodb_table.mod.name}"
  }
}

resource "layer0_deploy" "mod" {
  name    = "${var.deploy_name}"
  content = "${data.template_file.deploy.rendered}"
}

resource "layer0_service" "mod" {
  name        = "${var.service_name}"
  environment = "${var.environment_id}"
  deploy      = "${layer0_deploy.mod.id}"
  scale       = "${var.scale}"
}
