class SessionController < ApplicationController
  def create
    client_id = "tester1"
    client_key = "testtest"

    # TODO: request to router
    Session.create({
                     owner: client_id,
                     router_addr: "localhost:16283",
                     client_id: client_id,
                     client_key: client_key,
                     expires_at: Time.current.since(30.minutes)
                   })

    redirect_to controller: :user, action: :show
  end

  def destroy; end
end
