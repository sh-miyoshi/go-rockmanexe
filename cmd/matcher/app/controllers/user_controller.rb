require "securerandom"

class UserController < ApplicationController
  include Login
  before_action :set_login_user, except: [:index, :new, :create]

  # top page
  def index; end

  def show
    # Set session info
    @own_session = Session.find_by(owner_id: @user.user_id)
    @guest_session = Session.where(guest_id: @user.user_id)
  end

  def new
    redirect_to '/' unless session[:user_id].present?
  end

  def create
    return redirect_to "/" unless session[:user_id].present?

    User.create!(
      name: params[:name],
      login_id: session[:user_id],
      user_id: SecureRandom.uuid
    )

    redirect_to user_show_path
  rescue StandardError => e
    Rails.logger.info("Failed to create user: #{e}")
    flash[:danger] = "ユーザーの作成に失敗しました。<br/>#{e}"
    redirect_to user_new_path
  end

  def destroy; end

  private

  def set_login_user
    @user = login_user
  end
end
