class UsersController < ApplicationController
  skip_before_action :set_login_user, only: %i[index new create]

  def index
  end

  def show
  end

  def new
    redirect_to "/" unless session[:user_id].present?
  end

  def create
    return redirect_to "/" unless session[:user_id].present?

    user_id = SecureRandom.uuid
    User.create!(
      user_id: user_id,
      name: params[:name],
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

  private

  def create_params
    params.permit(:name)
  end
end
