### Emgo

To build and install egc (Emgo compiler): 

    cd egc
    go install
  
To build examples, you first have to install ARM toolchain. You can install a package included in your OS distribution. In case of Debian/Ubuntu Linux:

	apt-get install gcc-arm-none-eabi

or go to [GNU ARM Embedded Toolchain](https://developer.arm.com/open-source/gnu-toolchain/gnu-rm) and download most recent toolchain as TAR archive (this is preffered version of toolchain, try use it before report a bug with compilation).

Next set required enviroment variables:

	export EGCC=path_to_arm_gcc            # eg. /usr/local/arm/bin/arm-none-eabi-gcc
	export EGLD=path_to_arm_linekr         # eg. /usr/local/arm/bin/arm-none-eabi-ld
	export EGAR=path_to_arm_archiver       # eg. /usr/local/arm/bin/arm-none-eabi-ar

	export EGROOT=path_to_egroot_directory # eg. $HOME/emgo/egroot
	export EGPATH=path_to_egpath_directory # eg. $HOME/emgo/egpath

Now you are ready to compile some example code. There are two directories that contain examples:

	[$EGPATH/src/stm32/examples](https://github.com/ziutek/emgo/tree/master/egpath/src/stm32/examples)
	[$EGPATH/src/nrf5/examples](https://github.com/ziutek/emgo/tree/master/egpath/src/nrf5/examples)

Use one that contains example for your MCU/devboard.

For example, to build blinky for STM32 Nucleo-F411RE board you need:

	cd $EGPATH/src/stm32/examples/nucleo-f411re/blinky
    ../build.sh

To program your MCU using binary built to run from SRAM:

	../load.sh      # This uses st-util.

or

	../load-oocd.sh # This uses openocd.

To load binary built to run from flash (this erases flash and programs it with new firmware):

	../load.sh flash

or

	../load-oocd.sh flash

To change this SRAM/flash build option you need to edit script.ld file and change the line:

	INCLUDE stm32/loadram

to

	INCLUDE stm32/loadflash

or vice versa. More editing is need for STM32F1xx series.

There are also scripts for [Black Magic Probe](https://github.com/blacksphere/blackmagic/wiki): load-bmp.sh, debug-bmp.sh.

#### Documentation

[Standard library](https://godoc.org/github.com/ziutek/emgo/egroot/src)

[Libraries for STM32, nRF5 and other](https://godoc.org/github.com/ziutek/emgo/egpath/src)

#### Resources

[YouTube](https://www.youtube.com/channel/UCAW4PLMDGO7_vY4sCG0jg6Q)

[Forum](https://groups.google.com/forum/#!forum/emgo)

