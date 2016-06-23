def webmockProxySetup
  # Download goproxy root certificate https://raw.githubusercontent.com/elazarl/goproxy/master/ca.pem
  ca = File.expand_path("../ca.pem", File.dirname(__FILE__))
  ENV['SSL_CERT_FILE'] = ca
  ENV['http_proxy'] = 'http://localhost:8080'
end

RSpec.configure do |config|
  config.expect_with :rspec do |expectations|
    expectations.include_chain_clauses_in_custom_matcher_descriptions = true
  end

  config.mock_with :rspec do |mocks|
    mocks.verify_partial_doubles = true
  end
end

webmockProxySetup()
