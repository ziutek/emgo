#!/bin/bash

set -e

#EGC='egc -O g'
EGC='egc'

rm -rf egroot/pkg/* egpath/pkg/*

list=$(find egroot/src egpath/src -type d)

for p in $list; do
	if [ -n "$(find $p -maxdepth 1 -type f -name '*.go')" ]; then
		cd $p
		rm -f *.elf *.bin *.sizes
		printf "%-44s   " ${p#*/*/}
		if $EGC; then
			echo OK
		else
			echo Err
		fi
		cd - >/dev/null
	fi
done

echo "--"
