# pgsummary
This project builds executables that can be used to do the following:
 - write a summary report to a json file
 - compare two summary json files

NOTE:  The tool expects to get database credentials form AWS Secret Manager

The tool was build with AWS RDS postgres in mind, but probably works pretty broadly because it's simple.  


## Installation

The pgsummary executables can be manually downloaded from the githbib project [releases](https://github.com/natemarks/pgsummary/releases)

To automate the process, there is an installation script for it in [pipeline-scripts](https://github.com/natemarks/pipeline-scripts) project.

```
# use v0.0.29 of the pipeline installaer script to install version 0.0.5 of pgsummary

bash -c 'curl "https://raw.githubusercontent.com/natemarks/pipeline-scripts/v0.0.29/scripts/install_pgsummary.sh" | bash -s --  -d ~/pgsummary_install -r 0.0.5'

  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  2871  100  2871    0     0  12592      0 --:--:-- --:--:-- --:--:-- 12592
Read parameters:
- release_version: 0.0.5
- directory: /Users/nmarks/pgsummary_install
x ./
x ./pgcompare
x ./pgreport
x ./version.txt
/Users/nmarks/pgsummary_install/0.0.5
❯ tree ~/pgsummary_install
/Users/nmarks/pgsummary_install
└── 0.0.5
    ├── pgcompare
    ├── pgreport
    ├── pgsummary_0.0.5_darwin_amd64.tar.gz
    └── version.txt

```
## Usage

Run pgreport to create a summary json file. The JSON is also dumped to stdout
```shell
pgreport \
-host=my.database.fqdn \
-dbname=postgres \
-secretId=my/secret/id/path \
# These are the keys in the secret JSON document
# we extract the values and use those
-secretUsernameKey=master_username \
-secretPasswordKey=master_password
{"level":"info","test_key":"test_value","time":"2021-11-18T09:00:19-05:00","message":"setting up the AWS Secret Manager client"}
{"level":"info","test_key":"test_value","time":"2021-11-18T09:00:19-05:00","message":"getting the secret doc from AWS SM"}
{"level":"info","test_key":"test_value","time":"2021-11-18T09:00:19-05:00","message":"Validating postgres credentials"}
{
  "hostName": "my.database.fqdn",
  "port": 5432,
  "databases": {
    "postgres": {
      "tables": {},
      "extensions": [
        "plpgsql"
      ]
    },
    "zzzsoups": {
      "tables": {
        "pg_stat_statements": {
          "rowCount": 4229,
          "columns": {
            "blk_read_time": "double precision",
            "blk_write_time": "double precision",
            "calls": "bigint",
            "dbid": "oid",
            "local_blks_dirtied": "bigint",
            "local_blks_hit": "bigint",
            "local_blks_read": "bigint",
            "local_blks_written": "bigint",
            "max_exec_time": "double precision",
            "max_plan_time": "double precision",
            "mean_exec_time": "double precision",
            "mean_plan_time": "double precision",
            "min_exec_time": "double precision",
            "min_plan_time": "double precision",
            "plans": "bigint",
            "query": "text",
            "queryid": "bigint",
            "rows": "bigint",
            "shared_blks_dirtied": "bigint",
            "shared_blks_hit": "bigint",
            "shared_blks_read": "bigint",
            "shared_blks_written": "bigint",
            "stddev_exec_time": "double precision",
            "stddev_plan_time": "double precision",
            "temp_blks_read": "bigint",
            "temp_blks_written": "bigint",
            "total_exec_time": "double precision",
            "total_plan_time": "double precision",
            "userid": "oid",
            "wal_bytes": "numeric",
            "wal_fpi": "bigint",
            "wal_records": "bigint"
          }
        },
        "soup_table": {
          "rowCount": 2,
          "columns": {
            "email": "character varying",
            "username": "character varying"
          }
        }
      },
      "extensions": [
        "plpgsql",
        "pg_stat_statements"
      ]
    }
  },
  "users": [
    {
      "name": "rdstopmgr",
      "attributes": ""
    },
    {
      "name": "rdsadmin",
      "attributes": "superuser, create database"
    },
    {
      "name": "postgres",
      "attributes": "create database"
    }
  ]
}%
```