#!/bin/zsh

for k in `seq 3 12`; do
    echo "======== K: ${k} ==============="
    train_emotions ${k} "/tmp/gmms/gmm" /home/do/go/src/github.com/bitterfly/emotions/wavs/{anger.wav,happiness.wav,sadness.wav}
    echo -e "Anger\t$(test_emotion "/tmp/gmms/gmm_k${k}" /home/do/Emotions/database_wavs/anger)" >> "/tmp/results/gmm_${k}"
    echo -e "Happiness\t$(test_emotion "/tmp/gmms/gmm_k${k}" /home/do/Emotions/database_wavs/happiness)" >> "/tmp/results/gmm_${k}"
    echo -e "Sadness\t$(test_emotion "/tmp/gmms/gmm_k${k}" /home/do/Emotions/database_wavs/sadness)" >> "/tmp/results/gmm_${k}"
done

