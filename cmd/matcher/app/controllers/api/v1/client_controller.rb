class Api::V1::ClientController < ApplicationController
  include ApiHelper

  def auth
    # TODO request auth

    render json: {val:"test"}
  end
end
