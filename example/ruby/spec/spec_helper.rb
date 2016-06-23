ca = File.expand_path("./ca.pem", File.dirname(__FILE__))
ENV['SSL_CERT_FILE'] = ca
require 'open-uri'
require 'openssl'

def goproxyCert
  url = "https://raw.githubusercontent.com/elazarl/goproxy/master/ca.pem"
  dir = File.expand_path(File.dirname(__FILE__))
  path = dir + '/ca.pem'
  open(path, 'wb') do |f|
    open(url, "r", {:ssl_verify_mode => OpenSSL::SSL::VERIFY_NONE}) do |data|
      f.write(data.read)
    end
  end
end

if !File.exists?(ca)
  goproxyCert()
  puts "Downloaded goproxy using ca certificate file."
  puts "Please retry."
  exit 0
end

require_relative '../sample'
ENV['http_proxy'] = 'http://localhost:8080'

RSpec.configure do |config|
  config.expect_with :rspec do |expectations|
    expectations.include_chain_clauses_in_custom_matcher_descriptions = true
  end

  config.mock_with :rspec do |mocks|
    mocks.verify_partial_doubles = true
  end
end
