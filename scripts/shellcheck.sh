#!/bin/sh
set -ex

# SC2034 is unused variable 
# some generated scripts contain 1 or more variables with aren't used
# sometimes.
shellcheck -e SC2034 -s sh samples/*.sh
shellcheck -e SC2034 -s bash samples/*.sh
shellcheck -e SC2034 -s dash samples/*.sh
shellcheck -e SC2034 -s ksh samples/*.sh
