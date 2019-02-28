#!/bin/bash

function parse_line {
    read end start _ name <<< "${1}"
    if [[ "${name}" =~ .*_tag_(.*)_tag.* ]]; then
        echo "${BASH_REMATCH[1]} ${start} ${end}"
    fi
}

function parse_log {
    while read -d $'\n' line; do
        parse_line "${line}"
    done < <(paste <(cut -d ' ' -f 1 "${1}" | tail -n +2 ; echo 42) "${1}")
}

function parse_tags {
    log="/tmp/mpv_log"

    grep --line-buffered '^Playing' | ts '%.s' > "${log}"
    parse_log "${log}"
}

mpv --fs "${1}" | parse_tags
