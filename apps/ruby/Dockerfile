FROM ruby:3.2-alpine

COPY . /app
WORKDIR /app

RUN gem install rack rackup

EXPOSE 3000

CMD ["rackup", "--port", "3000", "--host", "0.0.0.0"]
