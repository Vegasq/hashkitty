# HashKitty

![Unit Tests](https://github.com/vegasq/hashkitty/actions/workflows/unit.yml/badge.svg?branch=master)
![Functional Tests](https://github.com/vegasq/hashkitty/actions/workflows/functional.yml/badge.svg?branch=master)

Reimplementation of _some_ of the [HashCat](https://github.com/hashcat/hashcat) features in GO.

## What works

1. Attack modes:
   1. 0 - wordlist attack
   2. 9 - association attack
2. Rules supported to some degree (see hashkitty/rules)
3. Handful of algos (see hashkitty/algos)


## What does not work

1. Everything


## Usage

```bash
hashkitty -a 0 -m 99001 leftlist.txt wordlist.txt  --potfile-path potfile.txt --rules-file rules.txt
```

```
usage: hashkitty [LEFTLIST] [WORDLIST] [--attack-mode ATTACK-MODE] [--hash-type HASH-TYPE] [--rules-file RULES-FILE] [--potfile-path POTFILE-PATH]

positional arguments:
  LEFTLIST                                      Leftlist file location
  WORDLIST                                      Wordlist file location

flags:
  --attack-mode ATTACK-MODE, -a ATTACK-MODE     Attack Mode
  --hash-type HASH-TYPE, -m HASH-TYPE           Hash Type
  --rules-file RULES-FILE, -r RULES-FILE        Multiple rules applied to each word from wordlists
  --potfile-path POTFILE-PATH, -p POTFILE-PATH  Specific path to potfile
  --hex-salt                                    Assume salt is given in hex
  --remove                                      Enable removal of hashes once they are cracked
```

### Supported rules

|        Name         |  Function  |                      Description                      |  Example Rule  |  Input Word  |         Output Word        |  Note  |
|:-------------------:|:----------:|:-----------------------------------------------------:|:--------------:|:------------:|:--------------------------:|:------:|
|  Nothing            |  :         |  Do nothing (passthrough)                             |  :             |  p@ssW0rd    |  p@ssW0rd                  |        |
|  Lowercase          |  l         |  Lowercase all letters                                |  l             |  p@ssW0rd    |  p@ssw0rd                  |        |
|  Uppercase          |  u         |  Uppercase all letters                                |  u             |  p@ssW0rd    |  P@SSW0RD                  |        |
|  Capitalize         |  c         |  Capitalize the first letter and lower the rest       |  c             |  p@ssW0rd    |  P@ssw0rd                  |        |
|  Invert Capitalize  |  C         |  Lowercase first found character, uppercase the rest  |  C             |  p@ssW0rd    |  p@SSW0RD                  |        |
|  Toggle Case        |  t         |  Toggle the case of all characters in word.           |  t             |  p@ssW0rd    |  P@SSw0RD                  |        |
|  Toggle @           |  TN        |  Toggle the case of characters at position N          |  T3            |  p@ssW0rd    |  p@sSW0rd                  |  *     |
|  Reverse            |  r         |  Reverse the entire word                              |  r             |  p@ssW0rd    |  dr0Wss@p                  |        |
|  Duplicate          |  d         |  Duplicate entire word                                |  d             |  p@ssW0rd    |  p@ssW0rdp@ssW0rd          |        |
|  Duplicate N        |  pN        |  Append duplicated word N times                       |  p2            |  p@ssW0rd    |  p@ssW0rdp@ssW0rdp@ssW0rd  |        |
|  Reflect            |  f         |  Duplicate word reversed                              |  f             |  p@ssW0rd    |  p@ssW0rddr0Wss@p          |        |
|  Rotate Left        |  {         |  Rotate the word left.                                |  {             |  p@ssW0rd    |  @ssW0rdp                  |        |
|  Rotate Right       |  }         |  Rotate the word right                                |  }             |  p@ssW0rd    |  dp@ssW0r                  |        |
|  Append Character   |  $X        |  Append character X to end                            |  $1            |  p@ssW0rd    |  p@ssW0rd1                 |        |
|  Prepend Character  |  ^X        |  Prepend character X to front                         |  ^1            |  p@ssW0rd    |  1p@ssW0rd                 |        |
|  Truncate left      |  [         |  Delete first character                               |  [             |  p@ssW0rd    |  @ssW0rd                   |        |
|  Truncate right     |  ]         |  Delete last character                                |  ]             |  p@ssW0rd    |  p@ssW0r                   |        |
|  Delete @ N         |  DN        |  Delete character at position N                       |  D3            |  p@ssW0rd    |  p@sW0rd                   |  *     |
|  Extract range      |  xNM       |  Extract M characters, starting at position N         |  x04           |  p@ssW0rd    |  p@ss                      |  *     |
|  Omit range         |  ONM       |  Delete M characters, starting at position N          |  O12           |  p@ssW0rd    |  psW0rd                    |  *     |
|  Insert @ N         |  iNX       |  Insert character X at position N                     |  i4!           |  p@ssW0rd    |  p@ss!W0rd                 |  *     |
|  Overwrite @ N      |  oNX       |  Overwrite character at position N with X             |  o3$           |  p@ssW0rd    |  p@s$W0rd                  |  *     |
|  Truncate @ N       |  'N        |  Truncate word at position N                          |  '6            |  p@ssW0rd    |  p@ssW0                    |  *     |
|  Replace            |  sXY       |  Replace all instances of X with Y                    |  ss$           |  p@ssW0rd    |  p@$$W0rd                  |        |
|  Purge              |  @X        |  Purge all instances of X                             |  @s            |  p@ssW0rd    |  p@W0rd                    |        |
|  Duplicate first N  |  zN        |  Duplicate first character N times                    |  z2            |  p@ssW0rd    |  ppp@ssW0rd                |        |
|  Duplicate last N   |  ZN        |  Duplicate last character N times                     |  Z2            |  p@ssW0rd    |  p@ssW0rddd                |        |
|  Duplicate all      |  q         |  Duplicate every character                            |  q             |  p@ssW0rd    |  pp@@ssssWW00rrdd          |        |


    * Indicates that N starts at 0. For character positions other than 0-9 use A-Z (A=10)

This is a copypaste of supported rules from [Hashcat rule based attack](https://hashcat.net/wiki/doku.php?id=rule_based_attack).

## Why it even exists?

`Hashcat` ATM is lacking support for `argon2` hashes, so this is a simple dropin replacement.
