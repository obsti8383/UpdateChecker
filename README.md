# UpdateChecker
Checks for software release states on Windows systems using the online service http://vergrabber.kingu.pl/ for current version information

## Currently supported software
* Windows itself (only for Windows 10)
* 7-Zip
* Adobe Flash Player
* Adobe Acrobat Reader
* Google Chrome
* OpenVPN
* Mozilla Firefox
* Mozilla Thunderbird
* Java (8 only, due to vergrabber restrictions)
* TeamViewer
* VeraCrypt

## Informations
UpdateChecker downloads a JSON file from http://vergrabber.kingu.pl/ (vergrabber.json) and caches it for one day. This JSON file contains information about the current versions for Windows and several client and server programs. UpdateChecker parses this file, reads in all installed programs in Windows (via Registry keys) and compares the installed software base against the parsed results.
UpdateChecker generates a static HTML file that is saved to the filesystem (in the directory from which UpdateChecker is executed) and opens it in the default browser.

## Limitations
* this is a alpha version. Do not use for productive purposes. Do not rely on the output.
* there is no integrity check for the vergrabber.json file and it is transfered via http (which in combination is inherently insecure)
* only supports a limited set of software right now
* only detects software that is installed in Windows. Portable software is not supported.
* Reloading the webpage doesn't update the results, since this is only a static webpage generated once during execution of the tool (which is terminated after starting the browser)
