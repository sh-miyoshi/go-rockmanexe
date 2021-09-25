FROM golang:1.16 as router-builder

WORKDIR /app
COPY go.mod .
COPY go.sum .
COPY cmd/router cmd/router
COPY pkg pkg
RUN go build ./cmd/router

FROM ruby:2.7.4 as app-builder

ARG SECRET_KEY_BASE
ENV RAILS_ENV production
WORKDIR /app/cmd/matcher
COPY cmd/matcher/Gemfile Gemfile
COPY cmd/matcher/Gemfile.lock Gemfile.lock
RUN gem install bundler
RUN bundle config set --local disable_checksum_validation true
RUN bundle install
COPY cmd/matcher cmd/matcher
RUN rails assets:precompile
RUN rails db:migrate

FROM ruby:2.7.4

ARG SECRET_KEY_BASE
ENV RAILS_ENV production
WORKDIR /app
COPY --from=router-builder /app/router /app/router/router
COPY cmd/router/config.yaml router/config.yaml
COPY --from=app-builder /app/cmd/matcher/public matcher/public
COPY --from=app-builder /app/cmd/matcher/db/production.sqlite3 matcher/db/production.sqlite3
WORKDIR /app/router
RUN ./router --config=config.yaml &
WORKDIR /app/matcher
COPY cmd/matcher/Gemfile Gemfile
COPY cmd/matcher/Gemfile.lock Gemfile.lock
RUN gem install bundler
RUN bundle config set --local disable_checksum_validation true
RUN bundle install
COPY cmd/matcher matcher
RUN rails server -e production
