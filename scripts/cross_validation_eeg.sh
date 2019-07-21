#!/bin/zsh

batch_files_dir="${1}"
gmm_models_dir="${2}"
result_dir="${3}"

batch_files=$(find ${batch_files_dir} -type f)
for file in $(find ${batch_files_dir} -type f | sort); do
    echo Testing ${file}
    train_eeg "gmm" 0 ${gmm_models_dir}/gmm_$(basename ${file%.txt}) <(cat $(comm -23 <(echo ${batch_files} | sort) <(echo ${file})))
    test_eeg "gmm" 0 ${gmm_models_dir}/gmm_$(basename ${file%.txt}) <(cat ${file}) > ${result_dir}/result_$(basename ${file%.txt}).res 2> ${result_dir}/result_$(basename ${file%.txt}).err 
done