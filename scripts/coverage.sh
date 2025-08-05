#!/usr/bin/env bash

# Copyright 2025 The Toodofun Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http:www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


echo "mode: atomic" > coverage.txt

for d in $(find ./* -maxdepth 10 -type d); do
    if ls $d/*.go &> /dev/null; then
        go test  -coverprofile=profile.out -covermode=atomic $d
        if [ -f profile.out ]; then
            cat profile.out | grep -v "mode: " >> /tmp/coverage.txt
            rm profile.out
        fi
    fi
done

echo "coverage output: /tmp/coverage.txt"