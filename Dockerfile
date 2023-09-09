FROM public.ecr.aws/lambda/provided:al2
COPY sagemaker-lambda-proxy sagemaker-lambda-proxy
ENTRYPOINT [ "./sagemaker-lambda-proxy" ]