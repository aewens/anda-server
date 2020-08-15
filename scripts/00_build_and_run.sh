#!/bin/bash

# $1 = config file
# $2 = log file
# $3 = pid file

# Kill previous process
if [ -f $3 ]; then
    PID=$(cat $3)
    kill -9 $PID 2>/dev/null
    if [ $? -ne 0 ]; then
        echo "PID $PID is not running"
        rm $3
    fi
fi

# Get scope of script
DIR=`dirname "$0"`
SRC="$DIR/.."
BIN="$SRC/bin"

# Build the project to an executable
go build -o $BIN/anda.o
chmod +x $BIN/anda.o

# Run project using config file and logging output to file
nohup $BIN/anda.o -config $1 >> $2 2>&1 < /dev/null &

# Write PID to pid file
echo $! > $3
