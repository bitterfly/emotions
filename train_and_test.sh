#!/bin/zsh

wavsDir=$1
outputDir=$2

./make_test_data.sh ${wavsDir} ${outputDir}

gmmDir=$(echo ${outputDir}/gmms)

if [[ ! -d ${gmmDir} ]]; then
    mkdir ${gmmDir}
fi

trainDir=$(echo ${outputDir}/train)
testDir=$(echo ${outputDir}/test)

echo "=====TRAIN ${k}======"
train_emotions ${k} "${gmmDir}/gmm" -h ${trainDir}/happiness/* -s ${trainDir}/sadness/* -a ${trainDir}/anger/* -n ${trainDir}/neutral/* 

for k in `seq 2 12`; do
   echo "=====TEST ${k}======"    
    test_emotion "${gmmDir}/gmm_k${k}" -h ${testDir}/happiness/* -s ${testDir}/sadness/* -a ${testDir}/anger/* -n ${testDir}/neutral/* > ${outputDir}/result_k${k}
done

