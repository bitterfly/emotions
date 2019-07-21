#!/bin/zsh

if [ $# -ne 3 ]; then
    echo "usage: <wav-dir> <output-dir> <batch-num>"
    exit 1
fi

wav_dir="${1}"
output_dir="${2}"
batchnum="${3}"

if [ ! -d ${output_dir} ]; then
    mkdir ${output_dir}
fi

rm -rf ${output_dir}/*

csvs=$(find "${wav_dir}" -type f -name "*.csv")
negative=$(echo  ${csvs} | egrep "(anger|sadness)")
positive=$(echo  ${csvs} | grep "happiness")
neutral=$(echo  ${csvs} | grep "neutral")

echo -e "negative: $(echo ${negative} |wc -l)"
echo -e "positive: $(echo ${positive} |wc -l)"
echo -e "neutral: $(echo ${neutral} |wc -l)"

negative_batch=$(echo -e ${negative} | wc -l | awk -v batchnum=${batchnum} '{print int($0/batchnum)}')
positive_batch=$(echo ${positive} | wc -l | awk  -v batchnum=${batchnum} '{print int($0/batchnum)}')
neutral_batch=$(echo ${neutral} | wc -l | awk  -v batchnum=${batchnum} '{print int($0/batchnum)}')

for i in $(seq 1 ${batchnum}); do
    batch=$(echo ${negative} | shuf -n ${negative_batch})
    paste <(yes "eeg-negative" | head -n ${negative_batch}) <(echo ${batch}) >> ${output_dir}/batch_${i}.txt
    anger=$(comm -23 <(echo ${negative}|sort) <(echo ${batch} | sort))

    batch=$(echo ${positive} | shuf -n ${positive_batch})
    paste <(yes "eeg-positive" | head -n ${positive_batch}) <(echo ${batch}) >> ${output_dir}/batch_${i}.txt
    happiness=$(comm -23 <(echo ${positive}|sort) <(echo ${batch} | sort))

    batch=$(echo ${neutral} | shuf -n ${neutral_batch})
    paste <(yes "eeg-neutral" | head -n ${neutral_batch}) <(echo ${batch}) >> ${output_dir}/batch_${i}.txt
    neutral=$(comm -23 <(echo ${neutral}|sort) <(echo ${batch} | sort))

    sort -o ${output_dir}/batch_${i}.txt ${output_dir}/batch_${i}.txt
done
