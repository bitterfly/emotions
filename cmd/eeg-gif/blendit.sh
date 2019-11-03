#!/bin/zsh

cdir="$(readlink -f "$(dirname "${0}")")"

outputDir="${1}"
for f in "${outputDir}"/*.png; do 
    cp "${f}" "/tmp/foo.png"
    blender -b "${cdir}/colour-mix.blend" -o "/tmp/bla" -f 1 &> /dev/null
    mv -f "/tmp/bla0001.png" "${f}"
done

ffmpeg -r 5 -i ${outputDir}/'α_%05d.png' -c:v libx264 ${outputDir}/α.mp4
ffmpeg -r 5 -i ${outputDir}/'β_%05d.png' -c:v libx264 ${outputDir}/β.mp4
ffmpeg -r 5 -i ${outputDir}/'γ_%05d.png' -c:v libx264 ${outputDir}/γ.mp4
ffmpeg -r 5 -i ${outputDir}/'δ_%05d.png' -c:v libx264 ${outputDir}/δ.mp4