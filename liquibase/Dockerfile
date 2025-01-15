FROM liquibase/liquibase
RUN lpm add mysql --global

COPY ./liquibase/liquibase.docker.properties /liquibase/changelog/liquibase.docker.properties
COPY ./liquibase/ticketing-system/sql /liquibase/changelog

ENTRYPOINT [ "sh", "-c", "sleep 10; liquibase --defaults-file=/liquibase/changelog/liquibase.docker.properties update" ]