
### Index Difference from MySQL

- The CREATE INDEX supports the UNIQUE index and does not support FULLTEXT and SPATIAL indexes.

- The index_col_name supports the length option with a maximum length limit of 3072 bytes. The
length limit does not change depending on the storage engine, and character set used when building the table.
This is because TiDB does not use storage engines like InnoDB and MyISAM, and only provides
syntax compatibility with MySQL for the storage engine options when creating tables. Similarly, TiDB
uses the utf8mb4 character set, and only provides syntax compatibility with MySQL for the character
set options when creating tables. For more information, see Compatibility with MySQL.

- The index_col_name supports the index sorting options of ASC and DESC. The behavior of sorting
options is similar to MySQL, and only syntax parsing is supported. All the internal indexes are stored
in ascending order ascending order. For more information, see CREATE INDEX Syntax.

- The index_option supports KEY_BLOCK_SIZE, index_type and COMMENT. The COMMENT supports
a maximum of 1024 characters and does not support the WITH PARSER option.

- The index_type supports BTREE and HASH only in MySQL syntax, which means the index type is
independent of the storage engine option in the creating table statement. For example, in MySQL,
when you use CREATE INDEX on a table using InnoDB, it only supports the BTREE index, while TiDB
supports both BTREE and HASH indexes.

- The CREATE INDEX does not support the algorithm_option and lock_option in MySQL.

- TiDB supports at most 512 columns in a single table. The corresponding number limit in InnoDB is
1017, and the hard limit in MySQL is 4096. For more details, see Limits on Table Column Count and
Row Size.

### Features that are different from MySQL

- Auto-increment ID：The auto-increment ID feature in TiDB is only guaranteed to be automatically incremental and unique but
is not guaranteed to be allocated sequentially. Currently, TiDB is allocating IDs in batches. If data is inserted
into multiple TiDB servers simultaneously, the allocated IDs are not sequential.

- Built-in functions：TiDB supports most of the MySQL built-in functions, but not all. See TiDB SQL Grammar for the supported
functions.

- 