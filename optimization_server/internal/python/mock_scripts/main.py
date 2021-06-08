import sys, errno, getopt, time

def main(argv):
    try:
        if len(argv) == 0:
            return

        if argv[0]=="error":    
            sys.exit(errno.EFAULT)
        if argv[0]=="success":
            sys.exit(0)
        if argv[0]=="sleep":
            time.sleep(float(argv[1])/1000.0)
            return
        if argv[0]=="print":
            print(argv[1])
            return
    
    except Exception as e:
        print('Exception:', e)
        sys.exit(errno.EFAULT)

if __name__ == "__main__":
   main(sys.argv[1:])
