require 'net/http'
require 'json'

# Helper module to maintain logic to access the Profile Services
module RecoEngineHelper

  # Get recommended valence helper, used to call the /getRecommendedValence present in RecommendationEngine
  def get_recommended_valence_helper(url)
    begin
      host_details = zookeeper_helper(url)
      host = host_details[0]
      port = host_details[1].to_s
      uri = URI('http://'+host+':'+port+'/getRecommendedValence')
      http = Net::HTTP.new(host, port)
      req = Net::HTTP::Get.new(uri.path, {'Content-Type' =>'application/json',
        'Authorization' => request.headers[:Authorization]})
      http.request(req)
    rescue => e
      raise e
    end
  end

  # Zookeeper handler to retrieve the auth services host and port
  def zookeeper_helper(url)
    # z = Zookeeper.new("149.165.170.7:2181")
    # host_details= z.get(:path => url)
    # host_details=host_details[:data]
    # host_details=host_details.split(":")
    host_details=[ENV["RECOENGINE_SERVICE_PORT_8001_TCP_ADDR"], ENV["RECOENGINE_SERVICE_PORT_8001_TCP_PORT"]]
    host_details
  end

end
