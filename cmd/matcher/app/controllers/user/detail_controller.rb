class User::DetailController < ApplicationController
  include Login
  before_action :set_login_user

  def edit; end

  def update; end

  private

  def set_login_user
    @user = login_user
  end
end
