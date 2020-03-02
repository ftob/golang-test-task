FROM grafana/grafana:latest

USER root

COPY ./etc/grafana/datasources /etc/grafana/provisioning/datasources
COPY ./etc/grafana/dashboards /etc/grafana/provisioning/dashboards

EXPOSE 3000