# HashKitty

[Tests](https://github.com/vegasq/hashkitty/actions/workflows/go.yml/badge.svg?branch=master)

Reimplementation of _some_ of the HashCat (https://github.com/hashcat/hashcat) features in GO.

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

## Why it even exists?

`Hashcat` ATM is lacking support for `argon2` hashes, so this is a simple dropin replacement.
