#/bin/bash
#
# Compile and run labctl
#
# Add a labctl alias in your shell. E.g. in ~/.zshrc add the following:
#  alias labctl='/home/kellerza/go/src/github.com/labctl/labctl/labctl.sh'
#
set -e
pwd=$PWD
src=${0%/*}
cd $src
go build --race
cd $pwd
$src/labctl "$@"
