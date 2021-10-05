#!/bin/bash

function clear() {
    # First argument is the Data directory
    find $1 -type f -not -name 'empty.txt' -delete
    echo "Old data removed"
}

# Current dir of this script
current_dir="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

# Dir containing data from the previous experiment
data_dir=${current_dir}/../Data

clear $data_dir