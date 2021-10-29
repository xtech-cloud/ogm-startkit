# *************************************
#
# OpenGM
#
# *************************************

FROM alpine:3.14

MAINTAINER XTech Cloud "xtech.cloud"

ENV container docker
ENV MSA_MODE release

EXPOSE 18899

ADD bin/ogm-startkit /usr/local/bin/
RUN chmod +x /usr/local/bin/ogm-startkit

CMD ["/usr/local/bin/ogm-startkit"]
