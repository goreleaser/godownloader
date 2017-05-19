#!/bin/sh
set -e

git_clone_or_update() {
  giturl=$1
  gitrepo=${giturl##*/}   # foo.git
  gitrepo=${gitrepo%.git} # foo
  if [ ! -d "$gitrepo" ]; then
    git clone "$giturl"
  else
    (cd "$gitrepo" && git pull > /dev/null)
  fi
}

git_clone_or_update https://github.com/client9/posixshell.git 
cd posixshell

echo "package main"
echo ""
echo 'const shellfn = `'
cat \
  license.sh \
  is_command.sh \
  uname_os.sh \
  uname_arch.sh \
  uname_os_check.sh \
  uname_arch_check.sh \
  untar.sh \
  mktmpdir.sh \
  http_download.sh \
  github_api.sh \
  github_last_release.sh \
  hash_sha256.sh \
  license_end.sh | \
  grep -v '^#' |grep -v ' #' | tr -s '\n'

echo '`'

