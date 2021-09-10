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
        owner_id: @user.user_id,
        guest_id: params[:guest_id],
        expires_at: Time.current.since(30.minutes),
        router_addr: Settings.router.data_addr
      }
    )

    if session.invalid?
      Rails.logger.info("Invalid request was specified: #{session.errors.messages}")
      flash[:danger] = "不正な値が入力されました。#{session.errors.messages}"
      return redirect_to session_new_path
    end

    # TODO: request to router
    # add owner client
    # add guest client
    # add route fot clients

    session.owner_client_id = "tester1"
    session.owner_client_key = "testtest"
    session.guest_client_id = "tester2"
    session.guest_client_key = "testtest"
    session.route_id = "nagnagnklrbhm"

    session.save!

    redirect_to controller: :user, action: :show
  rescue => e
      Rails.logger.error("Failed to create session: #{e}")
      flash[:danger] = "セッション情報の作成に失敗しました。#{e}"
      return redirect_to session_new_path
  end

  def new; end

  def destroy; end

  private

  def set_login_user
    @user = login_user
  end
end
