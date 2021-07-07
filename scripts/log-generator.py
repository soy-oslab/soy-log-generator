#!python3
import random
import time
import argparse
import logging
import os
import sys
import tqdm

parser = argparse.ArgumentParser(description="generic log generator")

parser.add_argument('--source-file', dest='src', type=str,
                    help="source log file")
parser.add_argument('--target-file', dest='dst', type=str,
                    help="target log file")
parser.add_argument('--flush-ratio', dest='flush_ratio', metavar='N',
                    default=0.1, type=float, help="every N seconds flush")
parser.add_argument('--max-sleep', dest='interval', metavar='N',
                    default=1, type=int, help="sleep max N seconds")
parser.add_argument('--append', dest='use_append', nargs='?', default=1,
                    help="using append manner based to open file")

args = parser.parse_args()

if not os.path.isfile(args.src):
    logging.error(f"you must specify the valid \
            source log file ({args.src})")
    sys.exit(-1)

src = open(args.src, "r")
try:
    if args.use_append is None:
        dst = open(args.dst, 'a')
    else:
        dst = open(args.dst, "w")
except OSError:
    logging.error(f"you must specify the valid \
            target file directory ({args.dst})")
    sys.exit(-1)

_ = input("type any key to generate log file start")
now = time.time()
for line in tqdm.tqdm(src.readlines()):
    dst.write(line)
    if time.time() - now > args.flush_ratio:
        dst.flush()
        now = time.time()
    time.sleep(random.randint(0, args.interval))
src.close()
dst.close()
