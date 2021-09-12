class SessionController < ApplicationController
  include RouterApiRequester
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

    # request to router
    # add owner client
    url = "#{Settings.router.api_addr}/api/v1/client"
    res = router_request(url: url, method: "post")
    raise RouterApiRequester::RequestError, res[:body] unless res[:success]

    Rails.logger.debug("add owner client response: #{res[:body]}")
    session.owner_client_id = res[:body]["id"]
    session.owner_client_key = res[:body]["key"]

    # add guest client
    res = router_request(url: url, method: "post")
    raise RouterApiRequester::RequestError, res[:body] unless res[:success]

    Rails.logger.debug("add guest client response: #{res[:body]}")
    session.guest_client_id = res[:body]["id"]
    session.guest_client_key = res[:body]["key"]

    # add route fot clients
    url = "#{Settings.router.api_addr}/api/v1/route"
    res = router_request(url: url, method: "post", body: { clients: [session.owner_client_id, session.guest_client_id] })
    raise RouterApiRequester::RequestError, res[:body] unless res[:success]

    Rails.logger.debug("add route response: #{res[:body]}")
    session.route_id = res[:body]["id"]

    session.save!

    redirect_to controller: :user, action: :show
  rescue StandardError => e
    Rails.logger.error("Failed to create session: #{e}")
    flash[:danger] = "セッション情報の作成に失敗しました。#{e}"
    redirect_to session_new_path
  end

  def new; end

  def destroy
    # TODO: request to router
    # delete route
    # delete clients

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
