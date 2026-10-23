FROM python:3.10

RUN apt-get update && apt-get install -y --only-upgrade ca-certificates && update-ca-certificates -f

RUN apt-get update && apt-get install -y libsodium-dev tree
RUN pip install -U 'certifi>=2026.10.19' 'pip>=26.2.2' && pip install 'securesystemslib[crypto,pynacl]==0.25.0' 'tuf==v2.0.0'

ADD generate.py generate.sh /
CMD /generate.sh
