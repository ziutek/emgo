### Emgo

First of all, to try Emgo you need the [Go compiler](https://golang.org/) installed. The current Emgo compiler and whole process described below requires also some kind of Unix-like operating system. There is a chance that Windows with Cygwin can be used but this was not tested.

You can probably use `go get` to install Emgo but the preffered way is to clone this repository using the git command:

	git clone https://github.com/ziutek/emgo.git

Next you need to build and install egc (Emgo compiler): 

    cd emgo/egc
    go install

For now, Emgo supports only ARM Cortex-M based MCUs. To build code for Cortex-M architecture, you need to install ARM embedded toolchain. You have two options: install a package included in your OS distribution (in case of Debian/Ubuntu Linux):

	apt-get install gcc-arm-none-eabi

or better go to the [GNU ARM Embedded Toolchain website](https://developer.arm.com/open-source/gnu-toolchain/gnu-rm) and download most recent toolchain. This is preffered version of toolchain, try use it before report any bug with compilation.

Installed toolchain contains set of `arm-none-eabi-*` binaries. Find their location and set required enviroment variables:

	export EGCC=path_to_arm_gcc            # eg. /usr/local/arm/bin/arm-none-eabi-gcc
	export EGLD=path_to_arm_linekr         # eg. /usr/local/arm/bin/arm-none-eabi-ld
	export EGAR=path_to_arm_archiver       # eg. /usr/local/arm/bin/arm-none-eabi-ar

	export EGROOT=path_to_egroot_directory # eg. $HOME/emgo/egroot
	export EGPATH=path_to_egpath_directory # eg. $HOME/emgo/egpath

Load/debug helper scripts use also some other tools from the ARM toolchain (eg. `arm-none-eabi-objcopy`). If you downloaded the toolchain manually, you probably need also to add its bin direcotry to the `PATH` enviroment variable:

	export PATH=$PATH:path_to_arm_bin_dir  # eg. /usr/local/arm/bin

Now you are ready to compile some example code. There are two directories that contain examples:

[$EGPATH/src/stm32/examples](https://github.com/ziutek/emgo/tree/master/egpath/src/stm32/examples)

[$EGPATH/src/nrf5/examples](https://github.com/ziutek/emgo/tree/master/egpath/src/nrf5/examples)

Use one that contains example for your MCU/devboard.

For example, to build blinky for STM32 NUCLEO-F411RE board:

	cd $EGPATH/src/stm32/examples/nucleo-f411re/blinky
    ../build.sh

The first compilation may take some time because egc must process all required libraries and runtime. If everything went well you obtain `cortexm4f.elf` binary file.

Compilation can produce two kind of binaries: binaries that should be loaded to RAM or loaded to Flash of your MCU.

Loading to RAM is useful in case of small programs, during working on the code and debuging. Loading into RAM is faster, allows unlimited number of breakpoints, allows to modify constants and even the code from debuger and saves your Flash, which has big but limited number of erase cycles.

To run program loaded to RAM you must change MCU boot option. In case of most STM32 MCUs you simply need to set high BOOT0 and BOOT1 pins. Newer STM32 MCUs do not provide BOOT1 pin, insdead they require change some persistant bits that change the default booting behavior if BOOT0 is set high.

However, eventually your program should be loaded to Flash. Sometimes you have no other alternative: your program is too big to fit in RAM, your MCU does not provide easy way to run program from RAM (eg. nRF51), some bugs may only appear when program runs from Flash.

You need tools to load compiled binary to your MCU's RAM/Flash and allow to debug it. Such tools usually have a hardware part and a software part. In case of STM32 Nucleo or Discovery development boards the hardware part (ST-LINK programmer) is integrated with the board, so you only need the software part, which can be [OpenOCD](http://openocd.org) or [Texane's stlink](https://github.com/texane/stlink). You must install one of them or both before next steps (ensure that `openocd` and/or `st-util` binaries are on your `PATH`).

There is a set of scripts for any board in `example` directory that simplifies loding and debuging process. The `load-oocd.sh` script cah handle SWO output from ITM (Instrumentation Trace Macrocell) but needs [itmsplit](https://github.com/ziutek/itmsplit) to convert binary stream to readable messages. SWO is very useful for debuging and`fmt.Print*` functions by default use ITM trace port as standard output. Install `itmsplit` with the command:

	go get github.com/ziutek/itmsplit
	
and ensure that produced binary is in your `PATH`.

To program your MCU using Texane's stlink run:

	../load-stutil.sh

If you want to use OpenOCD, run:

	../load-oocd.sh

Some examples by default are configured to run from RAM. If you have problem to setup your board to run from RAM, edit `script.ld` file and change the line:

	INCLUDE stm32/loadram

to

	INCLUDE stm32/loadflash

and run `../build.sh`. More editing is need for STM32F1xx series: you additionally have to comment two lines:

	bootOffset = 0x1E0;
	ENTRY(bootRAM)

You can also load your program during debug session in gdb. Run `../debug-stutil.sh` or `../debug-oocd.sh` and next invoke `load` command.

There are also scripts for [Black Magic Probe](https://github.com/blacksphere/blackmagic/wiki): `load-bmp.sh`, `debug-bmp.sh`.

#### Documentation

[Standard library](https://godoc.org/github.com/ziutek/emgo/egroot/src)

[Libraries for STM32, nRF5 and other](https://godoc.org/github.com/ziutek/emgo/egpath/src)

#### Other resources

[YouTube](https://www.youtube.com/channel/UCAW4PLMDGO7_vY4sCG0jg6Q)

[Forum](https://groups.google.com/forum/#!forum/emgo)
