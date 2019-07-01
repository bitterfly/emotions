#!/bin/zsh

batch_files_dir="${1}"
gmm_models_dir="${2}"
result_dir="${3}"

batch_files=$(find ${batch_files_dir} -type f)
for file in $(find ${batch_files_dir} -type f | sort); do
    for k in $(seq 7 10); do
        echo Testing ${file} with k${k}
        train_emotions ${k} ${gmm_models_dir}/gmm_$(basename ${file%.txt}) <(cat $(comm -23 <(echo ${batch_files} | sort) <(echo ${file})))
        test_emotion ${gmm_models_dir}/gmm_$(basename ${file%.txt})_k${k} <(cat ${file}) > ${result_dir}/result_$(basename ${file%.txt})_k${k}
    done
done