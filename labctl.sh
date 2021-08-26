#/bin/bash
set -e
pwd=$PWD
src=${0%/*}
cd $src
go build --race
cd $pwd
$src/labctl "$@"
