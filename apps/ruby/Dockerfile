FROM ruby:2.6-alpine

COPY . /app
WORKDIR /app

RUN gem install rack

EXPOSE 3000

CMD ["ruby", "app.rb"]
