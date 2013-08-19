#General healer for tsuru PaaS

[![Build Status](https://travis-ci.org/globocom/tsuru-healer.png?branch=master)](https://travis-ci.org/globocom/tsuru-healer)

tsuru-healer detects whether a virtual machine - belonging to a Load Balancer -
is down, and fixes it by replacing the VM.

##Running

To be able to talk to tsuru one must obtain a token via `tsr token` command. Export this token
in the environment variable `TSURU_TOKEN`.
The next step is to run the binary, e.g.:

    $ ./healer tsuruapi.com:8080
