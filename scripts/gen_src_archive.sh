#!/usr/bin/env bash

PROJECT=generator
DEST="build/${PROJECT}.src.tar.gz"

# Check that git-archive-all is installed
if ! git-archive-all --help &> /dev/null; then
  echo "# This script uses 'git-archive-all', please install:"
  echo "# pip install git-archive-all"
  exit 1
fi

mkdir -p build
git-archive-all --prefix="$PROJECT" --force-submodules $DEST

echo "# $0: DONE!"
ls -lh $DEST
