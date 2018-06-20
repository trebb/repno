#! /usr/bin/env bash

target_dir=../messages
index_file="$target_dir/message-index.txt"
ui_file="$target_dir/messages-ui.txt"
mb_file="$target_dir/messages-mb.txt"
interleaved_file="$target_dir/messages-interleaved.txt"
known_index_file="$target_dir/known-message-index.txt"

mkdir -p "$target_dir"

sort -s -k1.1,1.22 *_ui.txt > "$ui_file"
sort -s -k1.1,1.22 *_mb.txt > "$mb_file"
sort -s -k1.1,1.22 "$ui_file" "$mb_file" > "$interleaved_file"

for i in ui mb; do
    echo "############## $i $i $i $i $i $i $i $i $i"
    cat *_$i.txt | grep -e ">>" | sed -e 's/55 AA  0 /\
55 AA  0 /g' | grep -ve ">>" | sort | uniq
    done > "$index_file"

< ../message-dictionary.txt awk \
  '/(UI|MB)\(/,/\)/{sub(/UI\(/, ""); sub(/MB\(/, ""); sub(/\)/, ""); sub(/^[ ]*/, ""); print}' |
    sort | uniq > "$known_index_file"

< "$index_file" awk \
  '$1 !~ /^####/ {freq[NF] = freq[NF] + 1}; END{for (i in freq) printf "%4d x len = %d\n", freq[i], i}'
