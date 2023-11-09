FROM centos:7
COPY dms-devops config.yml  /opt/apps/
COPY template /opt/apps/template
WORKDIR /opt/apps
ENTRYPOINT [ "sh", "-c", "./dms-devops" ]