#!/bin/zsh

photos=${1}
system_pictures=${2}

black=$(realpath "${system_pictures}/black2.png")

ln -sf "${black}" "/tmp/feh.png"

feh -Z -F -R 0.1 -Y "/tmp/feh.png" &
pid=$!
while read sp; do
	sleep 1
	ln -sf "${black}" "/tmp/feh.png"
	sleep 0.2
	ln -sf "${sp}" "/tmp/feh.png"
done < <(find "${system_pictures}/numbers" -type f -exec realpath {} \; | sort -r)

sleep 1

while read p; do
	ln -sf "${p}" "/tmp/feh.png"
	before=$(date +"%s")
	sleep 3
	echo -e "$(realpath ${p})\t${before}\t$(date +"%s")"
	ln -sf "${black}" "/tmp/feh.png"
	sleep 2
done < <(find ${photos} -type f -exec realpath {} \;)
kill ${pid}
rm "/tmp/feh.png"
