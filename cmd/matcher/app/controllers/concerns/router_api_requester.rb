require "net/http"

module RouterApiRequester
  extend ActiveSupport::Concern

  class RequestError < StandardError
  end

  private

  def router_request(url:, method:, body: nil)
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

    if 200 <= res.code.to_i && res.code.to_i < 300
      {
        success: true,
        code: res.code,
        body: JSON.parse(res.body)
      }
    else
      {
        success: false,
        code: res.code,
        body: res.body
      }
    end
  end
end
