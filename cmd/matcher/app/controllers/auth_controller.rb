class AuthController < ApplicationController
  def callback
    auth_info = request.env["omniauth.auth"]
    Rails.logger.debug("login info: #{auth_info}")
    session[:user_id] = auth_info[:uid]

    redirect_to "/user/show"
  end

  def failure
    @error_msg = request.params["message"]
  end

  def logout
    # TODO
  end
end
