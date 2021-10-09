require "securerandom"

class UserController < ApplicationController
  include Login
  before_action :set_login_user, except: [:index, :new, :create]

  # top page
  def index; end

  def show
    # Set session info
    @own_session = Session.find_by(owner_id: @user.user_id)
    @own_session.guest_name = user_name(@own_session.guest_id) if @own_session.present?
    @guest_sessions = Session.where(guest_id: @user.user_id)
    @guest_sessions&.each do |s|
      s.owner_name = user_name(s.owner_id)
    end
  end

  def new
    redirect_to "/" unless session[:user_id].present?
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

  def user_name(id)
    user = User.find_by(user_id: id)
    user.name
  end
end
