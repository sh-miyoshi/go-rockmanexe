class ApplicationController < ActionController::Base
  include Login

  before_action :set_login_user

  attr_reader :current_user

  private

  def set_login_user
    @current_user = login_user
  end
end
