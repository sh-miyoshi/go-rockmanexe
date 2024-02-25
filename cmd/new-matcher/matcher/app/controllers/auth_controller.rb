class AuthController < ApplicationController
  skip_before_action :set_login_user

  def callback
    auth_info = request.env["omniauth.auth"]
    Rails.logger.debug("login info: #{auth_info}")
    session[:user_id] = auth_info[:uid]

    redirect_to "/users/show"
  end

  def failure
    @error_msg = request.params["message"]
  end

  def logout
    reset_session
    redirect_to logout_url
  end

  private

  AUTH0_CONFIG = Rails.application.config.auth0

  def logout_url
    request_params = {
      returnTo: root_url,
      client_id: AUTH0_CONFIG["auth0_client_id"]
    }

    URI::HTTPS.build(
      host: AUTH0_CONFIG["auth0_domain"],
      path: "/v2/logout",
      query: request_params.to_query
    ).to_s
  end
end
