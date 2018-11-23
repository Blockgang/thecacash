#FROM mysql:8

# theca: no volume please (yet)
#VOLUME /var/lib/mysql

#COPY docker-entrypoint.sh /usr/local/bin/

# theca sql files insert
#COPY *.sql /docker-entrypoint-initdb.d/

#RUN ln -s usr/local/bin/docker-entrypoint.sh /
#ENTRYPOINT ["docker-entrypoint.sh"]

#EXPOSE 3306
#CMD ["mysqld"]

FROM mysql/mysql-server:8.0

MAINTAINER Yanick Gubler <yanick.gubler@hsr.ch>

#ADD events.sql /docker-entrypoint-initdb.d
#ADD create.sql /docker-entrypoint-initdb.d

WORKDIR /root
