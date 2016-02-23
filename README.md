### Emgo

To build and install egc: 

    cd egc
    go install
  
To build some example, you first have to set two enviroment variables:

	export EGROOT=PATH_TO_EGROOT_DIRECTORY # eg. $HOME/emgo/egroot
	export EGPATH=PATH_TO_EGPATH_DIRECTORY # eg. $HOME/emgo/egpath

Next go to example directory and build it:

	cd $EGPATH/src/stm32/examples/nucleo-f411re/blinky
    ../build.sh

To program your MCU using binary built to run from SRAM:

	../load.sh      # This uses st-util

or

	../load-oocd.sh # This uses openocd

To load binary built to run from flash (this erases flash and programs it with new firmware):

	../load.sh flash

or

	../load-oocd.sh flash


[Home page](https://sites.google.com/site/embeddedgo/)
