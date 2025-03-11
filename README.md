# file-dedup

[日本語](/README-ja.md)

A macOS tool for detecting and suggesting duplicate files for deletion.

## Overview

This tool takes a CSV file containing SHA256 hash values as input, detects duplicate files, and suggests candidates for deletion.
Based on filename length, it keeps the file with the shortest name and suggests other files for deletion.

## Usage

First, run the following command to generate hash values for each file and create a CSV file.
Note: We use Pipe Viewer to display execution time.
You can install it using Homebrew:
```bash
brew install pv
```

```bash
find . -type f -print0 | tee >(pv -l -s $(find . -type f | wc -l) > /dev/null) | xargs -0 sha256sum | awk -F'  ' 'BEGIN {OFS=","} {gsub(/"/, "\"\"", $2); print "\"" $2 "\",\"" $1 "\""}' > hashes.csv
```

```bash
./file-dedup -csv <path_to_csv_file> [-out <output_filename>] [-debug]
```

After running this script, you can delete files using terminal commands.
Before deletion, please verify the files to be deleted using the -debug option of this script.
```bash
pv -l duplicates.txt | while read file; do
    rm "$file"
    echo "Deleted: $file"
done
```

### Options

- `-csv`: Required. Path to the CSV file for duplicate checking
- `-out`: Optional. Name of the output file (default: duplicates.txt)
- `-debug`: Optional. Detailed output mode

### Input CSV File Format

The CSV file must be in the following format:
- No header
- Column 1: Filename
- Column 2: SHA256 hash value

### Output Format

#### Normal Mode
Only outputs the file paths of deletion candidates.

#### Debug Mode
Outputs in the following format:
```
File to keep: [file_path]
Deletion candidates:
  [file_path1]
  [file_path2]
---
```

## How it works

1. Groups files with identical hash values from the input CSV file
2. Selects files within each group based on the following priorities:
   - Prioritizes manually named files (e.g., `vacation2023.jpg`)
   - Auto-generated camera filenames (starting with `DSC`, `IMG`) are lower priority (can be customized in `createFileGroup()` > `isAutoGenerated()` > `patterns`)
   - Within the same category, selects the file with the shortest filename
3. Keeps the selected file and outputs others as deletion candidates