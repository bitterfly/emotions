#!/bin/bash

cdir="$(dirname "$(readlink -f "${0}")")"
vid_dir="/tmp/generated_videos"

function info {
    echo "${@}" 1>&2
}

function make_name {
    if [[ ! -r "${2}" ]]; then
        info "${2} does not exist!"
        exit 1
    fi

    mkdir -p "${vid_dir}"
    if [[ -f "${3}" ]]; then
        hash="$(echo "$(sha256sum "${2}" ; sha256sum "${3}"; echo "${1}")" | sha256sum | cut -c1-8)"
    else
        hash="$(echo "$(sha256sum "${2}" ; echo "${1}")" | sha256sum | cut -c1-8)"
    fi
    echo "${vid_dir}/${hash}.mp4"
}

function tag {
    read name
    tagged="${name%.mp4}_tag_${1}_tag.mp4"
    info "tagging"
    info "    ${name}"
    info " -> ${tagged}"

    ln -sf "${name}" "${tagged}"
    echo "${tagged}"
}

function audio {
    # first argument is the image, second argument is the audio
    name="$(make_name 42 "${1}" "${2}")"

    if [[ -f "${name}" ]]; then
        echo "${name}"
        return 0    # already generated
    fi

    info "rendering audio with image"
    info "    ${1}"
    info "  , ${2}"
    info " -> ${name}"
    
    ffmpeg \
        -r 25 \
        -stream_loop -1 \
        -i "${1}" \
        -i "${2}" \
        -c:v libx264 \
        -tune stillimage \
        -c:a aac \
        -b:a 192k \
        -pix_fmt yuv420p \
        -shortest \
        -preset veryslow \
        -t $(soxi -D "${2}") \
        "${name}" 2>/tmp/ffmpeg_log

    if [[ $? -ne 0 ]]; then
        info "video creation failed for ${name} :("
        rm -rf "${name}"
        exit 1
    fi

    echo "${name}"
}

function video {
    name="$(make_name "42" "${1}")"

    if [[ -f "${name}" ]]; then
        echo "${name}"
        return 0    # already generated
    fi

    info "copying"
    info "    ${1}"
    info " -> ${name}"

    cp -rf "${1}" "${name}" || exit 1

    echo "${name}"
}

function image {
    name="$(make_name "${1}" "${2}")"

    if [[ -f "${name}" ]]; then
        echo "${name}"
        return 0    # already generated
    fi

    info "rendering"
    info "    ${2}"
    info " -> ${name}"
    
    # framerate is 1/duration
    ffmpeg \
        -r $(perl -e "print 1/${1}") \
        -i "${2}" \
        -preset veryslow \
        -c:v libx264 \
        -vf fps=25 \
        -pix_fmt yuv420p \
        "${name}" 2> /tmp/ffmpeg_log
    if [[ $? -ne 0 ]]; then
        info "video creation failed for ${name} :("
        rm -rf "${name}"
        exit 1
    fi

    echo "${name}"
}

function countdown {
    IFS=$'\n'
    for sp in $(find "${system_data}/numbers" -type f -exec realpath {} \; | sort -r); do
        image 0.2 "${black}"
        image 0.8 "${sp}"
    done
    image 0.2 "${black}"
}

function countup {
    IFS=$'\n'
    for sp in $(find "${system_data}/green_numbers" -type f -exec realpath {} \; | sort); do
        image 1 "${sp}" | tag "audio ${1}"
    done
    image 0.2 "${green}"
}

photos=${1}
videos=${2}
texts=$(realpath ${3})
system_data=${4}

black=$(realpath "${system_data}/black.png")
green=$(realpath "${system_data}/green.png")
part1=$(realpath "${system_data}/part1.png")
part2=$(realpath "${system_data}/part2.png")
description1=$(realpath "${system_data}/description1.png")
description2=$(realpath "${system_data}/description2.png")
beep=$(realpath "${system_data}/beep.wav")

audio "${green}" "${beep}" | tag start

image 0.5 "${green}"
image 2   "${part1}"
image 0.8 "${green}"
image 8   "${description1}"

countdown

IFS=$'\n'
for p in $(find ${photos} -type f -exec realpath {} \; | shuf); do
    image 3 "${p}" | tag "picture $(basename ${p})"
    image 2 "${black}"
done

exit 1


image 0.5 "${green}"
image   2 "${part2}"
image   1 "${green}"
image  15 "${description2}"
image 0.2 "${black}"

countdown

i=1
IFS=$'\n'
for vf in $(find  ${videos} -type f -exec realpath {} \; | shuf); do
    video "${vf}" | tag "video $(basename ${vf})"
    countup $(basename ${vf})
    image 0.5 ${green}
    image 12 ${texts}/${i}.png | tag "text ${i}"
    image 0.5 ${green}
    i=$((i+1))
done
