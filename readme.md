# Nested Set Model

>Nested Set Model (or Modified Preorder Tree) represents nested sets (trees or hierarchies) in relational databases.
>See [wikipedia](https://en.wikipedia.org/wiki/Nested_set_model).

## Model in short

The nested model:
1. assigns two numbers for each node as left and right, that
2. left of each node is less than all its children's left, and
3. right of each node is greater than all its children's right.

Numbers could be assigned according to a preorder tree traversal, which visits each node twice, and assigns numbers at both visits.

- Then querying all descendants of a node could be efficiency as:
```sql
SELECT child.id, child.node, child.lft, child.rgt
    FROM `nested` parent, `nested` child
    WHERE child.lft BETWEEN parent.left AND parent.right 
    AND parent.id=@node_id
```

- And querying the path from root to any node could be:
```sql
SELECT parent.id, parent.node, parent.lft, parent.rgt
    FROM `nested` parent, `nested` child
    WHERE child.lft BETWEEN parent.lft AND parent.rgt
    AND child.id=@node_id
    ORDER BY parent.lft
```

- To query immediate children of a node efficiently, another column `depth` is used to record depth of the node, and the querying could be:
```sql
SELECT child.id, child.node, child.lft, child.rgt
    FROM `nested` parent, `nested` child
    WHERE child.lft BETWEEN parent.lft AND parent.rgt
    AND child.depth=parent.depth+1
```

The nested model is also suitable for trees with more than one root, forests.

### Test clothing categories
Test data used in `nested_test.go` are collected from previous [wikipedia article](https://en.wikipedia.org/wiki/Nested_set_model).

```
- Clothing
    - Men's
        - Suits
            - Slacks
            - Jackets
    - Women's
        - Dresses
            - Evening Gowns
            - Sun Dresses
        - Skirts
        - Blouses
```

Traversal and number nodes as:

![traversing](https://upload.wikimedia.org/wikipedia/commons/thumb/b/b5/Clothing-hierarchy-traversal-2.svg/523px-Clothing-hierarchy-traversal-2.svg.png)

## Demo Chinese division data representation

Store Chinese division data with nested sets:

1. build a division tree from raw data,
2. assign left and right number for divisions by preorder tree traversal,
3. generate sql inserting queries

Data collected from [中国行政区划数据](https://github.com/modood/Administrative-divisions-of-China). Initial inserting SQL in `division.sql` are generated with `build.go`:

```sh
$ cd division && go run build.go   # generates data inserting sql 
```

## Use as dependency

1. create new table as in `createtable.sql` with your table name;
2. initialize table as in `division/build.go`, or
3. call `Add...()` continually as in `TestInserting()`;
4. call `SetTableName()` in your `init()`;
