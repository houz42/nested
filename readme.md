# Nested Sets Model

>Nested set model represents nested sets (trees or hierarchies) in relational databases.
>See [wikipedia](https://en.wikipedia.org/wiki/Nested_set_model).

## Model in short



## Demo Chinese division data representation

Data collected from [中国行政区划数据](https://github.com/modood/Administrative-divisions-of-China). Initial inserting SQL in `division.sql` are generated with `build.go`:

```sh
$ cd division && go run build.go   # generates data inserting sql 
```

## Use as a dependence

1. create new table as in `createtable.sql` with your table name;
2. call `SetTableName()` in your `init()`;
3. initialize table as in `division/build.go`, or
4. call `Add...()` continually as in `TestInserting()`