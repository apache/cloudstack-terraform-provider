#
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
#

#!/bin/bash

# This script rewrites [GH-nnnn]-style references in the CHANGELOG.md file to
# be Markdown links to the given github issues.
#
# This is run during releases so that the issue references in all of the
# released items are presented as clickable links, but we can just use the
# easy [GH-nnnn] shorthand for quickly adding items to the "Unrelease" section
# while merging things between releases.

set -e

if [[ ! -f CHANGELOG.md ]]; then
  echo "ERROR: CHANGELOG.md not found in pwd."
  echo "Please run this from the root of the terraform provider repository"
  exit 1
fi

if [[ `uname` == "Darwin" ]]; then
  echo "Using BSD sed"
  SED="sed -i.bak -E -e"
else
  echo "Using GNU sed"
  SED="sed -i.bak -r -e"
fi

PROVIDER_URL="https:\/\/github.com\/terraform-providers\/terraform-provider-cloudstack\/issues"

$SED "s/GH-([0-9]+)/\[#\1\]\($PROVIDER_URL\/\1\)/g" -e 's/\[\[#(.+)([0-9])\)]$/(\[#\1\2))/g' CHANGELOG.md

rm CHANGELOG.md.bak
