module Login
  extend ActiveSupport::Concern

  def login_user
    # redirect top page if not logged in
    # redirect_to '/' unless session[:user_id].present?

    # Userがいなければ作成ページにリダイレクト

    # debug
    User.find_by(login_id: "login_id_1")
  end
end
