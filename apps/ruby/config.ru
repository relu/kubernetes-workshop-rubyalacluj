require 'rack'

run do |env|
  name = ENV.fetch('NAME', __FILE__)
  [200, {'content-type' => 'text/html'}, ["This is %s" % name]]
end
