# doubleopenconfigurer

This will automatically configure meta-doubleopen to the project

```console
$ make
all                            performs clean build and test
build                          Generate the windows and linux builds for sep
test                           downloads the meta-doubleopen git repo and configures 2 files in test folder workdir/conf/ bblayers.conf and local.conf
del                            Delete the contents of test (named: proj) folder
cp                             Copy Sample conf files to test (named: proj) folder in workdir/conf/
clean                          Copy back the files to original state in test folder
git                            commits and push the changes if commit msg m is given without spaces ex m=added_files
help                           Show this help
```
