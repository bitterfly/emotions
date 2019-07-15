#!/bin/zsh
beep_file="${1}"
videos_dir="${2}"
output_file="${3}"

if [[ $# != 3 ]]; then
    echo "<beep-file> <video-dir> <output-file>"
fi

emotional_files=$(find "${videos_dir}" -type f)
i=1

before=$(date +"%s.%N")
aplay "${beep_file}" &>/dev/null
echo -e "start\t${p}\t${before}\t$(date +"%s.%N")" > "${output_file}"
echo "Start." 1>&2

rec="0"
while true; do
    read -s -k key
    if [[ "${key}" == " " ]]; then
        if [ "${rec}" -eq "0" ]; then
            rec="1"
            echo "Recording..." 1>&2
            before=$(date +"%s.%N")
        else
            rec="0"
            echo "Done." 1>&2
            echo -e "audio\t${p}\t${before}\t$(date +"%s.%N")" >> "${output_file}"
        fi
    elif [[ "${key}" == $'\n' ]]; then
        if [[ ${i} > $(echo "${emotional_files}" | wc -l) ]]; then
            echo "No more videos in playlist" 1>&2
        else
            file=$(echo "${emotional_files}" | sed -n ${i}p)
            i=$((i+1))
            before=$(date +"%s.%N")
            mpv --fs ${file} &>/dev/null
            echo -e "video\t${p}\t${before}\t$(date +"%s.%N")" >> "${output_file}"
        fi
    fi
done