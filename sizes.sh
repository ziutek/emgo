#!/bin/sh

getFlash() {
	echo $(expr $1 + $2)
}

getBSS() {
	echo $(expr $3)
}

sumFlash=0
sumBSS=0

for elf in $(find * -name 'cortexm*.elf'); do
	size=$(arm-none-eabi-size $elf |tail -1)
	flash=$(getFlash $size)
	bss=$(getBSS $size)
	sumFlash=$(expr $sumFlash + $flash)
	sumBSS=$(expr $sumBSS + $bss)
	printf '%-60s %6d %5d\n' $(dirname $elf) $flash $bss
done

printf '* %65d %5d\n' $sumFlash $sumBSS