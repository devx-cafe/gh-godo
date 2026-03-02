#!/usr/bin/env bash

set -eou pipefail

PREFIX="🍰  "
echo "$PREFIX Running $(basename $0)"

# GitHub CLI Dependencies
set +e
gh auth status >/dev/null 2>&1
AUTH_OK=$?
set -e
if [ $AUTH_OK -ne 0 ]; then
  echo "$PREFIX ⚠️ Not logged into GitHub CLI"
  echo "$PREFIX    This is not looking good  — we want GitHub CLI to work!"
else
  echo "$PREFIX ✅ GitHub Authentication is working smooth!"
  gh extension install devx-cafe/gh-tt
  echo "$PREFIX ✅ Installed the TakT gh cli extension from devx-cafe/gh-tt "
  gh alias import .devcontainer/.gh_alias.yml --clobber
  echo "$PREFIX ✅ Installed the gh shorthand aliases"    
  
fi

git config --global --add safe.directory $(pwd)
echo "$PREFIX ✅ Setting up safe git repository to prevent dubious ownership errors"

git config --local --get include.path | grep -e ../.gitconfig >/dev/null 2>&1 || git config --local --add include.path ../.gitconfig
echo "$PREFIX ✅ Setting up git configuration to support .gitconfig in repo-root"

# Install Go dependencies if go.mod exists
if [ -f "go.mod" ]; then
    echo "$PREFIX Installing Go dependencies (go mod download)..."
    go mod download

    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
    sh -s -- -b $(go env GOPATH)/bin latest
    echo "$PREFIX ✅ Installing golangci-lint"

fi


echo "$PREFIX ✅ SUCCESS"
exit 0
