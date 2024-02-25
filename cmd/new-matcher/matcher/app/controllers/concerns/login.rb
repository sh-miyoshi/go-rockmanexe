module Login
  extend ActiveSupport::Concern

  private

  def login_user
    # redirect top page if not logged in
    return redirect_to "/" unless session[:user_id].present?

    user = User.find_by(login_id: session[:user_id])

    # redirect user create page if not in matching server
    return redirect_to "/users/new" if user.nil?

    user
  end
end
