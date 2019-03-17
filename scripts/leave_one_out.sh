#!/bin/zsh

eeg_directory=${1}
result_file=${2}
keyword=${3}

echo "" > ${result_file}

negative=$(find ${eeg_directory} -type f -name "negative${keyword}*" -exec readlink -f {} \;| shuf)
positive=$(find ${eeg_directory} -type f -name "positive${keyword}*" -exec readlink -f {} \;| shuf)
neutral=$(find ${eeg_directory} -type f -name "neutral*" -exec readlink -f {} \; | shuf)

c=$(echo ${negative} | wc -l)
len=$(perl -e "print int(${c} * 0.7 + 0.5)")
negative_train=$(echo ${negative} | head -n "${len}")
negative_test=$(echo ${negative} | tail -n +$((len+1)))


c=$(echo ${positive} | wc -l)
len=$(perl -e "print int(${c} * 0.7 + 0.5)")
positive_train=$(echo ${positive} | head -n "${len}")
positive_test=$(echo ${positive} | tail -n +$((len+1)))


c=$(echo ${neutral} | wc -l)
len=$(perl -e "print int(${c} * 0.7 + 0.5)")
neutral_train=$(echo ${neutral} | head -n "${len}")
neutral_test=$(echo ${neutral} | tail -n +$((len+1)))

echo net-train ${negative_train}
echo neg-test ${negative_test}

echo pos-train ${positive_train}
echo pos-test ${positive_test}

echo neu-train ${neutral_train}
echo neu-test ${neutral_test}


echo -e "Whole file\teeg-negative\teeg-neutral\teeg-possitive" >> ${result_file}
train_eeg 1 /tmp/foo --eeg-positive $(echo ${positive_train}) --eeg-negative $(echo ${negative_train}) --eeg-neutral $(echo ${neutral_train})
test_eeg 1 /tmp/foo --eeg-positive $(echo ${positive_test}) --eeg-negative $(echo ${negative_test}) --eeg-neutral $(echo ${neutral_test}) 2>/dev/null >> ${result_file}
echo "" >> ${result_file}

for dur in `seq 200 200 5000`; do 
    echo -e "${dur}\teeg-negative\teeg-neutral\teeg-possitive" >> ${result_file}
    
    train_eeg ${dur} /tmp/foo --eeg-positive $(echo ${positive_train}) --eeg-negative $(echo ${negative_train}) --eeg-neutral $(echo ${neutral_train})
    test_eeg ${dur} /tmp/foo --eeg-positive $(echo ${positive_test}) --eeg-negative $(echo ${negative_test}) --eeg-neutral $(echo ${neutral_test}) 2>/tmp/log >> ${result_file}
    echo "" >> ${result_file}
done