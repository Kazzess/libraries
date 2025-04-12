# Common libraries

## How to add new library

1. Create folder ./<name>
2. write down a lot of code in ./<name>
3. cd ./<name>
4. go mod init github.com/Kazzess/libraries/<name>
5. go mod tidy
6. add to Makefile to variable `NAMES` your new library name
