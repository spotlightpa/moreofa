#!/bin/bash

set -eu -o pipefail

# Get the directory that this script file is in
THIS_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

cd "$THIS_DIR"

function _default() {
	# shellcheck disable=SC2119
	api
}

function _die() {
	echo >&2 "Fatal: ${*}"
	exit 1
}

function _installed() {
	hash "$1" >/dev/null 2>&1
}

function _git-xargs() {
	local PATTERN=$1
	shift
	git ls-files --exclude="$PATTERN" -ciz | xargs -0 -I _ "$@"
}

function help() {
	local SCRIPT=$0
	cat <<EOF
Usage

	$SCRIPT <task> <args>

Tasks:

EOF
	compgen -A function | grep -e '^_' -v | sort | xargs printf ' - %s\n'
	exit 2
}

function sql() {
	set -x
	format:sql
	sql:sqlc
	set +x
}

function sql:sqlc() {
	_installed sqlc || _die "sqlc not installed"
	sqlc generate
	sqlc compile
	sqlc vet
}

function db:migrate() {
	dbmate -d sql/migrations/ -u sqlite:./comments.db --no-dump-schema "$@"
}

function test() {
	set -x
	test:backend
	test:misc
	set +x
}

function test:backend() {
	go test -race -v ./...
}

function test:misc() {
	_git-xargs '*.sh' shellcheck _
	go mod tidy -diff
	sqruff lint sql
}

function format() {
	set -x
	format:go
	format:sh
	format:sql
	set +x
}

function format:go() {
	gofmt -s -w .
}

function format:sh() {
	_git-xargs '*.sh' shfmt -w _
}

function format:sql() {
	sqruff fix sql
}

function db:copy-prod() {
	litestream restore -o comments.db s3://moreofa-backup.data.spotlightpa.org/comments.db
}

# shellcheck disable=SC2120
function api() {
	# shellcheck disable=SC1091
	[[ -f .env ]] && echo "Using .env file" && source .env
	go run . "$@"
}

function check-deps() {
	_installed shellcheck || echo "install https://www.shellcheck.net"
	_installed shfmt || echo "install https://github.com/mvdan/sh"
	_installed sqlc || echo "install https://sqlc.dev"
	_installed sqruff
	_installed dbmate
}

TIMEFORMAT="Task completed in %1lR"
time "${@:-_default}"
