EXEC = mergeSub

all:
	@go build -o $(EXEC) -v

clean:
	@rm -rfv $(EXEC)

run:
	./$(EXEC) -i "sub/University.Of.Laughs.2004.DVDRip.DivX.CD1-AXiNE.chs.srt;sub/University.Of.Laughs.2004.DVDRip.DivX.CD2-AXiNE.chs.srt" -o output.srt -t 00:59:00,300
