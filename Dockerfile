FROM golang:1.17.6 as builder

WORKDIR /go/src/github.com/gavinmcnair/javatruststoreinit

COPY . .

RUN make javatruststoreinit

FROM gcr.io/distroless/static-debian11
MAINTAINER Gavin McNair

ARG git_repository="Unknown"
ARG git_commit="Unknown"
ARG git_branch="Unknown"
ARG built_on="Unknown"

LABEL git.repository=$git_repository
LABEL git.commit=$git_commit
LABEL git.branch=$git_branch
LABEL build.on=$built_on

COPY --from=builder /go/src/github.com/gavinmcnair/javatruststoreinit/bin/linux/javatruststoreinit .

CMD [ "/javatruststoreinit" ]
