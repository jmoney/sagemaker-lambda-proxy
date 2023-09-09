FROM public.ecr.aws/lambda/provided:al2
COPY main main
ENTRYPOINT [ "./main" ]