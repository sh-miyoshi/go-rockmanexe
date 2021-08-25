class UserController < ApplicationController
  before_action :set_login_user, except: :index

  def index; end

  def show; end

  def create; end

  def update; end

  def destroy; end

  private

  def set_login_user
    login_id = "login_id_1" # debug
    @user = User.find_by(login_id: login_id)
  end
end
