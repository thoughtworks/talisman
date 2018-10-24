#/bin/bash
docker image build --rm -t build_talisman .
docker container run --rm -v `pwd`:/go/src/github.com/thoughtworks/talisman -it --name build_talisman build_talisman
# need to find a good way to make sure that the binaries and the checksum file are not owned by root
# sudo chown $USER:$USER talisman_{l,d,w}* 
# sudo chown $USER:$USER checksums
