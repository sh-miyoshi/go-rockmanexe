class UsersController < ApplicationController
  skip_before_action :set_login_user, only: %i[index new create]
  before_action :check_access, only: %i[show destroy edit update]

  def index
  end

  def show
    @own_session = Session.find_by(owner_id: @current_user.user_id)
    @guest_sessions = Session.where(guest_id: @current_user.user_id)
  end

  def new
    redirect_to "/" unless session[:user_id].present?
  end

  def create
    return redirect_to "/" unless session[:user_id].present?

    user_id = SecureRandom.uuid
    User.create!(
      user_id: user_id,
      name: create_params[:name],
      login_id: session[:user_id]
    )

    redirect_to user_path(user_id)
  rescue StandardError => e
    Rails.logger.info("Failed to create user: #{e}")
    flash[:danger] = "ユーザーの作成に失敗しました。<br/>#{e}"
    redirect_to new_user_path
  end

  def destroy
  end

  def edit
  end

  def update
    msg = @current_user.update(update_params) ? "更新できました" : "更新に失敗しました"
    redirect_to edit_user_path, notice: msg
  end

  private

  def create_params
    params.permit(:name)
  end

  def update_params
    params.require(:user).permit(:user_id, :name)
  end

  def check_access
    redirect_to "/" if @current_user.user_id != params[:id]
  end
end
