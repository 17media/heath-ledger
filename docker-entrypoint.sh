#!/bin/sh
if [ $DEVELOPEMENT = "1" ];
then
    echo Development mode.
    go get github.com/codegangsta/gin
    gin run
else
    ./app
fi
