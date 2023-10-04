#!/bin/sh
set -e

# available endpoints:
SIGN_UP="auth/sign-up"   # POST
SIGN_IN="auth/sign-in"   # GET
LIST_NOTES="api/notes"   # GET
NEW_NOTE="api/notes/new" # POST

host="http://127.0.0.1"
port="8080"
endpoint=""
method="" # GET / POST
header=""
body=""
username="anon"
passwd="qwerty"
modes=""
token=""

_printHelp() {
	printf "%s\n" \
"Usage:
  -H | --host      (defalut is 127.0.0.1)
  -P | --port      (default is 8080)
  -u | --username  Set username
  -p | --passwd    Set password
  -t | --token     Set JWT token
  -n | --note      Set a note title
  -d | --note-desc Set a note description

  Set your modes (you may do a chain a modes) by
  supplying one of four arguments (without dashes):
    sign_up    (requires username and password)
    sign_in    (requires username and password)
    new_note   (requires JWT token)
    list_notes (requires JWT token)

  'sign_in' mode also sets the token variable,
  so you'll be able to proceed with token dependent
  modes.

  Example:
    $(basename "$0") -P 8000 -t \"qwerasdfzxcv\" \\
      -n \"My new note\" -d \"Just testing notes\" new_note"
}

[ $# -lt 1 ] && {
	_printHelp
	exit 1
}

_fatal() {
	printf "%s\n" "$*"
	exit 1
}

_sendRequest() {
	curl -w '\nStatus code: %{http_code}\n' \
		-X "${method}" "${url}" \
		-H "${header}" -d "${body}"
}

_getToken() {
	endpoint="$SIGN_IN"
	url="${host}:${port}/${endpoint}"
	method="GET"
	header=""
	body="{\"username\":\"${username}\",\"password\":\"${passwd}\"}"

	response="$(_sendRequest)"
	printf "%s\n" "$response"

	status="$(printf "%s\n" "$response" | tail -c4)"
	[ "$status" -ne 200 ] && _fatal "Login failed"

	token="$(printf "%s" "$response" | head -n -1 | \
		rev | cut -d'"' -f2 | rev)"
			printf "%s\n" "Token has been set ($token)"
}

_signIn() {
	{ [ -z "$username" ] || [ -z "$passwd" ]; } && \
		_fatal "Must have non-empty 'username' and " \
			"password for a 'sign_in' mode"

	_getToken
}

_signUp() {
	{ [ -z "$username" ] || [ -z "$passwd" ]; } && \
		_fatal "Must have non-empty 'username' and " \
			"password for a 'sign_in' mode"

	endpoint="$SIGN_UP"
	url="${host}:${port}/${endpoint}"
	method="POST"
	header=""
	body="{\"username\":\"${username}\",\"password\":\"${passwd}\"}"

	response="$(_sendRequest)"
	printf "%s\n" "$response"

	status="$(printf "%s\n" "$response" | tail -c4)"
	if [ "$status" -ne 200 ]; then
		_fatal "Unable to register as $username." "
Maybe username is taken?"
	else
		echo "Successfully registered!"
	fi
}

_noteNew() {
	[ -z "$note_title" ] && \
		_fatal "Must have a note defined to create one"

	endpoint="$NEW_NOTE"
	url="${host}:${port}/${endpoint}"
	method="POST"
	header="Authorization: Bearer $token"
	body="{\"title\":\"${note_title}\",\"description\":\"${note_desc}\"}"

	response="$(_sendRequest)"
	printf "%s\n" "$response"

	status="$(printf "%s\n" "$response" | tail -c4)"
	if [ "$status" -ne 200 ]; then
		_fatal "Unable to create a note"
	else
		echo "Successfully created a new note!"
	fi
}

_noteList() {
	endpoint="$LIST_NOTES"
	url="${host}:${port}/${endpoint}"
	method="GET"
	header="Authorization: Bearer $token"
	body=""

	response="$(_sendRequest)"
	printf "%s\n" "$response"

	status="$(printf "%s\n" "$response" | tail -c4)"
	[ "$status" -ne 200 ] && \
		_fatal "Unable to get a list of notes for $username"
}

while [ $# -gt 0 ]; do
	case $1 in
	"-h" | "-help" | "--help")
		_printHelp
		exit 0
		;;
	"-u" | "--username")
		{ [ -z "$2" ] || [ "$(printf "%s\n" "$2" | \
			head -c1)" = "-" ]; } && {
			_fatal "No username provided"
		}

		username="$2"
		shift 2
		;;
	"-p" | "--pwd" | "--passwd" | "--password")
		{ [ -z "$2" ] || [ "$(printf "%s\n" "$2" | \
			head -c1)" = "-" ]; } && {
			_fatal "No password provided"
		}

		passwd="$2"
		shift 2
		;;
	"-H" | "--host")
		{ [ -z "$2" ] || [ "$(printf "%s\n" "$2" | \
			head -c1)" = "-" ]; } && {
			_fatal "No host provided"
		}

		host="$2"
		shift 2
		;;
	"-P" | "--port")
		{ [ -z "$2" ] || [ "$(printf "%s\n" "$2" | \
			head -c1)" = "-" ]; } && {
			_fatal "No port provided"
		}

		port="$2"
		shift 2
		;;
	"-T" | "--token")
		{ [ -z "$2" ] || [ "$(printf "%s\n" "$2" | \
			head -c1)" = "-" ]; } && {
			_fatal "No token provided"
		}

		token="$2"
		shift 2
		;;

	"-n" | "--note")
		{ [ -z "$2" ] || [ "$(printf "%s\n" "$2" | \
			head -c1)" = "-" ]; } && {
			_fatal "Note must have at least the title"
		}

		note_title="$2"

		shift 2
		;;
	"-d" | "--note-desc")
		{ [ -z "$2" ] || [ "$(printf "%s\n" "$2" | \
			head -c1)" = "-" ]; } && {
			_fatal "No note description provided"
		}

		note_desc="$2"

		shift 2
		;;
	*)
		[ "$1" = "sign_up" ] && {
			modes="${modes}sign_up "
			shift; continue
		}
		[ "$1" = "sign_in" ] && {
			modes="${modes}sign_in "
			shift; continue
		}
		[ "$1" = "new_note" ] && {
			modes="${modes}new_note "
			shift; continue
		}
		[ "$1" = "list_notes" ] && {
			modes="${modes}list_notes "
			shift; continue
		}

		_fatal "Unknown argument - $1"
	esac
done

[ -z "$modes" ] && _fatal "No mode has been set, exiting..."

printf "%s\n" "Going to connect to ${host}:${port}..."
sleep 1

for mode in $modes; do
	case $mode in
		sign_up)
			_signUp
			;;
		sign_in)
			_signIn
			;;
		new_note)
		   	[ -z "$token" ] && \
		   		_fatal "Token is required for 'new_note' mode"
			_noteNew
			;;
		list_notes)
		   	[ -z "$token" ] && \
		   		_fatal "Token is required for 'list_notes' mode"
			_noteList
	esac
done
