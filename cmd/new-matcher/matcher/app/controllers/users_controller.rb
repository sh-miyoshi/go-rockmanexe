class UsersController < ApplicationController
  skip_before_action :set_login_user, only: %i[index new create]

  def index
  end

  def show
  end

  def new
  end

  def create
  end

  def destroy
  end
end
