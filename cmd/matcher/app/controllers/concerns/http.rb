require "net/http"

module HTTP
  extend ActiveSupport::Concern

  class RequestError < StandardError
  end

  def http_request(url:, method:, body: nil)
    uri = URI.parse(url)
    http = Net::HTTP.new(uri.host, uri.port)
    http.use_ssl = uri.scheme === "https"

    headers = { "Content-Type" => "application/json" } if body.present?
    case method.downcase
    when "get"
      res = http.get(uri.path, body&.to_json, headers)
    when "post"
      res = http.post(uri.path, body&.to_json, headers)
    else
      raise RequestError, "invalid method #{method} was specified"
    end

    if 200 <= res.code && res.code < 300
      {
        code: res.code,
        body: JSON.parse(res.body)
      }
    else
      {
        code: res.code,
        body: res.body
      }
    end
  end
end
