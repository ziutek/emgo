#!/bin/bash

set -e

#EGC='egc -O g'
EGC='egc'

rm -rf egroot/pkg/* 

list=$(find egroot/src -type d)

for p in $list; do
	if [ -n "$(find $p -maxdepth 1 -type f -name '*.go' |grep -v '/_')" ]; then
		cd $p
		printf "%-44s   " ${p#*/*/}
		if $EGC; then
			echo OK
		else
			echo Err
		fi
		cd - >/dev/null
	fi
done

rm -rf egpath/pkg/*

list=$(find egpath/src/stm32/examples -type d)

for p in $list; do
	if [ -n "$(find $p -maxdepth 1 -type f -name '*.go' |grep -v '/_')" ]; then
		cd $p
		if [ -x ../build.sh ]; then
			rm -f *.elf *.bin *.sizes
			printf "%-44s   " ${p#*/*/}
			if ../build.sh; then
				echo OK
			else
				echo Err
			fi
		fi
		cd - >/dev/null
	fi
done

echo "--"
