FROM gcr.io/distroless/static:nonroot
WORKDIR /app

COPY gha_register_build_actions_app /app

ENTRYPOINT ["/app/gha_register_build_actions_app"]