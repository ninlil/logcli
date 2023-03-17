#!/bin/bash

pwd
echo "...args = $*"
echo ""

rm xyzzy
echo "to error " >&2
echo ""

date
exit 1