FROM ruby:3.0-alpine

COPY . /app
WORKDIR /app

RUN gem install rack webrick

EXPOSE 3000

CMD ["ruby", "app.rb"]
