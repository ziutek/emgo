#!/bin/sh

set -e

cd ../../stm32

rm -rf hal/raw hal/irq

cd o
for target in *; do
	cd $target
	for pkg in *; do
		halpkg=../../hal/raw/$pkg
		mkdir -p $halpkg
		cd $pkg
		for f in *.go; do
			half=../$halpkg/$target--$f
			echo "// +build $target" >$half
			echo >>$half
			cat $f >>$half
		done
		cd ..
	done
	cd ..
done
cd ../hal

mv raw/irq .
