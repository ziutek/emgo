### Emgo

To build and install egc: 

    cd egc
    go install
  
To build some example:

    cd egpath/src/stm32/examples/nucleo-f411re/rtc
    ../build.sh

To load binary build for running from SRAM:

	../load.sh      # This uses st-util

or

	../load-oocd.sh # This uses openocd

To load binary build for running from flash (this erases flash and program it with new firmware):

	../load.sh flash

or

	../load-oocd.sh flash


[Home page](https://sites.google.com/site/embeddedgo/)
