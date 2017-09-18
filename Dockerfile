FROM cassandra:3.0.14

ADD ./bin /bin
COPY CaOps.yaml /etc/CaOps.yaml

ENV JOLOKIA_VERSION 1.3.7
ENV JOLOKIA_JVM_AGENT /usr/share/java/jolokia-jvm-agent.jar
ENV JOLOKIA_PORT 8778
ENV CAOPS_API_PORT 8080
ENV CAOPS_GOSSIP_PORT 7942

ENV CASSANDRA_LISTEN_ADDRESS auto
ENV CASSANDRA_CLUSTER_NAME CaOpsDev

RUN apt-get update
RUN apt-get install -y --no-install-recommends ca-certificates wget supervisor
RUN wget -O $JOLOKIA_JVM_AGENT "https://repo1.maven.org/maven2/org/jolokia/jolokia-jvm/$JOLOKIA_VERSION/jolokia-jvm-$JOLOKIA_VERSION-agent.jar"
RUN echo "JVM_OPTS=\"\$JVM_OPTS -javaagent:$JOLOKIA_JVM_AGENT=port=$JOLOKIA_PORT,host=0.0.0.0,discoveryEnabled=false\"" >> /etc/cassandra/cassandra-env.sh
RUN sed -ri 's/^listen_address:/#listen_address:/g' /etc/cassandra/cassandra.yaml
RUN sed -ri 's/^# listen_interface:/listen_interface:/g' /etc/cassandra/cassandra.yaml

RUN mkdir -p /var/log/supervisor

ENV SUPERVISOR_CASSANDRA /etc/supervisor/conf.d/cassandra.conf
RUN echo "[program:cassandra]" > $SUPERVISOR_CASSANDRA ; \
	echo "command = /usr/sbin/cassandra -f" >> $SUPERVISOR_CASSANDRA ;\
	echo "user = cassandra" >> $SUPERVISOR_CASSANDRA

ENV SUPERVISOR_CAOPS /etc/supervisor/conf.d/caops.conf
RUN echo "[program:caops]" > $SUPERVISOR_CAOPS ; \
	echo "command = /bin/CaOps --config /etc/CaOps.yaml serve" >> $SUPERVISOR_CAOPS ;\
	echo "user = cassandra" >> $SUPERVISOR_CAOPS

EXPOSE $JOLOKIA_PORT $CAOPS_API_PORT $CAOPS_GOSSIP_PORT

CMD ["/usr/bin/supervisord", "-n", "-c", "/etc/supervisor/supervisord.conf"]