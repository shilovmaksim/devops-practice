import csv, sys, errno, getopt, time, random
#import numpy as np

def readFile(name):
    data = []
    count = 0
    with open(name, 'r') as csvfile:
        reader = csv.reader(csvfile)
        for row in reader:
            if count==0:
                count=1
            else:
                data.append(int(row[0]))
    return data

def writeFile(name, data):
    with open(name, 'w') as f:
        writer = csv.writer(f, delimiter=',', quotechar='"', quoting=csv.QUOTE_MINIMAL)
        writer.writerow(['value'])
        for i in data:
            writer.writerow([i])

def main(argv):
    inputfile = []
    outputfile = 'def_output.csv'
    try:
        opts, args = getopt.gnu_getopt(argv,"hi:o:",["ofile="])
    except getopt.GetoptError:
        print('test.py -o <outputfile>')
        sys.exit(2)
   
    for opt, arg in opts:
        if opt == '-h':
            print('test.py [<inputfile>] -o <outputfile>')
            sys.exit()
        elif opt in ("-o", "--ofile"):
            outputfile = arg

    inputfile = args
    print('Input file is', inputfile)
    print('Output file is', outputfile)

    try:
        data1 = readFile(inputfile[0])
        data2 = readFile(inputfile[1])
    
        if len(data1)!=len(data2):
            sys.exit(errno.ENOKEY)

        time.sleep(random.randint(0,2000)/1000.0)
    
        writeFile(outputfile, data1+data2)        
    except Exception as e:
        print('Exception:', e)
        sys.exit(errno.EFAULT)

if __name__ == "__main__":
   main(sys.argv[1:])
