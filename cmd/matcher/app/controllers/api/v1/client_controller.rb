class Api::V1::ClientController < ApplicationController
  include ApiHelper
  skip_before_action :verify_authenticity_token

  def auth
    # TODO: request auth

    begin
      req = JSON.parse(request.body.read, { symbolize_names: true })
    rescue StandardError
      return response_bad_request
    end
    client_id = req[:client_id]
    client_key = req[:client_key]

    client = Client.find_by(client_id: client_id)
    return response_bad_request if client.nil?
    return response_bad_request if client_key != client.client_key

    render json: { session_id: client.session_id }
  rescue StandardError => e
    Rails.logger.error(e)
    response_internal_server_error
  end
end
