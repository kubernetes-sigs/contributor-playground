FROM alpine:3.5

ADD output/cce-cloud-controller-manager /usr/local/bin/cce-cloud-controller-manager

CMD ["cce-cloud-controller-manager"]