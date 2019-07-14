#!/bin/zsh

cdir="$(readlink -f "$(dirname "${0}")")"

outputDir="${1}"
for f in "${outputDir}"/*.png; do 
    cp "${f}" "/tmp/foo.png"
    blender -b "${cdir}/colour-mix.blend" -o "/tmp/bla" -f 1 &> /dev/null
    mv -f "/tmp/bla0001.png" "${f}"
done

montage -label '%f' -pointsize 20  "${outputDir}"/eeg-positive* \
          -geometry '512x512+4+4>'  "${outputDir}"/positive.png

montage -label '%f' -pointsize 20 "${outputDir}"/eeg-negative* \
          -geometry '512x512+4+4>'  "${outputDir}"/negative.png

montage -label '%f' -pointsize 20 "${outputDir}"/eeg-neutral* \
          -geometry '512x512+4+4>'  "${outputDir}"/neutral.png


montage "${outputDir}"/positive.png "${outputDir}"/negative.png "${outputDir}"/neutral.png \
          -geometry '1024x1024+4+4>'  "${outputDir}"/all.png