#!/bin/bash

# $1 = config file

# Get scope of script
DIR=`dirname "$0"`

SRC="$DIR/.."
BIN="$SRC/bin"
ETC="$SRC/etc"

PID_FILE="$ETC/anda.pid"
LOG_FILE="$ETC/anda.log"

# Kill previous process
if [ -f $PID_FILE ]; then
    PID=$(cat $PID_FILE)
    kill -9 $PID 2>/dev/null
    if [ $? -ne 0 ]; then
        echo "PID $PID is not running"
        rm $PID_FILE
    fi
fi

# Build the project to an executable
go build -o $BIN/anda.o
chmod +x $BIN/anda.o

# Run project using config file and logging output to file
nohup $BIN/anda.o -config $1 >> $LOG_FILE 2>&1 < /dev/null &

# Write PID to pid file
echo $! > $PID_FILE
