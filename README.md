CloudStack Terraform Provider
=============================

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 1.0.x
-	[Go](https://golang.org/doc/install) 1.20+ (to build the provider plugin)

See wiki: https://github.com/apache/cloudstack-terraform-provider/wiki

Installing from Github Release
------------------------------

User can install the CloudStack Terraform Provider using the [Github Releases](https://github.com/apache/cloudstack-terraform-provider/releases) with the installation steps below.

Replace the `RELEASE` version with the version you're trying to install and use.

The valid `ARCH` options are:

- linux_amd64
- linux_386
- linux_arm64
- linux_arm
- darwin_amd64
- darwin_arm64
- freebsd_amd64
- freebsd_386
- freebsd_arm64
- freebsd_arm

Steps for installation:

```
RELEASE=0.5.0
ARCH=darwin_arm64
mkdir -p ~/.terraform.d/plugins/local/cloudstack/cloudstack/${RELEASE}/${ARCH}
wget "https://github.com/apache/cloudstack-terraform-provider/releases/download/v${RELEASE}/cloudstack-terraform-provider_${RELEASE}_${ARCH}.zip"
unzip cloudstack-terraform-provider_${RELEASE}_${ARCH}.zip -d cloudstack-terraform-provider_${RELEASE}
mv cloudstack-terraform-provider_${RELEASE}/cloudstack-terraform-provider_v${RELEASE} ~/.terraform.d/plugins/local/cloudstack/cloudstack/${RELEASE}/${ARCH}/terraform-provider-cloudstack_v${RELEASE}
```

To use the locally installed provider, please use the following in your main.tf etc, and then run `terraform init`:

```
terraform {
  required_providers {
    cloudstack = {
      source = "local/cloudstack/cloudstack"
      version = "0.5.0"
    }
  }
}

provider "cloudstack" {
  # Configuration options
}
```

Note: this can be used when users are not able to install using the Terraform registry.

Installing from Terrafrom registry
----------------------------------
To install the CloudStack provider, copy and paste the below code into your Terraform configuration. Then, run terraform init.
```sh
terraform {
  required_providers {
    cloudstack = {
      source = "cloudstack/cloudstack"
      version = "0.5.0"
    }
  }
}

provider "cloudstack" {
  # Configuration options
}
```

User hitting installation issue using registry can install using the local install method.

Documentation
-------------

For more details on how to use the provider, click [here](website/) or visit https://registry.terraform.io/providers/cloudstack/cloudstack/latest/docs

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.16+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

Clone repository to: `$GOPATH/src/github.com/apache/cloudstack-terraform-provider`

```sh
$ mkdir -p $GOPATH/src/github.com/apache; cd $GOPATH/src/github.com/apache
$ git clone git@github.com:apache/cloudstack-terraform-provider
```

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/apache/cloudstack-terraform-provider
$ make build
$ ls $GOPATH/bin/terraform-provider-cloudstack
```
Once the build is ready, you have to copy the binary into Terraform locally (version appended).

On Linux and Mac this path is at ~/.terraform.d/plugins,
On Windows at %APPDATA%\terraform.d\plugins,

```sh
$  cd ~
$  mkdir -p ~/.terraform.d/plugins/localdomain/provider/cloudstack/0.4.0/linux_amd64
$  cp $GOPATH/bin/terraform-provider-cloudstack ~/.terraform.d/plugins/localdomain/provider/cloudstack/0.4.0/linux_amd64
```

Testing the Provider
--------------------

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests you will need to run the CloudStack Simulator. Please follow these steps to prepare an environment for running the Acceptance tests:

```sh
docker pull apache/cloudstack-simulator

or pull it with a particular build tag

docker pull apache/cloudstack-simulator:4.17.2.0

docker run --name simulator -p 8080:5050 -d apache/cloudstack-simulator

or

docker run --name simulator -p 8080:5050 -d apache/cloudstack-simulator:4.17.2.0
```

When Docker started the container you can go to http://localhost:8080/client and login to the CloudStack UI as user `admin` with password `password`. It can take a few minutes for the container is fully ready, so you probably need to wait and refresh the page for a few minutes before the login page is shown.

Once the login page is shown and you can login, you need to provision a simulated data-center:

```sh
docker exec -it cloudstack-simulator python /root/tools/marvin/marvin/deployDataCenter.py -i /root/setup/dev/advanced.cfg
```

If you refresh the client or login again, you will now get passed the initial welcome screen and be able to go to your account details and retrieve the API key and secret. Export those together with the URL:

```sh
$ export CLOUDSTACK_API_URL=http://localhost:8080/client/api
$ export CLOUDSTACK_API_KEY=r_gszj7e0ttr_C6CP5QU_1IV82EIOtK4o_K9i_AltVztfO68wpXihKs2Tms6tCMDY4HDmbqHc-DtTamG5x112w
$ export CLOUDSTACK_SECRET_KEY=tsfMDShFe94f4JkJfEh6_tZZ--w5jqEW7vGL2tkZGQgcdbnxNoq9fRmwAtU5MEGGXOrDlNA6tfvGK14fk_MB6w
```

In order for all the tests to pass, you will need to create a new (empty) project in the UI called `terraform`. When the project is created you can run the Acceptance tests against the CloudStack Simulator by simply running:

```sh
$ make testacc
```

Sample Terraform configuration when testing locally
------------------------------------------------------------
Below is an example configuration to initialize provider and create a Virtual Machine instance

```sh
$ cat provider.tf
terraform {
  required_providers {
    cloudstack = {
      source = "localdomain/provider/cloudstack"
      version = "0.4.0"
    }
  }
}

provider "cloudstack" {
  # Configuration options
  api_url    = "${var.cloudstack_api_url}"
  api_key    = "${var.cloudstack_api_key}"
  secret_key = "${var.cloudstack_secret_key}"
}

resource "cloudstack_instance" "web" {
  name             = "server-1"
  service_offering = "Small Instance"
  network_id       = "df5fc279-86d5-4f5d-b7e9-b27f003ca3fc"
  template         = "616fe117-0c1c-11ec-aec4-1e00610002a9"
  zone             = "2b61ed5d-e8bd-431d-bf52-d127655dffab"
}
```
## History

This codebase relicensed under APLv2 and donated to the Apache CloudStack
project under an [IP
clearance](https://github.com/apache/cloudstack/issues/5159) process and
imported on 26th July 2021.

## License

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License. You may obtain a copy of the
License at <http://www.apache.org/licenses/LICENSE-2.0>
