#! /usr/bin/env bash

resultfile=all-messages.txt

for i in ui mb; do
    echo "############## $i $i $i $i $i $i $i $i $i"
    cat *_$i.txt | grep -e ">>" | sed -e 's/55 AA  0 /\
55 AA  0 /g' | grep -ve ">>" | sort | uniq
    done > "$resultfile"
