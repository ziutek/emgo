#!/bin/bash

set -e

rm -rf egroot/pkg/* egpath/pkg/*


list=$(find egroot/src egpath/src/stm32 -type d)

for p in $list; do
	if [ -n "$(find $p -maxdepth 1 -type f -name '*.go')" ]; then
		cd $p
		printf "%-48s   " ${p#*/*/}
		if egc; then
			echo OK
		else
			echo Err
		fi
		cd - >/dev/null
	fi
done

echo "--"
