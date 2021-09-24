class SessionController < ApplicationController
  include Login
  before_action :set_login_user

  def create
    if params[:guest_id].blank? || User.where(user_id: params[:guest_id]).empty?
      Rails.logger.info("No such user: #{params[:guest_id]}")
      flash[:danger] = "ユーザー #{params[:guest_id]} は存在しません。"
      return redirect_to session_new_path
    end

    session = Session.new(
      {
        session_name: params[:name],
        router_addr: Settings.router.data_addr,
        owner_id: @user.user_id,
        owner_client_id: SecureRandom.uuid,
        owner_client_key: SecureRandom.uuid,
        guest_id: params[:guest_id],
        guest_client_id: SecureRandom.uuid,
        guest_client_key: SecureRandom.uuid,
        session_id: SecureRandom.uuid,
      }
    )

    Client.create(
      {
        client_id: session.owner_client_id,
        session_id: session.id,
      }
    )
    Client.create(
      {
        client_id: session.guest_client_id,
        session_id: session.id,
      }
    )

    session.save!

    redirect_to controller: :user, action: :show
  rescue StandardError => e
    Rails.logger.error("Failed to create session: #{e}")
    flash[:danger] = "セッション情報の作成に失敗しました。#{e}"
    redirect_to session_new_path
  end

  def new; end

  def destroy
    Session.destroy(params["id"])

    redirect_to user_show_path
  rescue StandardError => e
    Rails.logger.error("Failed to delete session: #{e}")
    flash[:danger] = "セッション削除に失敗しました。#{e}"
    redirect_to user_show_path
  end

  private

  def set_login_user
    @user = login_user
  end
end
