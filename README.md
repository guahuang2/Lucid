# Lucid

## To run a program
go run lucid.go yourFileName.golite
## To see the lexer output 
go run lucid.go -lex yourFileName.golite
## To see the parser ast output
go run lucid.go -ast yourFileName.golite
## To see the iloc output 
go run lucid.go -iloc test1.golite
go run lucid.go -iloc benchmarks/benchmarks/Twiddleedee/Twiddleedee.golite
## To see the iloc output
go run lucid.go -o benchmarks/benchmarks/Twiddleedee/Twiddleedee.golite

Note: make sure that you are under lucid project by running:
cd proj/lucid
