#!/usr/bin/env bash
set -x

export LOGRUS_LEVEL="debug"
export TEST_OPTIONS='-v'

make e2e-tests
