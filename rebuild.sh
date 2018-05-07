#!/bin/bash

./clean.sh

list=$(find egroot/src -type d)

for p in $list; do
	if [ -n "$(find $p -maxdepth 1 -type f -name '*.go' |egrep -v '/_|/os|linux')" ]; then
		cd $p
		printf "%-44s   " ${p#*/*/}
		result=$(egc $@ 2>&1)
		if [ "$result" ]; then
			echo "Error:"
			echo
			echo "$result"
			echo
		else
			echo OK
		fi
		cd - >/dev/null
	fi
done

list=$(find egpath/src/*/examples -type d) 

for p in $list; do
	if [ -n "$(find $p -maxdepth 1 -type f -name '*.go' |grep -v '/_')" ]; then
		cd $p
		if [ -x ../build.sh ]; then
			printf "%-44s   " ${p#*/*/}
			result=$(../build.sh $@ 2>&1)
			if [ "$result" ]; then
				overflow=$(
					echo $result |grep "region \`.*' overflowed by .* bytes" \
					|sed "s/.*region \`\(.*\)' overflowed by \(.*\) b.*/\1 overflow: \2 B/g"
				)
				if [ "$overflow" ]; then
					echo "$overflow"
				else
					echo "Error:"
					echo
					echo "$result"
					echo
				fi
			else
				echo OK
			fi
		fi
		cd - >/dev/null
	fi
done

echo "--"
