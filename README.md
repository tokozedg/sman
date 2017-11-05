# Sman
[![Build Status](https://travis-ci.org/tokozedg/sman.svg?branch=master)](https://travis-ci.org/tokozedg/sman)

***
A command-line snippet manager in Go

[![asciicast](https://asciinema.org/a/2e04fxybyyo5ubjk42mk4yuop.png)](https://asciinema.org/a/2e04fxybyyo5ubjk42mk4yuop)

## Install

```shell
bash -c "$(curl https://raw.githubusercontent.com/tokozedg/sman/master/install.sh)"
```

## Building

* Build with Go
```bash
go get -v github.com/tokozedg/sman
```

* Add to your rc:
```bash
[ -f $GOPATH/src/github.com/tokozedg/sman/sman.rc  ] && source $GOPATH/src/github.com/tokozedg/sman/sman.rc
```

* Optionally copy demo snippets dir or create yours:
```bash
cp -r $GOPATH/src/github.com/tokozedg/sman/snippets ~/
```

## Snippets Examples

```yaml
#~/snippets/shell.yml

smtp:server: # snippet name
  do: exec # copy or exec
  desc: smtp server in python. Prints mails to stdout
  command: python -m smtpd -n -c DebuggingServer localhost:1025

```

```yaml
#~/snippets/shell.yml
# Sman will ask input for placeholder <<port>>
tcpdump:port:
  do: copy
  desc: listen traffic on port
  command: tcpdump -nqt -s 0 -A -i eth0 port <<port>>
```

```yaml
#~/snippets/shell.yml
# To execute multiline commands, separate lines by semi-colon
curl:upload:
  do: exec
  command: >
    gpg -c <<file>>;
    curl --upload-file <<file>>.gpg https://transfer.sh/<<file>>.gpg
```

* You can export command to a separate file located at: `SMAN_SNIPPET_DIR/<Snippet File Name>/<Snippet Name>`



## Placeholders

```
<<name(option1,option2)#Description>>
```
* Include placeholder anywhere within snippet command
* Name is the only mandatory field
* You can have multiple placeholders with the same name. After input all of them will be replaced
* Use `\` to escape comma in options

## Usage Examples

### Run snippet

```bash
s run [-f <FILE>]  [-t <TAG>] <SNIPPET> [placeholder values...] [-cxyp]
```
```bash
~|⇒ s run -f shell curl:upload test.tar.gz -x
----
gpg -c test.tar.gz; curl --upload-file test.tar.gz.gpg https://transfer.sh/test.tar.gz.gpg
----
Execute Snippet? [Y/n]:
```

```bash
~|⇒ s run curl:ip
----
curl canhazip.com
----
Execute Snippet? [Y/n]:
```

### Show snippet

```bash
s show [-f <FILE>] [-t <TAG>] <SNIPPET>
```

### List and search snippets
```bash
s ls [-f <FILE>] [-t <TAG>] [<PATTERN>]
```

* Pattern is matched against snippet name, command and description


## Fuzzy search file and snippet name:
```bash
# `r` is alias for `run`
# matches file `mysql` and snippet `database:dump`

~|⇒ s r -f sql dmp
----
mysqldump -u[user] --lock-tables=[lock] -p[pass] -h [host] [database] > [database].sql
----
[user]:
```

## Config
```bash
# Append history can be useful to avoid re-entering all placeholders when you need to change single parameter.
export SMAN_APPEND_HISTORY=false
# Snippet directory
export SMAN_SNIPPET_DIR="~/snippets"
# Ask confirmation before executing
export SMAN_EXEC_CONFIRM=true
```

## vim-sman

Install vim plugin for better snippets colors:

*  [Pathogen](https://github.com/tpope/vim-pathogen)
  * `git clone https://github.com/tokozedg/vim-sman.git ~/.vim/bundle/vim-sman`
*  [vim-plug](https://github.com/junegunn/vim-plug)
  * `Plug 'tokozedg/vim-sman'`
*  [NeoBundle](https://github.com/Shougo/neobundle.vim)
  * `NeoBundle 'tokozedg/vim-sman'`
*  [Vundle](https://github.com/VundleVim/Vundle.vim)
  * `Plugin 'tokozedg/vim-sman'`

## Contributing

If you'd like to contribute, please fork the repository and make changes as
you'd like. Pull requests are warmly welcome, especially if you make a good snippet file.

