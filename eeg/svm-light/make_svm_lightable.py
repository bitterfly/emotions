import sys
import os

if __name__ == "__main__":
    name  = sys.argv[1]
    if os.path.basename(name)[:3] == "pos":
        fclass = "+1"
    else:
        fclass = "-1"

    with open(name, 'r') as f:
        for line in f:
            args = line.split(',')
            print("{}".format(fclass), end=" ")
            for i in range(len(args)-1):
                print("{}:{}".format(i+1, args[i]), end =" ")
            print("")