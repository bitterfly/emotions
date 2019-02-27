#!/bin/zsh
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

ln -sf "${green}" "/tmp/feh.png"
feh -Z -F -R 0.1 -Y "/tmp/feh.png" &
pid=$!

# =============== PICTURES =====================

sleep 0.5
ln -sf "${part1}" "/tmp/feh.png"
sleep 2
ln -sf "${green}" "/tmp/feh.png"
sleep 0.8
ln -sf "${description1}" "/tmp/feh.png"
sleep 8


echo "Pictures start: $(date +"%s")"
ln -sf "${black}" "/tmp/feh.png"
while read sp; do
	sleep 0.2
	ln -sf "${sp}" "/tmp/feh.png"
	sleep 0.8
	ln -sf "${black}" "/tmp/feh.png"
done < <(find "${system_data}/numbers" -type f -exec realpath {} \; | sort -r)

sleep 0.5

while read p; do
	ln -sf "${p}" "/tmp/feh.png"
	before=$(date +"%s")
	sleep 3
	echo -e "picture\t${p}\t${before}\t$(date +"%s")"
	ln -sf "${black}" "/tmp/feh.png"
	sleep 2
done < <(find ${photos} -type f -exec realpath {} \;)
echo "Pictures end: $(date +"%s")"
ln -sf "${green}" "/tmp/feh.png"

# ================= VIDEOS ===========================
aplay "$(realpath ${system_data}/pitch.wav)" &
ln -sf "${part2}" "/tmp/feh.png"
sleep 2
ln -sf "${green}" "/tmp/feh.png"
sleep 1
ln -sf "${description2}" "/tmp/feh.png"
sleep 15

ln -sf "${black}" "/tmp/feh.png"
while read sp; do
	sleep 0.2
	ln -sf "${sp}" "/tmp/feh.png"
	sleep 0.8
	ln -sf "${black}" "/tmp/feh.png"
done < <(find "${system_data}/numbers" -type f -exec realpath {} \; | sort -r)
sleep 0.5

i=1
echo "Videos start: $(date +"%s")"
while read vf; do
	before=$(date +"%s")
	mpv --no-terminal --fs --ontop ${vf}
	echo -e "video\t${vf}\t${before}\t$(date +"%s")"
	
	ln -sf "${green}" "/tmp/feh.png"	
	wmctrl -a "feh"

	sleep 0.5
	before=$(date +"%s")	
	while read sp; do
		ln -sf "${sp}" "/tmp/feh.png"
		sleep 1
	done < <(find "${system_data}/green_numbers" -type f -exec realpath {} \; | sort)
	echo -e "audio\t${vf}\t${before}\t$(date +"%s")"
	
	ln -sf "${green}" "/tmp/feh.png"	
	sleep 0.5
	before=$(date +"%s")
	ln -sf "${texts}/${i}.png" "/tmp/feh.png"
	sleep 10
	echo -e "neutral\t${texts}/${i}\t${before}\t$(date +"%s")"
	ln -sf "${green}" "/tmp/feh.png"
	sleep 0.5

	i=$((i+1))
done < <(find  ${videos} -type f -exec realpath {} \; | shuf)

kill ${pid}
rm "/tmp/feh.png"

echo "Videos end: $(date +"%s")"