FROM ubuntu:18.04

# a lib (not sure which?) needed for api request
RUN apt update && apt install -y curl 
COPY blockchain_sync /tmp/
CMD ["/tmp/blockchain_sync"]
#CMD ["bash"]
