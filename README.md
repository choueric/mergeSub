# mergeSub

Merge srt subtitle files.

# Usage

```sh
mergeSub -i cd1.srt;cd2.srt -o output.srt -t 00:59:00,300
```

* -i: input srt files. Use `;` to divide them.
* -o: output merged srt file.
* -t: offset between srt files. The format is same as the timecode in srt.

# TODO
- one SRT item may have empty line.
- detect dos and unix formats and add option to specify one format in the output.
