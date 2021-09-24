class Api::V1::SessionController < ApplicationController
  def show
    # TODO request auth

    session = Session.find_by(session_id: params[:session_id])
    return response_not_found(class_name: "session") if session.nil?

    render json: {
      id: session.session_id,
      owner_user_id: session.owner_id,
      owner_client_id: session.owner_client_id,
      guest_user_id: session.guest_id,
      guest_client_id: session.guest_client_id,
    }
  end
end
