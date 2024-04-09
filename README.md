# report-sitory
`report-sitory` is a tool which provides statistical report for a git repository.  
The report includes such numbers as lines of code, changed files and commits per contributor  

---
## Example
```
✗ report-sitory --repository=. --extensions='.go,.md' --order-by=lines
Name                   Lines Commits Files
Joe Tsai               12154 92      49
colinnewell            130   1       1
Roger Peppe            59    1       2
A. Ishikawa            36    1       1
Tobias Klauser         33    1       2
178inaba               11    2       4
Kyle Lemons            11    1       1
Dmitri Shuralyov       8     1       2
ferhat elmas           7     1       4
```

---
## Usage

Build the tool:
```
(cd ./cmd/report-sitory && go build .)
```
Write it to `GOPATH/bin`:
```
go install ./cmd/report-sitory/...
```

Add `GOPATH/bin` to `PATH` to use the tool from any path:
```
export PATH=$GOPATH/bin:$PATH
```  

Now simply call 
```
report-sitory --repository *path_to_repo*
```  
By default the tool tries to calculate stats in current directory  

---
## Interface  
`report-sitory` writes the output to stdout.  
The progress of the calculations is written to stderr  
Output format:  
`tabular`:
```
Name         Lines Commits Files
Joe Tsai     64    3       2
Ross Light   2     1       1
ferhat elmas 1     1       1
```

`csv`:
```
Name,Lines,Commits,Files
Joe Tsai,64,3,2
Ross Light,2,1,1
ferhat elmas,1,1,1
```

`json`:
```
[{"name":"Joe Tsai","lines":64,"commits":3,"files":2},{"name":"Ross Light","lines":2,"commits":1,"files":1},{"name":"ferhat elmas","lines":1,"commits":1,"files":1}]
```

`json-lines`:
```
{"name":"Joe Tsai","lines":64,"commits":3,"files":2}
{"name":"Ross Light","lines":2,"commits":1,"files":1}
{"name":"ferhat elmas","lines":1,"commits":1,"files":1}
```

---
## Stats

* Number of code lines
* Number of commits
* Number of changed files

All statistics are calculated for a specific commit  

---
## Flags

**--repository** — path to git repository  
**--revision** — commit hash, default is HEAD  
**--order-by** — the key for ordering the results. May be "lines" (default), "commits", "files"  
**--use-committer** — (bool) use committer instead of author in calculations  
**--format** — output format, one of "tabular" (default), "csv", "json", "json-lines"  
**--extensions** — set of file extensions narrowing down the list of files in calculations  
**--languages** — set of programming languages narrowing down the list of files in calculations  
**--exclude** — set of Glob-patterns which excludes all the files with a matching name from calculations  
**--restrict-to** — set of Glob-patterns which excludes all the files with a name not matching any of the patterns from it  