# mergeSub

Merge srt subtitle files.

# Usage

```sh
mergeSub -i cd1.srt;cd2.srt -o output.srt -t 00:59:00,300
```

* -i: input srt files. Use `;` to divide them.
* -o: output merged srt file.
* -t: offset between srt files. The format is same as the timecode in srt.
* -f: specify output file format, unix or dos.
* -d: verbose debug output.

# TODO
- [x] Text section of a SRT item may have empty line.
- [x] detect dos and unix formats and add option to specify the format (dos or unix)
  for the output file.
