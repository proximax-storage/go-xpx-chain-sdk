## Generate flatbuffers transactions

`flatc --go schema_aggregate_transaction.fbs`

Docs: http://google.github.io/flatbuffers/

You can install `flatc` manually

```
git clone https://github.com/google/flatbuffers.git
cd flatbuffers
cmake -G "Unix Makefiles"
make
cd -
```
You can use syntax like:

```$xslt
./flatbuffers/flatc --go *.fbs
```


