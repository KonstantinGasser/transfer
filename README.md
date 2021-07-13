# USB transfer


## Getting Started
1. Download the ZIP-File from `https://github.com/KonstantinGasser/transfer`
2. in there you will find two files per operating system (windows and MacOs). The files with the suffix of `-concurrent` process all plugged in USB-Sticks concurrently (these versions can be used to speed up things when using a big USB-Hub - not required for only 2 - 10 sticks unless wanted)


## On Windows
2. Unzip the file and move the `transfer-windows.exe` to a new folder on your desktop.
3. put the target/source data into the same new folder the executable now is.
4. open the `powerShell` from the search and navigate to the new folder on the desktop. Use `cd path\to\Desktop\folder-name`
5. run the executable with `./transfer-windows.exe`. You will see text appearing in the `powerShell` telling you what is going on.

## On Mac
2. Unzip the file and move the `transfer-darwin` to a new folder on your desktop.
3. put the target/source data into the same new folder the executable now is.
4. open the `terminal` from the spotlight and navigate to the new folder on the desktop. Use `cd ~/Desktop/folder-name`
5. run the executable with `./transfer-darwin`. You will see text appearing in the `terminal` telling you what is going on.


## Process
1. before you run the program make sure there are no USB-Sticks plugged in, so the program can initialize properly
2. after the program has initialized itself you will see a text, and you can start inserting USB-Sticks.
3. the program will listen for new USB-Sticks and informs you if a new one has been detected. 
4. once a USB-Stick has is detected, the target data will be copied to the USB-Stick and once finished again the program will tell you which USB-Stick is done.
5. the programs default for concurrent USB-Sticks is set to two - which means you can process two USB-Sticks at the same time. to change this, run the program with
`./transfer-darwin -concurrent <the-number-of-usb-sticks-you-want-to-process>`

## Terminate program
Hit keys `CTR-C` (in some cases this can take up to 3 seconds but you will be notified once the progam has stopped)
