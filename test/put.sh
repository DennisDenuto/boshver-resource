#!/bin/bash

set -e

source $(dirname $0)/helpers.sh

it_can_put_and_set_first_version() {
  local repo=$(init_repo)

  local src=$(mktemp -d $TMPDIR/put-src.XXXXXX)
  echo 1.2 > $src/some-new-file

  # cannot push to repo while it's checked out to a branch
  git -C $repo checkout refs/heads/master

  put_uri $repo $src some-new-file | jq -e "
    .version == {number: \"1.2\"}
  "

  # switch back to master
  git -C $repo checkout master

  test -e $repo/some-file
  test "$(cat $repo/some-file)" = 1.2
}

it_can_put_and_set_same_version() {
  local repo=$(init_repo)

  set_version $repo 1.2

  local src=$(mktemp -d $TMPDIR/put-src.XXXXXX)
  echo 1.2 > $src/some-new-file

  # cannot push to repo while it's checked out to a branch
  git -C $repo checkout refs/heads/master

  put_uri $repo $src some-new-file | jq -e "
    .version == {number: \"1.2\"}
  "

  # switch back to master
  git -C $repo checkout master

  test -e $repo/some-file
  test "$(cat $repo/some-file)" = 1.2
}

it_can_put_and_set_over_existing_version() {
  local repo=$(init_repo)

  set_version $repo 0.0

  local src=$(mktemp -d $TMPDIR/put-src.XXXXXX)
  echo 1.2 > $src/some-new-file

  # cannot push to repo while it's checked out to a branch
  git -C $repo checkout refs/heads/master

  put_uri $repo $src some-new-file | jq -e "
    .version == {number: \"1.2\"}
  "

  # switch back to master
  git -C $repo checkout master

  test -e $repo/some-file
  test "$(cat $repo/some-file)" = 1.2
}



run it_can_put_and_set_first_version
run it_can_put_and_set_same_version
run it_can_put_and_set_over_existing_version
