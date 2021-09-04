class UserController < ApplicationController
  include Login
  before_action :set_login_user, except: :index

  # top page
  def index; end

  def show; end

  def new
    
  end

  def create; end

  def destroy; end

  private

  def set_login_user
    @user = login_user()
  end
end
