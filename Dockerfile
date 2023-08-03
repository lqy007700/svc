FROM alpine
ADD svc /svc

ENV KUBECONFIG=/root/.kube/config

ENTRYPOINT [ "/svc" ]