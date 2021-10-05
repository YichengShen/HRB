from os import walk, path
from contextlib import ExitStack
import numpy as np


DATA_DIR = "./Data"

_, _, filenames = next(walk(DATA_DIR), (None, None, []))

# Note that file names contain "Throuput"
# Not "Throughput"
throughput_fnames = [f for f in filenames if "Throuput" in f]

# latency_fnames = [f for f in filenames if "Latency" in f]

with ExitStack() as stack:
    files = [stack.enter_context(open(path.join(DATA_DIR, fname)))
             for fname in throughput_fnames]
    mean_throughput = np.mean([float(f.read().strip()) for f in files])

print("Mean Throughput")
print(mean_throughput)