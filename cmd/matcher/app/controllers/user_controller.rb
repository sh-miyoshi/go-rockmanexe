class UserController < ApplicationController
  include Login
  before_action :set_login_user, except: :index

  # top page
  def index; end

  def show; end

  def create; end

  def destroy; end

  private

  def set_login_user
    login_id = login_user_id
    @user = User.find_by(login_id: login_id)
  end
end
