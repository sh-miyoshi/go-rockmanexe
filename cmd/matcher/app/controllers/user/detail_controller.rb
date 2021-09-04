class User::DetailController < ApplicationController
  include Login
  before_action :set_login_user

  def show; end

  def update; end

  private

  def set_login_user
    login_id = login_user_id
    @user = User.find_by(login_id: login_id)
  end
end
