FROM ubuntu:latest as user
RUN useradd -u 10001 scratchuser

FROM scratch
COPY substitution /substition
COPY --from=user /etc/passwd /etc/passwd
ENTRYPOINT ["/substition"]
