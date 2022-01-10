class User::DetailController < ApplicationController
  include Login
  before_action :set_login_user

  def edit; end

  def update
    @user.name = params[:user][:name]
    msg = @user.save ? "User was successfully updated." : "Failed to update user."
    redirect_to user_detail_edit_path, notice: msg
  end

  private

  def set_login_user
    @user = login_user
  end
end
