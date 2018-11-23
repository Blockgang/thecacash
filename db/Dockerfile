#FROM mysql/mysql-server:8.0
# ==> ERROR 2059 (HY000): Authentication plugin 'caching_sha2_password' cannot be loaded: /usr/lib/mysql/plugin/caching_sha2_password.so: cannot open shared object file: No such file or directory

FROM mysql/mysql-server:5.7

ADD events.sql /docker-entrypoint-initdb.d
ADD create.sql /docker-entrypoint-initdb.d
