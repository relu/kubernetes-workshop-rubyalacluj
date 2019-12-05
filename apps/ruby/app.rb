require 'rack'

app = Proc.new do |env|
  name = ENV.fetch('NAME', __FILE__)
  ['200', {'Content-Type' => 'text/html'}, ["This is %s" % name]]
end

Rack::Handler::WEBrick.run app
