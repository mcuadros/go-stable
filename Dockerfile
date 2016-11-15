FROM busybox 
MAINTAINER Máximo Cuadros <mcuadros@gmail.com>

ADD cli/stable/stable /usr/local/bin/
ENTRYPOINT ["stable"] 