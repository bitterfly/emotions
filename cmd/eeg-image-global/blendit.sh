#!/bin/zsh

cdir="$(readlink -f "$(dirname "${0}")")"

outputDir="${1}"
for f in "${outputDir}"/*.png; do 
    cp "${f}" "/tmp/foo.png"
    blender -b "${cdir}/colour-mix.blend" -o "/tmp/bla" -f 1 &> /dev/null
    mv -f "/tmp/bla0001.png" "${f}"
done

montage  "${outputDir}"/* -geometry '512x512+4+4>'  "${outputDir}"/all.png
