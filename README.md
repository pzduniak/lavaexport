# lavaexport

Tool for backing up contents of your Lavaboom account. In order to create a
copy of your account's contents, all you have to do is run the tool, fill in
the requested data and wait until the process is complete.

**Important:** This program does not use pagination, therefore I recommend
running it on a computer that has at least twice as much memory as there is
data on your Lavaboom account. Also, because of that, it might seem that the
progress freezes when the program switches to a new object type. Don't worry - 
it's fetching the data from Lavaboom's API in the background.

## Installation

Either download the binaries provided in the "Releases" tab on GitHub or run
the following command to build it from source:

```bash
go get github.com/pzduniak/lavaexport
```
