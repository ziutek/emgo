### Emgo

To try Emgo you first need to download it. You can probably use `go get` but preffered way is to clone this repository:

	git clone https://github.com/ziutek/emgo.git

Next you need to build and install egc (Emgo compiler): 

    cd emgo/egc
    go install
  
All examples are for ARM Cortex-M based MCUs. To build them, you need to install ARM toolchain. You have two options: install a package included in your OS distribution (in case of Debian/Ubuntu Linux):

	apt-get install gcc-arm-none-eabi

or better go to the [GNU ARM Embedded Toolchain website](https://developer.arm.com/open-source/gnu-toolchain/gnu-rm) and download most recent toolchain (this is preffered version of toolchain, try use it before report any bug with compilation).

Installed toolchain contains set of `arm-none-eabi-*` binaries. Find their location and set required enviroment variables:

	export EGCC=path_to_arm_gcc            # eg. /usr/local/arm/bin/arm-none-eabi-gcc
	export EGLD=path_to_arm_linekr         # eg. /usr/local/arm/bin/arm-none-eabi-ld
	export EGAR=path_to_arm_archiver       # eg. /usr/local/arm/bin/arm-none-eabi-ar

	export EGROOT=path_to_egroot_directory # eg. $HOME/emgo/egroot
	export EGPATH=path_to_egpath_directory # eg. $HOME/emgo/egpath

	export PATH=$PATH:path_to_arm_bin_dir  # eg. /usr/local/arm/bin

Now you are ready to compile some example code. There are two directories that contain examples:

[$EGPATH/src/stm32/examples](https://github.com/ziutek/emgo/tree/master/egpath/src/stm32/examples)

[$EGPATH/src/nrf5/examples](https://github.com/ziutek/emgo/tree/master/egpath/src/nrf5/examples)

Use one that contains example for your MCU/devboard.

For example, to build blinky for STM32 NUCLEO-F411RE board:

	cd $EGPATH/src/stm32/examples/nucleo-f411re/blinky
    ../build.sh

First compilation may take some time because `egc` must process all required libraries and runtime. If everything went well you obtain `cortexm4.elf` binary.

Compilation can produce two kind of binaries: binaries that should be loaded to RAM or to Flash of your MCU.

Load into RAM is useful in case of small programs, during working on the code and debuging. Loading into RAM is faster, allows unlimited number of breakpoints, allows to modify constants and even the code from debuger and saves your Flash, which has big but limited number of erase cycles.

To run program loaded to RAM you must change MCU boot option. In case of most STM32 MCUs you simply need to set high BOOT0 and BOOT1 pins.

But eventually your program should be loaded to Flash. Sometimes you simply can not load to RAM: program is too big, your MCU does not provide easy way to run program loaded to RAM (eg. nRF51). At last, some bugs may only appear when program runs from Flash.

At this point you need some tools to load compiled binary to your MCU's RAM/Flash and allow to debug it. Such tools usually have a hardware part and a software part. In case of STM32 Nucleo or Discovery development boards the hardware part (ST-LINK programmer) is integrated with the board, so you only need the software part, which can be [OpenOCD](http://openocd.org) or [Texane's stlink](https://github.com/texane/stlink).

To program your MCU using binary built to run from RAM:

	../load.sh      # This uses st-util.

or

	../load-oocd.sh # This uses openocd.

To load binary built to run from Flash (this erases Flash and programs it with new firmware):

	../load.sh flash

or

	../load-oocd.sh flash

To change this RAM/Flash build option you need to edit `script.ld` file and change the line:

	INCLUDE stm32/loadram

to

	INCLUDE stm32/loadflash

or vice versa. More editing is need for STM32F1xx series.

In case of nRF51 load scripts have no `flash` option: binary is always loaded to Flash.

You can also load you program during debug session in gdb. Try `../debug.sh` or `../debug-oocd.sh`.

There are also scripts for [Black Magic Probe](https://github.com/blacksphere/blackmagic/wiki): `load-bmp.sh`, `debug-bmp.sh`.

#### Documentation

[Standard library](https://godoc.org/github.com/ziutek/emgo/egroot/src)

[Libraries for STM32, nRF5 and other](https://godoc.org/github.com/ziutek/emgo/egpath/src)

#### Resources

[YouTube](https://www.youtube.com/channel/UCAW4PLMDGO7_vY4sCG0jg6Q)

[Forum](https://groups.google.com/forum/#!forum/emgo)
