FROM scratch

VOLUME  /flogo/flogo-cli
CMD ["/bin/true"]
COPY cli/ /flogo/flogo-cli/cli
COPY flogo/ /flogo/flogo-cli/flogo
COPY tools/ /flogo/flogo-cli/tools
COPY util/ /flogo/flogo-cli/util
COPY README.md /flogo/flogo-cli/
