class AuthController < ApplicationController
  def callback
    auth_info = request.env['omniauth.auth']
    session[:userinfo] = auth_info['extra']['raw_info']
    Rails.logger.debug("login info: #{auth_info}")

    redirect_to '/user/show'
  end

  def failure
    @error_msg = request.params['message']
  end

  def logout
    # TODO
  end
end
