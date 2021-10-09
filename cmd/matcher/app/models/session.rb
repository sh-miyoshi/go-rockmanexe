class Session < ApplicationRecord
  attr_writer :guest_name, :owner_name

  validates :session_name, presence: true
  validates :owner_id, presence: true
  validates :guest_id, presence: true

  def guest_name
    @guest_name
  end

  def owner_name
    @owner_name
  end
end
