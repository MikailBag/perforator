#!/usr/bin/env bash

set -euxo pipefail

mkdir ~/src

(cd ~/src && tar xf ~/code.tgz)

(cd ~/src && ./ya test --show-extra-progress -DCI=github ./perforator)

