import sys
import xml.etree.ElementTree as ET
import time
import datetime
import math

if __name__ == "__main__":
    log = sys.argv[1]

    beep_end = 0
    last_audio_file = ""
    last_audio_file_start = 0
    last_audio_file_end = 0
    audio = []

    with open(log, "r") as l:
        for line in l:
            spl = line.strip().split(" ")
            if spl[0] == "start":
                beep_end = float(spl[2])
            elif spl[0] == "audio":
                if spl[1] != last_audio_file:
                    # start new audio file
                    if last_audio_file_start != 0:
                        print("%s %f %f" % (last_audio_file, last_audio_file_start - beep_end, last_audio_file_end - last_audio_file_start))
                    last_audio_file_start = float(spl[2])
                    last_audio_file = spl[1]
                last_audio_file_end = float(spl[-1])
        print("%s %f %f" % (last_audio_file, last_audio_file_start - beep_end, last_audio_file_end - last_audio_file_start))
        