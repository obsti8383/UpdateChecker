# UpdateChecker Overview
Checks for software release states on Windows systems.

Only finds software installed via regular Windows Installer (no "portable" software)

UpdateChecker is using https://vergrabber.kingu.pl/ to fetch the current version of installed software and uses the Bootstrap framework (https://getbootstrap.com/) to show the results as a static webpage in your standard browser.


Currently supports:
* Windows 10
* Mozilla Firefox
* Google Chrome
* OpenVPN
* Adobe Flash Player
* Adobe Acrobat Reader
* 7-Zip
* TeamViewer
* Mozilla Thunderbird
* VeraCrypt
* Java 8
* LibreOffice (in preparation)


UpdateChecker ist Open Source (GPL 3.0), doesn't track you and is ad-free.

# Installation
Just unzip the provided ZIP-File to a location that fits your needs.
The UpdateChecker directory contains the following files and directories:
* UpdateChecker.exe: The executable
* /static: Contains the html, css and js files required from the Bootstrap framework (https://getbootstrap.com/) to show the results as static webpage
* main.html: The webpage template that is used for generating the results webpage

Only there after first start of UpdateChecker.exe:
* updatechecker_result.html:  This is the webpage with the results based on the main.html template; used for display in the standard browser
* vergrabber.json: This is a json file that contains the current version of the software packages. Will be updated when UpdateChecker.exe is started, but only once a day (if started more than once a day this cached version is used)
* UpdateChecker.log: Log output, check for errors if something doesn't work as expected or no Webpage is opened in your browser

# Usage
Just start UpdateChecker.exe (you don't need administrative rights) and wait a second. UpdateChecker fetches the current versions from https://vergrabber.kingu.pl/ and thereafter verifies your installed software and shows the results in your standard web browser.

![ResultsScreenshot](./graphics/result1.jpg)

The column "Status" shows you the state of the listed software installation:
* Outdated: Updates are available. It is recommended to install the recent version to be sure you have the latest security patches applied.
* Up-to-date: The most recent version is installed.
* Unkown status: This software is not known to UpdateChecker
