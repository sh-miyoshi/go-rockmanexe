class SessionsController < ApplicationController
  def new
  end

  def create
    guest = User.find_by(user_id: create_params[:guest_id])
    if guest.nil?
      Rails.logger.info("No such user: #{create_params[:guest_id]}")
      flash[:danger] = "ユーザー #{create_params[:guest_id]} は存在しません。"
      return redirect_to new_session_path
    end

    ActiveRecord::Base.transaction do
      owner_client_id = SecureRandom.uuid
      guest_client_id = SecureRandom.uuid

      session =
        Session.new(
          session_id: SecureRandom.uuid,
          router_addr: Settings.router.data_addr,
          owner_id: @current_user.user_id,
          guest_id: guest.user_id,
          name: create_params[:name],
          owner_client_id:,
          guest_client_id:
        )

      Client.create!(
        client_id: owner_client_id,
        client_key: SecureRandom.uuid,
        session_id: session.session_id
      )

      Client.create!(
        client_id: guest_client_id,
        client_key: SecureRandom.uuid,
        session_id: session.session_id
      )

      session.save!
      redirect_to user_path(@current_user.user_id)
    rescue StandardError => e
      Rails.logger.info("Failed to create session: #{e}")
      flash[:danger] = "セッションの作成に失敗しました。<br/>#{e}"
      redirect_to new_session_path
    end
  end

  def destroy
    s = Session.find_by(session_id: params[:session_id])
    begin
      s.destroy_with_clients!
    rescue StandardError => e
      Rails.logger.info("Failed to create session: #{e}")
      flash[:danger] = "セッションの削除に失敗しました。<br/>#{e}"
    end
    redirect_to user_path(@current_user.user_id)
  end

  private

  def create_params
    params.permit(:name, :guest_id)
  end
end
