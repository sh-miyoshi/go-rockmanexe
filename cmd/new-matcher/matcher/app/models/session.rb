class Session < ApplicationRecord
  validates :name, length: { in: 3..20 }
  validates :owner_id, presence: true
  validates :guest_id, presence: true

  def owner_name
    User.find_by(user_id: owner_id)&.name
  end

  def owner_client_key
    Client.find_by(client_id: owner_client_id)&.client_key
  end

  def guest_name
    User.find_by(user_id: guest_id)&.name
  end

  def guest_client_key
    Client.find_by(client_id: guest_client_id)&.client_key
  end

  def destroy_with_clients!
    ActiveRecord::Base.transaction do
      Client.find_by(client_id: owner_client_id)&.destroy!
      Client.find_by(client_id: guest_client_id)&.destroy!
      destroy!
    end
  end
end
