CloudStack Terraform Provider
=============================

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.8 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/apache/cloudstack-terraform-provider`

```sh
$ mkdir -p $GOPATH/src/github.com/apache; cd $GOPATH/src/github.com/apache
$ git clone git@github.com:apache/cloudstack-terraform-provider
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/apache/cloudstack-terraform-provider
$ make build
```

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.8+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make bin
...
$ $GOPATH/bin/cloudstack-terraform-provider
...
```

Testing the Provider
--------------------

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests you will need to run the CloudStack Simulator. Please follow these steps to prepare an environment for running the Acceptance tests:

```sh
$ docker pull svanharmelen/simulator:4.12.0.0
$ docker run -d -p 8080:8080 --name cloudstack svanharmelen/simulator:4.12.0.0
```

When Docker started the container you can go to http://localhost:8080/client and login to the CloudStack UI as user `admin` with password `password`. It can take a few minutes for the container is fully ready, so you probably need to wait and refresh the page for a few minutes before the login page is shown.

Once the login page is shown and you can login, you need to provision a simulated data-center:

```sh
$ docker exec -ti cloudstack python /root/tools/marvin/marvin/deployDataCenter.py -i /root/setup/dev/advanced.cfg
```

If you refresh the client or login again, you will now get passed the initial welcome screen and be able to go to your account details and retrieve the API key and secret. Export those together with the URL:

```sh
$ export CLOUDSTACK_API_URL=http://localhost:8080/client/api
$ export CLOUDSTACK_API_KEY=r_gszj7e0ttr_C6CP5QU_1IV82EIOtK4o_K9i_AltVztfO68wpXihKs2Tms6tCMDY4HDmbqHc-DtTamG5x112w
$ export CLOUDSTACK_SECRET_KEY=tsfMDShFe94f4JkJfEh6_tZZ--w5jqEW7vGL2tkZGQgcdbnxNoq9fRmwAtU5MEGGXOrDlNA6tfvGK14fk_MB6w
```

In order for all the tests to pass, you will need to create a new (empty) project in the UI called `terraform`. When the project is created you can run the Acceptance tests against the CloudStack Simulator by simply runnning:

```sh
$ make testacc
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
