require 'openssl'
require './github'

client = Github::Client.new
p client.apis
